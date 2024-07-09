package article

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mahesh-143/blog-platform-api/db"
)

type Handler struct{}

type Articles struct {
	Id        string    `json:"article_id"`
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

	var articles []Articles

	for rows.Next() {
		var article Articles
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

	var articles []Articles

	for rows.Next() {
		var article Articles
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

func (h *Handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	_, err := db.Connection.Exec("DELETE FROM articles WHERE article_id = $1", r.PathValue("id"))
	if err != nil {
		http.Error(w, "Failed to delete article: "+err.Error(), http.StatusInternalServerError)
		return
	} else {
		http.Error(w, "Article successfully", http.StatusAccepted)
	}
}
