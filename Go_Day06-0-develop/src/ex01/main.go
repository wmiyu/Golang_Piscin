package main

// These are the libraries we are going to use
// Both "fmt" and "net" are part of the Go standard library
import (
	// "fmt" has methods for formatted I/O operations (like printing to the console)
	"database/sql"
	"fmt"
	// The "net/http" library has methods to implement HTTP clients and servers
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	auth "github.com/abbot/go-http-auth"
)


type dbStore struct {
	db *sql.DB
}

func Secret(user, realm string) string {
	if user == "admin" {
			// password is "hello"
			return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
	}
	return ""
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handlerHello).Methods("GET")
	authenticator := auth.NewBasicAuthenticator("example.com", Secret)
	// http.HandleFunc("/admin", authenticator.Wrap(handle))
	r.HandleFunc("/admin", authenticator.Wrap(adminHandle)).Methods("GET")

	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/assets").Handler(staticFileHandler).Methods("GET")

	r.HandleFunc("/hello", handlerHello).Methods("GET")
	r.HandleFunc("/doc", getDocHandler).Methods("GET")
	r.HandleFunc("/bird", getBirdHandlers).Methods("GET")
	r.HandleFunc("/bird", createBirdHandler).Methods("POST")

	return r
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "pgpwd4habr"
	dbname   = "bird_wiki"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")

	InitStore(&dbStore{db: db})
	r := newRouter()
	http.ListenAndServe(":8888", r)
}
