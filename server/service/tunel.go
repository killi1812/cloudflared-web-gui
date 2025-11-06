package service

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"github.com/google/uuid"
	"github.com/killi1812/cloudflared-web-gui/app"
	"github.com/killi1812/cloudflared-web-gui/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	_CLOUDFLARED = "cloudflared"
	_TUNNEL      = "tunnel"
	_OUTPUT      = "--output=json"
)

type ITunnelSrv interface {
	Create(name string) (*model.Tunnel, error)
	Start(uuid uuid.UUID) (*model.Tunnel, error)
	Stop(uuid uuid.UUID) (*model.Tunnel, error)
	List() ([]model.Tunnel, error)
	Delete(uuid uuid.UUID) error
}

func NewTunelSrv() ITunnelSrv {
	var service ITunnelSrv
	app.Invoke(func(db *gorm.DB, logger *zap.SugaredLogger) {
		service = TunelSrv{
			db:     db,
			logger: logger,
		}
	})

	return service
}

type TunelSrv struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

// Create implements ITunnelSrv.
// runs and parses ❯ cloudflared tunnel create [name]
func (t TunelSrv) Create(name string) (*model.Tunnel, error) {
	cmd := exec.Command(_CLOUDFLARED, _TUNNEL, "create", _OUTPUT, name)
	data, err := cmd.Output()
	if err != nil {
		checkErr(err)
		t.logger.Errorf("Error running the command,err = %w", err)
		return nil, err
	}

	var tunnel model.Tunnel
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&tunnel)
	if err != nil {
		t.logger.Errorf("Error decoding data, err = %w", err)
		return nil, err
	}

	return &tunnel, nil
}

// Delete implements ITunnelSrv.
func (t TunelSrv) Delete(uuid uuid.UUID) error {
	cmd := exec.Command(_CLOUDFLARED, _TUNNEL, "delete", _OUTPUT, uuid.String())
	err := cmd.Run()
	if err != nil {
		checkErr(err)
		t.logger.Errorf("Error running the command,err = %w", err)
		return err
	}

	return nil
}

// List implements ITunnelSrv.
// runs and parses ❯ cloudflared tunnel list
func (t TunelSrv) List() ([]model.Tunnel, error) {
	cmd := exec.Command(_CLOUDFLARED, _TUNNEL, "list", _OUTPUT)
	data, err := cmd.Output()
	if err != nil {
		checkErr(err)
		t.logger.Errorf("Error running the command,err = %w", err)
		return nil, err
	}

	var list []model.Tunnel
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&list)
	if err != nil {
		t.logger.Errorf("Error decoding data, err = %w", err)
		return nil, err
	}

	return list, nil
}

// Start implements ITunnelSrv.
func (t TunelSrv) Start(uuid uuid.UUID) (*model.Tunnel, error) {
	panic("unimplemented")
}

// Stop implements ITunnelSrv.
func (t TunelSrv) Stop(uuid uuid.UUID) (*model.Tunnel, error) {
	panic("unimplemented")
}

func checkErr(err error) {
	nerr, ok := err.(*exec.Error)
	if ok {
		zap.S().Errorf("%+v", nerr)
	}

	exerr, ok := err.(*exec.ExitError)
	if ok {
		zap.S().Errorf("%+v", exerr.Error)
	}
}
