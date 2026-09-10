package api

//go:generate go tool oapi-codegen -config oapi-codegen.yaml openapi.yaml

// OpenAPIからGoコードを自動生成する。
//
// oapi-codegenを直接実行せず、
//
//   go tool oapi-codegen
//
// を使用する。
//
// oapi-codegenのバージョンはgo.modで管理されるため、
//
// - 開発者ごとのPATH設定
// - グローバルインストール
// - CIへの個別インストール
//
// が不要になる。
//
// API変更時は:
//
// 1. openapi.yamlを変更
// 2. go generate ./...
// 3. 生成されたinterfaceに合わせてHandlerを実装
//
// の順番で進める。
//
// api/generated配下は自動生成コードなので
// 手動で編集しない。
