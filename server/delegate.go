package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/themisir/databoard/domain"
	"github.com/themisir/databoard/query"
)

type queryDelegate struct {
	db    *sql.DB
	query *query.Query
	first bool
}

type mutationDelegate struct {
	db    *sql.DB
	query *query.Query
}

type queryResponseData struct {
	Columns []string            `json:"columns"`
	Data    []map[string]string `json:"data"`
}

// Create method delegate that handles requests by executing db queries
// and responding with returned data. If first set to true only first
// row will be responded or 404 status will be responded in case of empty
// result returned from database.
func Query(db *sql.DB, query *query.Query, first bool) MethodDelegate {
	return &queryDelegate{db: db, query: query, first: first}
}

// Create method delegate that handles requests by executing db queries
// and responding with 204 No Content.
func Mutation(db *sql.DB, query *query.Query) MethodDelegate {
	return &mutationDelegate{db: db, query: query}
}

func (d *queryDelegate) Handle(w http.ResponseWriter, ctx *domain.Context, values map[string]string) {
	result, err := d.query.Query(d.db, ctx, values)
	if err != nil {
		log.Printf("Failed to execute query: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if d.first {
		if len(result.Rows) > 0 {
			json.NewEncoder(w).Encode(result.Rows[0])
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		return
	}

	data := queryResponseData{
		Columns: result.Columns,
		Data:    result.Rows,
	}
	json.NewEncoder(w).Encode(data)
}

func (d *mutationDelegate) Handle(w http.ResponseWriter, ctx *domain.Context, values map[string]string) {
	_, err := d.query.Exec(d.db, ctx, values)
	if err != nil {
		log.Printf("Failed to execute query: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
