# Personal Finance Tracker

A Go-based personal finance tracker service built as a small HTTP API for recording income and expenses, calculating totals and reviewing monthly summaries.

## Features

- Add transactions with type, amount, category, description and date
- List and view individual transactions
- Delete transactions
- Review totals and net balance by month
- Filter transactions and summaries by category
- Persist to a JSON file on disk
- Docker-ready deployment

## Run locally

```bash
go run .
```

Then open:

- http://localhost:8080/
- http://localhost:8080/health
- http://localhost:8080/version

## Example API calls

```bash
curl -X POST http://localhost:8080/transactions \
  -H 'Content-Type: application/json' \
  -d '{"type":"income","amount":2500,"category":"salary","description":"Monthly pay","date":"2026-11-01"}'

curl http://localhost:8080/summary?month=2026-11
```
