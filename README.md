# 🏠 Property Transactions Microservice

A microservice for managing income and expense transactions by user and property.

---

## 🚀 Quick Start

### 1. Run Docker

```bash
cd cmd/property-transactions/
docker compose up -d
```

### 2. Run DB Migrations (Goose)

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir ./cmd/property-transactions/migrations/ clickhouse "clickhouse://myuser:mypassword@localhost:9000/default" up
```

---

## 🔌 API Endpoints

### ➕ POST `/property_transactions/v1/user/{userID}/`

**Description:** Add a new transaction for a user.

**Request Body:**

```json
{
  "propertyID": "property-123",
  "amount": 100.50,
  "date": 1712486400
}
```

**Response:**

```json
{
  "success": true,
  "message": "Transaction added successfully"
}
```

---

### 📥 GET `/property_transactions/v1/user/{userID}/property/{propertyID}/`

**Description:** Get all transactions for a specific property of a user.

**Response:**

```json
{
  "success": true,
  "data": {
    "transactions": [
      {
        "amount": 200.0,
        "date": 1712486400
      }
    ]
  }
}
```

---

### 💰 GET `/property_transactions/v1/user/{userID}/property/{propertyID}/balance/`

**Description:** Get current balance for a specific property.

**Response:**

```json
{
  "success": true,
  "data": {
    "balance": 1500.75
  }
}
```

---

### 📊 GET `/property_transactions/v1/user/{userID}/property/{propertyID}/monthly_report/`

**Description:** Get a detailed monthly report for a specific property.

**Response:**

```json
{
  "success": true,
  "data": {
    "monthlyBalance": {
      "startingCash": 1000.0,
      "records": [
        {
          "id": 1,
          "type": "income",
          "amount": 500.0,
          "total": 1500.0
        },
        {
          "id": 2,
          "type": "expense",
          "amount": -200.0,
          "total": 1300.0
        }
      ],
      "endCash": 1300.0
    }
  }
}
```

---

## 📂 Project Structure

```
cmd/property-transactions/
├── docker-compose.yml
├── property-transactions.go
├── migrations/
└── app/
```

---

## 🧾 Transaction Types

```go
type TransactionType string

const (
    TransactionTypeIncome  TransactionType = "income"
    TransactionTypeExpense TransactionType = "expense"
)
```

---

## 📅 2025 | Built with ❤️ by Meital


---

## 🧪 Example `curl` Commands

### ➕ Add Transaction (POST)

```bash
curl -X POST http://localhost/property_transactions/v1/user/3f9e3a47-7a91-4873-b27f-2b56e9cb06f0/   -H "Content-Type: application/json"   -d '{
    "propertyID": "4c8d1e8d-39ea-4df2-872a-8e4e45b0a119",
    "amount": 100.50,
    "date": '"$(date +%s)"'
  }'
```

---

### 📥 Get Transactions (GET)

```bash
curl http://localhost/property_transactions/v1/user/3f9e3a47-7a91-4873-b27f-2b56e9cb06f0/property/4c8d1e8d-39ea-4df2-872a-8e4e45b0a119/
```

---

### 💰 Get Balance (GET)

```bash
curl http://localhost/property_transactions/v1/user/3f9e3a47-7a91-4873-b27f-2b56e9cb06f0/property/4c8d1e8d-39ea-4df2-872a-8e4e45b0a119/balance/
```

---

### 📊 Get Monthly Report (GET)

```bash
curl http://localhost/property_transactions/v1/user/3f9e3a47-7a91-4873-b27f-2b56e9cb06f0/property/4c8d1e8d-39ea-4df2-872a-8e4e45b0a119/monthly_report/
```
