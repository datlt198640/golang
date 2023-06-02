package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func main() {
	// Data Source Name Properties
	dsn := mysql.Config{
		User:                 "root",
		Passwd:               "Ltd20198640troyita@",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "sakila",
		AllowNativePasswords: true,
	}
	// Get a database handle
	var err error
	db, err = sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		log.Fatal(err)
		// return
	}
	defer db.Close()

	// check connection is active or not
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	// Performs CRUD
	actorID, err := addActor("JOE", "BERRY")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added actor: %v\n", actorID)

	actors, err := GetActor(201)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Actor found: %v\n", actors)

	rowsUpdated, err := updateActor("JAMES", 201)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total actor affected: %d\n", rowsUpdated)

	rowsDeleted, err := deleteActor(201)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total actor deleted: %d\n", rowsDeleted)

	// Perform transaction
	ctx := context.Background()
	actorID, err = txActor(ctx, "JOEN", "BERRY")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Actor ID: %v\n", actorID)

	// Perform procedure
	actors, err = GetActorSP("Joe", "Simple")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Actor found: %v\n", actors)
}

// PROCEDURE STORED
func GetActorSP(first_name string, last_name string) ([]Actor, error) {
	var actors []Actor

	result, err := db.Query("CALL addActor (?, ?)", first_name, last_name)
	if err != nil {
		return nil, fmt.Errorf("GetActorSP: %v", err)
	}
	defer result.Close()

	// Loop through rows
	for result.Next() {
		var acts Actor
		if err := result.Scan(&acts.actor_id, &acts.first_name, &acts.last_name); err != nil {
			return nil, fmt.Errorf("GetActorSP: %v", err)
		}
		actors = append(actors, acts)

		if err := result.Err(); err != nil {
			return nil, fmt.Errorf("GetActorSP: %v", err)
		}
	}
	return actors, nil
}

// TRANSACTION
func txActor(ctx context.Context, firstname string, lastname string) (int64, error) {

	// Begin Transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("Adding Actor Failed: %v", err)
	}
	// Defer a rollback in case of failure
	defer tx.Rollback()

	// Check if name exists
	var actID int64
	if err = tx.QueryRowContext(ctx, "SELECT actor_id from actor where first_name = ? and last_name = ?",
		firstname, lastname).Scan(&actID); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Actor does not exist")
		} else {
			return 0, fmt.Errorf("txActor: %v", err)
		}
	}
	// Rollback if actor exists
	if actID > 0 {
		if err = tx.Rollback(); err != nil {
			return 0, fmt.Errorf("txActor: %v", err)
		}
		fmt.Println("Actor already exist: ", actID)
		fmt.Println("*** Transaction Rolling Back ***")
		return actID, nil
	}

	// Create a new row
	result, err := tx.ExecContext(ctx, "INSERT INTO actor (first_name, last_name) VALUES (?, ?)",
		firstname, lastname)
	if err != nil {
		return 0, fmt.Errorf("txActor: %v", err)
	}

	// Get the ID of the actor just inserted
	NewActorID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("txActor: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("txActor: %v", err)
	} else {
		fmt.Println("New Actor Created: ", NewActorID)
		fmt.Println("*** Transaction Commited ***")
	}

	return NewActorID, nil
}

// CRUD
type Actor struct {
	actor_id   int64
	first_name string
	last_name  string
}

func addActor(firstname, lastname string) (int64, error) {
	result, err := db.Exec("INSERT INTO actor(first_name, last_name) VALUES(?, ?)", firstname, lastname)
	if err != nil {
		return 0, fmt.Errorf("addActor: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addActor: %v", err)
	}
	return id, nil
}

func GetActor(actorID int32) ([]Actor, error) {
	var actors []Actor

	result, err := db.Query("SELECT actor_id, first_name, last_name FROM actor WHERE actor_id = ?",
		actorID)
	if err != nil {
		return nil, fmt.Errorf("GetActor %v: %v", actorID, err)
	}
	defer result.Close()

	// Loop through rows
	for result.Next() {
		var acts Actor
		if err := result.Scan(&acts.actor_id, &acts.first_name, &acts.last_name); err != nil {
			return nil, fmt.Errorf("GetActor %v: %v", actorID, err)
		}
		actors = append(actors, acts)

		if err := result.Err(); err != nil {
			return nil, fmt.Errorf("GetActor %v: %v", actorID, err)
		}
	}
	return actors, nil
}

func updateActor(firstname string, actorid int32) (int64, error) {
	result, err := db.Exec("UPDATE actor SET first_name = ? WHERE actor_id = ?",
		firstname, actorid)
	if err != nil {
		return 0, fmt.Errorf("updateActor: %v", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("updateActor: %v", err)
	}
	return id, nil
}

func deleteActor(actorid int32) (int64, error) {
	result, err := db.Exec("DELETE FROM actor WHERE actor_id = ?",
		actorid)
	if err != nil {
		return 0, fmt.Errorf("deleteActor: %v", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("deleteActor: %v", err)
	}
	return id, nil
}
