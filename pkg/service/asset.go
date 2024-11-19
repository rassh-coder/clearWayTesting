package service

import (
	"clearWayTest/pkg/models"
	"clearWayTest/pkg/repository"
	"database/sql"
	"errors"
)

type AssetService struct {
	repos repository.Asset
}

func NewAsset(repos repository.Asset) *AssetService {
	return &AssetService{repos: repos}
}

func (s *AssetService) SaveAsset(asset *models.Asset) error {
	err := s.repos.SaveAsset(asset)
	if err != nil {
		return err
	}

	return nil
}

func (s *AssetService) GetAsset(name string, uid uint) (*models.Asset, error) {
	asset, err := s.repos.GetAsset(name)
	if err != nil {
		return nil, err
	}

	if asset.UID != uid {
		return nil, errors.New("forbidden")
	}

	return asset, nil
}

func (s *AssetService) GetList(limit, offset int) (*[]models.Asset, error) {
	assets, err := s.repos.GetList(limit, offset)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (s *AssetService) DeleteByName(name string, uid uint) error {
	asset, err := s.GetAsset(name, uid)
	if err != nil {
		return err
	}

	if asset != nil {
		err = s.repos.DeleteByName(name)
		if err != nil {
			return err
		}
	} else {
		return sql.ErrNoRows
	}

	return nil
}
