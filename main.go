package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB

type NewsFeed struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Article   string `json:"article"`
	Category  string `json:"category"`
	TimeStamp int    `json:"timeStamp"`
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

	//Routing - Gorilla-Mux + HandleFunc(S)
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		println("The server is working")
	})

	r.HandleFunc("/articles", getAllArticles).Methods("GET")
	r.HandleFunc("/article/{id}", getArticleHandler).Methods("GET")
	r.HandleFunc("/article", createArticleHandler).Methods("POST")
	r.HandleFunc("/news/update/:id", UpdateNewsPiece).Methods("PUT")
	r.HandleFunc("/news//delete/:id", DeleteNewsPiece).Methods("DELETE")

	serverErr := http.ListenAndServe(":8000", r)
	if serverErr != nil {
		log.Fatal(serverErr)
	}
}

// the logic of handlers for all queries
func getAllArticles(w http.ResponseWriter, r *http.Request) {
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
	row, _ := db.Query("SELECT * FROM articles WHERE id = ?", id)
	var news NewsFeed
	err := row.Scan(&news.Id, &news.Title, &news.Article, &news.Category)
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
	result, err := db.Exec("INSERT INTO articles (title, article, category, createdAt) VALUES (?,?,?,?)", newsArticle.Title, newsArticle.Article, newsArticle.Category, newsArticle.TimeStamp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	//newsArticle.Id, err = result.Id + 1
	//if err != nil {
	//	log.Fatal(err)
	//}
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
	_, err = db.Exec("UPDATE articles SET title = ?, article = ?, category = ?, timeStamp = ? WHERE id = ?", newsArticle.Title, newsArticle.Article, newsArticle.Category, newsArticle.TimeStamp, params["id"])
	if err != nil {
		log.Fatal(err)
	}
	newsArticle.Id, err = strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(newsArticle)
}

func DeleteNewsPiece(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	_, err := db.Exec("DELETE FROM articles WHERE id = ?", params["id"])
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode("News article deleted successfully")
}
