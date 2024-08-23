package controllers

import (
	"apiA/commands"
	"apiA/services"
	"github.com/gofiber/fiber/v3"
	"log"
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
		return err
	}
	err = obj.transferService.Transfer(command)
	if err != nil {
		log.Println(err)
		return err
	}

	// todo loop waiting for transaction status is COMPLETED

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "transfer successfully",
		"refId":   command.RefID,
	})
}
