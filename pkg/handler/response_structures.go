package handler

import "clearWayTest/pkg/models"

type ErrorResponse struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type AssetList struct {
	Data   []models.Asset `json:"data"`
	Status string         `json:"status"`
}
