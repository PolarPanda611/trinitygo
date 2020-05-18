package model

import (
	"github.com/PolarPanda611/trinitygo/crud/model"
)

// Direction 运行方向
type Direction string

const (
	// Up 上行列车
	Up Direction = "Up"
	// Down 下行列车
	Down Direction = "Down"
)

// Station  metro station info
// 地铁线路站点信息
type Station struct {
	model.Model
	Code          string    `json:"code" gorm:"type:varchar(50);index;not null;unique;"`
	Name          string    `json:"name" gorm:"type:varchar(50);"`
	Description   string    `json:"description"`
	NextStationID int64     `json:"next_station_id,string" `
	NextStation   *Station  `json:"next_station"`
	LineID        int64     `json:"line_id,string"`
	Line          *Line     `json:"line"`
	Direction     Direction `json:"direction" gorm:"type:varchar(50);default:'Up';"`
}
