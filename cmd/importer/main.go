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
	// 1. загрузили .env
	if err := godotenv.Load(); err != nil {
		log.Println("no .env, read OS env vars")
	}
	dsn := os.Getenv("DATABASE_URL")
	log.Println(">> DSN:", dsn) // ← покажет, что реально прочиталось

	// 2. подключаемся к БД
	log.Println(">> opening DB …")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil { // ← проверяем сразу
		log.Fatal("ping:", err)
	}
	log.Println(">> DB connected")

	// 3. тянем эмодзи
	log.Println(">> fetching Emojihub …")
	rawEmojis, err := service.FetchAll()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(">> got %d emojis\n", len(rawEmojis))

	// Convert raw emojis to model emojis
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
