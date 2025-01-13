package userGinRouter

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/pedramktb/schwarzit-probearbeit/internal/datasource"
	"github.com/pedramktb/schwarzit-probearbeit/internal/dtos"
	ginRouter "github.com/pedramktb/schwarzit-probearbeit/internal/gin"
	"github.com/pedramktb/schwarzit-probearbeit/internal/logging"
	"github.com/pedramktb/schwarzit-probearbeit/internal/types"
)

type r struct {
	datasource.Getter[types.User]
	datasource.Querier[types.User]
	datasource.Saver[types.User]
	datasource.Deleter[types.User]
}

func create(
	getter datasource.Getter[types.User],
	querier datasource.Querier[types.User],
	saver datasource.Saver[types.User],
	deleter datasource.Deleter[types.User],
) *r {
	return &r{
		getter,
		querier,
		saver,
		deleter,
	}
}

// @Summary Create a user
// @Description Create a user
// @Tags user
// @Accept json
// @Produce json
// @Param user body SaveUser true "User"
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse "Bad Request Error"
// @Failure 401 {object} ErrorResponse "Unauthorized Error"
// @Failure 403 {object} ErrorResponse "Forbidden Error"
// @Failure 404 {object} ErrorResponse "Not Found Error"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/users [post]
func (r *r) Create(c *gin.Context) {
	if isAdmin := c.GetBool(string(logging.CtxUserIsAdmin)); !isAdmin {
		ginRouter.ErrorResponse(c, types.ErrForbidden)
		return
	}

	userDTO := dtos.SaveUser{}
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		ginRouter.ErrorResponse(c, errors.CombineErrors(types.ErrBadRequest, err))
		return
	}

	if user, err := r.Saver.Save(c.Request.Context(), userDTO.ToCreateUser()); err != nil {
		ginRouter.ErrorResponse(c, err)
	} else {
		c.JSON(http.StatusOK, dtos.FromUser(&user))
	}
}

// @Summary Query users
// @Description Query users
// @Tags user
// @Security Bearer
// @Accept json
// @Produce json
// @Param params query UserQueryParams false "Query Parameters"
// @Success 200 {array} []User
// @Failure 400 {object} ErrorResponse "Bad Request Error"
// @Failure 401 {object} ErrorResponse "Unauthorized Error"
// @Failure 403 {object} ErrorResponse "Forbidden Error"
// @Failure 404 {object} ErrorResponse "Not Found Error"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/users [get]
func (r *r) Query(c *gin.Context) {
	paramsDTO := dtos.UserQueryParams{}
	if err := c.ShouldBindQuery(&paramsDTO); err != nil {
		ginRouter.ErrorResponse(c, errors.CombineErrors(types.ErrBadRequest, err))
		return
	}

	if users, err := r.Querier.Query(c, paramsDTO.ToQueryParams()); err != nil {
		ginRouter.ErrorResponse(c, err)
	} else {
		userDTOs := make([]dtos.User, len(users))
		for i, user := range users {
			userDTOs[i] = dtos.FromUser(&user)
		}
		c.JSON(http.StatusOK, userDTOs)
	}
}

// @Summary Get a user
// @Description Get a user by id
// @Tags user
// @Security Bearer
// @Produce json
// @Param id path string true "User ID"
// @Failure 400 {object} ErrorResponse "Bad Request Error"
// @Failure 401 {object} ErrorResponse "Unauthorized Error"
// @Failure 403 {object} ErrorResponse "Forbidden Error"
// @Failure 404 {object} ErrorResponse "Not Found Error"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/users/{id} [get]
func (r *r) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ginRouter.ErrorResponse(c, errors.CombineErrors(types.ErrInvalidID, err))
		return
	}

	if user, err := r.Getter.Get(c, id); err != nil {
		ginRouter.ErrorResponse(c, err)
	} else {
		c.JSON(http.StatusOK, dtos.FromUser(&user))
	}
}

