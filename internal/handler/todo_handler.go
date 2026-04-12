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
	CategoryID  *uint  `json:"category_id"`
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

func (h *TodoHandler) bindAndValidate(c *gin.Context, input interface{}) bool {
	if err := c.ShouldBindJSON(input); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}
		c.Error(utils.ValidationError(validationErrors))
		return false
	}
	return true
}

// GetAllTodos godoc
// @Summary      Ambil semua todo
// @Description  Mengambil semua todo milik user dengan pagination dan filter
// @Tags         todos
// @Produce      json
// @Security     BearerAuth
// @Param        page      query int    false "Halaman ke berapa (default: 1)"
// @Param        limit     query int    false "Jumlah data per halaman (default: 10)"
// @Param        search    query string false "Cari berdasarkan judul"
// @Param        completed query bool   false "Filter berdasarkan status"
// @Success      200 {object} domain.Response
// @Failure      401 {object} domain.Response
// @Router       /todos [get]
func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.Error(utils.Unauthorized("User tidak terautentikasi"))
		return
	}
	var query domain.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(utils.BadRequest("ID tidak valid", "ID harus berupa angka"))
		return
	}
	result, err := h.service.GetAllTodos(userID, query)
	if err != nil {
		c.Error(utils.InternalServerError("Gagal mengambil data todo"))
		return
	}
	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Data todo berhasil diambil",
		result,
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
		c.Error(utils.Unauthorized("User tidak terautentikasi"))
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(utils.BadRequest("ID tidak valid", "ID harus berupa angka"))
		return
	}
	todo, err := h.service.GetTodoByID(uint(id), userID)
	if err != nil {
		c.Error(utils.NotFound("Todo tidak ditemukan"))
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
		c.Error(utils.Unauthorized("User tidak terautentikasi"))
		return
	}

	var input CreateTodoInput
	if !h.bindAndValidate(c, &input) {
		return
	}

	todo, err := h.service.CreateTodo(userID, input.Title, input.Description, input.CategoryID)
	if err != nil {
		c.Error(utils.BadRequest("Gagal membuat todo", err.Error()))
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
		c.Error(utils.Unauthorized("User tidak terautentikasi"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(utils.BadRequest("ID tidak valid", "ID harus berupa angka"))
		return
	}

	var input UpdateTodoInput
	if !h.bindAndValidate(c, &input) {
		return
	}

	todo, err := h.service.UpdateTodo(uint(id), userID, input.Title, input.Description, input.Completed)
	if err != nil {
		c.Error(utils.BadRequest("Gagal mengupdate todo", err.Error()))
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
		c.Error(utils.Unauthorized("User tidak terautentikasi"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(utils.BadRequest("ID tidak valid", "ID harus berupa angka"))
		return
	}

	err = h.service.DeleteTodo(uint(id), userID)
	if err != nil {
		c.Error(utils.NotFound("Gagal menghapus todo"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(
		"Todo berhasil dihapus",
		nil,
	))
}
