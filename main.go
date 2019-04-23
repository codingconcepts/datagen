package main

import (
	"database/sql"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/codingconcepts/datagen/internal/pkg/parse"
	"github.com/codingconcepts/datagen/internal/pkg/runner"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	driver := flag.String("driver", "", "name of the database driver to use [postgres|mysql]")
	script := flag.String("script", "", "the full or relative path to your script file")
	conn := flag.String("conn", "", "the database connection string")
	flag.Parse()

	if *script == "" || *driver == "" || *conn == "" {
		flag.Usage()
		os.Exit(2)
	}

	db := mustConnect(*driver, *conn)
	defer db.Close()

	runner := runner.New(db)

	file, err := os.Open(*script)
	if err != nil {
		log.Fatalf("error reading script file: %v", err)
	}
	defer file.Close()

	blocks, err := parse.Blocks(file)
	if err != nil {
		log.Fatalf("error reading blocks from script file: %v", err)
	}

	for _, block := range blocks {
		for i := 0; i < block.Repeat; i++ {
			if err = runner.Run(block); err != nil {
				log.Fatalf("error running block: %v", err)
			}
		}
	}
}

func mustConnect(driver, connStr string) *sql.DB {
	conn, err := sql.Open(driver, connStr)
	if err != nil {
		log.Fatalf("error opening connection: %d", err)
	}

	if err = conn.Ping(); err != nil {
		log.Fatalf("error checking connection: %v", err)
	}

	return conn
}
