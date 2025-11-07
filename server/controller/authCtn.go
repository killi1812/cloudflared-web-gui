package controller

import (
	"net/http"

	"github.com/killi1812/cloudflared-web-gui/app"
	"github.com/killi1812/cloudflared-web-gui/dto"
	"github.com/killi1812/cloudflared-web-gui/service"
	"github.com/killi1812/cloudflared-web-gui/util/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthCtn struct {
	auth   service.IAuthService
	logger *zap.SugaredLogger
}

func NewAuthCtn() app.Controller {
	var controller *AuthCtn

	// Use the mock service for testing
	app.Invoke(func(loginService service.IAuthService, logger *zap.SugaredLogger) {
		// create controller
		controller = &AuthCtn{
			auth:   loginService,
			logger: logger,
		}
	})

	return controller
}

func (ctn *AuthCtn) RegisterEndpoints(api *gin.RouterGroup) {
	// create a group with the name of the router
	group := api.Group("/auth")

	// register Endpoints
	group.POST("/login", ctn.login)
	group.POST("/refresh", auth.Protect(), ctn.refreshToken)
	group.POST("/logout", ctn.logout)
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticates a user and returns access and refresh tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			loginDto	body		dto.LoginDto	true	"Login credentials"
//	@Success		200			{object}	dto.TokenDto
//	@Router			/auth/login [post]
func (ctn *AuthCtn) login(c *gin.Context) {
	var loginDto dto.LoginDto

	if err := c.BindJSON(&loginDto); err != nil {
		ctn.logger.Errorf("Invalid login request err = %+v", err)
		return
	}

	accessToken, err := ctn.auth.Login(loginDto.Username, loginDto.Password)
	if err != nil {
		ctn.logger.Errorf("Login failed err = %+v", err)
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.TokenDto{
		AccessToken: accessToken,
	})
}

// Refresh godoc
//
//	@Summary		Refresh Access Token
//	@Description	Generates a new access token using a valid refresh token
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	dto.TokenDto
//	@Router			/auth/refresh [post]
func (ctn *AuthCtn) refreshToken(c *gin.Context) {
	tokenStr := c.Request.Header.Get("Authorization")
	token, err := ctn.auth.RefreshTokens(tokenStr)
	if err != nil {
		ctn.logger.Errorf("Refresh failed err = %w", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.TokenDto{
		AccessToken: token,
	})
}

// Refresh godoc
//
//	@Summary		Refresh Access Token
//	@Description	Generates a new access token using a valid refresh token
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	dto.TokenDto
//	@Router			/auth/logout [post]
func (ctn *AuthCtn) logout(c *gin.Context) {
	_, claims, err := auth.ParseToken(c.Request.Header.Get("Authorization"))
	if err != nil {
		ctn.logger.Errorf("Logout failed err = %w", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = ctn.auth.Logout(claims.ID)
	if err != nil {
		ctn.logger.Errorf("Logout failed err = %w", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
