package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
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
	Start(uuid uuid.UUID) error
	Stop(uuid uuid.UUID) error
	Restart(uuid uuid.UUID) error

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
		service = &TunnelSrv{
			db:         db,
			logger:     logger,
			tunnelProc: make(map[uuid.UUID]*os.Process),
		}
	})

	return service
}

type TunnelSrv struct {
	db         *gorm.DB
	logger     *zap.SugaredLogger
	tunnelProc map[uuid.UUID]*os.Process // tunnelPid is a map with [Key] tunnel id and [Value] *os.process
}

// Info implements ITunnelSrv.
func (t *TunnelSrv) Info(uuid uuid.UUID) (*model.Tunnel, error) {
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
func (t *TunnelSrv) AddConn(uuid uuid.UUID, domain string) (*model.Tunnel, error) {
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
func (t *TunnelSrv) RemoveConn(uuid uuid.UUID) (*model.Tunnel, error) {
	// TODO: see with cloudflared api
	panic("unimplemented")
}

// Restart implements ITunnelSrv.
func (t *TunnelSrv) Restart(uuid uuid.UUID) error {
	panic("unimplemented")
}

// Start implements ITunnelSrv.
// runs and parses ❯ cloudflared tunnel run [tunnel id]
func (t *TunnelSrv) Start(uuid uuid.UUID) error {
	_, ok := t.tunnelProc[uuid]
	if ok {
		zap.S().Infof("Tunnel uuid = %s already running", uuid)
		return errors.New("tunnel already running")
	}

	cmd := exec.Command(_CLOUDFLARED, _TUNNEL, _OUTPUT, "run", uuid.String())
	err := cmd.Start()
	if err != nil {
		checkErr(err)
		t.logger.Errorf("Error running the command = %s, err = %w", cmd.String(), err)
		return err
	}
	t.tunnelProc[uuid] = cmd.Process

	return nil
}

// Stop implements ITunnelSrv.
func (t *TunnelSrv) Stop(uuid uuid.UUID) error {
	proc, ok := t.tunnelProc[uuid]
	if !ok {
		t.logger.Errorf("Procces running a tunnel %s not found", uuid.String())
		return errors.New("Procces not found")
	}

	err := proc.Kill()
	if err != nil {
		t.logger.Errorf("Failed to kill procces = %+v, err = %v", *proc, err)
		return err
	}

	_, err = proc.Wait()
	if err != nil {
		t.logger.Errorf("Failed to Wait for procces = %+v, err = %v", *proc, err)
		return err
	}

	delete(t.tunnelProc, uuid)
	return nil
}

// Create implements ITunnelSrv.
// runs and parses ❯ cloudflared tunnel create [name]
func (t *TunnelSrv) Create(name string) (*model.Tunnel, error) {
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
func (t *TunnelSrv) Delete(uuid uuid.UUID) error {
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
func (t *TunnelSrv) List() ([]model.Tunnel, error) {
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
