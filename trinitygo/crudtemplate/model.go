package crudtemplate

func init() {
	_templates["/domain/model/%v.go"] = genModel()
}

func genModel() string {
	return `
package model

import "github.com/PolarPanda611/trinitygo/crud/model"

// {{.ModelName}} model for {{.ModelName}}
type {{.ModelName}} struct {
	model.Model
	// to add your customize param inside here
	Code string ` + "`" + `json:"code" gorm:"type:varchar(50);index;not null;unique"` + "`" + `
}
	
`
}
