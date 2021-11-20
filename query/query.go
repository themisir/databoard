package query

import (
	"database/sql"
	"fmt"
	"text/template"

	"github.com/themisir/databoard/domain"
)

type Query struct {
	name   string
	tmpl   *template.Template
	params []Parameter
}

type StringMap = map[string]string

type QueryResult struct {
	Columns []string
	Rows    []StringMap
}

// Prepare new query from give details. This function will parse given query
// using go templates and return Query pointer for executing parsed query.
// If parsing fails the function will return error.
func New(name string, query string) (*Query, error) {
	tmpl, err := template.New("query." + name).Parse(query)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query template: %s", err)
	}
	return &Query{
		name:   name,
		tmpl:   tmpl,
		params: []Parameter{},
	}, nil
}

// Add parameter
func (q *Query) AddParam(param Parameter) {
	q.params = append(q.params, param)
}

// Execute query on given database object using db.Exec with provided values.
func (q *Query) Exec(db *sql.DB, ctx *domain.Context, values StringMap) (sql.Result, error) {
	query, err := ctx.Transform(q.tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s", err)
	}

	args, err := q.prepareArgs(values)
	if err != nil {
		return nil, err
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute query: %s\n\n%s", err, query)
	}
	return result, nil
}

// Execute query on given database object using db.Query with provided values.
// Returns query result and column information.
func (q *Query) Query(db *sql.DB, ctx *domain.Context, values StringMap) (*QueryResult, error) {
	query, err := ctx.Transform(q.tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %s", err)
	}

	args, err := q.prepareArgs(values)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute query: %s\n\n%s", err, query)
	}
	defer rows.Close()

	result, err := q.scanRows(rows)
	if err != nil {
		return nil, fmt.Errorf("Failed to read row: %s\n\n%s", err, query)
	}

	return result, nil
}

func (q *Query) scanRows(rows *sql.Rows) (*QueryResult, error) {
	result := []map[string]string{}
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		scans := make([]interface{}, len(fields))
		for i := range scans {
			scans[i] = &scans[i]
		}
		if err := rows.Scan(scans...); err != nil {
			return nil, err
		}
		row := make(map[string]string)
		for i, v := range scans {
			value := ""
			if v != nil {
				value = fmt.Sprintf("%v", v)
			}
			row[fields[i]] = value
		}
		result = append(result, row)
	}

	return &QueryResult{
		Columns: fields,
		Rows:    result,
	}, nil
}

func (q *Query) prepareArgs(values StringMap) ([]interface{}, error) {
	args := make([]interface{}, 0)
	for _, param := range q.params {
		if value, ok := values[param.Name()]; ok {
			transformedValue, err := param.Transform(value)
			if err != nil {
				return nil, fmt.Errorf("failed to transform parameter '%s': %s", param.Name(), err)
			}
			args = append(args, transformedValue)
		} else if !param.IsOptional() {
			return nil, fmt.Errorf("parameter '%s' is missing", param.Name())
		}
	}
	return args, nil
}
