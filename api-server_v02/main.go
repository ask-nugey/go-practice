package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// User はユーザー情報を表す構造体。
type User struct {
	ID   int    `json:"id"`   // JSONエンコード時のフィールド名を指定
	Name string `json:"name"` // JSONエンコード時のフィールド名を指定
}

// グローバル変数の宣言
var (
	users  = []User{} // 保存しているユーザー情報のスライス
	nextID = 1        // 次に追加されるユーザーに割り当てるID
	mu     sync.Mutex // usersの排他制御のためのmutex
)

// addUser は新しいユーザーを追加するエンドポイントのハンドラ
func addUser(w http.ResponseWriter, r *http.Request) {
	// POSTメソッドのみ受け付ける
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusBadRequest)
		return
	}

	// リクエストボディからUserをデコード
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ユーザー情報にIDを割り当てて保存
	mu.Lock()                // 排他制御の開始
	u.ID = nextID            // 新しいIDを割り当て
	nextID++                 // 次のIDのインクリメント
	users = append(users, u) // ユーザーの追加
	mu.Unlock()              // 排他制御の終了

	// 追加されたユーザー情報をレスポンスとして返す
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

// getUser は指定されたIDのユーザー情報を取得するエンドポイントのハンドラ
func getUser(w http.ResponseWriter, r *http.Request) {
	// GETメソッドのみ受け付ける
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusBadRequest)
		return
	}

	// クエリパラメータからIDを取得
	id := r.URL.Query().Get("id")
	for _, u := range users {
		if fmt.Sprint(u.ID) == id {
			json.NewEncoder(w).Encode(u)
			return
		}
	}

	// 一致するユーザーが見つからなかった場合のエラーレスポンス
	http.Error(w, "User not found", http.StatusNotFound)
}

// getAllUsers は全てのユーザー情報を取得するエンドポイントのハンドラ
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	// GETメソッドのみ受け付ける
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusBadRequest)
		return
	}

	// 保存している全てのユーザー情報をレスポンスとして返す
	json.NewEncoder(w).Encode(users)
}

func main() {
	// 各エンドポイントとハンドラ関数の関連付け
	http.HandleFunc("/add-user", addUser)
	http.HandleFunc("/get-user", getUser)
	http.HandleFunc("/get-all-users", getAllUsers)

	// HTTPサーバの起動
	http.ListenAndServe(":8080", nil)
}
