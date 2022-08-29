package implementations

import (
	"github.com/timofef/technopark_subd_sem_project/models"
	repository "github.com/timofef/technopark_subd_sem_project/repository/interfaces"
	usecase "github.com/timofef/technopark_subd_sem_project/usecase/interfaces"
)

type ServiceUsecase struct {
	serviceRepo repository.ServiceRepository
}

func NewServiceUsecase(serviceR repository.ServiceRepository) usecase.ServiceUsecase {
	return &ServiceUsecase{serviceRepo: serviceR}
}

func (s *ServiceUsecase) GetStatus() (*models.Status, error) {
	status, err := s.serviceRepo.GetStatus()

	return status, err
}

func (s *ServiceUsecase) Clear() error {
	err := s.serviceRepo.Clear()

	return err
}
