package handlers

import (
	"fmt"
	"net/http"
	"ysearch/api"
	"ysearch/views"
)

func (h *Handler) HandleCardModal(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	ctx := r.Context()

	id := q.Get("id")

	sql := `SELECT source_id, name, "desc", humanreadablecardtype, race, atk, def, archetype, attribute, level, linkval, linkmarkers, scale, banlistinfo  FROM all_cards_search WHERE source_id = $1`
	fmt.Println(id)

	card, err := api.SingleSearchCards(ctx, h.db.DB, sql, id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	views.SingleCardModal(card).Render(r.Context(), w)

}
