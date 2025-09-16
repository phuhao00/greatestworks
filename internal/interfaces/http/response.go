package http

import (
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
)

// APIResponse 统一API响应结构
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// APIError API错误信息
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Field   string      `json:"field,omitempty"`
}

// Meta 元数据信息
type Meta struct {
	Page       int   `json:"page,omitempty"`
	PageSize   int   `json:"page_size,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// PaginationRequest 分页请求
type PaginationRequest struct {
	Page     int `form:"page" json:"page" binding:"min=1"`
	PageSize int `form:"page_size" json:"page_size" binding:"min=1,max=100"`
}

// GetPagination 获取分页参数
func (p *PaginationRequest) GetPagination() (int, int) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.Page, p.PageSize
}

// GetOffset 获取偏移量
func (p *PaginationRequest) GetOffset() int {
	page, pageSize := p.GetPagination()
	return (page - 1) * pageSize
}

// GetLimit 获取限制数量
func (p *PaginationRequest) GetLimit() int {
	_, pageSize := p.GetPagination()
	return pageSize
}

// 响应构建器

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}, message ...string) {
	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	
	response := APIResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: getRequestID(c),
	}
	
	c.JSON(http.StatusOK, response)
}

// SuccessResponseWithMeta 带元数据的成功响应
func SuccessResponseWithMeta(c *gin.Context, data interface{}, meta *Meta, message ...string) {
	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	
	response := APIResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now(),
		RequestID: getRequestID(c),
	}
	
	c.JSON(http.StatusOK, response)
}

// CreatedResponse 创建成功响应
func CreatedResponse(c *gin.Context, data interface{}, message ...string) {
	msg := "Created successfully"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	
	response := APIResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: getRequestID(c),
	}
	
	c.JSON(http.StatusCreated, response)
}

// NoContentResponse 无内容响应
func NoContentResponse(c *gin.Context, message ...string) {
	msg := "No content"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	
	response := APIResponse{
		Success:   true,
		Message:   msg,
		Timestamp: time.Now(),
		RequestID: getRequestID(c),
	}
	
	c.JSON(http.StatusNoContent, response)
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, statusCode int, code, message string, details ...interface{}) {
	apiError := &APIError{
		Code:    code,
		Message: message,
	}
	
	if len(details) > 0 {
		apiError.Details = details[0]
	}
	
	response := APIResponse{
		Success:   false,
		Error:     apiError,
		Timestamp: time.Now(),
		RequestID: getRequestID(c),
	}
	
	c.JSON(statusCode, response)
}

// BadRequestResponse 400错误响应
func BadRequestResponse(c *gin.Context, message string, details ...interface{}) {
	ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", message, details...)
}

// UnauthorizedResponse 401错误响应
func UnauthorizedResponse(c *gin.Context, message ...string) {
	msg := "Unauthorized"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", msg)
}

// ForbiddenResponse 403错误响应
func ForbiddenResponse(c *gin.Context, message ...string) {
	msg := "Forbidden"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", msg)
}

// NotFoundResponse 404错误响应
func NotFoundResponse(c *gin.Context, message ...string) {
	msg := "Resource not found"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", msg)
}

// ConflictResponse 409错误响应
func ConflictResponse(c *gin.Context, message string, details ...interface{}) {
	ErrorResponse(c, http.StatusConflict, "CONFLICT", message, details...)
}

// ValidationErrorResponse 422验证错误响应
func ValidationErrorResponse(c *gin.Context, field, message string, details ...interface{}) {
	apiError := &APIError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Field:   field,
	}
	
	if len(details) > 0 {
		apiError.Details = details[0]
	}
	
	response := APIResponse{
		Success:   false,
		Error:     apiError,
		Timestamp: time.Now(),
		RequestID: getRequestID(c),
	}
	
	c.JSON(http.StatusUnprocessableEntity, response)
}

// InternalServerErrorResponse 500错误响应
func InternalServerErrorResponse(c *gin.Context, message ...string) {
	msg := "Internal server error"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", msg)
}

// ServiceUnavailableResponse 503错误响应
func ServiceUnavailableResponse(c *gin.Context, message ...string) {
	msg := "Service unavailable"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	ErrorResponse(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", msg)
}

// 辅助函数

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	return c.GetHeader("X-Request-ID")
}

// BindAndValidate 绑定并验证请求数据
func BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		BadRequestResponse(c, "Invalid request data", err.Error())
		return false
	}
	return true
}

// BindQueryAndValidate 绑定并验证查询参数
func BindQueryAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		BadRequestResponse(c, "Invalid query parameters", err.Error())
		return false
	}
	return true
}

// BindURIAndValidate 绑定并验证URI参数
func BindURIAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindUri(obj); err != nil {
		BadRequestResponse(c, "Invalid URI parameters", err.Error())
		return false
	}
	return true
}

// GetIDParam 获取ID参数
func GetIDParam(c *gin.Context) string {
	return c.Param("id")
}

// GetPlayerIDParam 获取玩家ID参数
func GetPlayerIDParam(c *gin.Context) string {
	return c.Param("player_id")
}

// ValidateID 验证ID参数
func ValidateID(c *gin.Context, paramName string) (string, bool) {
	id := c.Param(paramName)
	if id == "" {
		BadRequestResponse(c, fmt.Sprintf("%s is required", paramName))
		return "", false
	}
	return id, true
}

// CreateMeta 创建元数据
func CreateMeta(page, pageSize int, total int64) *Meta {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return &Meta{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

// HandleError 处理错误
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	
	// 根据错误类型返回不同的响应
	switch {
	case isNotFoundError(err):
		NotFoundResponse(c, err.Error())
	case isValidationError(err):
		BadRequestResponse(c, err.Error())
	case isConflictError(err):
		ConflictResponse(c, err.Error())
	default:
		InternalServerErrorResponse(c, err.Error())
	}
}

// 错误类型检查函数

func isNotFoundError(err error) bool {
	// TODO: 实现具体的错误类型检查
	return false
}

func isValidationError(err error) bool {
	// TODO: 实现具体的错误类型检查
	return false
}

func isConflictError(err error) bool {
	// TODO: 实现具体的错误类型检查
	return false
}