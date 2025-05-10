package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/balqadishaPRO/Emoji-Hub/internal/service"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
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
	emojis, err := service.FetchAll()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(">> got %d emojis\n", len(emojis))

	// 4. prepare
	log.Println(">> preparing statement …")
	stmt, err := db.PrepareContext(context.Background(), `
        INSERT INTO emoji(id,name,category,"group",html_code,unicode)
        VALUES(gen_random_uuid(),$1,$2,$3,$4,$5)
        ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Fatal("prepare:", err)
	}

	// 5. insert
	for i, e := range emojis {
		if _, err := stmt.Exec(e.Name, e.Category, e.Group,
			pq.StringArray(e.HtmlCode), pq.StringArray(e.Unicode)); err != nil {
			log.Println("insert:", err)
		}
		if i%100 == 0 {
			log.Printf(">> inserted %d/%d\n", i, len(emojis))
		}
	}
	log.Println(">> import done:", len(emojis))
}
