package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Task 構造体 (DBとJSONの両方に対応)
type Task struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

var db *sqlx.DB

func main() {
	// 1. DB接続設定
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
	os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the database!")

	// 2. エンドポイントの設定
	// すべての /tasks へのリクエストをこの関数で受け取り、中でGETとPOSTを振り分けます
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		// CORS設定 (ブラウザからのアクセスを許可)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// プリフライトリクエスト(OPTIONS)への対応
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case http.MethodPost:
			createTaskHandler(w, r)
		case http.MethodDelete: // 追加：削除リクエストの振り分け
			deleteTaskHandler(w, r)
		case http.MethodPut:
			updateTaskHandler(w, r)
		default:
			getTasksHandler(w, r)
		}
	})

	// 接続確認用のルート
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "API Server is running.")
	})

	fmt.Println("Server starting at :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// --- 以下、ハンドラー関数 (main関数の外に定義します) ---

// タスク取得 (GET)
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := []Task{}
	err := db.Select(&tasks, "SELECT * FROM tasks ORDER BY created_at DESC")
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// タスク作成 (POST)
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id, created_at"
	err := db.QueryRow(query, t.Title, t.Description).Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}
// タスク削除 (DELETE)
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// URLのクエリパラメータから "id" を取得 (例: /tasks?id=5)
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 成功（中身は空）を返す
}

// 追加：タスク更新 (PUT)
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }

    // ステータスを 'DONE' に更新
    _, err := db.Exec("UPDATE tasks SET status = 'DONE' WHERE id = $1", id)
    if err != nil {
        log.Printf("Database error: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}
