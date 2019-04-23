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
	db      *sql.DB
	funcs   template.FuncMap
	helpers map[string]interface{}
	context map[interface{}][]interface{}
}

func New(db *sql.DB) *runner {
	r := runner{
		db:      db,
		context: map[interface{}][]interface{}{},
	}

	r.funcs = template.FuncMap{
		"s":    random.String,
		"i":    random.Int,
		"d":    random.Date,
		"f32":  random.Float32,
		"f64":  random.Float64,
		"uuid": func() string { return uuid.New().String() },
		"set":  random.Set,
		"ref":  r.reference,
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

			key := ct.Name()
			if b.Name != "" {
				key = b.Name + "_" + key
			}
			r.context[key] = append(r.context[key], values[i])
		}
	}
	return nil
}

func (r *runner) reference(key interface{}) interface{} {
	value, ok := r.context[key]
	if !ok {
		log.Fatalf("key %v not found in context", key)
	}

	return value[rand.Intn(len(value))]
}
