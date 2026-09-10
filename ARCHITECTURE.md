# MUSEUM Architecture

## 1. このドキュメントについて

このドキュメントでは、MUSEUMのソフトウェア設計について説明します。

対象読者は以下です。

* MUSEUMへ新しく参加した開発者
* Go / Next.js / PostgreSQLを初めて触る人
* Clean Architectureに慣れていない人
* 「どこに何を書けばいいのか」を知りたい人

単にディレクトリ構成を説明するだけではなく、

> なぜこの構成なのか

> 新しい機能を追加するとき、どこを変更すればいいのか

まで理解できることを目的としています。

---

# 2. MUSEUMとは

MUSEUMは、

> 自分だけの展示会をつくるSNS

です。

一般的なSNSのようなタイムラインを中心とするのではなく、ユーザー自身のプロフィールを「Museum」として表現します。

基本構造は次の通りです。

```text
User
 │
 ▼
Museum
 │
 ▼
Exhibition
 │
 ▼
Exhibit
```

例えば、

```text
Haruki's Museum

├── THINGS I MADE
│   ├── CinemaAI
│   ├── Illustration
│   └── Application
│
├── SUMMER IN KYOTO
│   ├── Photo
│   ├── Photo
│   └── Video
│
└── THINGS I LOVE
    ├── Movie
    ├── Book
    └── Music
```

のように、自分自身を展示していきます。

---

# 3. 技術スタック

MUSEUMでは以下を使用します。

```text
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
PostgreSQL
```

ローカル環境ではDocker Composeを利用します。

```text
Browser
   │
   ▼
Next.js
   │
   ▼
Go API
   │
   ▼
PostgreSQL 18
```

---

# 4. Backendの設計方針

Backendでは、役割ごとにコードを分離します。

大きく5つに分かれています。

```text
handler
   │
   ▼
usecase
   │
   ▼
domain
   ▲
   │
repository
   ▲
   │
infrastructure
```

それぞれの責務は次の通りです。

| Layer          | 役割                    |
| -------------- | --------------------- |
| handler        | HTTPを受け取る             |
| usecase        | アプリケーションの処理を実行する      |
| domain         | MUSEUMの概念を表現する        |
| repository     | データ操作のルールを定義する        |
| infrastructure | PostgreSQLなど具体的な技術を扱う |

重要なのは、

> MUSEUMのルールと、PostgreSQLやHTTPなどの技術を分離する

ことです。

---

# 5. Domain

例：

```text
internal/domain/exhibition.go
```

```go
type Exhibition struct {
    ID          string
    MuseumID    string
    Title       string
    Description string
    CreatedAt   time.Time
}
```

DomainはMUSEUMそのものの概念を表します。

そのため、

```go
json:"title"
```

や、

```go
db:"title"
```

などは基本的に持たせません。

Domainは、

```text
HTTP
PostgreSQL
JSON
OpenAPI
pgx
sqlc
```

などを知らなくてよい状態を目指します。

---

# 6. UseCase

UseCaseは、

> ユーザーがMUSEUMで何をするか

を表します。

例えば、

```text
CreateExhibition
GetExhibition
CreateExhibit
MoveExhibit
PublishExhibition
```

などです。

例：

```text
internal/usecase/create_exhibition.go
```

```go
func (u *CreateExhibitionUseCase) Execute(
    ctx context.Context,
    input CreateExhibitionInput,
) (domain.Exhibition, error) {

    // タイトルを整形する。
    title := strings.TrimSpace(input.Title)

    // MUSEUMのルールをチェックする。
    if title == "" {
        return domain.Exhibition{},
            errors.New("exhibition title is required")
    }

    exhibition := domain.Exhibition{
        ID:          u.idGenerator.Generate(),
        MuseumID:    input.MuseumID,
        Title:       title,
        Description: input.Description,
        CreatedAt:   time.Now(),
    }

    // Repositoryを通して保存する。
    if err := u.repository.Save(
        ctx,
        exhibition,
    ); err != nil {
        return domain.Exhibition{}, err
    }

    return exhibition, nil
}
```

UseCaseは、

```text
HTTP Request
JSON
PostgreSQL
SQL
```

を直接扱いません。

---

# 7. Repository

Repositoryは、

