package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq" // Postgresドライバをインポート
)

func main() {
	// 1. 環境変数から接続情報を取得（docker-compose.ymlで設定した変数名）
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// 2. 接続文字列を作成
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname)

	// 3. DBへ接続
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the database!")

	// 接続確認用のAPI
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Connected to DB and Server is running!")
	})

	fmt.Println("Server starting at :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
