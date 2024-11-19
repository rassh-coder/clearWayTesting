package handler

import (
	"clearWayTest/pkg/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (h *Handler) UploadAsset(w http.ResponseWriter, r *http.Request) {
	var asset models.Asset
	w.Header().Set("Content-Type", "application/json")
	assetName := r.PathValue("assetName")

	if r.Method != http.MethodPost {
		h.NewNotFound(w)
		return
	}

	if assetName == "" {
		h.NewNotFound(w)
		return
	}

	uid := r.Context().Value(ctxtUidKey)

	if uid == "" || uid == 0 {
		errRes := ErrorResponse{Error: "unauthorized"}
		h.NewError(w, &errRes, http.StatusUnauthorized)
		return
	}

	data, errRead := io.ReadAll(r.Body)
	if errRead != nil {
		errRes := ErrorResponse{Error: "can't read body"}
		h.NewError(w, &errRes, http.StatusBadRequest)
	}
	defer r.Body.Close()

	asset.Name = assetName
	asset.UID = uid.(uint)
	asset.Data = &data
	err := h.Services.Asset.SaveAsset(&asset)
	if err != nil {
		errRes := ErrorResponse{Error: err.Error()}
		h.NewError(w, &errRes, http.StatusInternalServerError)
		return
	}

	sucRes := StatusResponse{Status: "ok"}
	jsonSucRes, err := json.Marshal(&sucRes)
	if err != nil {
		errRes := ErrorResponse{Error: err.Error()}
		h.NewError(w, &errRes, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonSucRes)
	if err != nil {
		http.Error(w, "can't write request", http.StatusInternalServerError)
		return
	}
	return
}

func (h *Handler) DownloadAsset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.NewNotFound(w)
		return
	}

	assetName := r.PathValue("assetName")
	if assetName == "" {
		h.NewNotFound(w)
		return
	}

	uid := r.Context().Value(ctxtUidKey)
	if uid == nil || uid == "" || uid == 0 {
		errRes := ErrorResponse{Error: "unauthorized"}
		h.NewError(w, &errRes, http.StatusUnauthorized)
		return
	}

	asset, err := h.Services.Asset.GetAsset(assetName, uid.(uint))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errRes := ErrorResponse{Error: "asset is not found"}
			h.NewError(w, &errRes, http.StatusBadRequest)
			return
		}

		if err.Error() == "forbidden" {
			errRes := ErrorResponse{Error: err.Error()}
			h.NewError(w, &errRes, http.StatusForbidden)
			return
		}

		errRes := ErrorResponse{Error: err.Error()}
		h.NewError(w, &errRes, http.StatusInternalServerError)
		return
	}

	if asset == nil {
		errRes := ErrorResponse{Error: "asset is not found"}
		h.NewError(w, &errRes, http.StatusBadRequest)
		return
	}

	if asset.Data != nil {
		_, err = w.Write(*asset.Data)
		if err != nil {
			errRes := ErrorResponse{Error: fmt.Sprintf("can't write body: %s", err)}
			h.NewError(w, &errRes, http.StatusInternalServerError)
			return
		}
		return
	}

	errRes := ErrorResponse{Error: "can't get asset"}
	h.NewError(w, &errRes, http.StatusInternalServerError)
	return
}

func (h *Handler) AssetGetList(w http.ResponseWriter, r *http.Request) {
	var list ListStructure
	if r.Method != http.MethodPost {
		h.NewNotFound(w)
		return
	}

	err := DecodeJSONBody(w, r, &list)
	if err != nil {
		errRes := ErrorResponse{Error: "invalid json"}
		h.NewError(w, &errRes, http.StatusBadRequest)
		return
	}

	assets, errService := h.Services.Asset.GetList(list.Limit, list.Offset)
	res := AssetList{
		Data:   *assets,
		Status: "ok",
	}
	if errService != nil {
		if errors.Is(errService, sql.ErrNoRows) {
			assetJson, err := json.Marshal(&res)
			if err != nil {
				errRes := ErrorResponse{Error: fmt.Sprintf("can't marshal json: %s", err.Error())}
				h.NewError(w, &errRes, http.StatusInternalServerError)
				return
			}
			_, err = w.Write(assetJson)
			if err != nil {
				errRes := ErrorResponse{Error: err.Error()}
				h.NewError(w, &errRes, http.StatusInternalServerError)
				return
			}
			return
		}

		errRes := ErrorResponse{Error: "can't fetch asset list"}
		h.NewError(w, &errRes, http.StatusInternalServerError)
		return
	}

	assetJson, errJson := json.Marshal(&res)
	if errJson != nil {
		errRes := ErrorResponse{Error: fmt.Sprintf("can't marshal json: %s", errJson.Error())}
		h.NewError(w, &errRes, http.StatusInternalServerError)
		return
	}

	_, errWrite := w.Write(assetJson)
	if errWrite != nil {
		errRes := ErrorResponse{Error: errWrite.Error()}
		h.NewError(w, &errRes, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AssetDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.NewNotFound(w)
		return
	}

	assetName := r.PathValue("assetName")
	if assetName == "" {
		h.NewNotFound(w)
		return
	}

	uid := r.Context().Value(ctxtUidKey)
	if uid == nil || uid == 0 || uid == "" {
		errRes := ErrorResponse{Error: "unauthorized"}
		h.NewError(w, &errRes, http.StatusUnauthorized)
		return
	}

	err := h.Services.Asset.DeleteByName(assetName, uid.(uint))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errRes := ErrorResponse{Error: "asset is not found"}
			h.NewError(w, &errRes, http.StatusBadRequest)
			return
		}

		if err.Error() == "forbidden" {
			errRes := ErrorResponse{err.Error()}
			h.NewError(w, &errRes, http.StatusForbidden)
			return
		}

		errRes := ErrorResponse{Error: err.Error()}
		h.NewError(w, &errRes, http.StatusInternalServerError)
		return
	}

	res := StatusResponse{Status: "ok"}
	jsonRes, err := json.Marshal(&res)
	if err != nil {
		errRes := ErrorResponse{Error: "can't marshal res"}
		h.NewError(w, &errRes, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonRes)
	if err != nil {
		errRes := ErrorResponse{Error: err.Error()}
		h.NewError(w, &errRes, http.StatusInternalServerError)
		return
	}
}
