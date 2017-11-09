package actions

import (
	"io"
	"os"
	"fmt"
	"time"
	"errors"
	"regexp"
	"strings"
	"io/ioutil"
	"recipes/models"
	"recipes/helpers"
	"golang.org/x/net/html"
	"github.com/gobuffalo/buffalo"
)

func UploadHandler(c buffalo.Context) error {
	dsn := c.Request().FormValue("device")

	// fetch device
	var device models.Device
	device, err := device.ByDSN(dsn)

	if err != nil {
		msg := fmt.Sprintf("BAD_DEV. Device with dsn '%v' not found.", dsn)

		return errors.New(msg)
	}

	// prepare raw log
	var rawLog models.RawLog
	if err := rawLog.CreateFromRequest(c.Request(), device); err != nil {
		return err
	}

	serverTime := time.Now()
	clientTime := c.Request().FormValue("client-date-time")

	pattern := "^\\d{2}/\\d{2}/\\d{4} \\d{2}:\\d{2}:\\d{2}$"
	match, _ := regexp.MatchString(pattern, clientTime)
	parsedTime, err := time.Parse("02/01/2006 15:04:05", clientTime)

	if match == false || err != nil {
		msg := fmt.Sprintf("Client date time has invalid format. Expected format: dd/MM/YYYY HH:mm:ss.")

		return errors.New(msg)
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
	if err != nil {
		return err
	}

	defer file.Close()

	// for html logs dailyQuota is bigger then for other types
	if fileHeaders.Header.Get("Content-Type") == "text/html" {
		dailyQuota += 350500
	}

	if usedQuota >= dailyQuota {
		return errors.New("reject - you've reached the daily quota")
	}

	usedQuota += fileHeaders.Size
	device.DailyTraffic = usedQuota

	if err := device.UpdateState(); err != nil {
		return errors.New("could not save device info")
	}

	cwd, _ := os.Getwd()
	targetDir := fmt.Sprintf("%v/tmp/files/device_%v", cwd, device.ID)

	stat, err := os.Stat(targetDir)
	if err != nil || !stat.IsDir() {
		err := os.MkdirAll(targetDir, os.ModePerm)

		if err != nil {
			return err
		}
	}

	// Upload File
	targetFile, err := os.OpenFile(targetDir+"/"+fileHeaders.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer targetFile.Close()
	io.Copy(targetFile, file)

	// @todo: Read only 250 kb
	// Lines: 145-148

	var content string
	fc, err := ioutil.ReadFile(targetFile.Name())
	if err != nil {
		return err
	}

	content = string(fc)

	// @todo: Check & Adjust Member's account
	// Lines: 150-163

	member := new(models.Member)

	var log []helpers.DeltaLog

	if member.AccountName() != "basic" {
		content, err := helpers.CleanContentWithTidy(content)
		if err != nil {
			return err
		}

		content = helpers.StripTags(content)

		dom, err := html.Parse(strings.NewReader(content))

		if err != nil {
			return err
		}

		log = helpers.Parse(dom, device.ID, log)

		st := new(models.LogsDelta)
		st.Persist(log)
	}

	return c.Render(200, r.JSON(map[string]interface{}{
		"log": log,
	}))
}
