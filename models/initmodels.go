package models

import (
	"admin/confs"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)
var DB *gorm.DB
// 初始化Gorm
func InitMysql()  (err error){
	conn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",confs.AdminConf.MysqlConf.User, confs.AdminConf.MysqlConf.Password, confs.AdminConf.MysqlConf.Host, confs.AdminConf.MysqlConf.DataBase)
	db, err := gorm.Open("mysql", conn)
	if err != nil {
		return
	}

	DB = db
	//if confs.AdminConf.MysqlConf.Debug {
	//	DB.LogMode(true)
	//	DB.SetLogger(logs.Info(os.Stdout, "\r\n", 0))
	//}
	DB.SingularTable(true)

	DB.DB().SetMaxOpenConns(confs.AdminConf.MysqlConf.MaxOpenConns)
	DB.DB().SetConnMaxLifetime(time.Duration(confs.AdminConf.MysqlConf.MaxLifetime) * time.Second)
	DB.DB().SetMaxIdleConns(confs.AdminConf.MysqlConf.MaxIdleConns)
	Migration()
	return
}
func Migration() {
	fmt.Println(DB.AutoMigrate(new(Menu)).Error)
	fmt.Println(DB.AutoMigrate(new(User)).Error)
	fmt.Println(DB.AutoMigrate(new(RoleMenu)).Error)
	fmt.Println(DB.AutoMigrate(new(Role)).Error)
	fmt.Println(DB.AutoMigrate(new(UserRole)).Error)
}


