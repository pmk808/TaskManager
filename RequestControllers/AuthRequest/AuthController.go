package AuthRequest

import (
	"net/http"
	"taskmanager/RequestControllers/AuthRequest/dto"
	"taskmanager/RequestControllers/AuthRequest/interfaces"
	"taskmanager/RequestControllers/httpSetup/jwt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type authController struct {
	jwtManager *jwt.JWTManager
	logger     *logrus.Logger
}

func NewAuthController(jwtManager *jwt.JWTManager, logger *logrus.Logger) interfaces.AuthController {
	return &authController{
		jwtManager: jwtManager,
		logger:     logger,
	}
}

func (c *authController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/token", c.GenerateToken)
}

// GenerateToken godoc
// @Summary Generate JWT token
// @Description Generates a JWT token for client authentication
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.GenerateTokenRequest true "Client credentials"
// @Success 200 {object} dto.GenerateTokenResponse
// @Failure 400 {object} dto.GenerateTokenResponse
// @Router /auth/token [post]
func (c *authController) GenerateToken(ctx *gin.Context) {
	var request dto.GenerateTokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.logger.WithError(err).Error("Invalid token request")
		ctx.JSON(http.StatusBadRequest, dto.GenerateTokenResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	token, err := c.jwtManager.GenerateToken(request.ClientName, request.ClientID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to generate token")
		ctx.JSON(http.StatusInternalServerError, dto.GenerateTokenResponse{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.GenerateTokenResponse{
		Success: true,
		Token:   token,
		Message: "Token generated successfully",
	})
}
