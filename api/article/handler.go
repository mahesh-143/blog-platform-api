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
	Id          int       `json:"article_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	Category_id int       `json:"category_id"`
	Category    string    `json:"name"`
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

func (h *Handler) GetArticlesByCategory(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Connection.Query(
		"SELECT a.* FROM articles a JOIN article_categories ac ON a.article_id = ac.article_id JOIN categories c ON ac.category_id = c.id WHERE c.name = $1;", r.PathValue("category"))
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
	rows, err := db.Connection.Query("SELECT * FROM articles WHERE article_id = $1 " + r.PathValue("id"))
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

	// begin transaction

	tx, err := db.Connection.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// insert article

	err = tx.QueryRow("INSERT INTO articles (title, content) VALUES ($1, $2) RETURNING article_id, created_at", newArticle.Title, newArticle.Content).Scan(&newArticle.Id, &newArticle.CreatedAt)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create article :"+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Inserted article")

	// insert category
	err = tx.QueryRow("INSERT INTO categories (name) VALUES ($1) RETURNING category_id", newArticle.Category).Scan(&newArticle.Category_id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create article :"+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Inserted category")

	// insert relation
	_, err = tx.Exec("INSERT INTO article_categories (article_id, category_id) VALUES ($1, $2) RETURNING category_id", newArticle.Category_id, newArticle.Id)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create article :"+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Inserted into article_categories")

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// send response
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

	articleId := r.URL.Query().Get("id")
	if articleId == "" {
		http.Error(w, "Article ID  is required", http.StatusBadRequest)
		return
	}

	// begin transaction
	tx, err := db.Connection.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// update Article
	err = tx.QueryRow("UPDATE articles SET title = $1, content = $2 WHERE article_id = $3 RETURNING article_id, created_at", updatedArticle.Title, updatedArticle.Content, r.PathValue("id")).Scan(&updatedArticle.Id, &updatedArticle.CreatedAt)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update the article : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update the article : "+err.Error(), http.StatusInternalServerError)
	}

	response := map[string]interface{}{"message": "Article updated!", "article": updatedArticle}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
