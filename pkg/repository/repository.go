package repository

import (
	"clearWayTest/pkg/models"
	"github.com/jackc/pgx/v5"
)

type Authorization interface {
	GetUserByLogin(login string) (*models.User, error)
	SaveSession(session *models.Session) (string, error)
	GetSession(sessionId string) (*models.Session, error)
}

type Asset interface {
	SaveAsset(asset *models.Asset) error
	GetAsset(name string) (*models.Asset, error)
	GetList(limit, offset int) (*[]models.Asset, error)
	DeleteByName(name string) error
}

type Repository struct {
	Authorization
	Asset
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		Authorization: NewAuthorization(db),
		Asset:         NewAsset(db),
	}
}
