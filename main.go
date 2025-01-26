package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Author struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Book struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	AuthorID      int    `json:"author_id"`
	PublishedYear int    `json:"published_year"`
	Genre         string `json:"genre"`
}

type Member struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	JoinDate string `json:"join_date"`
}

func main() {
	// Database connection
	connStr := "postgresql://pgx_demo:tmp_pwd@localhost:5432/pgx_demo"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// Insert initial data
	_, err = pool.Exec(context.Background(), "INSERT INTO authors (name, email) VALUES ($1, $2) ON CONFLICT DO NOTHING", "J.K. Rowling", "jk.rowling@codeheim.io")
	if err != nil {
		log.Fatalf("Error inserting author: %v", err)
	}

	// Transaction
	tx, err := pool.Begin(context.Background())
	if err != nil {
		log.Fatalf("Error starting transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), "INSERT INTO authors (name, email) VALUES ($1, $2)", "George R.R. Martin", "george.martin@codeheim.io")
	if err != nil {
		log.Fatalf("Error inserting author: %v", err)
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO books (title, author_id, published_year, genre) VALUES ($1, $2, $3, $4)", "Harry Potter", 1, 1997, "Fantasy")
	if err != nil {
		log.Fatalf("Error inserting book: %v", err)
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO members (name, email) VALUES ($1, $2)", "John Doe", "john.doe@example.com")
	if err != nil {
		log.Fatalf("Error inserting member: %v", err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}

	// Query all authors
	rows, err := pool.Query(context.Background(), "SELECT id, name, email FROM authors")
	if err != nil {
		log.Fatalf("Error querying authors: %v", err)
	}
	defer rows.Close()

	var authors []Author
	for rows.Next() {
		var author Author
		if err := rows.Scan(&author.ID, &author.Name, &author.Email); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		authors = append(authors, author)
	}
	fmt.Println("Authors:", authors)

	// Query specific book by title
	var book Book
	err = pool.QueryRow(context.Background(), "SELECT id, title, author_id, published_year, genre FROM books WHERE title=$1", "Harry Potter").Scan(&book.ID, &book.Title, &book.AuthorID, &book.PublishedYear, &book.Genre)
	if err != nil {
		log.Fatalf("Error querying book: %v", err)
	}
	fmt.Println("Book Details:", book)

	// Update an author's name
	authorID := 1
	newName := "J.K. Rowling Updated"
	_, err = pool.Exec(context.Background(), "UPDATE authors SET name = $1 WHERE id = $2", newName, authorID)
	if err != nil {
		log.Fatalf("Error updating author: %v", err)
	}
	fmt.Println("Author updated successfully")

	// Delete a book by ID
	bookID := 1
	_, err = pool.Exec(context.Background(), "DELETE FROM books WHERE id = $1", bookID)
	if err != nil {
		log.Fatalf("Error deleting book: %v", err)
	}
	fmt.Println("Book deleted successfully")
}
