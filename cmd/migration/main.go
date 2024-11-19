package main

import (
	"clearWayTest/config"
	"clearWayTest/pkg/repository"
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5"
	"io"
	"log"
	"os"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config err: %s", err)
	}

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Can't open db connect: %s", err)
	}

	defer db.Close(context.Background())

	action := flag.String("action", "", "")
	flag.Parse()
	if string(*action) == "up" {
		Up(db)
		fmt.Println("Migration up successfully!")
	}

	if string(*action) == "down" {
		Down(db)
		fmt.Println("Migration down successfully!")
	}
}

func Up(dbConn *pgx.Conn) {
	dir, err := os.Open("./cmd/migration/up")

	if err != nil {
		log.Fatalf("Can't open user migration file: %s", err)
	}

	defer func() {
		err := dir.Close()
		if err != nil {
			log.Fatalf("Can't close migration dir: %s", err)
		}
	}()

	files, err := dir.Readdir(-1)

	for _, file := range files {
		f, err := os.Open(fmt.Sprintf("%s/%s", dir.Name(), file.Name()))

		if err != nil {
			log.Fatalf("Can't open migration file: %s", err)
		}

		res, err := io.ReadAll(f)

		resStr := string(res)
		_, err = dbConn.Exec(context.Background(), resStr)

		if err != nil {
			log.Fatalf("Can't exec migration: %s", err)
		}

		err = f.Close()
		if err != nil {
			log.Printf("Can't close migration file: %s", err)
		}
	}
}

func Down(dbConn *pgx.Conn) {
	dir, err := os.Open("./cmd/migration/down")

	if err != nil {
		log.Fatalf("Can't open user migration file: %s", err)
	}

	defer func() {
		err := dir.Close()
		if err != nil {
			log.Fatalf("Can't close migration dir: %s", err)
		}
	}()

	files, err := dir.Readdir(-1)

	for _, file := range files {
		f, err := os.Open(fmt.Sprintf("%s/%s", dir.Name(), file.Name()))

		if err != nil {
			log.Fatalf("Can't open migration file: %s", err)
		}

		res, err := io.ReadAll(f)

		resStr := string(res)
		_, err = dbConn.Exec(context.Background(), resStr)

		if err != nil {
			log.Fatalf("Can't exec migration: %s", err)
		}

		err = f.Close()
		if err != nil {
			log.Fatalf("Can't close migration file: %s", err)
		}
	}
}
