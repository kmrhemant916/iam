package helpers

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/kmrhemant916/iam/models"
	"gorm.io/gorm"
)

var scripts = []string{
	"scripts/permission.sql",
	"scripts/role.sql",
	"scripts/rolePermission.sql",
}

func InitialiseServices(db *gorm.DB) {
	Createtable(db)
	// DatabaseMigration(db)
}

func DatabaseMigration(db *gorm.DB) {
	for _, script := range scripts {
		scriptPath,err := GetAbsPath(script)
		if err != nil {
			panic("Error in path")
		}
		ExecuteDatabaseScript(scriptPath, db)
	}
}

func Createtable(db *gorm.DB) {
	db.AutoMigrate(models.Permission{})
	db.AutoMigrate(models.Role{})
	db.AutoMigrate(models.RolePermission{})
	db.AutoMigrate(models.User{})
	db.AutoMigrate(models.Organization{})
}

func GetAbsPath(p string) (string, error) {
	absPath, err := filepath.Abs(p)
	return absPath, err
}

func ExecuteDatabaseScript(p string, db *gorm.DB) {
	data, err := os.Open(p)
	if err != nil {
		panic("Error loading sql script")
	}
	defer data.Close()
    scanner := bufio.NewScanner(data)
    for scanner.Scan() {
        sql := scanner.Text()
        if err := db.Exec(sql).Error; err != nil {
            panic(err)
        }
    }
    if err := scanner.Err(); err != nil {
        panic(err)
    }
}