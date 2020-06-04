package migrate

import (
	"http/infra/db"

	"http/domain/model"
)

func Migrate() {
	db.DB.AutoMigrate(&model.User{})
}
