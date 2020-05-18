package db

import "github.com/jinzhu/gorm"

// DB instance
var DB *gorm.DB

// SequenceResult select nextval result
type SequenceResult struct {
	Nextval string
}
