package actions

import (
	"os"
	"fmt"
	"time"
	"errors"
	"regexp"
	"recipes/models"
	"github.com/gobuffalo/buffalo"
	"mime/multipart"
	"recipes/helpers"
	"strings"
	"golang.org/x/net/html"
)

func UploadHandler(c buffalo.Context) error {
	dsn := c.Request().FormValue("device")

	// fetch device
	var device models.Device
	device, err := device.ByDSN(dsn)

	if err != nil {
		msg := fmt.Sprintf("BAD_DEV. Device with dsn '%v' not found.", dsn)

		return stringError(c, errors.New(msg))
	}

	// prepare raw log
	var rawLog models.RawLog
	if err := rawLog.CreateFromRequest(c.Request(), device); err != nil {
		return stringError(c, err)
	}

	serverTime := time.Now()
	clientTime := c.Request().FormValue("client-date-time")

	pattern := "^\\d{2}/\\d{2}/\\d{4} \\d{2}:\\d{2}:\\d{2}$"
	match, _ := regexp.MatchString(pattern, clientTime)
	parsedTime, err := time.Parse("02/01/2006 15:04:05", clientTime)

	if match == false || err != nil {
		msg := fmt.Sprintf("Client date time has invalid format. Expected format: dd/MM/YYYY HH:mm:ss.")

		return stringError(c, errors.New(msg))
	}

	deviceSettings, err := device.GetSettings()

	// set correct Device GMT
	if !device.Gmt.Valid {
		device.SetGMT(parsedTime, serverTime)
	}

	// save Device client version
	device.SetClientVersion(
		c.Request().FormValue("client-ver"),
	)

	// update device log dates
	device.SetLogDates()

	// reset daily traffic if today != lastLogDate
	ty, tm, td := serverTime.Date()
	ly, lm, ld := device.LastLogDate.Time.Date()

	if ty != ly || tm != lm || td != ld {
		device.DailyTraffic = 0
	}

	dailyQuota := deviceSettings.DailyQuota
	usedQuota := device.DailyTraffic

	c.Request().ParseMultipartForm(102400)
	file, fileHeaders, err := c.Request().FormFile("file")
	defer file.Close()

	if err != nil {
		return stringError(c, err)
	}

	// for html logs dailyQuota is bigger then for other types
	if fileHeaders.Header.Get("Content-Type") == "text/html" {
		dailyQuota += 350500
	}

	if usedQuota >= dailyQuota {
		return stringError(c, errors.New("REJECT - you've reached the daily quota"))
	}

	usedQuota += fileHeaders.Size
	device.DailyTraffic = usedQuota

	if err := device.UpdateState(); err != nil {
		return stringError(c, errors.New("could not save device info"))
	}

	target, err := device.HandleStorage()

	// return Human error if it fails
	if err != nil {
		return stringError(c, err)
	}

	content, err := storeUploadedFile(target, fileHeaders, file)

	// return Human error if it fails
	if err != nil {
		return stringError(c, err)
	}

	// @todo: Check & Adjust Member's account
	// Lines: 150-163

	member := new(models.Member)

	if member.AccountName() != "basic" {
		chl, ech := make(chan []helpers.DeltaLog), make(chan error)
		go parseHtmlFile(content, device, chl, ech)
		log, err := <-chl, <-ech

		if err != nil {
			return stringError(c, err)
		}

		st := new(models.LogsDelta)
		st.Persist(log)
	}

	return c.Render(200, r.String("OK"))
}

func stringError(c buffalo.Context, err error) error {
	return c.Render(500, r.String(err.Error()))
}

func parseHtmlFile(content string, device models.Device, chl chan []helpers.DeltaLog, ech chan error) {
	defer close(chl)
	defer close(ech)

	var log []helpers.DeltaLog

	content, err := helpers.CleanContentWithTidy(content)
	if err != nil {
		ech <- err
		return
	}

	content = helpers.StripTags(content)
	dom, err := html.Parse(strings.NewReader(content))

	if err != nil {
		ech <- err
		return
	}

	chl <- helpers.Parse(dom, device.ID, log)
}

func storeUploadedFile(targetDir string, fileHeaders *multipart.FileHeader, file multipart.File) (string, error) {
	// Upload File
	tf, err := os.OpenFile(targetDir+"/"+fileHeaders.Filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	defer tf.Close()

	if err != nil {
		return "", err
	}

	// @todo: Read only 250 kb
	// Lines: 145-148
	fc := make([]byte, fileHeaders.Size)

	if _, err = file.Read(fc); err != nil {
		return "", err
	}

	content := string(fc)
	_, err = tf.WriteString(content)

	if err != nil {
		return "", err
	}

	return content, nil
}
