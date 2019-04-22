package main

import (
	"bytes"
	"database/sql"
	"flag"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/codingconcepts/datagen/internal/pkg/parse"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	funcs = template.FuncMap{
		"s":     randomString,
		"i":     randomInt,
		"d":     randomDate,
		"f32":   randomFloat32,
		"f64":   randomFloat64,
		"rep":   repeat,
		"ref":   reference,
		"join":  strings.Join,
		"times": times,
	}

	helperArgs = map[string]interface{}{
		"times_1":      make([]struct{}, 1),
		"times_10":     make([]struct{}, 10),
		"times_100":    make([]struct{}, 100),
		"times_1000":   make([]struct{}, 1000),
		"times_10000":  make([]struct{}, 10000),
		"times_100000": make([]struct{}, 100000),
	}

	ascii = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

const (
	dateFormat = "2006-01-02 15:04:05Z07:00"
)

var (
	context = map[interface{}][]interface{}{}
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
			if err = run(db, block); err != nil {
				log.Fatalf("error running block: %v", err)
			}
		}
	}
}

func run(db *sql.DB, b parse.Block) error {
	personTmpl := template.Must(template.New("block").Funcs(funcs).Parse(b.Body))

	buf := &bytes.Buffer{}
	if err := personTmpl.Execute(buf, helperArgs); err != nil {
		return errors.Wrap(err, "executing template")
	}

	rows, err := db.Query(buf.String())
	if err != nil {
		return errors.Wrap(err, "executing query")
	}

	for rows.Next() {
		columnTypes, err := rows.ColumnTypes()
		if err != nil {
			return errors.Wrap(err, "getting columns types from result")
		}

		values := make([]interface{}, len(columnTypes))
		for i, ct := range columnTypes {
			switch ct.DatabaseTypeName() {
			case "UUID":
				values[i] = reflect.New(reflect.TypeOf("")).Interface()
			default:
				values[i] = reflect.New(ct.ScanType()).Interface()
			}
		}

		if err = rows.Scan(values...); err != nil {
			return errors.Wrap(err, "scanning columns")
		}

		for i, ct := range columnTypes {
			values[i] = reflect.ValueOf(values[i]).Elem()
			context[ct.Name()] = append(context[ct.Name()], values[i])
		}
	}
	return nil
}

func randomString(min, max int, prefix string) string {
	var length int
	if min >= max {
		length = min
	} else {
		length = between(min, max) - len(prefix)
	}

	output := make([]rune, length)
	for i := 0; i < length; i++ {
		output[i] = ascii[rand.Intn(len(ascii))]
	}

	return prefix + string(output)
}

func randomInt(min, max int) int {
	return between(min, max)
}

func randomDate(minStr, maxStr string) string {
	min, err := time.Parse(time.RFC3339, minStr)
	if err != nil {
		log.Fatalf("invalid min date: %v", err)
	}
	max, err := time.Parse(time.RFC3339, maxStr)
	if err != nil {
		log.Fatalf("invalid max date: %v", err)
	}

	diff := between64(min.Unix(), max.Unix())
	return time.Unix(diff, 0).Format(dateFormat)
}

func randomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func repeat(times int, input string, separator string) string {
	return "repeated"
}

func reference(key interface{}) interface{} {
	value, ok := context[key]
	if !ok {
		log.Fatalf("key %v not found in context", key)
	}

	return value[rand.Intn(len(value))]
}

func times(i int) []struct{} {
	return make([]struct{}, i)
}

func between(min, max int) int {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min) + min
}

func between64(min, max int64) int64 {
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	return rand.Int63n(max-min) + min
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
