package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var db *sql.DB

type NewsFeed struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Article   string    `json:"article"`
	Category  string    `json:"category"`
	TimeStamp time.Time `json:"timeStamp"`
}

// the logic of handlers for all queries

func HomePage(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
func getAllArticles(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newsArticles []NewsFeed
	rows, err := db.Query("SELECT id, title, article, category, createdAt FROM articles")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var newsArticle NewsFeed
		err := rows.Scan(&newsArticle.Id, &newsArticle.Title, &newsArticle.Article, &newsArticle.Category, &newsArticle.TimeStamp)
		if err != nil {
			log.Fatal(err)
		}
		newsArticles = append(newsArticles, newsArticle)
	}
	json.NewEncoder(w).Encode(newsArticles)
}
func getArticleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	var newsArticle NewsFeed
	err := db.QueryRow("SELECT id, title, article, category, createdAt FROM articles WHERE id = ?", id).Scan(&newsArticle.Id, &newsArticle.Title, &newsArticle.Article, &newsArticle.Category, &newsArticle.TimeStamp)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(newsArticle)
}
func createArticleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newsArticle NewsFeed
	err := json.NewDecoder(r.Body).Decode(&newsArticle)
	if err != nil {
		log.Fatal(err)
	}
	//newsArticle.TimeStamp = time.Now()
	_, e := db.Exec("INSERT INTO articles (title, article, category) VALUES (?,?,?)", newsArticle.Title, newsArticle.Article, newsArticle.Category)
	if e != nil {
		log.Fatal(e)
	}
	json.NewEncoder(w).Encode(newsArticle)
}
func UpdateNewsPiece(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var newsArticle NewsFeed
	err := json.NewDecoder(r.Body).Decode(&newsArticle)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("UPDATE articles SET title = ?, article = ?, category = ? WHERE id = ?", newsArticle.Title, newsArticle.Article, newsArticle.Category, params["id"])
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(newsArticle)
}
func DeleteNewsPiece(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	_, err := db.Exec("DELETE FROM articles WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode("News article deleted successfully")
}

func handleRequests() {
	//Routing - Gorilla-Mux + HandleFunc(S)
	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/", HomePage)
	myRouter.HandleFunc("/articles", getAllArticles).Methods("GET")
	myRouter.HandleFunc("/article", createArticleHandler).Methods("POST")
	myRouter.HandleFunc("/article/{id}", getArticleHandler).Methods("GET")
	myRouter.HandleFunc("/article/update/{id}", UpdateNewsPiece).Methods("PUT")
	myRouter.HandleFunc("/article/delete/{id}", DeleteNewsPiece).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", myRouter))

}

func main() {

	//Capture connection properties.
	//cfg := mysql.Config{
	//	User:                 "root",
	//	Passwd:               "brickDBmaria9",
	//	Net:                  "tcp",
	//	Addr:                 "127.0.0.1:3306",
	//	DBName:               "news",
	//	AllowNativePasswords: true,
	//}

	// Get a database handle
	var err error
	db, err = sql.Open("mysql", "root:brickDBmaria9@tcp(127.0.0.1:3306)/news")
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	handleRequests()
}
