package runner

import (
	"bytes"
	"database/sql"
	"log"
	"math/rand"
	"reflect"
	"text/template"
	"time"

	"github.com/codingconcepts/datagen/internal/pkg/random"

	"github.com/google/uuid"

	"github.com/codingconcepts/datagen/internal/pkg/parse"
	"github.com/pkg/errors"
)

var logFatalf = log.Fatalf

type Runner struct {
	db           *sql.DB
	funcs        template.FuncMap
	helpers      map[string]interface{}
	context      map[string][]map[string]interface{}
	contextGroup map[rowKey]map[string]interface{}

	dateFormat string
}

type rowKey struct {
	groupType interface{}
	groupID   int
}

func New(db *sql.DB, opts ...Option) *Runner {
	r := Runner{
		db:           db,
		context:      map[string][]map[string]interface{}{},
		contextGroup: map[rowKey]map[string]interface{}{},
	}

	for _, opt := range opts {
		opt(&r)
	}

	r.funcs = template.FuncMap{
		"string": random.String,
		"int":    random.Int,
		"date":   random.Date(r.dateFormat),
		"float":  random.Float,
		"uuid":   func() string { return uuid.New().String() },
		"set":    random.Set,
		"ref":    r.reference,
		"row":    r.row,
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
func (r *Runner) Run(b parse.Block) error {
	tmpl, err := template.New("block").Funcs(r.funcs).Parse(b.Body)
	if err != nil {
		return errors.Wrap(err, "parsing template")
	}

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

func (r *Runner) scan(b parse.Block, rows *sql.Rows) error {
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
			values[i] = r.prepareValue(reflect.ValueOf(values[i]).Elem())
			curr[ct.Name()] = values[i]
		}
		r.context[b.Name] = append(r.context[b.Name], curr)
	}

	return nil
}

// prepareValue ensures that data being read out of the database following
// a scan is in the correct format for being re-inserted into the database
// during follow-up queries.
func (r *Runner) prepareValue(v reflect.Value) interface{} {
	switch v.Type() {
	case reflect.TypeOf(time.Time{}):
		t := v.Interface().(time.Time)
		return t.Format(r.dateFormat)
	default:
		return v
	}
}

func (r *Runner) reference(key string, column string) interface{} {
	rows, ok := r.context[key]
	if !ok {
		logFatalf("key %v not found in context", key)
		return nil // Break out early for tests.
	}

	value, ok := rows[rand.Intn(len(rows))][column]
	if !ok {
		logFatalf("key %v not found in context", key)
		return nil // Break out early for tests.
	}

	return value
}

func (r *Runner) row(key string, column string, group int) interface{} {
	groupKey := rowKey{groupType: key, groupID: group}

	// Check if we've scanned this row before.
	row, ok := r.contextGroup[groupKey]
	if ok {
		value, ok := row[column]
		if !ok {
			logFatalf("key %v not found in context", key)
			return nil // Break out early for tests.
		}
		return value
	}

	// Get a random item from the row context and cache it for the next read.
	var randomValue map[string]interface{}
	for _, v := range r.context {
		randomValue = v[rand.Intn(len(v))]
		break
	}

	r.contextGroup[groupKey] = randomValue

	value, ok := randomValue[column]
	if !ok {
		logFatalf("key %v not found in context", key)
		return nil // Break out early for tests.
	}

	return value
}
