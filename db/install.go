package db

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/PolarPanda611/trinitygo/util"

	"github.com/jinzhu/gorm"

	_ "github.com/lib/pq" // pg
	uuid "github.com/satori/go.uuid"
)

// NilLogger nil logger
type NilLogger struct{}

// Print nil logger do noothing
func (l *NilLogger) Print(v ...interface{}) {

}

// DefaultInstallGORM default install gorm
func DefaultInstallGORM(
	debug bool,
	singular bool,
	dbType string,
	tablePrefix string,
	server string,
	maxIdleConn int,
	maxOpenConn int,
) *gorm.DB {
	db, err := gorm.Open(dbType, server)
	if err != nil {
		log.Fatal("db connect build failed")
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return fmt.Sprintf("%v%v", tablePrefix, defaultTableName)
	}
	db.SetLogger(&NilLogger{})
	db.LogMode(debug)
	db.SingularTable(singular)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampAndUUIDForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.DB().SetMaxIdleConns(maxIdleConn)
	db.DB().SetMaxOpenConns(maxOpenConn)
	return db
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampAndUUIDForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		userIDInterface, _ := scope.Get("user_id")
		userID, _ := userIDInterface.(int64)
		nowTime := util.GetCurrentTime()
		if createTimeField, ok := scope.FieldByName("CreatedTime"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}
		if createUserIDField, ok := scope.FieldByName("CreateUserID"); ok {
			if createUserIDField.IsBlank {
				createUserIDField.Set(userID)
			}
		}
		if idField, ok := scope.FieldByName("ID"); ok {
			idField.Set(util.GenerateSnowFlakeID(int64(rand.Intn(100))))
		}
		if modifyTimeField, ok := scope.FieldByName("UpdatedTime"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
		if updateUserIDField, ok := scope.FieldByName("UpdateUserID"); ok {
			if updateUserIDField.IsBlank {
				updateUserIDField.Set(userID)
			}
		}

		if updateDVersionField, ok := scope.FieldByName("DVersion"); ok {
			if updateDVersionField.IsBlank {
				updateDVersionField.Set(uuid.NewV4().String())
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		userID, _ := scope.Get("user_id")
		if attrs, ok := scope.InstanceGet("gorm:update_attrs"); ok {
			updateAttrs, _ := attrs.(map[string]interface{})
			updateAttrs["updated_time"] = util.GetCurrentTime()
			updateAttrs["update_user_id"] = userID
			updateAttrs["d_version"] = uuid.NewV4().String()
			scope.InstanceSet("gorm:update_attrs", updateAttrs)
		}
	}

}

// deleteCallback will set `DeletedOn` where deleting
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		userID, ok := scope.Get("user_id")
		if !ok {
			userID = nil
		}
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}
		deletedAtField, hasDeletedAtField := scope.FieldByName("deleted_time")
		deleteUserIDField, hasDeleteUserIDField := scope.FieldByName("DeleteUserID")
		dVersionField, hasDVersionField := scope.FieldByName("d_version")

		if !scope.Search.Unscoped && hasDeletedAtField && hasDVersionField && hasDeleteUserIDField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v,%v=%v,%v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedAtField.DBName),
				scope.AddToVars(util.GetCurrentTime()),
				scope.Quote(deleteUserIDField.DBName),
				scope.AddToVars(userID),
				scope.Quote(dVersionField.DBName),
				scope.AddToVars(uuid.NewV4().String()),
				util.AddExtraSpaceIfExist(scope.CombinedConditionSql()),
				util.AddExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				util.AddExtraSpaceIfExist(scope.CombinedConditionSql()),
				util.AddExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}