// @Summary Update a user
// @Description Update a user by id
// @Tags user
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body SaveUser true "User"
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse "Bad Request Error"
// @Failure 401 {object} ErrorResponse "Unauthorized Error"
// @Failure 403 {object} ErrorResponse "Forbidden Error"
// @Failure 404 {object} ErrorResponse "Not Found Error"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/users/{id} [put]
func (r *r) Update(c *gin.Context) {
	if isAdmin := c.GetBool(string(logging.CtxUserIsAdmin)); !isAdmin {
		ginRouter.ErrorResponse(c, types.ErrForbidden)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ginRouter.ErrorResponse(c, errors.CombineErrors(types.ErrInvalidID, err))
		return
	}

	userDTO := dtos.SaveUser{}
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		ginRouter.ErrorResponse(c, errors.CombineErrors(types.ErrBadRequest, err))
		return
	}

	if user, err := r.Saver.Save(c.Request.Context(), userDTO.ToUpdateUser(id)); err != nil {
		ginRouter.ErrorResponse(c, err)
	} else {
		c.JSON(http.StatusOK, dtos.FromUser(&user))
	}
}

// @Summary Patch a user
// @Description Patch a user by id
// @Tags user
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body PatchUser true "User"
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse "Bad Request Error"
// @Failure 401 {object} ErrorResponse "Unauthorized Error"
// @Failure 403 {object} ErrorResponse "Forbidden Error"
// @Failure 404 {object} ErrorResponse "Not Found Error"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/users/{id} [patch]
func (r *r) Patch(c *gin.Context) {
	if isAdmin := c.GetBool(string(logging.CtxUserIsAdmin)); !isAdmin {
		ginRouter.ErrorResponse(c, types.ErrForbidden)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ginRouter.ErrorResponse(c, errors.CombineErrors(types.ErrInvalidID, err))
		return
	}

	userDTO := dtos.PatchUser{}
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		ginRouter.ErrorResponse(c, errors.CombineErrors(types.ErrBadRequest, err))
		return
	}

	user, err := r.Getter.Get(c.Request.Context(), id)
	if err != nil {
		ginRouter.ErrorResponse(c, err)
		return
	}

	user.ApplyPatch(userDTO.ToUserPatch())

	if user, err := r.Saver.Save(c.Request.Context(), user); err != nil {
		ginRouter.ErrorResponse(c, err)
	} else {
		c.JSON(http.StatusOK, dtos.FromUser(&user))
	}
}

// @Summary Delete a user
// @Description Delete a user by id
// @Tags user
// @Security Bearer
// @Produce json
// @Param id path string true "User ID"
// @Success 200
// @Failure 400 {object} ErrorResponse "Bad Request Error"
// @Failure 401 {object} ErrorResponse "Unauthorized Error"
// @Failure 403 {object} ErrorResponse "Forbidden Error"
// @Failure 404 {object} ErrorResponse "Not Found Error"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/users/{id} [delete]
func (r *r) Delete(c *gin.Context) {
	if isAdmin := c.GetBool(string(logging.CtxUserIsAdmin)); !isAdmin {
		ginRouter.ErrorResponse(c, types.ErrForbidden)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ginRouter.ErrorResponse(c, errors.CombineErrors(types.ErrInvalidID, err))
		return
	}

	if err := r.Deleter.Delete(c, id); err != nil {
		ginRouter.ErrorResponse(c, err)
	} else {
		c.Status(http.StatusOK)
	}
}

// @Summary Get me (user)
// @Description Get me as a user
// @Tags user
// @Security Bearer
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse "Bad Request Error"
// @Failure 401 {object} ErrorResponse "Unauthorized Error"
// @Failure 403 {object} ErrorResponse "Forbidden Error"
// @Failure 404 {object} ErrorResponse "Not Found Error"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/users/me [get]
func (r *r) GetMe(c *gin.Context) {
	aid, _ := c.Get(string(logging.CtxUserID))
	id, _ := aid.(uuid.UUID)

	if user, err := r.Getter.Get(c.Request.Context(), id); err != nil {
		ginRouter.ErrorResponse(c, err)
	} else {
		c.JSON(http.StatusOK, dtos.FromUser(&user))
	}
}

