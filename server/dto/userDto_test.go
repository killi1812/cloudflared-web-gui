package dto_test

import (
	"testing"

	"github.com/killi1812/cloudflared-web-gui/dto"
	"github.com/killi1812/cloudflared-web-gui/model"
	"github.com/killi1812/cloudflared-web-gui/util/cerror"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserDto_ToModel(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name    string
		dto     dto.UserDto
		want    *model.User
		wantErr error
	}{
		{
			name: "Valid DTO to Model",
			dto: dto.UserDto{
				Uuid: validUUID.String(),
				Role: "user",
			},
			want: &model.User{
				Uuid: validUUID,
				Role: model.ROLE_USER,
			},
			wantErr: nil,
		},
		{
			name: "Invalid UUID",
			dto: dto.UserDto{
				Uuid: "not-a-uuid",
				Role: "firma",
			},
			want:    nil,
			wantErr: cerror.ErrBadUuid,
		},
		{
			name: "Invalid Role",
			dto: dto.UserDto{
				Uuid: validUUID.String(),
				Role: "invalid_role",
			},
			want:    nil,
			wantErr: cerror.ErrUnknownRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dto.ToModel()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				// Compare individual fields as GORM Model has unexported fields
				assert.Equal(t, tt.want.Uuid, got.Uuid)
				assert.Equal(t, tt.want.Role, got.Role)
			}
		})
	}
}

func TestUserDto_FromModel(t *testing.T) {
	userUUID := uuid.New()
	userModel := &model.User{
		Uuid: userUUID,
		Role: model.ROLE_USER,
	}

	expectedDto := dto.UserDto{
		Uuid: userUUID.String(),
		Role: string(model.ROLE_USER),
	}

	var gotDto dto.UserDto
	gotDto = gotDto.FromModel(userModel)

	assert.Equal(t, expectedDto, gotDto)
}
