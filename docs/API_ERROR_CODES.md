# API 錯誤碼定義文件

> **文檔位置**：`docs/API_ERROR_CODES.md`  
> **GitHub 連結**：https://github.com/victorbuild/go-gin-restful-api/blob/master/docs/API_ERROR_CODES.md

## 📋 設計原則

### 1. HTTP Status Code（HTTP 狀態碼）
- **用途**：表達錯誤的**層級**（使用者端錯誤 vs 伺服器端錯誤）
- **標準**：遵循 HTTP 標準（RFC 7231）
- **範圍**：400-599

### 2. Error Code（應用程式錯誤碼）
- **用途**：表達錯誤的**具體類型**，用於前端做不同處理
- **範圍**：
  - `1000-1999`：使用者端錯誤（Client Errors）
  - `4000-4999`：伺服器端錯誤（Server Errors）

### 3. 使用規則
- **HTTP Status Code**：必須使用，表達錯誤層級
- **Error Code**：必須使用，表達錯誤類型
- **Message**：提供人類可讀的錯誤訊息
- **Trace ID**：每次回應都包含，用於追蹤錯誤（在 HTTP Header `X-Trace-ID` 中）

---

## 📊 HTTP Status Code 與 Error Code 對應表

### 400 Bad Request（使用者端請求錯誤）

| Error Code | 錯誤名稱 | 說明 | 使用場景 |
|-----------|---------|------|---------|
| `1001` | `CodeInvalidInput` | JSON 格式錯誤或無效輸入 | 請求的 JSON 格式不正確、無效的參數值 |
| `1002` | `CodeMissingRequiredFields` | 缺少必填欄位 | 請求缺少必填欄位（目前未使用，保留） |

**範例回應**：
```json
{
  "status": "error",
  "message": "Invalid input",
  "error_code": 1001
}
```

---

### 401 Unauthorized（未授權）

| Error Code | 錯誤名稱 | 說明 | 使用場景 |
|-----------|---------|------|---------|
| `1004` | `CodeInvalidCredentials` | 帳號密碼錯誤 | 登入時帳號或密碼錯誤 |
| `1006` | `CodeTokenMissing` | Token 缺失 | 請求缺少 Authorization header |
| `1007` | `CodeTokenInvalid` | Token 無效或過期 | Token 格式錯誤、已過期、或已被撤銷 |
| `1008` | `CodeRefreshTokenInvalid` | Refresh Token 無效或過期 | Refresh Token 無效、過期、或已被撤銷 |

**範例回應**：
```json
{
  "status": "error",
  "message": "Invalid token",
  "error_code": 1007
}
```

**前端處理建議**：
- `1006`：提示使用者「請先登入」
- `1007`：自動嘗試使用 Refresh Token 更新 Access Token
- `1008`：清除本地 Token，導向登入頁

---

### 403 Forbidden（禁止存取）

| Error Code | 錯誤名稱 | 說明 | 使用場景 |
|-----------|---------|------|---------|
| `1010` | `CodeForbidden` | 權限不足 | 使用者已登入，但沒有足夠權限存取該資源 |

**範例回應**：
```json
{
  "status": "error",
  "message": "Forbidden",
  "error_code": 1010
}
```

---

### 404 Not Found（找不到資源）

| Error Code | 錯誤名稱 | 說明 | 使用場景 |
|-----------|---------|------|---------|
| `1009` | `CodeUserNotFound` | 使用者不存在 | 查詢的使用者 ID 不存在 |

**範例回應**：
```json
{
  "status": "error",
  "message": "User not found",
  "error_code": 1009
}
```

---

### 409 Conflict（資源衝突）

| Error Code | 錯誤名稱 | 說明 | 使用場景 |
|-----------|---------|------|---------|
| `1003` | `CodeEmailExists` | Email 已存在 | 註冊時 Email 已被使用 |
| `1013` | `CodeEmailInUse` | Email 已被使用 | 更新使用者時 Email 已被其他使用者使用 |

**範例回應**：
```json
{
  "status": "error",
  "message": "Email already registered",
  "error_code": 1003
}
```

---

### 415 Unsupported Media Type（不支援的媒體類型）