> Domainを保存・取得するために必要な操作

を定義します。

例：

```text
internal/repository/exhibition_repository.go
```

```go
type ExhibitionRepository interface {
    Save(
        ctx context.Context,
        exhibition domain.Exhibition,
    ) error

    FindByID(
        ctx context.Context,
        id string,
    ) (domain.Exhibition, error)
}
```

ここでは、

```text
PostgreSQLを使う
```

とは書いていません。

そのためUseCaseから見ると、

```text
UseCase
   │
   ▼
ExhibitionRepository
```

だけが見えます。

実際に、

```text
PostgreSQL
Memory
Fake Repository
```

のどれが使われるかはUseCaseには関係ありません。

これによりテストもしやすくなります。

---

# 8. Infrastructure

Infrastructureでは具体的な技術を扱います。

例えば、

```text
PostgreSQL
UUID
外部API
Object Storage
```

などです。

Exhibitionの場合、

```text
internal/infrastructure/exhibition/
```

に、

```text
postgres_repository.go
uuid_generator.go
sqlcgen/
```

があります。

例えば、

```go
func (r *PostgresRepository) FindByID(
    ctx context.Context,
    id string,
) (domain.Exhibition, error) {

    row, err :=
        r.queries.GetExhibitionByID(
            ctx,
            id,
        )

    if err != nil {
        return domain.Exhibition{}, err
    }

    return domain.Exhibition{
        ID:          row.ID,
        MuseumID:    row.MuseumID,
        Title:       row.Title,
        Description: row.Description,
        CreatedAt:   row.CreatedAt,
    }, nil
}
```

ここで初めて、

```text
PostgreSQL
sqlc
pgx
```

が登場します。

---

# 9. Handler

HandlerはHTTPとUseCaseをつなぎます。

```text
HTTP Request

     ↓

Handler

     ↓

UseCase
```

Handlerでは、

```text
HTTP Request → UseCase Input
```

と、

```text
Domain → HTTP Response
```

の変換を担当します。

ビジネスロジックはHandlerに書きません。

例えば、

```text
タイトルが空なら作成できない
```

というルールはHandlerではなくUseCase側に置きます。

---

# 10. OpenAPI

MUSEUMではAPI仕様を、

```text
backend/api/openapi.yaml
```

で管理します。

OpenAPIを、

> HTTP APIのSource of Truth

として扱います。

例えば、

```yaml
/exhibitions/{id}:
  get:
    operationId: getExhibition
```

と定義すると、`oapi-codegen`によってGoコードを生成します。

```text
openapi.yaml

      ↓

oapi-codegen

      ↓

generated/server.gen.go
```

生成コードには、

```text
Request型
Response型
Path Parameter
Route
Server Interface
```

などが含まれます。

そのため、これらを毎回手作業で実装する必要がありません。

---

# 11. sqlc

SQLからGoコードを生成するために`sqlc`を使用します。

SQLは、

```text
backend/db/query/
```

に置きます。

例えば、

```sql
-- name: GetExhibitionByID :one

SELECT
    id,
    museum_id,
    title,
    description,
    created_at
FROM exhibitions
WHERE id = $1
LIMIT 1;
```

から、

```go
GetExhibitionByID(...)
```

というGoコードが生成されます。

流れは、

```text
SQL

 ↓

sqlc

 ↓

Go
```

です。

これによって、

```text
QueryRow
Scan
SQL Parameter Struct
DB Model
```

などを手作業で書く量を減らします。

---

# 12. Migration

Database SchemaはMigrationで管理します。

```text
backend/db/migrations/
```

例えば、

```text
000001_create_exhibitions.up.sql
000001_create_exhibitions.down.sql
```

です。

`up.sql`はSchemaを進めます。

```sql
CREATE TABLE exhibitions (...);
```

`down.sql`は元に戻します。

```sql
DROP TABLE exhibitions;
```

Databaseを直接手作業で変更するのではなく、

> Schema変更はMigrationとしてGit管理する

のが原則です。

---

# 13. Code Generation

MUSEUMでは、繰り返し発生するコードをできるだけ自動生成します。

現在は、

```text
OpenAPI
   │
   └── oapi-codegen
          ↓
       HTTP Code


SQL
 │
 └── sqlc
       ↓
    Database Code
```