// @Summary Update me (user)
// @Description Update me as a user
// @Tags user
// @Security Bearer
// @Accept json
// @Produce json
// @Param user body SaveUser true "User"
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse "Bad Request Error"
// @Failure 401 {object} ErrorResponse "Unauthorized Error"
// @Failure 403 {object} ErrorResponse "Forbidden Error"
// @Failure 404 {object} ErrorResponse "Not Found Error"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/users/me [put]
func (r *r) UpdateMe(c *gin.Context) {
	aid, _ := c.Get(string(logging.CtxUserID))
	id, _ := aid.(uuid.UUID)

	userDTO := dtos.SaveUser{}
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		ginRouter.ErrorResponse(c, errors.CombineErrors(types.ErrBadRequest, err))
		return
	}

	if isAdmin := c.GetBool(string(logging.CtxUserIsAdmin)); !isAdmin && userDTO.IsAdmin {
		ginRouter.ErrorResponse(c, types.ErrForbidden)
		return
	}

	if user, err := r.Saver.Save(c.Request.Context(), userDTO.ToUpdateUser(id)); err != nil {
		ginRouter.ErrorResponse(c, err)
	} else {
		c.JSON(http.StatusOK, dtos.FromUser(&user))
	}
}

// @Summary Patch me (user)
// @Description Patch me as a user
// @Tags user
// @Security Bearer
// @Accept json
// @Produce json
// @Param user body PatchUser true "User"
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse "Bad Request Error"
// @Failure 401 {object} ErrorResponse "Unauthorized Error"
// @Failure 403 {object} ErrorResponse "Forbidden Error"
// @Failure 404 {object} ErrorResponse "Not Found Error"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/users/me [patch]
func (r *r) PatchMe(c *gin.Context) {
	aid, _ := c.Get(string(logging.CtxUserID))
	id, _ := aid.(uuid.UUID)

	userDTO := dtos.PatchUser{}
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		ginRouter.ErrorResponse(c, errors.CombineErrors(types.ErrBadRequest, err))
		return
	}

	if isAdmin := c.GetBool(string(logging.CtxUserIsAdmin)); !isAdmin && userDTO.IsAdmin != nil {
		ginRouter.ErrorResponse(c, types.ErrForbidden)
		return
	}

	user, err := r.Getter.Get(c.Request.Context(), id)
	if err != nil {
		ginRouter.ErrorResponse(c, err)
		return
	}

	user.ApplyPatch(userDTO.ToUserPatch())

	if user, err := r.Saver.Save(c.Request.Context(), user); err != nil {
		ginRouter.ErrorResponse(c, err)
	} else {
		c.JSON(http.StatusOK, dtos.FromUser(&user))
	}
}

// @Summary Delete me (user)
// @Description Delete me as a user
// @Tags user
// @Security Bearer
// @Produce json
// @Success 200
// @Failure 400 {object} ErrorResponse "Bad Request Error"
// @Failure 401 {object} ErrorResponse "Unauthorized Error"
// @Failure 403 {object} ErrorResponse "Forbidden Error"
// @Failure 404 {object} ErrorResponse "Not Found Error"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/users/me [delete]
func (r *r) DeleteMe(c *gin.Context) {
	aid, _ := c.Get(string(logging.CtxUserID))
	id, _ := aid.(uuid.UUID)

	if err := r.Deleter.Delete(c.Request.Context(), id); err != nil {
		ginRouter.ErrorResponse(c, err)
	} else {
		c.Status(http.StatusOK)
	}
}
