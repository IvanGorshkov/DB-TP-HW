package service

import "github.com/IvanGorshkov/DB-TP-HW/internal/app/models"


type ServiceRepository interface {
	getStatus() (*models.Status, error)
}


