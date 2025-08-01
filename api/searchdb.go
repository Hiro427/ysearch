package api

import (
	"context"
	"fmt"
	"log"
	"ysearch/types"

	// "log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SearchCards(ctx context.Context, db *pgxpool.Pool, sql string, args []any) ([]types.Card, error) {
	fmt.Println("SearchCards called")
	fmt.Println("SQL:", sql)
	fmt.Println("Args:", args)

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		fmt.Println("Query error:", err)
		return nil, err // Return the error instead of ignoring it
	}
	defer rows.Close() // Always close rows

	fmt.Println("Query completed")

	var cards []types.Card
	for rows.Next() {
		var c types.Card
		if err := rows.Scan(
			&c.SourceID,
			&c.Name,
			&c.Desc,
			&c.HumanReadableCardType,
			&c.Race,
			&c.Atk,
			&c.Def,
			&c.Archetype,
			&c.Attribute,
			&c.Level,
			&c.Linkval,
			&c.Linkmarkers,
			&c.Scale,
			&c.Banlistinfo,
		); err != nil {
			log.Println("Scan error:", err)
			return nil, err // Return scan errors
		}
		cards = append(cards, c)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		fmt.Println("Rows iteration error:", err)
		return nil, err
	}

	fmt.Printf("Found %d cards\n", len(cards))
	return cards, nil
}
