package account

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	service *AccountService
	logger  *log.Logger
}

func NewAccountController(service *AccountService, logger *log.Logger) *AccountController {
	return &AccountController{
		service: service,
		logger:  logger,
	}
}

func (c *AccountController) CreateAccount(ctx *gin.Context) {
	userIDStr := ctx.PostForm("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqCtx := ctx.Request.Context()
	reqCtx, cancel := context.WithTimeout(reqCtx, 5*time.Second)
	defer cancel()

	account, err := c.service.CreateAccount(reqCtx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

func (c *AccountController) GetAccountByAccountNumber(ctx *gin.Context) {
	idNumber := ctx.Param("id_number")

	reqCtx := ctx.Request.Context()
	reqCtx, cancel := context.WithTimeout(reqCtx, 5*time.Second)
	defer cancel()

	account, err := c.service.GetAccountByIDNumber(reqCtx, idNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *AccountController) GetAccountBalance(ctx *gin.Context) {
	idNumber := ctx.Param("id_number")

	reqCtx := ctx.Request.Context()
	reqCtx, cancel := context.WithTimeout(reqCtx, 5*time.Second)
	defer cancel()

	balance, err := c.service.GetAccountBalance(reqCtx, idNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (c *AccountController) GetAllAccounts(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	reqCtx, cancel := context.WithTimeout(reqCtx, 5*time.Second)
	defer cancel()

	accounts, err := c.service.GetAllAccounts(reqCtx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (c *AccountController) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/accounts")
	{
		group.POST("", c.CreateAccount)
		group.GET("/:id_number", c.GetAccountByAccountNumber)
		group.GET("/:id_number/balance", c.GetAccountBalance)
		group.GET("", c.GetAllAccounts)
	}
}
