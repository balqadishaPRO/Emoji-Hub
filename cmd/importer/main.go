package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/balqadishaPRO/Emoji-Hub/internal/model"
	"github.com/balqadishaPRO/Emoji-Hub/internal/repo"
	"github.com/balqadishaPRO/Emoji-Hub/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env, read OS env vars")
	}
	dsn := os.Getenv("DATABASE_URL")
	log.Println(">> DSN:", dsn)

	log.Println(">> opening DB …")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("ping:", err)
	}
	log.Println(">> DB connected")

	log.Println(">> fetching Emojihub …")
	rawEmojis, err := service.FetchAll()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(">> got %d emojis\n", len(rawEmojis))

	emojis := make([]model.Emoji, len(rawEmojis))
	for i, raw := range rawEmojis {
		emojis[i] = model.Emoji{
			ID:       raw.ID,
			Name:     raw.Name,
			Category: raw.Category,
			Group:    raw.Group,
			HtmlCode: raw.HtmlCode,
			Unicode:  raw.Unicode,
		}
	}

	repository, err := repo.New(dsn)
	if err != nil {
		log.Fatal("repo:", err)
	}

	svc := &service.EmojiService{Repo: repository}
	if err := svc.ImportEmojis(context.Background(), emojis); err != nil {
		log.Fatal("import:", err)
	}

	log.Println(">> import done:", len(emojis))
}
