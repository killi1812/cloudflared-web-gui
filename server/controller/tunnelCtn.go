package controller

import (
	"errors"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/killi1812/cloudflared-web-gui/app"
	"github.com/killi1812/cloudflared-web-gui/dto"
	"github.com/killi1812/cloudflared-web-gui/service"
	"github.com/killi1812/cloudflared-web-gui/util/auth"
	"github.com/killi1812/cloudflared-web-gui/util/cerror"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type (
	nameDto   struct{ name string }
	domainDto struct {
		Domain string `json:"domain" binding:"required`
	}
)

// NewImageCnt creates a new controller for images.
func NewTunnelCtn() app.Controller {
	var controller *TunnelCtn
	app.Invoke(func(logger *zap.SugaredLogger, srv service.ITunnelSrv, dns service.IDnsSrv) {
		controller = &TunnelCtn{
			Logger:    logger,
			TunnelSrv: srv,
			DndSrv:    dns,
		}
	})
	return controller
}

type TunnelCtn struct {
	Logger    *zap.SugaredLogger
	TunnelSrv service.ITunnelSrv
	DndSrv    service.IDnsSrv
}

// RegisterEndpoints registers the image manipulation endpoints.
func (cnt *TunnelCtn) RegisterEndpoints(router *gin.RouterGroup) {
	grp := router.Group("/tunnel", auth.Protect())
	grp.GET("", cnt.getTunnels)
	grp.POST("", cnt.createTunnel)

	grp.DELETE("/:id", cnt.deleteTunnel)
	grp.GET("/:id", cnt.getInfo)

	grp.POST("/dns/:id", cnt.createDnsRecord)

	grp.PUT("/:id/start", cnt.startTunnel)
	grp.PUT("/:id/stop", cnt.stopTunnel)
	grp.PUT("/:id/restart", cnt.restartTunnel)
}

// getTunnels godoc
//
//	@Summary		Get a list of all tunnels
//	@Description	returns a list of all tunnels
//	@Tags			tunnel
//	@Produce		json
//	@Success		200	{object}	[]dto.TunnelDto	"List of tunnels"
//	@Router			/tunnel [get]
func (ctn *TunnelCtn) getTunnels(c *gin.Context) {
	list, err := ctn.TunnelSrv.List()
	if err != nil {
		ctn.Logger.Errorf("Error retrieving tunnels, err %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var resp []dto.TunnelDto = make([]dto.TunnelDto, len(list))
	var wg sync.WaitGroup

	for i, tnl := range list {
		wg.Go(func() {
			resp[i].FromModel(tnl)
			dnsRecords, err := ctn.DndSrv.GetDnsRecords(tnl.Id)
			if err != nil {
				if !errors.Is(err, cerror.ErrZoneIdNotSet) && !errors.Is(err, cerror.ErrCloudflaredApiKeyNotSet) {
					ctn.Logger.Warnln(err)
				} else {
					ctn.Logger.Errorf("Error Dns Records for tunnel %s, err %v", tnl.Id, err)
				}
			} else {
				resp[i].DnsRecords.FromModel(dnsRecords)
			}
		})
	}

	wg.Wait()
	c.AbortWithStatusJSON(http.StatusOK, resp)
}

// getTunnels godoc
//
//	@Summary		creates a new tunnel
//	@Description	creates a new tunnel with given name and returns it
//	@Tags			tunnel
//	@Produce		json
//	@Success		201		{object}	model.Tunnel	"Newly created tunnel"
//	@Param			name	body		nameDto			true	"tunnel name"
//	@Router			/tunnel [post]
func (ctn *TunnelCtn) createTunnel(c *gin.Context) {
	var req nameDto

	err := c.BindJSON(&req)
	if err != nil {
		ctn.Logger.Errorf("Error body format, err = %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tunnel, err := ctn.TunnelSrv.Create(req.name)
	if err != nil {
		ctn.Logger.Errorf("Error creating a tunnel, err = %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, tunnel)
}

// getTunnels godoc
//
//	@Summary		Get a list of all tunnels
//	@Description	returns a list of all tunnels
//	@Tags			tunnel
//	@Produce		json
//	@Success		204	"Tunnel deleted"
//	@Param			id	path	string	true	"tunnel id"
//	@Router			/tunnel/{id} [delete]
func (ctn *TunnelCtn) deleteTunnel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ctn.Logger.Error("Error id not found in a form")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		ctn.Logger.Errorf("Error parsing uuid, id = %s", id)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = ctn.TunnelSrv.Delete(uuid)
	if err != nil {
		ctn.Logger.Errorf("Error deleting a tunnel, err = %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// createDnsRecord godoc
//
//	@Summary		Creates a dns record on the tunnel
//	@Description	Creates a new dns record on the tunnel
//	@Tags			tunnel
//	@Produce		json
//	@Success		201	{object}	model.Tunnel	"Tunnel dns record created"
//	@Param			id	path		string			true	"tunnel id"
//	@Param			id	body		domainDto		true	"dns domain"
//	@Router			/tunnel/dns/{id} [post]
func (ctn *TunnelCtn) createDnsRecord(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ctn.Logger.Error("Error id not found in a form")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		ctn.Logger.Errorf("Error parsing uuid, id = %s", id)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var req domainDto
	err = c.BindJSON(&req)
	if err != nil {
		ctn.Logger.Errorf("Error body format, err = %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctn.Logger.Debugf("req: %+v", req)

	tunnel, err := ctn.TunnelSrv.AddConn(uuid, req.Domain)
	if err != nil {
		ctn.Logger.Errorf("Error deleting a tunnel, err = %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, tunnel)
}

// getInfo godoc
//
//	@Summary		Creates a dns record on the tunnel
//	@Description	Creates a new dns record on the tunnel
//	@Tags			tunnel
//	@Produce		json
//	@Success		200	{object}	model.Tunnel	"Tunnel dns record created"
//	@Param			id	path		string			true	"tunnel id"
//	@Router			/tunnel/{id} [get]
func (ctn *TunnelCtn) getInfo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ctn.Logger.Error("Error id not found in a form")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		ctn.Logger.Errorf("Error parsing uuid, id = %s", id)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tunnel, err := ctn.TunnelSrv.Info(uuid)
	if err != nil {
		ctn.Logger.Errorf("Error getting tunnel info, id = %s , err = %v", uuid, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var resp dto.TunnelDto
	resp.FromModel(*tunnel)

	dnsRecords, err := ctn.DndSrv.GetDnsRecords(tunnel.Id)
	if err != nil {
		ctn.Logger.Errorf("Error retrieving tunnel dns records, id = %s, err = %v", tunnel.Id, err)
		if !errors.Is(err, cerror.ErrZoneIdNotSet) && !errors.Is(err, cerror.ErrCloudflaredApiKeyNotSet) {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else {
		resp.DnsRecords.FromModel(dnsRecords)
	}

	c.AbortWithStatusJSON(http.StatusOK, resp)
}

// startTunnel godoc
//
//	@Summary		starts a tunnel
//	@Description	starts a tunnel as system proccess
//	@Tags			tunnel
//	@Produce		json
//	@Success		204	"Tunnel started"
//	@Param			id	path	string	true	"tunnel id"
//	@Router			/tunnel/{id}/start [put]
func (ctn *TunnelCtn) startTunnel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ctn.Logger.Error("Error id not found in a form")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		ctn.Logger.Errorf("Error parsing uuid, id = %s", id)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = ctn.TunnelSrv.Start(uuid)
	if err != nil {
		ctn.Logger.Errorf("Error starting a tunnel, err = %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// stopTunnel godoc
//
//	@Summary		stops a tunnel
//	@Description	stops a tunnel running as system proccess
//	@Tags			tunnel
//	@Produce		json
//	@Success		204	"Tunnel stopped"
//	@Param			id	path	string	true	"tunnel id"
//	@Router			/tunnel/{id}/start [put]
func (ctn *TunnelCtn) stopTunnel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ctn.Logger.Error("Error id not found in a form")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		ctn.Logger.Errorf("Error parsing uuid, id = %s", id)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = ctn.TunnelSrv.Stop(uuid)
	if err != nil {
		ctn.Logger.Errorf("Error stopping a tunnel, err = %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// restartTunnel godoc
//
//	@Summary		restarts a tunnel
//	@Description	restarts a tunnel with zero downtime
//	@Tags			tunnel
//	@Produce		json
//	@Success		204	"Tunnel restarted"
//	@Param			id	path	string	true	"tunnel id"
//	@Router			/tunnel/{id}/restart [put]
func (ctn *TunnelCtn) restartTunnel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		ctn.Logger.Error("Error id not found in a form")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		ctn.Logger.Errorf("Error parsing uuid, id = %s", id)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = ctn.TunnelSrv.Restart(uuid)
	if err != nil {
		ctn.Logger.Errorf("Error stopping a tunnel, err = %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
