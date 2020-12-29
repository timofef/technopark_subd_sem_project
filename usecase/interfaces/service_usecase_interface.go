package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type ServiceUsecase interface {
	GetStatus() (*models.Status, error)
	Clear() (error)
}
