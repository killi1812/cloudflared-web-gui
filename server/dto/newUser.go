package dto

import (
	"github.com/killi1812/cloudflared-web-gui/model"
	"github.com/killi1812/cloudflared-web-gui/util/cerror"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NewUserDto struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username" binding:"required,min=2,max=100"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin user superadmin"`
}

// ToModel create a model from a dto
func (dto NewUserDto) ToModel() (*model.User, error) {
	role, err := model.StrToUserRole(dto.Role)
	if err != nil {
		zap.S().Error("Failed to parse role = %+v, err = %+v", dto.Role, err)
		return nil, cerror.ErrUnknownRole
	}
	if dto.Uuid != "" {
		_, err := uuid.Parse(dto.Uuid)
		if err != nil {
			zap.S().Errorf("Failed to parse uuid = %s, err = %+v", dto.Uuid, err)
			return nil, cerror.ErrBadUuid
		}
	}

	return &model.User{
		Uuid:     uuid.New(),
		Username: dto.Username,
		Role:     role,
	}, nil
}

// FromModel returns a dto from model struct
func (NewUserDto) FromModel(m *model.User) NewUserDto {
	dto := NewUserDto{
		Uuid:     m.Uuid.String(),
		Username: m.Username,
		Role:     string(m.Role),
	}
	return dto
}
