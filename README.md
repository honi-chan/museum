# MUSEUM

> 自分という人間を、展示して残す。

MUSEUMは、**自分だけの展示会をつくるSNS**です。

一般的なSNSのように「今日何をしたか」をタイムラインへ投稿するのではなく、自分のプロフィールそのものをひとつのMuseumとしてつくっていきます。

```text id="p2imk0"
User
  ↓
Museum
  ↓
Exhibition
  ↓
Exhibit
```

写真、作品、思い出、好きなもの、考えていたこと。

残しておきたいものを自分で選び、自分で並べ、自分だけのMuseumを少しずつつくっていきます。

---

## Concept

Instagramが、

> 今日の自分を見せる場所

だとしたら、MUSEUMは、

> 自分という人間を編集して残す場所

です。

MUSEUMでは、投稿の新しさや人気ではなく、**その人が「何を残したいと思ったか」**を大切にします。

そのため、一般的なSNSにある以下の要素をプロダクトの中心には置きません。

* タイムライン
* いいね数
* フォロワー数
* バズ
* レコメンドアルゴリズム

MUSEUMの主役はコンテンツではなく、**その人自身が編集したMuseum**です。

---

## Product Structure

MUSEUMは4つの基本概念から構成されます。

```text id="lpsv22"
User

 └── Museum
      │
      ├── Exhibition
      │    ├── Exhibit
      │    ├── Exhibit
      │    └── Exhibit
      │
      └── Exhibition
           ├── Exhibit
           └── Exhibit
```

### Museum

ユーザー自身を表現する空間です。

プロフィールページというより、**その人自身の美術館**として扱います。

### Exhibition

Museumの中につくられる「展示室」です。

例えば、

```text id="3tvxkq"
つくったもの
京都の夏
好きなもの
昔好きだったもの
持っているもの
仕事
人
場所
考えていたこと
```

など、ユーザー自身がテーマを決めます。

### Exhibit

Exhibitionの中に置かれる展示物です。

写真や作品など、ユーザーが残したいものを展示します。

---

## Technology

```text id="h6f04r"
Frontend
Next.js
   │
   │ HTTP / JSON
   ▼
Backend
Go
   │
   │ SQL
   ▼
Database
PostgreSQL 18
```

ローカル開発環境にはDocker Composeを使用します。

API仕様はOpenAPIをSource of Truthとして管理しています。

```text id="t5yk9u"
backend/api/openapi.yaml
        │
        ├── oapi-codegen
        │      ↓
        │     Go
        │
        └── openapi-typescript
               ↓
           TypeScript
```

BackendとFrontendでAPI型を別々に管理しません。

---

## Architecture

Backendでは以下の構成を採用しています。

```text id="uv45kj"
HTTP
 │
 ▼
Handler
 │
 ▼
UseCase
 │
 ▼
Repository
 │
 ▼
Infrastructure
 │
 ▼
PostgreSQL
```

MUSEUMのDomainと、HTTP・PostgreSQLなどの技術詳細を分離します。

詳しくは [`ARCHITECTURE.md`](./ARCHITECTURE.md) を参照してください。

---

## Design

MUSEUMのUIは一般的なSNS Dashboardではなく、**静かな現代美術館**をイメージしています。

基本原則は、

* 展示物を主役にする
* 大きな余白を使う
* 情報を詰め込みすぎない
* 過剰なカードUIを避ける
* 過剰な角丸やShadowを避ける
* 色より余白・Typography・配置で階層をつくる
* 非対称なレイアウトも許容する
* 美術館のキャプションや図録のようなTypographyを使う

です。

詳細は [`DESIGN.md`](./DESIGN.md) を参照してください。

---

## Directory

```text id="8tddbe"
MUSEUM/
│
├── README.md
├── PRODUCT.md
├── DESIGN.md
├── ARCHITECTURE.md
│
├── backend/
│   ├── api/
│   ├── db/
│   ├── cmd/
│   └── internal/
│
├── frontend/
│   ├── app/
│   ├── components/
│   └── lib/
│
└── compose.yaml
```

それぞれのドキュメントの役割は次の通りです。

| Document        | 内容             |
| --------------- | -------------- |
| README.md       | プロジェクトの入口      |
| PRODUCT.md      | MUSEUMとは何か     |
| DESIGN.md       | MUSEUMをどう見せるか  |
| ARCHITECTURE.md | MUSEUMをどう実装するか |

---

## Setup

### Database

プロジェクトルートで実行します。

```bash id="4sjgfk"
make db-up
```

Migrationを実行します。

```bash id="u9l3mq"
make db-migrate-up
```

---

### Generate

OpenAPI / sqlcなどのコード生成を実行します。

```bash id="22xvru"
make generate
```

生成されたコードは直接編集しません。

変更する場合はSource of Truthを変更してから再生成します。

---

### Backend

```bash id="ekspvb"
cd backend
```

Database URLを設定します。

```bash id="11nqmd"
export DATABASE_URL="postgres://museum:museum@localhost:15432/museum?sslmode=disable"
```

起動します。

```bash id="1tr4ns"
go run ./cmd/server
```

Backend API:

```text id="jnc43e"
http://localhost:8080
```

---

### Frontend

```bash id="gt8grh"
cd frontend
```

API型を生成します。

```bash id="t2frzj"
npm run generate:api
```

起動します。

```bash id="j0bixv"
npm run dev
```

Frontend:

```text id="9s21tf"
http://localhost:3000
```

---

## Test

Backend:

```bash id="94iym7"
cd backend

go test ./...
```

Frontend:

```bash id="qlpzpl"
cd frontend

npm run lint
```

---

## Development Rules

### generatedコードを直接編集しない

以下は自動生成されます。

```text id="n1cjcw"
backend/api/generated/
backend/internal/infrastructure/**/sqlcgen/
frontend/lib/api/generated/
```

APIを変更する場合：

```text id="mymz0u"
backend/api/openapi.yaml
```

SQLを変更する場合：

```text id="sjh7gd"
backend/db/query/
```

を変更します。

---

### 重複する作業は自動化する

同じ変更を複数箇所で繰り返すようになった場合は、

```text id="6hzwnk"
共通化
OpenAPI
Code Generation
sqlc
Middleware
共通Component
```

などによって変更箇所を減らせないか検討します。

一方で、将来使うかもしれないという理由だけで抽象化は行いません。

必要になったタイミングでリファクタリングします。

---

## Status

現在はMVP開発中です。

最初に検証する価値は、

> 自分のMuseumをつくり、展示室をつくり、その中に自分の大切なものを並べる体験に価値があるか

です。

```text id="l94d8g"
Museum
   ↓
Exhibition
   ↓
Exhibit
   ↓
Visit
```

まずこの体験を完成させます。
