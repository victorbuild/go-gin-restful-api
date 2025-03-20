# Go Gin RESTful API

基於 Golang + Gin 架構的 RESTful API，包含 JWT 驗證、CRUD 操作、資料庫連接（PostgreSQL）、Kafka，主要以 user 資源為主，包含註冊、登入、`admin` 有 CRUD 會員 user 的權限。

## 使用技術

需要有這些環境，必且在 env 設定。

| 類別        | 技術與框架                |
|-----------|----------------------|
| **後端框架**  | Golang（Gin）          |
| **身分驗證**  | JWT                  |
| **資料庫**   | PostgreSQL（GORM ORM） |
| **快取機制**  | Redis                |
| **QUEUE** | Kafka、RabbitMQ       |
| **EMAIL** | SMTP                 |
| **日誌管理**  | ELK                  |
| **監測工具**  | Prometheus、Grafana   |
| **測試工具**  | K6（壓力測試）             |


## 快速開始

### 安裝 Go

請安裝GO，並且確認你的 Go 版本為 1.23 或更之後更新的版本。

### 安裝依賴
```
go mod tidy
```

### 啟動專案

執行以下指令

```
go run . 
```

## 專案結構

```
go-gin-restful-api/
│── config/           # 配置（如環境變數、DB config）
│   ├── database.go
│   ├── env.go
│   ├── jwt.go
│   ├── kafka.go
│   ├── rabbitmq.go
│   ├── redis.go
│   ├── smtp.go
│   ├── ...
│
│── controllers/      # 控制器（處理 HTTP 請求）
│   ├── admin_controller.go
│   ├── auth_controller.go
│   ├── user_controller.go
│
│── database/         # 資料庫相關
│   ├── postgresql.go
│
│── middlewares/         # 中介層 （JWT 驗證、日誌、錯誤處理）
│
│── models/           # 資料庫模型（定義 DB 結構）
│   ├── user.go
│
│── pkg/           # 作為可以復用的套件
│
│── repositories/     # 資料存取層（處理 DB 查詢，解耦模型與控制器）
│   ├── user_repository.go
│
│── routes/           # 定義路由
│   ├── user_router.go
│
│── services/         # 處理商業邏輯
│   ├── user_service.go
│
│── tests/            # 測試程式（K6壓測、單元測試、整合測試）
│
│── utils/            # 工具函式
│   ├── response.go
│
│── .env              # 環境變數（DB 連線、JWT Secret）
│── .env.example      # 環境範例
│── .gitignore        # 忽略不必要 git 的檔案
│── go.mod            # Go Module
│── go.sum            # 依賴鎖定檔 類似 PHP composer.lock
│── main.go
```

## 測試

執行壓力測試（K6），請確認是否安裝 K6

```sh
k6 run tests/load_test.js
```
