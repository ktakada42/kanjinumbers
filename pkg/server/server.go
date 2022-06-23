package server

import (
	"log"
	"net/http"

	"kanjinumbers/pkg/server/handler"
)

// Serve HTTPサーバを起動する
func Serve(addr string) {
	/* ===== URLマッピングを行う ===== */
	mappingURL()

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}

// get GETリクエストを処理する
func get(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodGet)
}

func mappingURL() {
	http.HandleFunc("/v1/number2kanji/", get(handler.HandleNumberToKanji))
	http.HandleFunc("/v1/kanji2number/", get(handler.HandleKanjiToNumber))
}

// httpMethod 指定したHTTPメソッドでAPIの処理を実行する
func httpMethod(apiFunc http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// CORS対応
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

		// プリフライトリクエストは処理を通さない
		if r.Method == http.MethodOptions {
			return
		}
		// 指定のHTTPメソッドでない場合はエラー
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method Not Allowed"))
			return
		}

		// 共通のレスポンスヘッダを設定
		w.Header().Add("Content-Type", "application/json")
		apiFunc(w, r)
	}
}
