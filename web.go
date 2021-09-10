package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var logger *log.Logger

func main() {
	var err error
	// Create database connection
	db, err = sql.Open("sqlite3", "./keys.db")
	handleError(err)

	args := os.Args[1:]
	// Show help if no arguments were given
	if len(args) == 0 {
		fmt.Println("Usage:")
		fmt.Println("  run start\t\tStarts Fiber server on port 8080")
		fmt.Println("  run add <user>\tAdds a new user and generates key for them")
		fmt.Println("  run create\t\tCreates new database file")
		fmt.Println("  run delete <user>\tDeletes a user")
		return
	}

	// Start Fiber server
	if args[0] == "start" {
		// Creates a logger for RTMP access
		logFile, err := os.OpenFile("auth.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		handleError(err)
		defer logFile.Close()
		logger = log.New(logFile, "", log.LstdFlags)

		// Serve index.html and static files
		app := fiber.New(fiber.Config{
			Views: html.New("./static", ".html"),
		})
		app.Static("/", "./static/files")
		app.Get("/", index)

		// Authorization and disconnect routes
		app.Get("/auth", auth)
		app.Get("/disconnect", disconnect)

		log.Fatal(app.Listen(":8080"))
	}

	// Add a new user
	if args[0] == "add" {
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
		generatedKey := stringWithCharset(10, charset, seededRand)
		prep, err := db.Prepare("INSERT INTO keys VALUES(?, ?)")
		handleError(err)
		_, err = prep.Exec(args[1], generatedKey)
		handleError(err)
		fmt.Println("Key for " + args[1] + " = " + generatedKey)
	}

	// Delete a user
	if args[0] == "delete" {
		prep, err := db.Prepare("DELETE FROM keys WHERE name=?")
		handleError(err)
		_, err = prep.Exec(args[1])
		handleError(err)
	}

	// Initializes database
	if args[0] == "create" {
		prep, err := db.Prepare("CREATE TABLE keys (name VARCHAR(64) PRIMARY KEY, key VARCHAR(128))")
		handleError(err)
		_, err = prep.Exec()
		handleError(err)
	}
}

// Serve info web
func index(c *fiber.Ctx) error {
	f, err := os.Open("results.csv")
	results := [][]string{}
	handleError(err)
	reader := csv.NewReader(f)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		handleError(err)
		results = append(results, record)
	}
	return c.Render("index", fiber.Map{"results": results})
}

// Authorization API
func auth(c *fiber.Ctx) error {
	// Get user and key from RTMP server
	name := c.Query("user")
	enteredKey := c.Query("key")
	var key string
	row := db.QueryRow("SELECT key FROM keys WHERE name=?", name)
	err := row.Scan(&key)
	if err == sql.ErrNoRows {
		logger.Println("Unknown user " + name + " attempted to connect")
		return c.SendStatus(fiber.StatusNotFound)
	} else {
		handleError(err)
		if enteredKey == key {
			logger.Println(name + " successfully connected")
			return c.SendStatus(fiber.StatusCreated)
		} else {
			logger.Println(name + " attempted to connect with wrong key (" + enteredKey + ")")
			return c.SendStatus(fiber.StatusNotFound)
		}
	}
}

// Log disconnects
func disconnect(c *fiber.Ctx) error {
	name := c.Query("user")
	logger.Println(name + " disconnected")
	return c.SendStatus(fiber.StatusOK)
}

// Handle unexpected errors
func handleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// Generate random authentication string
func stringWithCharset(length int, charset string, seededRand *rand.Rand) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
