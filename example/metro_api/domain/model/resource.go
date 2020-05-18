
package model

import "github.com/PolarPanda611/trinitygo/crud/model"

// Resource model for Resource
type Resource struct {
	model.Model
	// to add your customize param inside here
	Code string `json:"code" gorm:"type:varchar(50);index;not null;unique"`
}
	
