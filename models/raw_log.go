package models

import (
	"net"
	"net/http"
	"encoding/json"
	"github.com/markbates/pop/nulls"
)

type RawLog struct {
	ID       int          `json:"id" db:"id"`
	DeviceID int          `json:"device_id" db:"device_id"`
	Ip       uint32       `json:"ip" db:"ip"`
	Log      string       `json:"log" db:"log"`
	Error    nulls.String `json:"error" db:"error"`
	Time     nulls.Time   `json:"time" db:"time"`
	HasError bool         `json:"has_error" db:"hasError"`
}

func (r RawLog) CreateFromRequest(req *http.Request, device Device) error {
	row := map[string]interface{}{
		"POST": req.PostForm,
	}

	if files := req.MultipartForm; len(files.File) > 0 {
		row["FILES"] = files.File
	}

	var ip net.IP
	if ip = net.ParseIP(req.RemoteAddr); ip == nil {
		ip = net.ParseIP("89.28.82.88")
	}

	data, _ := json.Marshal(row)

	model := &RawLog{
		DeviceID: int(device.ID),
		Ip:       device.Ip2Int(ip),
		Log:      string(data),
	}

	return DB.Create(model)
}

// String is not required by pop and may be deleted
func (r RawLog) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// RawLogs is not required by pop and may be deleted
type RawLogs []RawLog

// String is not required by pop and may be deleted
func (r RawLogs) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}
