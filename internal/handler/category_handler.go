package handler

import (
	"net/http"
	"strconv"

	"todo-api/internal/domain"
	"todo-api/internal/service"
	"todo-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CategoryHandler struct {
	service  service.CategoryService
	validate *validator.Validate
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service:  service,
		validate: validator.New(),
	}
}

type CreateCategoryInput struct {
	Name  string `json:"name" validate:"required,min=3"`
	Color string `json:"color" validate:"omitempty"`
}

type UpdateCategoryInput struct {
	Name  string `json:"name" validate:"required,min=3"`
	Color string `json:"color" validate:"omitempty"`
}

// GetAllCategories godoc
// @Summary      Ambil semua kategori
// @Description  Mengambil semua kategori milik user
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Router       /categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.Error(utils.Unauthorized("User tidak terautentikasi"))
		return
	}

	categories, err := h.service.GetAllCategories(userID)
	if err != nil {
		c.Error(utils.InternalServerError("Gagal mengambil data kategori"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Data kategori berhasil diambil",
		categories,
	))
}

// GetCategoryByID godoc
// @Summary      Ambil kategori berdasarkan ID
// @Description  Mengambil satu kategori milik user berdasarkan ID
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Success      200 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Failure      404 {object} domain.Response
// @Router       /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.Error(utils.Unauthorized("User tidak terautentikasi"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(utils.BadRequest("ID tidak valid", "ID harus berupa angka"))
		return
	}

	category, err := h.service.GetCategoryByID(uint(id), userID)
	if err != nil {
		c.Error(utils.NotFound("Kategori tidak ditemukan"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Kategori berhasil ditemukan",
		category,
	))
}

// CreateCategory godoc
// @Summary      Buat kategori baru
// @Description  Membuat kategori baru untuk user yang sedang login
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateCategoryInput true "Data kategori"
// @Success      201 {object} domain.Response
// @Failure      400 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Router       /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.Error(utils.Unauthorized("User tidak terautentikasi"))
		return
	}

	var input CreateCategoryInput
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

	category, err := h.service.CreateCategory(userID, input.Name, input.Color)
	if err != nil {
		c.Error(utils.BadRequest("Gagal membuat kategori", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, domain.NewSuccessResponse(
		"Kategori berhasil dibuat",
		category,
	))
}

// UpdateCategory godoc
// @Summary      Update kategori
// @Description  Mengupdate kategori milik user berdasarkan ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Param        request body UpdateCategoryInput true "Data kategori"
// @Success      200 {object} domain.Response
// @Failure      400 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Router       /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.Error(utils.Unauthorized("User tidak terautentikasi"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(utils.BadRequest("ID tidak valid", "ID harus berupa angka"))
		return
	}

	var input UpdateCategoryInput
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

	category, err := h.service.UpdateCategory(uint(id), userID, input.Name, input.Color)
	if err != nil {
		c.Error(utils.BadRequest("Gagal mengupdate kategori", err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Kategori berhasil diupdate",
		category,
	))
}

// DeleteCategory godoc
// @Summary      Hapus kategori
// @Description  Menghapus kategori milik user berdasarkan ID
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Category ID"
// @Success      200 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Failure      404 {object} domain.Response
// @Router       /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.Error(utils.Unauthorized("User tidak terautentikasi"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(utils.BadRequest("ID tidak valid", "ID harus berupa angka"))
		return
	}

	err = h.service.DeleteCategory(uint(id), userID)
	if err != nil {
		c.Error(utils.NotFound("Gagal menghapus kategori"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Kategori berhasil dihapus",
		nil,
	))
}
