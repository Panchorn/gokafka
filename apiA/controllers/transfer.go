package controllers

import (
	"apiA/commands"
	"apiA/services"
	"github.com/labstack/echo/v4"
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
	command := commands.TransferCommand{}

	err := c.Bind(&command)
	if err != nil {
		logs.Error(err)
		return returnError(c, err)
	}

	start := time.Now()
	err = obj.transferService.Transfer(command)
	elapsed := time.Since(start)
	logs.Info("transfer process took " + elapsed.String())

	if err != nil {
		logs.Error(err)
		return returnError(c, err)
	}

	return returnSuccess(c, command.RefID)
}

func (obj transferController) TransferTransactions(c echo.Context) error {
	start := time.Now()

	transactions, err := obj.transferService.TransferTransactions()
	if err != nil {
		logs.Error(err)
		return returnError(c, err)
	}

	elapsed := time.Since(start)
	logs.Info("get transfer transactions process took " + elapsed.String())

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
