package handlers

import "ysearch/storage"

type Handler struct {
	db *storage.DB
}

func NewHandlers(db *storage.DB) *Handler {
	return &Handler{
		db: db,
	}
}
