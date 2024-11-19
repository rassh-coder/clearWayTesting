package service

import (
	"clearWayTest/pkg/models"
	"clearWayTest/pkg/repository"
)

type Authorization interface {
	SignIn(login, password, ip string) (string, error)
	GetSessionByToken(token string) (*models.Session, error)
}

type Asset interface {
	SaveAsset(asset *models.Asset) error
	GetAsset(name string, uid uint) (*models.Asset, error)
	GetList(limit, offset int) (*[]models.Asset, error)
	DeleteByName(name string, uid uint) error
}

type Service struct {
	Authorization
	Asset
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthorization(repos.Authorization),
		Asset:         NewAsset(repos.Asset),
	}
}
