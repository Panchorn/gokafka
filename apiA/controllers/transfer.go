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

func returnError(c fiber.Ctx, err error) error {
	return c.JSON(fiber.Map{
		"success": false,
		"message": "transfer error: " + err.Error(),
	})
}

func returnSuccess(c fiber.Ctx, id string) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "transfer successfully",
		"refId":   id,
	})
}
