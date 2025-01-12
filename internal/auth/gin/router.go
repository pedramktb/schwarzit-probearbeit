package authGinRouter

import (
	"net/http"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	authJWT "github.com/pedramktb/schwarzit-probearbeit/internal/auth/jwt"
	"github.com/pedramktb/schwarzit-probearbeit/internal/datasource"
	"github.com/pedramktb/schwarzit-probearbeit/internal/dtos"
	ginRouter "github.com/pedramktb/schwarzit-probearbeit/internal/gin"
	"github.com/pedramktb/schwarzit-probearbeit/internal/types"
)

type r struct {
	datasource.UserByEmailGetter
	userSaver datasource.Saver[types.User]
	jwt       *authJWT.JWT
}

func create(
	userByEmailGetter datasource.UserByEmailGetter,
	userSaver datasource.Saver[types.User],
	jwt *authJWT.JWT,
) *r {

	return &r{
		userByEmailGetter,
		userSaver,
		jwt,
	}
}

// @Summary Register
// @Description user registration
// @Tags auth
// @Produce json
// @Param user body RegisterUser true "Register User"
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/register [post]
func (r *r) Register(c *gin.Context) {
	var registerDTO dtos.RegisterUser
	if err := c.ShouldBindBodyWithJSON(&registerDTO); err != nil {
		ginRouter.ErrorResponse(c, errors.CombineErrors(types.ErrBadRequest, err))
		return
	}

	user, err := r.userSaver.Save(c.Request.Context(), registerDTO.ToUser())
	if err != nil {
		ginRouter.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.FromUser(&user))
}

// @Summary Login
// @Description user login
// @Tags auth
// @Produce json
// @Param login body LoginRequest true "Login Request"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/login [post]
func (r *r) Login(c *gin.Context) {
	var loginRequest dtos.LoginRequest
	if err := c.ShouldBindBodyWithJSON(&loginRequest); err != nil {
		ginRouter.ErrorResponse(c, errors.CombineErrors(types.ErrBadRequest, err))
		return
	}

	user, err := r.UserByEmailGetter.GetByEmail(c.Request.Context(), loginRequest.Email)
	if err != nil {
		ginRouter.ErrorResponse(c, err)
		return
	}

	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"admin": user.IsAdmin,
	}

	accessToken, err := r.jwt.GenerateAccessToken(claims)
	if err != nil {
		ginRouter.ErrorResponse(c, err)
		return
	}

	refreshToken, err := r.jwt.GenerateRefreshToken(claims)
	if err != nil {
		ginRouter.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// @Summary Refresh
// @Description refresh token
// @Tags auth
// @Produce json
// @Param Authorization header string true "Refresh Token"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/refresh [post]
func (r *r) Refresh(c *gin.Context) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

	claims, err := r.jwt.ValidateRefreshToken(token)
	if err != nil {
		ginRouter.ErrorResponse(c, errors.Wrap(types.ErrUnauthorized, "invalid refresh token"))
		return
	}

	accessToken, err := r.jwt.GenerateAccessToken(claims)
	if err != nil {
		ginRouter.ErrorResponse(c, err)
		return
	}

	refreshToken, err := r.jwt.GenerateRefreshToken(claims)
	if err != nil {
		ginRouter.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// @Summary AuthMiddleware
// @Description AuthMiddleware is the middleware for the authentication
// @Tags auth
// @Param Authorization header string true "Bearer Token"
func (r *r) AuthMiddleware(c *gin.Context) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

	claims, err := r.jwt.ValidateAccessToken(token)
	if err != nil {
		ginRouter.ErrorResponse(c, errors.Wrap(types.ErrUnauthorized, "invalid access token"))
		c.Abort()
		return
	}

	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		ginRouter.ErrorResponse(c, errors.Wrap(types.ErrUnauthorized, "invalid access token"))
		c.Abort()
		return
	}

	c.Set("user_id", userID)

	if admin := claims["admin"].(bool); admin {
		c.Set("admin", true)
	}

	c.Next()
}
