package migrate

import (
	"http/infra/db"

	"http/domain/model"
)

// Migrate migrate func to sync table structure and initial data
func Migrate() {
	db.DB.AutoMigrate(&model.User{})
}
