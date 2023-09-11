package main

//func CreateNewsPiece(w http.ResponseWriter, r *http.Request) {
//	return
//}
//
//func UpdateNewsPiece(w http.ResponseWriter, r *http.Request) {
//	return
//}

//func GetNewsPiece(w http.ResponseWriter, r *http.Request) {
//	id := r.URL.Query().Get("id")
//	row := db.QueryRow("SELECT * FROM news WHERE id = ?", id)
//
//	var news NewsFeed
//	err := row.Scan(&news.Id, &news.Title, &news.Article, &news.Category)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	json.NewEncoder(w).Encode(news)
//}

//func DeleteNewsPiece(w http.ResponseWriter, r *http.Request) {
//	return
//}
