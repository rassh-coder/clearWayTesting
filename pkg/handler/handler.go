package handler

import (
	"clearWayTest/pkg/service"
	"net/http"
)

type Handler struct {
	Services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{Services: services}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	apirouter := http.NewServeMux()
	apirouter.HandleFunc("/auth", h.LogIn)
	apirouter.Handle("/upload-asset/{assetName}", h.authMiddleware(http.HandlerFunc(h.UploadAsset)))
	apirouter.Handle("/asset/{assetName}", h.authMiddleware(http.HandlerFunc(h.DownloadAsset)))
	apirouter.Handle("/asset/get-list", h.authMiddleware(http.HandlerFunc(h.AssetGetList)))
	apirouter.Handle("/asset/delete/{assetName}", h.authMiddleware(http.HandlerFunc(h.AssetDelete)))

	router := http.NewServeMux()
	router.Handle("/api/", http.StripPrefix("/api", apirouter))

	return router
}
