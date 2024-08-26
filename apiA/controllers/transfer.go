package controllers

import (
	"apiA/commands"
	"apiA/services"
	"github.com/gofiber/fiber/v3"
	"log"
	"time"
)

type TransferController interface {
	Transfer(c fiber.Ctx) error
	TransferTransactions(c fiber.Ctx) error
}

type transferController struct {
	transferService services.TransferService
}

func NewTransferController(transferService services.TransferService) TransferController {
	return transferController{transferService}
}

func (obj transferController) Transfer(c fiber.Ctx) error {
	command := commands.TransferCommand{}

	err := c.Bind().JSON(&command)
	if err != nil {
		log.Println(err)
		return returnError(c, err)
	}

	start := time.Now()
	err = obj.transferService.Transfer(command)
	elapsed := time.Since(start)
	log.Printf("transfer process took %s", elapsed)

	if err != nil {
		log.Println(err)
		return returnError(c, err)
	}

	c.Status(fiber.StatusOK)
	return returnSuccess(c, command.RefID)
}

func (obj transferController) TransferTransactions(c fiber.Ctx) error {
	start := time.Now()

	transactions, err := obj.transferService.TransferTransactions()
	if err != nil {
		log.Println(err)
		return returnError(c, err)
	}

	elapsed := time.Since(start)
	log.Printf("get transfer transactions process took %s", elapsed)

	c.Status(fiber.StatusOK)
	return returnSuccessBody(c, transactions)
}

func returnError(c fiber.Ctx, err error) error {
	return c.JSON(fiber.Map{
		"success": false,
		"message": "error: " + err.Error(),
	})
}

func returnSuccess(c fiber.Ctx, id string) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "success",
		"refId":   id,
	})
}

func returnSuccessBody(c fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    data,
	})
}
