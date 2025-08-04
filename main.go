package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ysearch/handlers"
	"ysearch/storage"
	"ysearch/types"
)

func main() {
	var db *pgxpool.Pool
	conn := types.LoadSecret("DATABASE_URL")

	config, err := pgxpool.ParseConfig(conn)
	if err != nil {
		log.Fatal(err)
	}

	// Disable statement caching
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	db, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Println(err)
	}

	database := storage.NewStore(db)

	handler := handlers.NewHandlers(database)
	r := chi.NewRouter()

	r.Get("/", handler.HandleForm)
	r.Get("/search/cards", handler.HandleSearch)
	r.Get("/search/singlecard", handler.HandleCardModal)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
