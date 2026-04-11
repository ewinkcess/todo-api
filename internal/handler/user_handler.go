package handler

import (
	"net/http"

	"todo-api/internal/domain"
	"todo-api/internal/service"
	"todo-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service  service.UserService
	validate *validator.Validate
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validator.New(),
	}
}

type UpdateProfileInput struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

type UpdatePasswordInput struct {
	OldPassword string `json:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" validate:"required,min=8,strongpassword"`
}

// GetProfile godoc
// @Summary      Ambil profil user
// @Description  Mengambil data profil user yang sedang login
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Router       /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.Error(utils.Unauthorized("User tidak terautentikasi"))
		return
	}
	user, err := h.service.GetProfile(userID)
	if err != nil {
		c.Error(utils.NotFound("User tidak ditemukan"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Profil berhasil diambil",
		gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	))
}

// UpdateProfile godoc
// @Summary      Update profil user
// @Description  Mengupdate nama dan email user yang sedang login
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body UpdateProfileInput true "Data profil"
// @Success      200 {object} domain.Response
// @Failure      400 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Router       /users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.Error(utils.Unauthorized("User tidak terautentikasi"))
		return
	}

	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(utils.BadRequest("Format request tidak valid", err.Error()))
		return
	}

	if err := h.validate.Struct(input); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}
		c.Error(utils.ValidationError(validationErrors))
		return
	}
	user, err := h.service.UpdateProfile(userID, input.Name, input.Email)
	if err != nil {
		c.Error(utils.BadRequest("Gagal mengupdate profil", err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Profil berhasil diupdate",
		gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	))
}

// UpdatePassword godoc
// @Summary      Ganti password
// @Description  Mengganti password user yang sedang login
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body UpdatePasswordInput true "Data password"
// @Success      200 {object} domain.Response
// @Failure      400 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Router       /users/me/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.Error(utils.Unauthorized("User tidak terautentikasi"))
		return
	}

	var input UpdatePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(utils.BadRequest("Format request tidak valid", err.Error()))
		return
	}

	if err := h.validate.Struct(input); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}
		c.Error(utils.ValidationError(validationErrors))
		return
	}

	err := h.service.UpdatePassword(userID, input.OldPassword, input.NewPassword)
	if err != nil {
		c.Error(utils.BadRequest("Gagal mengganti password", err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Password berhasil diganti",
		nil,
	))
}
