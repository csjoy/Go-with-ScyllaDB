package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"net/http"
)

type status map[string]interface{}

type User struct {
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Password     string    `db:"password" json:"password"`
	Email        string    `db:"email" json:"email"`
	Phone        string    `db:"phone" json:"phone"`
	NewToken     string    `db:"new_token" json:"new_token"`
	RefreshToken string    `db:"refresh_token" json:"refresh_token"`
	UserType     string    `db:"user_type" json:"user_type"`
	UserID       string    `db:"user_id" json:"user_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

const createTable = `
CREATE TABLE IF NOT EXISTS users (
	first_name text,
	last_name text,
	password text,
	email text,
	phone text,
	new_token text,
	refresh_token text,
	user_type text,
	user_id text,
	created_at timestamp,
	updated_at timestamp,
	PRIMARY KEY((last_name), first_name, user_id, email, password)
)`

func main() {
	cluster := gocql.NewCluster("nodeX", "nodeY", "nodeZ")
	cluster.Keyspace = "chijwt"

	session, err := gocqlx.WrapSession(cluster.CreateSession())

	if err != nil {
		log.Fatal("unable to connect to scylla", err)
	}
	log.Println("Successfully connected")

	err = gocqlx.Session.Query(session, createTable, []string{}).ExecRelease()
	if err != nil {
		log.Fatal("failed to create table", err)
	}

	currentTime, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	newUser := User{
		FirstName:    "Prosenjit",
		LastName:     "Joy",
		Password:     "123456",
		Email:        "antohin.joy@gmail.com",
		Phone:        "01764277629",
		NewToken:     "will be generated later",
		RefreshToken: "will be generated later",
		UserType:     "ADMIN",
		UserID:       "1",
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}

	stmt, names := qb.Insert("users").Columns("first_name", "last_name", "password", "email", "phone", "user_type", "user_id").ToCql()

	err = gocqlx.Session.Query(session, stmt, names).BindStruct(newUser).ExecRelease()
	if err != nil {
		log.Fatal("failed to insert data:", err)
	}

	fmt.Println("Successfully inserted data")

	stmt, names = qb.Select("users").ToCql()
	var items []User
	err = gocqlx.Session.Query(session, stmt, names).Select(&items)
	if err != nil {
		log.Fatal("failed to read data:", err)
	}

	fmt.Println(items)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status{"success": "access granted for v1"})
	})

	if err := http.ListenAndServe(":5000", router); err != nil {
		log.Fatal("Error initializing server", err)
	}
}
