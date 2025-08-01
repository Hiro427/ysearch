package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"ysearch/api"
	"ysearch/views"
)

func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	// TODO: Split the code up into reusable functions
	// TODO: Add cap to limit of search results query

	name := strings.TrimSpace(q.Get("name"))
	cardType := strings.TrimSpace(q.Get("type"))
	effect := strings.TrimSpace(q.Get("effect"))

	races := filterEmpty(q["race"])
	banlistinfo := filterEmpty(q["banlist"])
	linkmarkers := filterEmpty(q["linkmarkers"])

	atk := strings.TrimSpace(q.Get("atk"))
	def := strings.TrimSpace(q.Get("def"))
	attribute := filterEmpty(q["attribute"])

	linkval := strings.TrimSpace(q.Get("linkval"))
	level := strings.TrimSpace(q.Get("level"))

	scale := strings.TrimSpace(q.Get("scale"))
	limit := strings.TrimSpace(q.Get("limit"))

	var where []string
	var args []interface{}
	argIdx := 1

	ArgumentMultiLike(name, "name", &argIdx, &args, &where)
	ArgumentMultiLike(effect, "\"desc\"", &argIdx, &args, &where)
	ArgumentMultiLike(cardType, "humanreadablecardtype", &argIdx, &args, &where)
	ArgumentNumComparison(atk, "atk", &argIdx, &args, &where)
	ArgumentNumComparison(def, "def", &argIdx, &args, &where)
	ArgumentAddArray(races, "race", &argIdx, &args, &where)
	ArgumentAddArray(attribute, "attribute", &argIdx, &args, &where)
	ArgumentNumComparison(linkval, "linkval", &argIdx, &args, &where)
	ArgumentNumComparison(scale, "scale", &argIdx, &args, &where)
	ArgumentNumComparison(level, "level", &argIdx, &args, &where)

	// Array conditions
	if len(banlistinfo) > 0 {
		where = append(where, "banlistinfo @> $"+strconv.Itoa(argIdx))
		args = append(args, banlistinfo)
		argIdx++
	}
	if len(linkmarkers) > 0 {
		where = append(where, "linkmarkers @> $"+strconv.Itoa(argIdx))
		args = append(args, linkmarkers)
		argIdx++
	}

	sql := `SELECT source_id, name, "desc", humanreadablecardtype, race, atk, def, archetype, attribute, level, linkval, linkmarkers, scale, banlistinfo FROM all_cards_search`

	if len(where) > 0 {
		sql += " WHERE " + strings.Join(where, " AND ")
	}

	sql += " ORDER by name"

	if limit != "" {
		sql += fmt.Sprintf(" LIMIT %s", limit)
	} else {
		sql += " LIMIT 40"
	}

	//For testing
	fmt.Println("Final SQL:", sql)
	fmt.Println("Args:", args)

	// Query the database
	cards, err := api.SearchCards(ctx, h.db.DB, sql, args)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	views.CardsList(cards).Render(r.Context(), w)
}

// Helper function to filter empty values
func filterEmpty(values []string) []string {
	var filtered []string
	for _, v := range values {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			filtered = append(filtered, trimmed)
		}
	}
	return filtered
}

func ArgumentNumComparison(val string, parameter string, argIdx *int, args *[]interface{}, where *[]string) {
	if val != "" {
		chars := []rune(val)
		a := string(chars[1:])
		switch s := string(chars[0]); s {
		case ">":
			*where = append(*where, fmt.Sprintf("%s > $", parameter)+strconv.Itoa(*argIdx))
			*args = append(*args, a)
		case "<":
			*where = append(*where, fmt.Sprintf("%s < $", parameter)+strconv.Itoa(*argIdx))
			*args = append(*args, a)
		default:
			*where = append(*where, fmt.Sprintf("%s = $", parameter)+strconv.Itoa(*argIdx))
			*args = append(*args, val)
		}
		*argIdx++
	}
}

func ArgumentAddArray(vals []string, parameter string, argIdx *int, args *[]interface{}, where *[]string) {
	if len(vals) > 0 {
		*where = append(*where, fmt.Sprintf("%s ~* $", parameter)+strconv.Itoa(*argIdx))
		regex := "^(" + strings.Join(vals, "|") + ")"
		*args = append(*args, regex)
		*argIdx++
	}
}

func ArgumentMultiLike(vals string, parameter string, argIdx *int, args *[]interface{}, where *[]string) {
	if vals != "" {
		if strings.Contains(vals, ",") {
			newVals := strings.Split(vals, ",")

			// var newArgs []string
			for _, a := range newVals {
				chars := []rune(a)
				if string(chars[0]) == "!" {
					*where = append(*where, fmt.Sprintf("%s NOT ILIKE $"+strconv.Itoa(*argIdx), parameter))
					*args = append(*args, "%"+string(chars[1:])+"%")
				} else {
					*where = append(*where, fmt.Sprintf("%s ILIKE $"+strconv.Itoa(*argIdx), parameter))
					*args = append(*args, "%"+a+"%")

				}
				*argIdx++
			}
			// *where = append(*where, "("+strings.Join(newArgs, " OR ")+")")

		} else {
			chars := []rune(vals)
			if string(chars[0]) == "!" {
				*where = append(*where, fmt.Sprintf("%s NOT ILIKE $"+strconv.Itoa(*argIdx), parameter))
				*args = append(*args, "%"+string(chars[1:])+"%")
			} else {
				*where = append(*where, fmt.Sprintf("%s ILIKE $"+strconv.Itoa(*argIdx), parameter))
				*args = append(*args, "%"+vals+"%")
			}
			*argIdx++
		}
	}
}
