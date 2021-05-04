package usecase

import (
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/errors"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/models"
	"github.com/IvanGorshkov/DB-TP-HW/internal/app/service"
)

type ServiceUsecase struct {
	serviceRepo service.ServiceRepository
}

func NewUserUsecase(repo service.ServiceRepository) service.ServiceUsecase {
	return &ServiceUsecase{
		serviceRepo: repo,
	}
}


func (su *ServiceUsecase) getStatus() (*models.Status, *errors.Error) {

	res, err := su.serviceRepo.getStatus()

	if err != nil {
		return nil, errors.UnexpectedInternal(err)
	}

	return res, nil
}