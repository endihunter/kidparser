package models

import (
	"time"
	"recipes/helpers"
	"github.com/markbates/pop"
)

type LogsDelta struct {
	ID       int    `db:"id"`
	DeviceID int    `db:"device_id"`
	LogType  string `db:"log_type"`
	Duration int    `db:"duration"`
	Name     string `db:"name"`
	Title    string `db:"title"`
	Date     string `db:"date"`
	Time     string `db:"time"`
}

func (l LogsDelta) Import(log helpers.DeltaLog) LogsDelta {
	return LogsDelta{
		DeviceID: log["device_id"].(int),
		LogType:  log["log_type"].(string),
		Duration: int(log["duration"].(float64)),
		Name:     log["name"].(string),
		Title:    log["title"].(string),
		Date:     log["time"].(time.Time).Format("2006-01-02"),
		Time:     log["time"].(time.Time).Format("15:04:05"),
	}
}

func (l *LogsDelta) Persist(logs []helpers.DeltaLog) error {
	return DB.Transaction(func(tx *pop.Connection) error {
		for _, v := range logs {
			model := l.Import(v)
			DB.Create(&model)
		}

		return nil
	})
}
