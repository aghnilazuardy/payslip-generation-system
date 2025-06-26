# payslip-generation-system
A scalable payroll system built using **Go** and **PostgreSQL** that supports attendance-based prorated salary, overtime payment, and reimbursement submissions.

## Installation

### Requirements

- Go 1.21+
- PostgreSQL
- Docker (for local dev)

### Setup

```bash
git clone https://github.com/your-org/payslip-generation-system.git
cd payslip-generation-system
cp .env.example .env
```

Update `.env` with your PostgreSQL credentials:

```env
DATABASE_DSN=postgres://user:password@localhost:5432/payslip?sslmode=disable
JWT_SECRET=your-jwt-secret
WORK_HOUR_START=9
WORK_HOUR_END=17
```

### Run with Docker

```bash
docker-compose up --build -d
```

### Apply Migrations

```bash
migrate -path migrations -database "$DATABASE_DSN" up
```

### Seed Database

```bash
go run cmd/seed/main.go
```

---

## Usage

### Login

```http
POST /login
{
  "username": "employee001",
  "password": "password"
}
```

Returns:

```json
{
  "token": "<JWT>"
}
```

Use `Bearer <token>` in `Authorization` header for all endpoints.

### Admin Endpoints

- `POST /admin/attendance-period`
- `POST /admin/payroll/run`
- `GET /admin/payslips`

### Employee Endpoints

- `POST /employee/attendance`
- `POST /employee/overtime`
- `POST /employee/reimbursement`
- `GET /employee/payslip`

---

## Testing

```bash
go test ./...
```

---
