package controller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/killi1812/cloudflared-web-gui/app"
	"github.com/killi1812/cloudflared-web-gui/service"
	"github.com/killi1812/cloudflared-web-gui/util/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TunnelCtn struct {
	Logger    *zap.SugaredLogger
	TunnelSrv service.ITunnelSrv
}

// NewImageCnt creates a new controller for images.
func NewTunnelCtn() app.Controller {
	var controller *TunnelCtn
	app.Invoke(func(logger *zap.SugaredLogger, srv service.ITunnelSrv) {
		controller = &TunnelCtn{
			Logger:    logger,
			TunnelSrv: srv,
		}
	})
	return controller
}

// RegisterEndpoints registers the image manipulation endpoints.
func (cnt *TunnelCtn) RegisterEndpoints(router *gin.RouterGroup) {
	grp := router.Group("/tunnel", auth.Protect())
	grp.GET("", cnt.getTunnels)
	grp.POST("", cnt.createTunnel)
	grp.DELETE("/:id", cnt.deleteTunnel)
}

// getTunnels godoc
//
//	@Summary		Get a list of all tunnels
//	@Description	returns a list of all tunnels
//	@Tags			tunnel
//	@Produce		json
//	@Success		200	{struct}	[]model.Tunnel	"List of tunnels"
//	@Router			/tunnel [get]
func (ctn *TunnelCtn) getTunnels(c *gin.Context) {
	list, err := ctn.TunnelSrv.List()
	if err != nil {
		ctn.Logger.Errorf("Error retrieving tunnels, err %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, list)
}

// getTunnels godoc
//
//	@Summary		creates a new tunnel
//	@Description	creates a new tunnel with given name and returns it
//	@Tags			tunnel
//	@Produce		json
//	@Success		201		{struct}	model.Tunnel	"Newly created tunnel"
//	@Param			name	body		string			true	"tunnel name"
//	@Router			/tunnel [post]
func (ctn *TunnelCtn) createTunnel(c *gin.Context) {
	var req struct{ name string }

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
