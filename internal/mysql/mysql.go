package mysql

import (
	"aiStudio/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(c *conf.Mysql) (err error) {
	DB, err = gorm.Open(mysql.Open(c.ToDSNString()))
	if err != nil {
		return err
	}
	return nil
}
