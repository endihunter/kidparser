package models

import "github.com/markbates/pop/nulls"

type DeviceSettings struct {
	ID           int        `json:"id" db:"id"`
	NumDevices   int        `json:"num_devices" db:"num_devices"`
	LogDays      int        `json:"log_days" db:"log_days"`
	Quota        int64      `json:"quota" db:"quota"`
	DailyQuota   int64      `json:"daily_quota" db:"daily_quota"`
	MaxAccId     nulls.Int  `json:"macc_id" db:"macc_id"`
	AccountID    int        `json:"account_id" db:"account_id"`
	AccountName  string     `json:"account_name" db:"name"`
	AccountTitle string     `json:"account_title" db:"name_title"`
	TariffID     nulls.Int  `json:"tariff_id" db:"tariff_id"`
	StartDate    nulls.Time `json:"start_date" db:"start_date"`
	EndDate      nulls.Time `json:"end_date" db:"end_date"`
	DayRemains   int        `json:"day_remains" db:"day_remains"`
	Reports      int        `json:"reports" db:"reports"`
}
