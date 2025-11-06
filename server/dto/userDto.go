package dto

import (
	"fmt"

	"github.com/killi1812/cloudflared-web-gui/model"
	"github.com/killi1812/cloudflared-web-gui/util/cerror"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserDto struct {
	Uuid        string `json:"uuid"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	PoliceToken string `json:"policeToken"`
}

func (dto UserDto) ToModel() (*model.User, error) {
	uuid, err := uuid.Parse(dto.Uuid)
	if err != nil {
		zap.S().Error("Failed to parse uuid = %s, err = %+v", dto.Uuid, err)
		return nil, cerror.ErrBadUuid
	}

	role, err := model.StrToUserRole(dto.Role)
	if err != nil {
		zap.S().Errorf("Failed to parse role = %+v, err = %+v", dto.Role, err)
		return nil, cerror.ErrUnknownRole
	}

	return &model.User{
		Uuid: uuid,
		Role: role,
	}, nil
}

// FromModel returns a dto from model struct
func (UserDto) FromModel(m *model.User) UserDto {
	dto := &UserDto{
		Uuid: m.Uuid.String(),
		Role: fmt.Sprint(m.Role),
	}
	return *dto
}
