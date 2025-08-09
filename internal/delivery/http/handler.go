package http

import (
	"net/http"
	"shop/internal/domain"
	"shop/internal/middleware"
	"shop/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Router         *gin.Engine
	userUsecase    *usecase.UserUsecase
	productUsecase *usecase.ProductUsecase
	orderUsecase   *usecase.OrderUsecase
	cartUsecase    *usecase.CartUsecase
}

func NewHandler(userUsecase *usecase.UserUsecase, productUsecase *usecase.ProductUsecase, orderUsecase *usecase.OrderUsecase, cartUsecase *usecase.CartUsecase) *Handler {
	handler := &Handler{
		Router:         gin.Default(),
		userUsecase:    userUsecase,
		productUsecase: productUsecase,
		orderUsecase:   orderUsecase,
		cartUsecase:    cartUsecase,
	}

	handler.setupRoutes()
	return handler
}

func (h *Handler) setupRoutes() {
	api := h.Router.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
	}

	// User routes
	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("/profile", h.getProfile)
	}

	// Product routes
	products := api.Group("/products")
	{
		products.GET("", h.getProducts)
		products.GET("/:id", h.getProduct)
	}

	// Admin product routes
	adminProducts := api.Group("/admin/products")
	adminProducts.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		adminProducts.POST("", h.createProduct)
		adminProducts.PUT("/:id", h.updateProduct)
		adminProducts.DELETE("/:id", h.deleteProduct)
	}

	// Order routes
	orders := api.Group("/orders")
	orders.Use(middleware.AuthMiddleware())
	{
		orders.POST("", h.createOrder)
		orders.GET("", h.getUserOrders)
		orders.GET("/:id", h.getOrder)
	}

	cart := api.Group("/cart")
	cart.Use(middleware.AuthMiddleware())
	{
		cart.POST("", h.createCart)
		//cart.GET("", h.getUserOrders)
	}
}

func (h *Handler) register(c *gin.Context) {
	var req domain.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUsecase.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *Handler) login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.userUsecase.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) getProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	user, err := h.userUsecase.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Product handlers
func (h *Handler) createProduct(c *gin.Context) {
	var req domain.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productUsecase.CreateProduct(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"product": product})
}

func (h *Handler) getProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.productUsecase.GetProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (h *Handler) getProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	products, err := h.productUsecase.GetProducts(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (h *Handler) updateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req domain.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productUsecase.UpdateProduct(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (h *Handler) deleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.productUsecase.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// Order handlers
func (h *Handler) createOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req domain.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderUsecase.CreateOrder(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order": order})
}

func (h *Handler) getOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderUsecase.GetOrder(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

func (h *Handler) getUserOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	orders, err := h.orderUsecase.GetUserOrders(userID.(uint), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *Handler) createCart(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req domain.RequestCart
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.cartUsecase.CreateCart(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"added": "ok"})
}
