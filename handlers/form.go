package handlers

import "net/http"
import "ysearch/views"

func (h *Handler) HandleForm(w http.ResponseWriter, r *http.Request) {
	views.Search().Render(r.Context(), w)
}
