package main

import (
	"encoding/json" // json変換用
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgresドライバをインポート
)

// Task 構造体 (DBのテーブルとJSONの両方に対応)
type Task struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"` // string から time.Time へ変更
}

var db *sqlx.DB // グローバル変数として保持

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
	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the database!")

	// 接続確認用
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "API Server is running. Access /tasks to see data.")
	})

	// エンドポイントの設定
	http.HandleFunc("/tasks", getTasksHandler) // GET /tasks
	fmt.Println("Server starting at :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// タスク一覧を取得して返すハンドラー
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	// 【重要】CORSヘッダーを追加（すべてのオリジンを許可するか、特定のポートを指定）
   	 w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
   	 w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
   	 w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 安全策：DBが接続されていない場合はエラーを返す
	if db == nil {
		http.Error(w, "Database connection is not initialized", http.StatusInternalServerError)
		return
	}

	tasks := []Task{}
	err := db.Select(&tasks, "SELECT * FROM tasks ORDER BY created_at DESC")
	if err != nil {
		log.Printf("Database error: %v", err) // ログに詳細を出す
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