| Error Code | 錯誤名稱 | 說明 | 使用場景 |
|-----------|---------|------|---------|
| `1000` | `CodeUnsupportedMediaType` | Content-Type 錯誤 | 請求的 Content-Type 不是 `application/json` |

**範例回應**：
```json
{
  "status": "error",
  "message": "Unsupported media type. Expected application/json",
  "error_code": 1000
}
```

---

### 500 Internal Server Error（伺服器內部錯誤）

| Error Code | 錯誤名稱 | 說明 | 使用場景 |
|-----------|---------|------|---------|
| `4001` | `CodeInternalError` | 伺服器內部錯誤（通用） | 未分類的伺服器錯誤、Token 解析錯誤、未知錯誤 |
| `4002` | `CodePasswordHashFail` | 密碼加密錯誤 | 密碼雜湊處理失敗 |
| `4003` | `CodeDatabaseError` | 資料庫錯誤 | 資料庫查詢、更新、刪除失敗 |
| `1005` | `CodeTokenGenerationFailed` | Token 產生失敗 | JWT Token 產生失敗 |
| `1011` | `CodeDeleteUserFailed` | 刪除使用者失敗 | 刪除使用者時發生錯誤 |
| `1012` | `CodeUpdateUserFailed` | 更新使用者失敗 | 更新使用者時發生錯誤 |

**範例回應**：
```json
{
  "status": "error",
  "message": "Internal server error",
  "error_code": 4001
}
```

**重要說明**：
- 500 錯誤的詳細資訊記錄在伺服器 log 中（包含 `Trace ID`）
- 使用者收到的是通用錯誤訊息，避免暴露系統細節
- 如需除錯，請提供 `Trace ID` 給工程師，工程師可在 log 中查找詳細錯誤資訊

---

## 🔍 錯誤追蹤流程

### 當使用者遇到錯誤時：

1. **使用者收到錯誤回應**，包含 `Trace ID`（在 HTTP Header `X-Trace-ID` 中）
   ```json
   {
     "status": "error",
     "message": "Internal server error",
     "error_code": 4001
   }
   ```
   HTTP Header: `X-Trace-ID: 550e8400-e29b-41d4-a716-446655440000`

2. **使用者提供 `Trace ID` 給工程師**

3. **工程師在 log 系統中搜尋 `Trace ID`**，找到詳細錯誤資訊：
   ```json
   {
     "timestamp": "2025-01-24T10:00:00Z",
     "trace_id": "550e8400-e29b-41d4-a716-446655440000",
     "api": "POST /v1/auth/login",
     "level": "error",
     "message": "failed to login: email=user@example.com, error=connection timeout"
   }
   ```

---

## 📝 錯誤碼範圍分配

### 使用者端錯誤（1000-1999）
- `1000-1099`：Auth 相關錯誤
- `1100-1999`：保留給未來擴充

### 伺服器端錯誤（4000-4999）
- `4000-4999`：Server 相關錯誤

---

## ✅ 使用檢查清單

新增錯誤時，請確認：

- [ ] HTTP Status Code 是否正確（符合 HTTP 標準）
- [ ] Error Code 是否在正確範圍內（1000-1999 或 4000-4999）
- [ ] Error Code 是否與其他錯誤不重複
- [ ] Message 是否清楚易懂
- [ ] 500 錯誤是否包含詳細資訊在 log 中
- [ ] 文件是否已更新

---

## 🔄 維護指南

### 新增錯誤碼時：

1. 在 `internal/util/errors.go` 中定義新的 Error Code 常數
2. 在 `internal/util/errors.go` 中新增對應的 Helper 函式（如需要）
3. 更新本文件（`docs/API_ERROR_CODES.md`）
4. 更新 API 文件（Swagger）

### 修改錯誤碼時：

1. 確認沒有其他地方使用該錯誤碼
2. 更新所有相關文件
3. 通知前端團隊（如有）

---

## 📚 參考資料

- [HTTP Status Codes (RFC 7231)](https://tools.ietf.org/html/rfc7231#section-6)
- [RESTful API 設計最佳實踐](https://restfulapi.net/)

---

**最後更新**：2025-01-24  
**維護者**：Victor

