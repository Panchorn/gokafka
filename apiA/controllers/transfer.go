package controllers

import (
	"apiA/commands"
	"apiA/services"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"logs"
	"net/http"
	"time"
)

type TransferController interface {
	Transfer(c echo.Context) error
	TransferTransactions(c echo.Context) error
}

type transferController struct {
	transferService services.TransferService
}

func NewTransferController(transferService services.TransferService) TransferController {
	return transferController{transferService}
}

func (obj transferController) Transfer(c echo.Context) error {
	requestID := getRequestID(c)
	command := commands.TransferCommand{}

	err := json.NewDecoder(c.Request().Body).Decode(&command)
	if err != nil {
		logs.Error(requestID, err)
		return returnError(c, err)
	}

	start := time.Now()
	err = obj.transferService.Transfer(c, command)
	elapsed := time.Since(start)
	logs.Info(requestID, "transfer process took "+elapsed.String())

	if err != nil {
		logs.Error(requestID, err)
		return returnError(c, err)
	}

	return returnSuccess(c, command.RefID)
}

func (obj transferController) TransferTransactions(c echo.Context) error {
	requestID := getRequestID(c)
	start := time.Now()

	transactions, err := obj.transferService.TransferTransactions(c)
	if err != nil {
		logs.Error(requestID, err)
		return returnError(c, err)
	}

	elapsed := time.Since(start)
	logs.Info(requestID, "get transfer transactions process took "+elapsed.String())

	return returnSuccessBody(c, transactions)
}

func returnError(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"success": false,
		"message": "error: " + err.Error(),
	})
}

func returnSuccess(c echo.Context, id string) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "success",
		"refId":   id,
	})
}

func returnSuccessBody(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "success",
		"data":    data,
	})
}

func logger(c echo.Context) *zap.Logger {
	//zap.driver
	ctx := logs.CreateLogContext(c)
	return ctx.Get(logs.RequestIDLogger).(*zap.Logger)
}

func getRequestID(c echo.Context) string {
	requestID := c.Response().Header().Get("X-Request-ID")
	c.Set(logs.RequestID, requestID)
	return requestID
}
