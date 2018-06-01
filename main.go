package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("starting server...")
	connStr := "postgres://rqdratsn:MNNiHU_roaqH-4igp47fQwZvq3FMLSBO@stampy.db.elephantsql.com:5432/rqdratsn"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open connection %s", err)
	}
	db := &database{conn: conn}

	http.Handle("/developers", withLogging(&handler{db: db}))
	http.HandleFunc("/healthz", healthz)
	http.ListenAndServe(":8080", nil)
}

func withLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("before")
		h.ServeHTTP(w, r)
		log.Println("after")
	})
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json;charset=utf-8")
	w.Write([]byte(`{"healthy": true}`)) // {"healthy" : true}
}

type database struct {
	conn *sql.DB
}

func (d *database) AllDevelopers() []developer {
	rows, err := d.conn.Query("SELECT id, name FROM developers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var devs []developer
	for rows.Next() {
		var dev developer
		if err := rows.Scan(&dev.ID, &dev.Name); err != nil {
			log.Fatal(err)
		}

		devs = append(devs, dev)

	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return devs
}

type handler struct {
	db *database
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json;charset=utf-8")
	ds := h.db.AllDevelopers()
	json.NewEncoder(w).Encode(developerList{Developers: ds})
}

type developer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type developerList struct {
	Developers []developer `json:"developers"`
}

// /developers

// {"developers": [{"id": 1, "name": Alice}, {"id: 2", "name": "Bob"}]}
