package types

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// TODO: Move to utils package
func LoadSecret(env string) string {
	var res string
	if os.Getenv(env) == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
		res = os.Getenv(env)
	}

	res = os.Getenv(env)

	return res
}

type Card struct {
	SourceID              int
	Name                  *string
	Desc                  *string
	HumanReadableCardType *string
	Race                  *string
	Atk                   *int
	Def                   *int
	Archetype             *string
	Scale                 *int
	Attribute             *string
	Level                 *int
	Linkval               *int
	Linkmarkers           *[]string
	Banlistinfo           *[]string
}
