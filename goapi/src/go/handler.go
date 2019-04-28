package swagger

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	database     *gorm.DB
	err          error
	currentValue int
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

func Get(params string) (string, error) {
	secret := Secret{}
	err := database.Where("hash=?", params).First(&secret).Error
	if err != nil {
		return "", err
	}
	if secret.ExpiresAt.Sub(time.Now().Local()) <= 0 {
		return "Link expired", nil
	}

	decryptData, err := Decrypt(CipherKey, secret.SecretText)

	if err == nil {
		tx := database.Begin()
		tx.Raw("SELECT `remaining_views` FROM secrets WHERE `hash` = ? LIMIT 1 FOR UPDATE", params).Row().Scan(&currentValue)
		currentValue++
		err := tx.Exec("UPDATE secrets SET `remaining_views` = ? WHERE `hash` = ? LIMIT 1", currentValue, params).Error
		if err != nil {
			return "", err
		}
		tx.Commit()
	}

	return decryptData, nil

}
