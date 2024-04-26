package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

var (
	db  *sql.DB
	rdb *redis.Client
	ctx = context.Background()
)

func setupDB() *sql.DB {
	// This connection string should be changed according to your database settings.
	connStr := "user=postgres dbname=Kali password=Anara sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	ensureDatabaseSetup(db)

	return db
}

func ensureDatabaseSetup(db *sql.DB) {
	// Create the products table if it does not exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255),
        description TEXT,
        price NUMERIC
    );`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating products table: %v", err)
	}

	// Check if the table is empty and populate it
	var count int
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count == 0 {
		// Populate the table with some sample data
		insertSampleData(db)
	}
}

func insertSampleData(db *sql.DB) {
	sampleProducts := []Product{
		{Name: "Widget A", Description: "A useful widget", Price: 19.99},
		{Name: "Widget B", Description: "Another great widget", Price: 29.99},
	}
	for _, p := range sampleProducts {
		_, err := db.Exec("INSERT INTO products (name, description, price) VALUES ($1, $2, $3)", p.Name, p.Description, p.Price)
		if err != nil {
			log.Fatalf("Error inserting sample data: %v", err)
		}
	}
}

func setupRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // or the appropriate host and port
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	// Check if the connection is successful
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	return rdb
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	// Get ID from query parameter
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// First check Redis
	val, err := rdb.Get(ctx, id).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(val))
		return
	}

	// If not in Redis, query the database
	var product Product
	err = db.QueryRow("SELECT id, name, description, price FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Description, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serialize the product and store in Redis
	jsonData, err := json.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rdb.Set(ctx, id, jsonData, 10*time.Minute) // Set with TTL

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func checkRedis() {
	// First check Redis
	val, err := rdb.Get(ctx, idQuery).Result()
	if err == nil {
		log.Println("Retrieved data from Redis")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(val))
		return
	} else {
		log.Printf("Data not found in Redis, error: %v", err)
	}

	// If not in Redis, query the database
	var product Product
	err = db.QueryRow("SELECT id, name, description, price FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Description, &product.Price)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serialize the product and store in Redis
	jsonData, err := json.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = rdb.Set(ctx, idQuery, jsonData, 10*time.Minute).Result() // Set with TTL
	if err != nil {
		log.Printf("Failed to set data in Redis: %v", err)
	} else {
		log.Println("Data stored in Redis successfully")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func main() {
	db = setupDB()
	defer db.Close()

	rdb = setupRedis()
	defer rdb.Close()

	http.HandleFunc("/product", getProduct)
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
