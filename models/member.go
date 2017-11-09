package models

import (
	"github.com/markbates/pop/nulls"
	"time"
)

type Member struct {
	ID                 int         `json:"int" db:"int"`
	Email              string      `json:"email" db:"email"`
	Password           string      `json:"-" db:"password"`
	Name               string      `json:"name" db:"name"`
	Role               string      `json:"-" db:"role"`
	AccountID          nulls.Int   `json:"account_id" db:"account_id"`
	MemberAccountID    nulls.Int   `json:"macc_id" db:"macc_id"`
	TariffID           nulls.Int64 `json:"tariff_id" db:"tariff_id"`
	StartDate          nulls.Time  `json:"start_date" db:"accountStartDate"`
	EndDate            nulls.Time  `json:"end_date" db:"accountEndDate"`
	NumDevices         nulls.Int   `json:"num_devices" db:"num_devices"`
	LogDays            nulls.Int   `json:"log_days" db:"log_days"`
	Quota              int         `json:"quota" db:"quota"`
	DailyQuota         int         `json:"daily_quota" db:"daily_quota"`
	Active             bool        `json:"active" db:"active"`
	LastLogin          nulls.Time  `json:"last_login" db:"last_login"`
	Joined             time.Time   `json:"joined" db:"joined"`
	Session            string      `json:"-" db:"session"`
	UnfinishedPurchase string      `json:"unfinished_purchase" db:"unfinishedPurchase"`
	ResetUrl           string      `json:"reset_url" db:"reset_url"`
	ResetUntil         time.Time   `json:"reset_time" db:"reset_time"`
	LangID             int         `json:"lang_id" db:"lang_id"`
}

func (m *Member) AccountName() string {
	return "premium"
}
