package handler

import (
	"net/http"
	"strconv"
	"todo-api/internal/domain"
	"todo-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TodoHandler struct {
	service  service.TodoService
	validate *validator.Validate
}

func NewTodoHandler(service service.TodoService) *TodoHandler {
	return &TodoHandler{
		service:  service,
		validate: validator.New(),
	}
}

type CreateTodoInput struct {
	Title       string `json:"title"       validate:"required,min=3"`
	Description string `json:"description" validate:"omitempty,min=3"`
}

type UpdateTodoInput struct {
	Title       string `json:"title"       validate:"required,min=3"`
	Description string `json:"description" validate:"omitempty,min=3"`
	Completed   bool   `json:"completed"`
}

func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetAllTodos godoc
// @Summary      Ambil semua todo
// @Description  Mengambil semua todo milik user yang sedang login
// @Tags         todos
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Router       /todos [get]
func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
			"Akses ditolak",
			"User tidak terautentikasi",
		))
		return
	}
	todos, err := h.service.GetAllTodos(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(
			"Gagal mengambil data todo",
			err.Error(),
		))
		return
	}
	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Data todo berhasil diambil",
		todos,
	))
}

// GetTodoByID godoc
// @Summary      Ambil todo berdasarkan ID
// @Description  Mengambil satu todo milik user berdasarkan ID
// @Tags         todos
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Todo ID"
// @Success      200 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Failure      404 {object} domain.Response
// @Router       /todos/{id} [get]
func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
			"Akses ditolak",
			"User tidak terautentikasi",
		))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(
			"ID tidak valid",
			"ID harus berupa angka",
		))
		return
	}
	todo, err := h.service.GetTodoByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.NewErrorResponse(
			"Todo tidak ditemukan",
			"Todo dengan ID tersebut tidak ada",
		))
		return
	}
	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Todo berhasil dibuat",
		todo,
	))
}

// CreateTodo godoc
// @Summary      Buat todo baru
// @Description  Membuat todo baru untuk user yang sedang login
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateTodoInput true "Data todo"
// @Success      201 {object} domain.Response
// @Failure      400 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Router       /todos [post]
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
			"Akses ditolak",
			"User tidak terautentikasi",
		))
		return
	}
	var input CreateTodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(
			"Format request tidak valid",
			err.Error(),
		))
		return
	}
	if err := h.validate.Struct(input); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(
			"Validasi gagal",
			validationErrors,
		))
		return
	}

	todo, err := h.service.CreateTodo(userID, input.Title, input.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(
			"Gagal memuat todo",
			err.Error(),
		))
		return
	}
	c.JSON(http.StatusCreated, domain.NewSuccessResponse(
		"Todo berhasil dibuat",
		todo,
	))
}

// UpdateTodo godoc
// @Summary      Update todo
// @Description  Mengupdate todo milik user berdasarkan ID
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int            true "Todo ID"
// @Param        request body UpdateTodoInput true "Data todo"
// @Success      200 {object} domain.Response
// @Failure      400 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Router       /todos/{id} [put]
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
			"Akses ditolak",
			"User tidak terautentikasi",
		))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(
			"ID tidak valid",
			"ID harus berupa angka",
		))
		return
	}
	var input UpdateTodoInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(
			"Format request tidak valid",
			err.Error(),
		))
		return
	}
	if err := h.validate.Struct(input); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(
			"Validasi gagal",
			validationErrors,
		))
		return
	}
	todo, err := h.service.UpdateTodo(uint(id), userID, input.Title, input.Description, input.Completed)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(
			"Gagal mengupdate todo",
			err.Error(),
		))
		return
	}
	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Todo berhasil diupdate",
		todo,
	))
}

// DeleteTodo godoc
// @Summary      Hapus todo
// @Description  Menghapus todo milik user berdasarkan ID
// @Tags         todos
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Todo ID"
// @Success      200 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Failure      404 {object} domain.Response
// @Router       /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(
			"Akses ditolak",
			"User tidak terautentikasi",
		))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(
			"ID tidak valid",
			"Id harus berupa angka",
		))
		return
	}

	err = h.service.DeleteTodo(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.NewErrorResponse(
			"Gagal menghapus todo",
			err.Error(),
		))
		return
	}
	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Todo berhasil dihapus",
		nil,
	))
}
