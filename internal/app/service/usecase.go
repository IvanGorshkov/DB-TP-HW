package service

import (
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
)


type ServiceUsecase interface {
	getStatus() (*models.Status, *errors.Error)
}
