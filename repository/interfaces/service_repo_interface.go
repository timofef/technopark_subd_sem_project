package interfaces

import "github.com/timofef/technopark_subd_sem_project/models"

type ServiceRepository interface {
	GetStatus() (*models.Status, error)
	Clear() (error)
}