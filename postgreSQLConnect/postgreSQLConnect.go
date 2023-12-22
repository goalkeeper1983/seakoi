package postgreSQLConnect

import (
	"fmt"
	"github.com/goalkeeper1983/seakoi/tools"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// host user password dbname port
func CreatePostgreSQLConnect(option ...string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", option[0], option[1], option[2], option[3], option[4])
	//dsn := "host=localhost user=your_username password=your_password dbname=your_dbname port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		tools.Log.Panic(tools.RunFuncName(), zap.Any("err", err.Error()))
	}
	if err = db.Exec("select 1").Error; err != nil {
		tools.Log.Panic(tools.RunFuncName(), zap.Any("err", err.Error()))
	}
	return db
}
