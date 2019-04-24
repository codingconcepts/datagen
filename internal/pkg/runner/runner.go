package runner

import (
	"bytes"
	"database/sql"
	"log"
	"math/rand"
	"reflect"
	"text/template"

	"github.com/google/uuid"

	"github.com/codingconcepts/datagen/internal/pkg/random"

	"github.com/codingconcepts/datagen/internal/pkg/parse"
	"github.com/pkg/errors"
)

type runner struct {
	db           *sql.DB
	funcs        template.FuncMap
	helpers      map[string]interface{}
	context      map[string][]map[string]interface{}
	contextGroup map[rowKey]map[string]interface{}
}

type rowKey struct {
	groupType interface{}
	groupID   int
}

func New(db *sql.DB) *runner {
	r := runner{
		db:           db,
		context:      map[string][]map[string]interface{}{},
		contextGroup: map[rowKey]map[string]interface{}{},
	}

	r.funcs = template.FuncMap{
		"s":    random.String,
		"i":    random.Int,
		"d":    random.Date,
		"f":    random.Float,
		"uuid": func() string { return uuid.New().String() },
		"set":  random.Set,
		"ref":  r.reference,
		"row":  r.row,
	}

	r.helpers = map[string]interface{}{
		"times_1":      make([]struct{}, 1),
		"times_10":     make([]struct{}, 10),
		"times_100":    make([]struct{}, 100),
		"times_1000":   make([]struct{}, 1000),
		"times_10000":  make([]struct{}, 10000),
		"times_100000": make([]struct{}, 100000),
	}

	return &r
}

// Run executes a given block, returning any errors encountered.
func (r *runner) Run(b parse.Block) error {
	tmpl := template.Must(template.New("block").Funcs(r.funcs).Parse(b.Body))

	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, r.helpers); err != nil {
		return errors.Wrap(err, "executing template")
	}

	rows, err := r.db.Query(buf.String())
	if err != nil {
		return errors.Wrap(err, "executing query")
	}

	return r.scan(b, rows)
}

func (r *runner) scan(b parse.Block, rows *sql.Rows) error {
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

		curr := map[string]interface{}{}
		for i, ct := range columnTypes {
			values[i] = reflect.ValueOf(values[i]).Elem()
			curr[ct.Name()] = values[i]
		}
		r.context[b.Name] = append(r.context[b.Name], curr)
	}

	return nil
}

func (r *runner) reference(key string, column string) interface{} {
	rows, ok := r.context[key]
	if !ok {
		log.Fatalf("key %v not found in context", key)
	}

	value, ok := rows[rand.Intn(len(rows))][column]
	if !ok {
		log.Fatalf("key %v not found in context", key)
	}

	return value
}

func (r *runner) row(key string, column string, i int) interface{} {
	groupKey := rowKey{groupType: key, groupID: i}

	// Check if we've scanned this row before.
	row, ok := r.contextGroup[groupKey]
	if ok {
		return row[column]
	}

	// Get a random item from the row context and cache it for the next read.
	var randomValue map[string]interface{}
	for _, v := range r.context {
		randomValue = v[rand.Intn(len(v))]
		break
	}

	r.contextGroup[groupKey] = randomValue

	return randomValue[column]
}
