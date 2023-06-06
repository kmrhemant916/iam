package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connection(host string, database string, password string, port string, username string) (*gorm.DB, error){
    connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
        username,
        password,
        host,
        port,
        database,
    )
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return db, nil
}
