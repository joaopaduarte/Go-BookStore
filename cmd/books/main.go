package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"books/internal/cli"
	"books/internal/service"
	"books/internal/web"

	_ "modernc.org/sqlite"
)

func main() {
	// Conexão com o banco de dados SQLite
	db, err := sql.Open("sqlite", "./books.db")
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Criação da tabela se não existir
	createTableQuery := `
			CREATE TABLE IF NOT EXISTS books (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				title TEXT NOT NULL,
				author TEXT NOT NULL,
				genre TEXT NOT NULL
			);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Erro ao criar a tabela: %v", err)
	}

	// Inicializando o serviço
	bookService := service.NewBookService(db)

	// Inicializando os handlers
	bookHandlers := web.NewBookHandlers(bookService)

	// Verifica se o CLI foi chamado com o comando "search" ou "simulate"
	if len(os.Args) > 1 && (os.Args[1] == "search" || os.Args[1] == "simulate") {
		bookCLI := cli.NewBookCLI(bookService)
		bookCLI.Run()
		return
	}

	// Criando o roteador com o novo servidor
	router := http.NewServeMux()

	// Configurando as rotas RESTful
	router.HandleFunc("GET /books", bookHandlers.GetBooks)
	router.HandleFunc("POST /books", bookHandlers.CreateBook)
	router.HandleFunc("GET /books/{id}", bookHandlers.GetBookByID)
	router.HandleFunc("PUT /books/{id}", bookHandlers.UpdateBook)
	router.HandleFunc("DELETE /books/{id}", bookHandlers.DeleteBook)

	// Iniciando o servidor
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))

}
