package dto_test

import (
	"testing"

	"github.com/killi1812/cloudflared-web-gui/dto"
	"github.com/killi1812/cloudflared-web-gui/model"
	"github.com/killi1812/cloudflared-web-gui/util/cerror"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUserDto_ToModel(t *testing.T) {
	userUUID := uuid.New()

	tests := []struct {
		name    string
		dto     dto.NewUserDto
		want    *model.User // Uuid will be generated, so we compare other fields
		wantErr error
	}{
		{
			name: "Valid DTO to Model - No UUID in DTO",
			dto: dto.NewUserDto{
				// Uuid: "", // Intentionally empty, ToModel should generate it
				Password: "password123",
				Role:     string(model.ROLE_USER),
			},
			want: &model.User{
				// Uuid will be checked for non-nil
				Role: model.ROLE_USER,
				// PasswordHash is set by service, not DTO
			},
			wantErr: nil,
		},
		{
			name: "Valid DTO to Model - With UUID in DTO",
			dto: dto.NewUserDto{
				Uuid:     userUUID.String(),
				Password: "securepass",
				Role:     string(model.ROLE_ADMIN),
			},
			want: &model.User{
				// Uuid will be checked for non-nil (ToModel generates a new one regardless of DTO input)
				Role: model.ROLE_ADMIN,
			},
			wantErr: nil,
		},
		{
			name: "Invalid Role",
			dto: dto.NewUserDto{
				Username: "Bad role",
				Role:     "nonexistent_role",
				Password: "password",
			},
			want:    nil,
			wantErr: cerror.ErrUnknownRole,
		},
		{
			name: "Invalid UUID in DTO",
			dto: dto.NewUserDto{
				Uuid:     "this-is-not-a-uuid",
				Username: "test",
				Role:     string(model.ROLE_USER),
				Password: "password",
			},
			want:    nil,
			wantErr: cerror.ErrBadUuid,
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
				assert.NotNil(t, got)
				assert.NotEqual(t, uuid.Nil, got.Uuid, "Generated UUID should not be nil")
				assert.Equal(t, tt.want.Role, got.Role)
			}
		})
	}
}

func TestNewUserDto_FromModel(t *testing.T) {
	userUUID := uuid.New()
	userModel := &model.User{
		Uuid:         userUUID,
		Username:     "test",
		PasswordHash: "somehash", // Not included in NewUserDto
		Role:         model.ROLE_SUPER_ADMIN,
	}

	expectedDto := dto.NewUserDto{
		Uuid:     userUUID.String(),
		Username: "test",
		Password: "", // Password is not part of FromModel for NewUserDto
		Role:     string(model.ROLE_SUPER_ADMIN),
	}

	gotDto := dto.NewUserDto{}.FromModel(userModel)

	assert.Equal(t, expectedDto, gotDto)
}
