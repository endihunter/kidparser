package models

import (
	"encoding/json"
	"time"
	"net"
	"encoding/binary"
	"github.com/markbates/pop/nulls"
	"math"
	"os"
	"fmt"
	"errors"
	"strings"
)

type Device struct {
	ID           int          `json:"id" db:"id"`
	Dsn          string       `json:"dsn" db:"dsn"`
	SectionID    nulls.Int    `json:"section_id" db:"wcSection_id"`
	MemberID     int          `json:"member_id" db:"member_id"`
	Name         string       `json:"name" db:"name"`
	LogDays      nulls.Int    `json:"log_days" db:"log_days"`
	Quota        nulls.Int64  `json:"quota" db:"quota"`
	DailyQuota   nulls.Int64  `json:"daily_quota" db:"daily_quota"`
	Gmt          nulls.Int    `json:"gmt" db:"gmt"`
	ClientVer    nulls.String `json:"client_ver" db:"clientVer"`
	DailyTraffic int64        `json:"daily_traffic" db:"trafficToday"`
	Time         time.Time    `json:"time" db:"time"`
	FirstLogDate nulls.Time   `json:"first_log" db:"firstLogDate"`
	LastLogDate  nulls.Time   `json:"last_log" db:"lastLogDate"`
}

func (d Device) UpdateState() error {
	toSave := &Device{
		ID:           d.ID,
		ClientVer:    d.ClientVer,
		FirstLogDate: d.FirstLogDate,
		LastLogDate:  d.LastLogDate,
		Gmt:          d.Gmt,
		DailyTraffic: d.DailyTraffic,
	}

	excluded := []string{"dsn", "wcSection_id", "member_id", "name", "time", "quota", "daily_quota"}

	return DB.Update(toSave, excluded...)
}

// Calculate the difference between Client & Server times
func (d Device) SetGMT(ct time.Time, st time.Time) {
	diff := math.Ceil(float64((ct.Unix() - st.Unix()) / 3600))
	d.Gmt = nulls.NewInt(int(diff))
}

func (d *Device) SetClientVersion(v string) {
	if !d.ClientVer.Valid {
		d.ClientVer = nulls.NewString(v)
	}
}

func (d *Device) SetLogDates() {
	if !d.FirstLogDate.Valid {
		d.FirstLogDate = nulls.NewTime(time.Now())
	}
	d.LastLogDate = nulls.NewTime(time.Now())
}

func (d *Device) ByDSN(dsn string) (item Device, err error) {
	query := DB.Where("dsn=?", dsn)
	err = query.First(&item)

	return
}

func (d *Device) Ip2Int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func (d *Device) Int2Ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

func (d *Device) HandleStorage(ch chan interface{}) {
	defer close(ch)

	cwd, _ := os.Getwd()
	targetDir := fmt.Sprintf("%v/storage/files/device_%v", cwd, d.ID)
	stat, err := os.Stat(targetDir)
	if err != nil || !stat.IsDir() {
		err := os.MkdirAll(targetDir, os.ModePerm)

		if err != nil {
			ch <- errors.New(fmt.Sprintf("could not create directory \"%v\"", strings.Replace(targetDir, cwd, "", -1)))
			return
		}
	}

	ch <- targetDir
}

// String is not required by pop and may be deleted
func (d Device) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

func (d Device) Member() (member Member, err error) {
	return Member{}, nil
}

func (d Device) GetSettings() (settings DeviceSettings, err error) {
	// This Fucking Query does not return a result,
	// so temporary we'll stick the settings to a fixed one
	// @todo: Fix the query
	settings = DeviceSettings{
		ID:           1,
		NumDevices:   1,
		LogDays:      9,
		Quota:        9437184,
		DailyQuota:   9437184,
		AccountID:    1,
		AccountName:  "Basic",
		AccountTitle: "basic",
		Reports:      0,
	}

	return settings, nil

	//query := DB.RawQuery("SELECT "+
	//	"`d`.`id`, "+
	//	"IF(`m`.`num_devices`, `m`.`num_devices`, `a`.`num_devices`) AS `num_devices`, "+
	//	"IF(`m`.`log_days`, `m`.`log_days`, `a`.`log_days`) AS `log_days`, "+
	//	"(IF(`m`.`quota`, `m`.`quota`, `a`.`quota`) * POW(1024, 2)) AS `quota`, "+
	//	"(IF(`m`.`daily_quota`, `m`.`daily_quota`, `a`.`daily_quota`) * POW(1024, 2)) AS `daily_quota`, "+
	//	"`ma`.`id` AS `macc_id`, IFNULL(`ma`.`account_id`, 1) AS `account_id`, `ma`.`tariff_id` AS `tariff_id`, `ma`.`start_date` AS `start_date`, "+
	//	"`ma`.`end_date` AS `end_date`, IF(`ma`.`end_date`, DATEDIFF(`ma`.`end_date`, NOW()), NULL) AS `day_remains`, "+
	//	"`a`.`name` AS `name`, `a`.`name_title` AS `name_title`, `a`.`reports` AS `reports` "+
	//	"FROM `devices` AS `d` "+
	//	"LEFT JOIN `members` AS `m` ON `d`.`member_id`=`m`.`id` "+
	//	"LEFT JOIN `memberAccounts` AS `ma` ON `ma`.`member_id`=`m`.`id` AND (SELECT `id` FROM `memberAccounts` WHERE `member_id` = `m`.`id` ORDER BY `id` DESC LIMIT 1)=`ma`.`id` "+
	//	"INNER JOIN `accounts` AS `a` ON `a`.`id` = IFNULL(`ma`.`account_id`, 1) "+
	//	"WHERE `d`.`id`=? LIMIT 1",
	//	d.ID,
	//)
	//
	//err = query.All(&settings)
	//
	//return
}

// Devices is not required by pop and may be deleted
type Devices []Device

// String is not required by pop and may be deleted
func (d Devices) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}
