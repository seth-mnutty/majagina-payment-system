# Payment System API

A lightweight Go payment management system for businesses to track incoming payments and store transaction logs in SQLite.

## Features

- Business management
- Payment transaction recording
- Transaction and audit logging
- Revenue summary
- SQLite storage
- REST API over HTTP

## Run

```bash
go mod tidy
go run .
```

The API runs on `http://localhost:8080`.

## Endpoints

### Health

```bash
curl http://localhost:8080/health
```

### Create business

```bash
curl -X POST http://localhost:8080/api/v1/businesses \
  -H "Content-Type: application/json" \
  -d '{"name":"Acme Labs","email":"billing@acmelabs.com"}'
```

### List businesses

```bash
curl http://localhost:8080/api/v1/businesses
```

### Create payment transaction

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "business_id":"biz-1234567890",
    "customer_name":"Jane Doe",
    "amount":299.99,
    "currency":"USD",
    "status":"paid",
    "reference":"INV-1001",
    "notes":"Monthly subscription"
  }'
```

### List transactions

```bash
curl "http://localhost:8080/api/v1/transactions?business_id=biz-1234567890"
```

### Summary

```bash
curl http://localhost:8080/api/v1/summary
```

## Data files

The app stores its SQLite database in `data/payment.db`.

## Notes

This implementation is intentionally simple and production-ready enough for a small business dashboard or internal payment API.
