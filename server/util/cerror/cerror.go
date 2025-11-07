package cerror

import (
	"errors"
	"fmt"

	"github.com/killi1812/cloudflared-web-gui/util/format"
)

var (
	ErrBadDateFormat           = fmt.Errorf("bad date format, should be %s", format.DateFormat)
	ErrBadDateTimeFormat       = fmt.Errorf("bad date and time format, should be %s", format.DateTimeFormat)
	ErrBadTimeFormat           = fmt.Errorf("bad time format, should be %s", format.TimeFormat)
	ErrBadUuid                 = errors.New("failed to parse uuid")
	ErrUnknownRole             = errors.New("unknown role")
	ErrInvalidCredentials      = errors.New("invalid email or password")
	ErrInvalidTokenFormat      = errors.New("invalid token format")
	ErrUserIsNil               = errors.New("user is nil")
	ErrBadRole                 = errors.New("role is not allowed")
	ErrCloudflaredApiKeyNotSet = errors.New("cloudflared api key not set")
	ErrZoneIdNotSet            = errors.New("zone id not set")
)
