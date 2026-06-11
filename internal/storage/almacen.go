package storage

import "modulo3_go/internal/models"

type Storage interface {
	CreateChofer(models.Chofer) (int, error)
	ListChoferes() ([]models.Chofer, error)
	GetChofer(int) (models.Chofer, error)
	UpdateChofer(int, models.Chofer) error
	DeleteChofer(int) error

	CreateMantenimiento(models.Mantenimiento) (int, error)
	ListMantenimientos() ([]models.Mantenimiento, error)
	GetMantenimiento(int) (models.Mantenimiento, error)
	UpdateMantenimiento(int, models.Mantenimiento) error
	DeleteMantenimiento(int) error

	Close() error
}
