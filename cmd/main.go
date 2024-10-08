package main

import (
	"log"
	"net/http"

	"github.com/mahesh-143/blog-platform-api/api/article"
	"github.com/mahesh-143/blog-platform-api/db"
)

func init() {
	db.InitDB()
}

func main() {

	article_handler := &article.Handler{}
	router := http.NewServeMux()
	router.HandleFunc("POST /api/articles/create", article_handler.CreateArticle)
	router.HandleFunc("GET /api/articles", article_handler.GetAll)
	router.HandleFunc("GET /api/articles/{category}", article_handler.GetArticlesByCategory)
	router.HandleFunc("GET /api/articles/{id}", article_handler.FindByID)
	router.HandleFunc("DELETE /api/articles/{id}", article_handler.DeleteByID)
	router.HandleFunc("PUT /api/articles/update/{id}", article_handler.UpdateArticle)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Printf("Starting server on port %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
