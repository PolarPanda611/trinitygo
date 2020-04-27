package model

//Changelog model Error
type Changelog struct {
	Logmodel
	Resource    string `json:"resource"  gorm:"type:varchar(50);not null;default:''"`     // resource name
	ResourceKey string `json:"resource_key"  gorm:"type:varchar(50);not null;default:''"` // resource key
	DVersion    string `json:"d_version" gorm:"type:varchar(50);not null;default:''"`     // current change d version
	Type        string `json:"type"  gorm:"type:varchar(50);not null;default:''"`         // operation type create , update , delete
	Column      string `json:"column"  gorm:"type:varchar(50);not null;default:''"`
	OldValue    string `json:"old_value"  gorm:"type:varchar(200);not null;default:''"` // old value before change
	NewValue    string `json:"new_value"  gorm:"type:varchar(200);not null;default:''"` // new value after change
	TraceID     string `json:"trace_id"  gorm:"type:varchar(50);not null;"`
}
