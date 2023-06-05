package database

import (
	"fmt"

	"github.com/kmrhemant916/iam/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	Config = "config/config.yaml"
)

func Connection() (*gorm.DB, error){
    var config utils.Config
	c, _:= config.ReadConf(Config)
    connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		c.Database.Username,
        c.Database.Password,
        c.Database.Host,
        c.Database.Port,
        c.Database.Name,
    )
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return db, nil
}