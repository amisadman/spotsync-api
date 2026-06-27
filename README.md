# 🚗 SpotSync – Smart Parking & EV Charging Reservation

SpotSync is a centralized backend API for busy airports and malls to manage parking zones and reserve limited EV charging spots, ensuring capacity limits are strictly enforced.

## 🔗 Project URLs
- **Live API Endpoint:** `https://spotsync-api.onrender.com` (Example Placeholder)
- **GitHub Repository:** `https://github.com/amisadman/spotsync-api`

---

## 🛠️ Technology Stack
- **Language:** Go (Golang) v1.26+
- **HTTP Framework:** Echo v5
- **ORM:** GORM (with PostgreSQL driver)
- **Database:** PostgreSQL (NeonDB)
- **Validator:** playground validator/v10
- **JWT:** golang-jwt/jwt/v5
- **Security:** bcrypt password hashing

---

## 🏛️ Clean Architecture & Layer Interactions

The project is structured according to Clean Architecture principles to separate concerns and ensure maintainability:

```
                  ┌──────────────────────┐
                  │        Client        │
                  └──────────┬───────────┘
                             │
                             ▼ (JWT Auth & Role Authorization)
                  ┌──────────────────────┐
                  │   HTTP Router/Echo   │
                  └──────────┬───────────┘
                             │
                             ▼ (DTO Binding & Validator)
                  ┌──────────────────────┐
                  │       Handlers       │
                  └──────────┬───────────┘
                             │
                             ▼ (Business Rules & Security)
                  ┌──────────────────────┐
                  │       Services       │
                  └──────────┬───────────┘
                             │
                             ▼ (GORM Transactions & Lockings)
                  ┌──────────────────────┐
                  │     Repositories     │
                  └──────────┬───────────┘
                             │
                             ▼
                  ┌──────────────────────┐
                  │   PostgreSQL DB      │
                  └──────────────────────┘
```

- **DTO (`dto/`):** Data Transfer Objects representing API request and response payloads.
- **Handler (`handler.go`):** HTTP layer which binds requests, validates inputs, and triggers services.
- **Service (`service.go`):** Core business logic containing calculations, bcrypt password hashing, and token signing.
- **Repository (`repository.go`):** Interacts with the database using GORM. Handles transactions and concurrency row locks.
- **Entity (`entity.go`):** Defines GORM models mapping to the database tables.

---

## ⚡ Concurrency-Safe Reservations (EV Spot Bottleneck)

To prevent over-allocation of limited EV spots during concurrent reservation attempts, the reservation layer utilizes a **database transaction** combined with **row-level locking** (`FOR UPDATE`) on the parking zone record:

1. **Locking Zone:** The transaction queries the target parking zone with an `UPDATE` lock strength, preventing concurrent transactions from updating the capacity or obtaining the lock until the current transaction commits/rolls back.
2. **Double Check:** Active reservations for the zone are counted within the same transaction.
3. **Validation:** If the active count is less than the zone's capacity, the reservation is created. Otherwise, a `409 Conflict` (ErrZoneFull) error is returned, aborting the transaction.

---

## 🚀 Setup & Running Locally

### 1. Prerequisites
- **Go** (Version 1.22 or higher)
- **PostgreSQL** database instance (e.g. NeonDB or local)

### 2. Environment Configuration
Create a `.env` file in the root directory:
```env
PORT=8080
DSN=postgresql://username:password@hostname:port/database?sslmode=require
JWT_SECRET=verysecretkey
```

### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Run the Server
Using Go:
```bash
go run cmd/main.go
```
Using Air (for hot-reload):
```bash
air
```

---

## 🌐 API Endpoints Specification

### 🔹 Authentication Module
- `POST /api/v1/auth/register` (Public) - Create user account. Supports `role` (`driver` or `admin`).
- `POST /api/v1/auth/login` (Public) - User login. Returns JWT token.
- `GET /api/v1/auth/me` (Authenticated) - Retrieve profile of the logged-in user.

### 🔹 Parking Zones Module
- `POST /api/v1/zones` (Admin Only) - Create a parking zone.
- `GET /api/v1/zones` (Public) - Retrieve all zones with their dynamically calculated `available_spots`.
- `GET /api/v1/zones/:id` (Public) - Retrieve details and availability for a single zone.

### 🔹 Reservations Module
- `POST /api/v1/reservations` (Authenticated) - Reserve a parking/EV spot. Safe against concurrent race conditions.
- `GET /api/v1/reservations/my-reservations` (Authenticated) - View own active and cancelled reservations.
- `DELETE /api/v1/reservations/:id` (Authenticated) - Cancel reservation (Drivers can cancel only their own; Admins can cancel any).
- `GET /api/v1/reservations` (Admin Only) - View all reservations in the system.
