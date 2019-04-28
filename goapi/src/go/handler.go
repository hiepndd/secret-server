package swagger

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	database *gorm.DB
	err      error
)

func Init() error {
	database, err = gorm.Open("mysql", "root:123456@/secret?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return err
	}
	database.AutoMigrate(Secret{})
	return nil
}

func Add(secret *Secret) error {
	err := database.Create(&secret).Error
	if err != nil {
		return err
	}

	return nil
}
