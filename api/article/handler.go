package article

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mahesh-143/blog-platform-api/db"
)

type Handler struct{}

type Article struct {
	Id        int       `json:"article_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Connection.Query("SELECT * FROM articles")
	if err != nil {
		http.Error(w, "Failed to get articles: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var articles []Article

	for rows.Next() {
		var article Article
		rows.Scan(
			&article.Id,
			&article.Title,
			&article.Content,
			&article.CreatedAt,
		)
		articles = append(articles, article)
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(articles); err != nil {
		http.Error(w, "Failed to encode articles to JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) FindByID(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Connection.Query("SELECT * FROM articles WHERE article_id =" + r.PathValue("id"))
	if err != nil {
		http.Error(w, "Failed to get articles: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var articles []Article

	for rows.Next() {
		var article Article
		rows.Scan(
			&article.Id,
			&article.Title,
			&article.Content,
			&article.CreatedAt,
		)
		articles = append(articles, article)
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(articles); err != nil {
		http.Error(w, "Failed to encode articles to JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var newArticle Article
	err := json.NewDecoder(r.Body).Decode(&newArticle)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Error decoding request body: ", err)
		return
	}
	defer r.Body.Close()
	err = db.Connection.QueryRow("INSERT INTO articles (title, content) VALUES ($1, $2) RETURNING article_id, created_at", newArticle.Title, newArticle.Content).Scan(&newArticle.Id, &newArticle.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to create article :"+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"message": "Article created!", "article": newArticle}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	_, err := db.Connection.Exec("DELETE FROM articles WHERE article_id = $1", r.PathValue("id"))
	if err != nil {
		http.Error(w, "Failed to delete article: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(w, "Article successfully", http.StatusAccepted)
}

func (h *Handler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	var updatedArticle Article
	err := json.NewDecoder(r.Body).Decode(&updatedArticle)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Error decoding request body: ", err)
		return
	}
	defer r.Body.Close()
	err = db.Connection.QueryRow("UPDATE articles SET title = $1, content = $2 WHERE article_id = $3 RETURNING article_id, created_at", updatedArticle.Title, updatedArticle.Content, r.PathValue("id")).Scan(&updatedArticle.Id, &updatedArticle.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to update the article : "+err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{"message": "Article updated!", "article": updatedArticle}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
