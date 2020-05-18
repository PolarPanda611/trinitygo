package model

import (
	"time"

	"github.com/PolarPanda611/trinitygo/crud/model"
)

// RunHistory  metro 运行历史
type RunHistory struct {
	model.Model
	LineID     int64      `json:"line_id"`
	Line       *Line      `json:"line"`
	MetroID    int64      `json:"metro_id"`
	Metro      *Metro     `json:"metro"`
	StationID  int64      `json:"station_id"`
	Station    *Station   `json:"station"`
	RecorderID int64      `json:"recorder_id"`
	Recorder   *User      `json:"recorder"`
	In         *time.Time `json:"in"`
	Out        *time.Time `json:"Out"`
	Code       string     `gorm:"-"`
}
