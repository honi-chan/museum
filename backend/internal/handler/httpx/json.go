package httpx

import (
	"encoding/json"
	"net/http"
)

// DecodeJSON はHTTP Request BodyのJSONを
// 指定された構造体へ変換する。
//
// 各Handlerで毎回、
//
//	json.NewDecoder(r.Body).Decode(...)
//
// と書かなくて済むように共通化する。
//
// Handler側は
// 「Requestを読み込む」という意図だけを
// 表現できるようになる。
func DecodeJSON(
	r *http.Request,
	dst any,
) error {
	return json.NewDecoder(
		r.Body,
	).Decode(dst)
}

// WriteJSON はJSONレスポンスを返す共通処理。
//
// 以下をまとめて担当する。
//
// 1. Content-Type設定
// 2. HTTP Status設定
// 3. JSON Encode
//
// Handlerごとの重複コードを減らす。
func WriteJSON(
	w http.ResponseWriter,
	status int,
	data any,
) error {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	return json.NewEncoder(
		w,
	).Encode(data)
}

// ErrorResponse はAPI共通のエラー形式。
//
// 全APIでエラー形式を統一することで、
// Frontend側も扱いやすくなる。
//
// 例:
//
//	{
//	  "error": "exhibition title is required"
//	}
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteError はJSON形式のエラーを返す。
//
// http.Error()を使うとレスポンス形式が
// 他のJSON APIと揃わなくなるため、
// MUSEUMではJSONに統一する。
func WriteError(
	w http.ResponseWriter,
	status int,
	err error,
) error {
	return WriteJSON(
		w,
		status,
		ErrorResponse{
			Error: err.Error(),
		},
	)
}
