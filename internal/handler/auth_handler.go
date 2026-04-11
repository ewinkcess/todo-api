package handler

import (
	"net/http"
	"strings"
	"todo-api/internal/domain"
	"todo-api/internal/service"
	"todo-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	service  service.AuthService
	validate *validator.Validate
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	v := validator.New()

	// Daftarkan custom validator untuk password kuat
	v.RegisterValidation("strongpassword", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		hasUpper := false
		hasLower := false
		hasNumber := false
		hasSpecial := false

		for _, char := range password {
			switch {
			case char >= 'A' && char <= 'Z':
				hasUpper = true
			case char >= 'a' && char <= 'z':
				hasLower = true
			case char >= '0' && char <= '9':
				hasNumber = true
			case strings.ContainsRune("!@#$%^&*", char):
				hasSpecial = true
			}
		}
		return hasUpper && hasLower && hasNumber && hasSpecial
	})
	return &AuthHandler{
		service:  service,
		validate: v,
	}
}

type RegisterInput struct {
	Name     string `json:"name"     validate:"required,min=3"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,strongpassword"`
}

type LoginInput struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Register godoc
// @Summary      Register user baru
// @Description  Mendaftarkan user baru dengan nama, email, dan password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterInput true "Data registrasi"
// @Success      201 {object} domain.Response
// @Failure      400 {object} domain.Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput
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

	user, err := h.service.Register(input.Name, input.Email, input.Password)
	if err != nil {
		c.Error(utils.BadRequest("Registrasi gagal", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, domain.NewSuccessResponse(
		"Registrasi berhasil",
		gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	))
}

// Login godoc
// @Summary      Login user
// @Description  Login dengan email dan password untuk mendapatkan JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginInput true "Data login"
// @Success      200 {object} domain.Response
// @Failure      400 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput

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
	token, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		c.Error(utils.Unauthorized("Email atau password salah"))
		return
	}
	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Login berhasil",
		gin.H{"token": token},
	))
}

func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " wajib diisi"
	case "email":
		return err.Field() + " harus format email yang valid"
	case "min":
		return err.Field() + " minimal " + err.Param() + " Karakter"
	case "strongpassword":
		return " harus mengandung karakter "
	default:
		return err.Field() + " tidak valid"
	}
}