という構成です。

生成は、

```bash
make generate
```

でまとめて実行します。

内部では、

```bash
cd backend
go generate ./...
```

が実行されます。

---

# 14. generatedファイルは編集しない

以下のようなディレクトリは自動生成です。

```text
api/generated/
```

```text
internal/infrastructure/exhibition/sqlcgen/
```

ここにあるコードは、

> 手動編集禁止

です。

変更したい場合は、

```text
OpenAPI生成コード
→ openapi.yamlを変更

sqlc生成コード
→ SQL / sqlc.yamlを変更
```

します。

そして、

```bash
make generate
```

します。

---

# 15. APIを追加するときの基本手順

例えば、

```text
GET /exhibitions/{id}
```

を追加するとします。

まずOpenAPIを変更します。

```text
backend/api/openapi.yaml
```

次にSQLを追加します。

```text
backend/db/query/exhibitions.sql
```

そして、

```bash
make generate
```

します。

その後、

```text
Repository
    ↓
UseCase
    ↓
Handler
```

を実装します。

つまり基本的な開発フローは、

```text
1. API仕様を書く
        ↓
2. SQLを書く
        ↓
3. make generate
        ↓
4. Repositoryを実装
        ↓
5. UseCaseを実装
        ↓
6. Handlerを実装
        ↓
7. Test
```

です。

---

# 16. 新しいDB Tableを追加するとき

例えば`exhibits`を追加するとします。

まずMigrationを作ります。

```text
db/migrations/

000002_create_exhibits.up.sql
000002_create_exhibits.down.sql
```

次にQueryを書く。

```text
db/query/exhibits.sql
```

その後、

```bash
make db-migrate-up
make generate
```

します。

---

# 17. Go初心者向け：packageとは

Goではディレクトリ単位でpackageを作ります。

例えば、

```text
internal/domain/
```

なら、

```go
package domain
```

です。

別packageを使用するときは、

```go
import "museum/internal/domain"
```

として、

```go
domain.Exhibition
```

のように使用します。

---

# 18. Go初心者向け：interfaceとは

MUSEUMではRepositoryでinterfaceをよく使います。

```go
type ExhibitionRepository interface {
    Save(...) error
}
```

これは、

> Saveという操作ができるもの

という「契約」です。

PostgreSQL実装が、

```go
func (r *PostgresRepository) Save(...) error
```

を持っていれば、このinterfaceを満たします。

Goでは、

```text
implements ExhibitionRepository
```

のような宣言は必要ありません。

必要なメソッドを持っていれば、自動的にinterfaceを満たします。

MUSEUMではさらに、

```go
var _ repository.ExhibitionRepository =
    (*PostgresRepository)(nil)
```

を書くことで、

> PostgresRepositoryがinterfaceを満たしているか

をコンパイル時に確認します。

---

# 19. Go初心者向け：context.Contextとは

Backendでは、

```go
ctx context.Context
```

が頻繁に登場します。

例えば、

```go
func (u *GetExhibitionUseCase) Execute(
    ctx context.Context,
    id string,
)
```

です。

HTTP Requestがキャンセルされた場合などに、その情報を、

```text
Handler
   ↓
UseCase
   ↓
Repository
   ↓
PostgreSQL
```

まで伝えるために使用します。

基本的には、

> Requestから受け取ったcontextを下のLayerへ渡していく

と覚えておけば大丈夫です。

---

# 20. Dependency Injection

`main.go`を見ると、

```go
repository :=
    exhibition.NewPostgresRepository(pool)

useCase :=
    usecase.NewCreateExhibitionUseCase(
        repository,
        idGenerator,
    )

handler :=
    handler.NewExhibitionHandler(
        useCase,
    )
```

のようになっています。

これはDependency Injectionと呼ばれる考え方です。

内部で勝手に、

```go
postgres.Connect()
```

するのではなく、必要なものを外から渡します。

```text
main

 ├── Repositoryを作る
 │
 ├── UseCaseへ渡す
 │
 └── Handlerへ渡す
```

`main.go`がアプリケーション全体を組み立てる場所になります。

---

# 21. エラーの境界

技術固有のエラーを上位Layerへ漏らしすぎないようにします。

例えばPostgreSQLでは、

