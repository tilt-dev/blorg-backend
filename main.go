package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/windmilleng/blorg-backend/golink"

	_ "github.com/lib/pq"
)

const USER = "blorger"

var dbAddr = flag.String("dbAddr", "localhost:26257", "address of the blorg database")
var db *sql.DB

func main() {
	flag.Parse()

	setupDatabase()

	r := mux.NewRouter()
	r.HandleFunc("/pong", Pong)
	r.HandleFunc("/golink/{name}", Golink)
	http.Handle("/", r)
	fmt.Println("Starting up on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupDatabase() {
	// Connect to the "bank" database.
	db2, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s@%s/golink?sslmode=disable", USER, *dbAddr))
	db = db2
	if err != nil {
		log.Fatal("error opening database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}
}

func Pong(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "pong")
}

func Golink(w http.ResponseWriter, req *http.Request) {
	gl := golink.NewGolink(db)
	vars := mux.Vars(req)
	name := vars["name"]

	switch req.Method {
	case "GET":
		link, err := gl.LinkFromName(name)
		if err != nil {
			log.Fatal("error getting link from name: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if link == "" {
			w.WriteHeader(http.StatusNotFound)
			// TODO(dmiller) make this JSON
			fmt.Fprint(w, "Link not found")
			return
		}
		j, _ := golink.LinkAsJSON(name, link)
		fmt.Fprint(w, j)

	case "PUT":
		payload, err := golink.ParseParams(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", err)
			return
		}
		err = gl.WriteLink(payload)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		j, _ := golink.LinkAsJSON(payload.Name, payload.Address)
		fmt.Fprint(w, j)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
