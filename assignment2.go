package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Anara"
	dbname   = "Kali"
)

func connectDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database.")
	return db
}

func createTask(db *sql.DB, name string) error {
	sqlStatement := `INSERT INTO tasks (name) VALUES ($1) ON CONFLICT (name) DO NOTHING;`
	_, err := db.Exec(sqlStatement, name)
	if err != nil {
		return err
	}
	fmt.Println("Task added successfully:", name)
	return nil
}

func readTasks(db *sql.DB) {
	rows, err := db.Query("SELECT id, name, completed FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var completed bool
		err = rows.Scan(&id, &name, &completed)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Task: %s, Completed: %v\n", id, name, completed)
	}
}

func updateTask(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE tasks SET completed = TRUE WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	fmt.Println("Task updated successfully:", id)
	return nil
}

func deleteTask(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	fmt.Println("Task deleted successfully:", id)
	return nil
}

func ensureTableExists(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		completed BOOLEAN NOT NULL DEFAULT FALSE
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal("Failed to create table:", err)
	}
	fmt.Println("Table ensured: tasks")
}

func main() {
	db := connectDB()
	defer db.Close()

	ensureTableExists(db)

	if err := createTask(db, "Read Go programming book"); err != nil {
		log.Println("Error creating task:", err)
	}

	readTasks(db)

	if err := updateTask(db, 1); err != nil {
		log.Println("Error updating task:", err)
	}

	if err := deleteTask(db, 2); err != nil {
		log.Println("Error deleting task:", err)
	}
}