```go
pgx.ErrNoRows
```

があります。

これをInfrastructureで、

```go
repository.ErrNotFound
```

へ変換します。

```text
PostgreSQL

pgx.ErrNoRows

      ↓

Infrastructure

repository.ErrNotFound

      ↓

UseCase / Handler

HTTP 404
```

こうすることでUseCaseがpgxに依存しません。

---

# 22. Backend Directory

現在の主要構成は次の通りです。

```text
backend/
│
├── api/
│   ├── openapi.yaml
│   ├── oapi-codegen.yaml
│   ├── generate.go
│   │
│   └── generated/
│       └── server.gen.go
│
├── db/
│   ├── generate.go
│   │
│   ├── migrations/
│   │   └── 000001_create_exhibitions.*
│   │
│   └── query/
│       └── exhibitions.sql
│
├── cmd/
│   └── server/
│       └── main.go
│
└── internal/
    │
    ├── domain/
    │   └── exhibition.go
    │
    ├── repository/
    │   ├── errors.go
    │   └── exhibition_repository.go
    │
    ├── usecase/
    │   ├── create_exhibition.go
    │   └── get_exhibition.go
    │
    ├── handler/
    │   └── exhibition_handler.go
    │
    └── infrastructure/
        │
        ├── database/
        │   └── postgres.go
        │
        └── exhibition/
            ├── postgres_repository.go
            ├── uuid_generator.go
            │
            └── sqlcgen/
```

---

# 23. Local Development

PostgreSQLを起動します。

```bash
make db-up
```

Migrationします。

```bash
make db-migrate-up
```

コード生成します。

```bash
make generate
```

テストします。

```bash
cd backend

go test ./...
```

Backendを起動します。

```bash
export DATABASE_URL="postgres://museum:museum@localhost:15432/museum?sslmode=disable"

go run ./cmd/server
```

---

# 24. 開発時の重要ルール

MUSEUMでは以下を基本ルールとします。

### Domainに技術詳細を持ち込まない

避ける：

```go
type Exhibition struct {
    Title string `json:"title" db:"title"`
}
```

DomainはMUSEUMの概念を表現します。

---

### Handlerにビジネスロジックを書かない

避ける：

```go
if title == "" {
    ...
}
```

のようなMUSEUM固有ルールをHandlerだけに実装すること。

UseCase / Domain側で判断します。

---

### SQLをGoコードに散らさない

避ける：

```go
db.Query(
    "SELECT * FROM exhibitions ...",
)
```

Queryは、

```text
db/query/
```

で管理し、sqlcを使用します。

---

### generatedを編集しない

避ける：

```text
api/generated/server.gen.go
```

を直接変更すること。

Source of Truthを変更して再生成します。

---

### 重複コードを増やさない

同じ処理が複数箇所に増えてきた場合は、そのままコピーせず、

```text
共通化
Code Generation
OpenAPI
sqlc
共通Middleware
共通Error Handling
```

などで変更箇所を減らせないか検討します。

ただし、将来必要になるかもしれないという理由だけで抽象化はしません。

実際に重複や変更コストが発生し始めたタイミングで導入します。

---

# 25. 設計思想

MUSEUMでは、

> 変更しやすいコード

を重視します。

ただし、

> Clean Architectureだから

> Design Patternだから

という理由だけで構造を複雑にしません。

問題が発生していない場所に無理にPatternを導入するのではなく、

```text
責務が大きくなった

変更箇所が増えた

同じコードが増えた

テストが難しくなった

特定技術への依存が強くなった
```

といった兆候が出たところを改善します。

---

# 26. 最初に覚えること

すべてを最初から理解する必要はありません。

新しく参加した場合は、まずこの流れだけ理解してください。

```text
HTTP Request

     ↓

Handler

     ↓

UseCase

     ↓

Repository

     ↓

Infrastructure

     ↓

PostgreSQL
```

そして、

```text
API変更
    ↓
openapi.yaml


DB Query変更
    ↓
db/query/*.sql


DB Schema変更
    ↓
db/migrations/


MUSEUMの処理変更
    ↓
usecase/


MUSEUMの概念変更
    ↓
domain/
```

という対応関係を覚えれば、MUSEUM Backendのコードを追えるようになります。
