package seed

import (
	"errors"
	"os"

	"github.com/killi1812/cloudflared-web-gui/dto"
	"github.com/killi1812/cloudflared-web-gui/model"
	"github.com/killi1812/cloudflared-web-gui/service"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	_PASSWORD_ENV = "SUPERADMIN_PASSWORD"
)

var suadmin *model.User

// CreateSuperAdmin creates a SuperAdmin user if one doesn't already exist.
// It reads the password from the SUPERADMIN_PASSWORD environment variable.
// The function will panic if required environment variables are missing or
// if user creation fails, as this is critical for application bootstrap.
func createSuperAdmin() error {
	userCrud := service.NewUserCrudService()

	// Check if SuperAdmin exists
	{
		_, err := userCrud.SearchUsersByName("superadmin")
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				zap.S().Infof("SuperAdmin not found, err %+v", err)
			} else {
				return err
			}
		} else {
			zap.S().Infoln("SuperAdmin found")
			zap.S().Infoln("Skipping superadmin creation")
			return nil
		}
	}

	zap.S().Infoln("Crating superadmin creation")
	password := os.Getenv(_PASSWORD_ENV)
	if password == "" {
		return errors.New("env variable is empty")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	dto := dto.NewUserDto{
		Username: "superadmin",
		Password: password,
		Role:     "superadmin",
	}
	newUser, err := dto.ToModel()
	if err != nil {
		return err
	}

	user, err := userCrud.Create(newUser, dto.Password)
	if err != nil {
		return err
	}
	suadmin = user
	zap.S().Infof("superadmin created, %+v\n", user)
	return nil
}
