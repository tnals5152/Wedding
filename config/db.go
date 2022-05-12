package config

import (
	"fmt"
	"os"
	"wedding/models"
	"wedding/utils"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GORMDB *gorm.DB

//DB연결 및 DB table migrate하는 함수(초기 메인에서 실행하여 테이블 변경 값을 실행시켜준다)
func ConnectDB() {
	var err error
	//사용자:비밀번호@tcp(ipAddress:port)
	dsnSet := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWD"),
		os.Getenv("GORMDB_HOST"), os.Getenv("GORMDB_PORT"),
		os.Getenv("DB_NAME"))
	GORMDB, err = gorm.Open(mysql.Open(dsnSet), &gorm.Config{
		// DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.Default.LogMode(logger.Info),
	})
	utils.IfErrorMakePanic(err, "can not connect GORMDB")

	migrateAllTable()
}

//여러 테이블 있을 때 migrate할 수 있게 함수로 따로 작성
func migrateAllTable() {
	GORMDB.AutoMigrate(&models.Comment{})
}
