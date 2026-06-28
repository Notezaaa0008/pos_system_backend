# POS System Backend

RESTful API สำหรับระบบ **Point of Sale (POS)** พัฒนาด้วย **Golang** โดยใช้ **Gin Framework**, **GORM**, และ **PostgreSQL** พร้อมระบบ Authentication, Authorization และโครงสร้างแบบ Modular เพื่อให้ง่ายต่อการขยายระบบในอนาคต

---

## Current Features

- JWT Authentication
- Access Token / Refresh Token
- Role-Based Authorization (RBAC)
- Role Management
- PostgreSQL Database
- GORM ORM
- Request Validation
- Email Service
- Cloudinary Image Upload
- Environment Configuration
- CORS Configuration
- Database Auto Migration

---

## Future Roadmap

- Product Module
- Category Module
- Inventory Management
- POS Transactions
- Order History
- Dashboard
- Reports
- Barcode Support
- Docker Deployment
- Swagger API Documentation
- Unit Testing

---

## Tech Stack

| Technology | Description        |
| ---------- | ------------------ |
| Golang     | Backend Language   |
| Gin        | HTTP Framework     |
| GORM       | ORM                |
| PostgreSQL | Database           |
| JWT        | Authentication     |
| Cloudinary | Image Storage      |
| Go-Mail    | Email Service      |
| Validator  | Request Validation |

---

## Project Structure

```text
cmd/
└── app/
    └── main.go

internal/
├── database/
├── middleware/
├── models/
├── module/
│   ├── auth/
│   ├── roles/
│   └── ...
├── routes/
└── utils/

pkg/
└── validator/

.env
go.mod
```

---

## Architecture

```text
Client
   │
Gin Router
   │
Middleware
   │
Controller
   │
Service
   │
Repository
   │
PostgreSQL
```

---

## Installation

Clone repository

```bash
git clone https://github.com/Notezaaa0008/pos_system_backend.git
```

Move into project

```bash
cd pos_system_backend
```

Install dependencies

```bash
go mod download
```

Create environment file

```bash
cp .env.example .env
```

Run application

```bash
go run cmd/app/main.go
```

---

## API Base URL

```
http://localhost:8080/api/v1
```

---

## Authentication

The API uses **JWT Authentication**.

Include the access token in every protected request.

```
Authorization: Bearer <access_token>
```

---

## Current Modules

- Authentication
- User
- Role

The project is designed in a modular architecture, allowing new modules such as Product, Order, Customer, Inventory, and Report to be added easily.

---

## Database

- PostgreSQL
- Auto Migration using GORM

---

## Security

- Password Hashing
- JWT Authentication
- Role-Based Authorization
- Request Validation
- Middleware Protection

---

## Author

Developed by **Notezaaa0008**

---

## License

This project is for educational and portfolio purposes.
