package repository

import (
	"clearWayTest/pkg/models"
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

type AssetRepository struct {
	db *pgx.Conn
}

func NewAsset(db *pgx.Conn) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) SaveAsset(asset *models.Asset) error {
	query := `
		INSERT INTO assets (name, uid, data) values (@name, @uid, @data)
	`
	args := pgx.NamedArgs{
		"name": asset.Name,
		"uid":  asset.UID,
		"data": asset.Data,
	}

	_, err := r.db.Exec(context.Background(), query, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *AssetRepository) GetAsset(name string) (*models.Asset, error) {
	var asset models.Asset
	query := `
		SELECT * FROM assets
		WHERE name = @name
	`

	args := pgx.NamedArgs{
		"name": name,
	}

	row := r.db.QueryRow(context.Background(), query, args)

	err := row.Scan(&asset.Name, &asset.UID, &asset.Data, &asset.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *AssetRepository) GetList(limit, offset int) (*[]models.Asset, error) {
	var assets []models.Asset
	var query string
	if limit != 0 {
		query = `
		SELECT * FROM assets
			LIMIT @limit
		    OFFSET @offset
	`
	} else {
		query = `
		SELECT * FROM assets
		OFFSET @offset
	`
	}
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := r.db.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var asset models.Asset
		err := rows.Scan(&asset.Name, &asset.UID, &asset.Data, &asset.CreatedAt)
		if err != nil {
			log.Printf("error fetching asset details")
			return &assets, err
		}
		assets = append(assets, asset)
	}
	return &assets, nil
}

func (r *AssetRepository) DeleteByName(name string) error {
	query := `
		DELETE FROM assets WHERE name = @name	
	`
	args := pgx.NamedArgs{
		"name": name,
	}

	e, err := r.db.Exec(context.Background(), query, args)
	if err != nil {
		return err
	}

	if e.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
