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
	Start(uuid uuid.UUID) (*model.Tunnel, error)
	Restart(uuid uuid.UUID) (*model.Tunnel, error)
	Stop(uuid uuid.UUID) (*model.Tunnel, error)

	AddConn(uuid uuid.UUID, domain string) (*model.Tunnel, error)
	RemoveConn(uuid uuid.UUID) (*model.Tunnel, error)

	Create(name string) (*model.Tunnel, error)
	Info(uuid uuid.UUID) (*model.Tunnel, error)
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

// Info implements ITunnelSrv.
func (t TunelSrv) Info(uuid uuid.UUID) (*model.Tunnel, error) {
	cmd := exec.Command(_CLOUDFLARED, _TUNNEL, "info", _OUTPUT, uuid.String())
	data, err := cmd.Output()
	if err != nil {
		checkErr(err)
		t.logger.Errorf("Error running the command = %s, err = %w", cmd.String(), err)
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

// AddConn implements ITunnelSrv.
// runs and parses ❯ cloudflared tunnel route dns [uuid] [domain]
func (t TunelSrv) AddConn(uuid uuid.UUID, domain string) (*model.Tunnel, error) {
	cmd := exec.Command(_CLOUDFLARED, _TUNNEL, "route", "dns", uuid.String(), domain)

	err := cmd.Run()
	if err != nil {
		checkErr(err)
		t.logger.Errorf("Error running the command = %s, err = %w", cmd.String(), err)
		return nil, err
	}

	tunnel, err := t.Info(uuid)
	if err != nil {
		t.logger.Errorf("Error retriving information about the tunnel, err = %v", err)
		return nil, err
	}

	return tunnel, nil
}

// RemoveConn implements ITunnelSrv.
func (t TunelSrv) RemoveConn(uuid uuid.UUID) (*model.Tunnel, error) {
	// TODO: see with cloudflared api
	panic("unimplemented")
}

// Restart implements ITunnelSrv.
func (t TunelSrv) Restart(uuid uuid.UUID) (*model.Tunnel, error) {
	panic("unimplemented")
}

// Start implements ITunnelSrv.
func (t TunelSrv) Start(uuid uuid.UUID) (*model.Tunnel, error) {
	panic("unimplemented")
}

// Stop implements ITunnelSrv.
func (t TunelSrv) Stop(uuid uuid.UUID) (*model.Tunnel, error) {
	panic("unimplemented")
}

// Create implements ITunnelSrv.
// runs and parses ❯ cloudflared tunnel create [name]
func (t TunelSrv) Create(name string) (*model.Tunnel, error) {
	cmd := exec.Command(_CLOUDFLARED, _TUNNEL, "create", _OUTPUT, name)
	data, err := cmd.Output()
	if err != nil {
		checkErr(err)
		t.logger.Errorf("Error running the command = %s, err = %w", cmd.String(), err)
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
		t.logger.Errorf("Error running the command = %s, err = %w", cmd.String(), err)
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
		t.logger.Errorf("Error running the command = %s, err = %w", cmd.String(), err)
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
