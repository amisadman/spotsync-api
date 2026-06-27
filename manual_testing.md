# SpotSync API – Manual Testing Guide

This guide contains complete request payloads, curl commands, and sample responses to test all the API endpoints of the **SpotSync** service.

You can run these tests against your local environment (`http://localhost:8080`) or your deployed service (e.g., `https://spotsync-api.onrender.com`). Replace `{{BASE_URL}}` with the appropriate server host.

---

## 🔹 Authentication Module

### 1. User Registration (Driver)
Creates a new account with the default role of `driver`.

- **Endpoint:** `POST {{BASE_URL}}/api/v1/auth/register`
- **Headers:** `Content-Type: application/json`

**Request Body:**
```json
{
  "name": "Jane Driver",
  "email": "jane.driver@spotsync.com",
  "password": "password123"
}
```

**cURL Command:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Driver", "email": "jane.driver@spotsync.com", "password": "password123"}'
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "name": "Jane Driver",
    "email": "jane.driver@spotsync.com",
    "role": "driver",
    "created_at": "2026-06-27T12:00:00Z",
    "updated_at": "2026-06-27T12:00:00Z"
  }
}
```

---

### 2. User Registration (Admin)
Creates an administrator account.

- **Endpoint:** `POST {{BASE_URL}}/api/v1/auth/register`
- **Headers:** `Content-Type: application/json`

**Request Body:**
```json
{
  "name": "Alex Admin",
  "email": "alex.admin@spotsync.com",
  "password": "adminPassword123",
  "role": "admin"
}
```

**cURL Command:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "Alex Admin", "email": "alex.admin@spotsync.com", "password": "adminPassword123", "role": "admin"}'
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 2,
    "name": "Alex Admin",
    "email": "alex.admin@spotsync.com",
    "role": "admin",
    "created_at": "2026-06-27T12:01:00Z",
    "updated_at": "2026-06-27T12:01:00Z"
  }
}
```

---

### 3. User Login
Authenticates and returns a JWT access token.

- **Endpoint:** `POST {{BASE_URL}}/api/v1/auth/login`
- **Headers:** `Content-Type: application/json`

**Request Body:**
```json
{
  "email": "jane.driver@spotsync.com",
  "password": "password123"
}
```

**cURL Command:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "jane.driver@spotsync.com", "password": "password123"}'
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJuYW1lIjoiSmFuZSBEcml2ZXIiLCJlbWFpbCI6ImphbmUuZHJpdmVyQHNwb3RzeW5jLmNvbSIsInJvbGUiOiJkcml2ZXIiLCJ0b2tlbl90eXBlIjoiYWNjcmVzcyIsImV4cCI6MTgwODgxNjAwMH0...",
    "user": {
      "id": 1,
      "name": "Jane Driver",
      "email": "jane.driver@spotsync.com",
      "role": "driver"
    }
  }
}
```

---

### 4. Get Current User Profile (`/me`)
Retrieves user profile claims using the Bearer JWT token.

- **Endpoint:** `GET {{BASE_URL}}/api/v1/auth/me`
- **Headers:** `Authorization: Bearer <your_jwt_token>`

**cURL Command:**
```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <your_jwt_token>"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "User details retrieved successfully",
  "data": {
    "id": 1,
    "name": "Jane Driver",
    "email": "jane.driver@spotsync.com",
    "role": "driver"
  }
}
```

---

## 🔹 Parking Zones Module

### 5. Create Parking Zone (Admin Only)
Requires admin privileges.

- **Endpoint:** `POST {{BASE_URL}}/api/v1/zones`
- **Headers:** 
  - `Content-Type: application/json`
  - `Authorization: Bearer <admin_jwt_token>`

**Request Body:**
```json
{
  "name": "Terminal 1 EV Charging",
  "type": "ev_charging",
  "total_capacity": 3,
  "price_per_hour": 6.50
}
```

**cURL Command:**
```bash
curl -X POST http://localhost:8080/api/v1/zones \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_jwt_token>" \
  -d '{"name": "Terminal 1 EV Charging", "type": "ev_charging", "total_capacity": 3, "price_per_hour": 6.50}'
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "Parking zone created successfully",
  "data": {
    "id": 1,
    "name": "Terminal 1 EV Charging",
    "type": "ev_charging",
    "total_capacity": 3,
    "price_per_hour": 6.50,
    "created_at": "2026-06-27T12:05:00Z",
    "updated_at": "2026-06-27T12:05:00Z"
  }
}
```

**Forbidden Response (403 Forbidden - when requested by a driver):**
```json
{
  "success": false,
  "message": "Access denied",
  "errors": "insufficient permissions"
}
```

---

### 6. Get All Parking Zones
Retrieves availability details. Publicly accessible.

- **Endpoint:** `GET {{BASE_URL}}/api/v1/zones`

**cURL Command:**
```bash
curl -X GET http://localhost:8080/api/v1/zones
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Parking zones retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Terminal 1 EV Charging",
      "type": "ev_charging",
      "total_capacity": 3,
      "available_spots": 3,
      "price_per_hour": 6.5,
      "created_at": "2026-06-27T12:05:00Z"
    }
  ]
}
```

---

### 7. Get Single Parking Zone
Publicly accessible.

- **Endpoint:** `GET {{BASE_URL}}/api/v1/zones/:id`

**cURL Command:**
```bash
curl -X GET http://localhost:8080/api/v1/zones/1
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Parking zone retrieved successfully",
  "data": {
    "id": 1,
    "name": "Terminal 1 EV Charging",
    "type": "ev_charging",
    "total_capacity": 3,
    "available_spots": 3,
    "price_per_hour": 6.5,
    "created_at": "2026-06-27T12:05:00Z"
  }
}
```

---

## 🔹 Reservations Module

### 8. Reserve Parking Spot
Creates an active reservation. Concurrency-safe.

- **Endpoint:** `POST {{BASE_URL}}/api/v1/reservations`
- **Headers:** 
  - `Content-Type: application/json`
  - `Authorization: Bearer <driver_jwt_token>`

**Request Body:**
```json
{
  "zone_id": 1,
  "license_plate": "EV-999-XYZ"
}
```

**cURL Command:**
```bash
curl -X POST http://localhost:8080/api/v1/reservations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <driver_jwt_token>" \
  -d '{"zone_id": 1, "license_plate": "EV-999-XYZ"}'
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "Reservation confirmed successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "zone_id": 1,
    "license_plate": "EV-999-XYZ",
    "status": "active",
    "created_at": "2026-06-27T12:10:00Z",
    "updated_at": "2026-06-27T12:10:00Z"
  }
}
```

**Conflict Response (409 Conflict - when capacity has been filled):**
```json
{
  "success": false,
  "message": "Reservation failed",
  "errors": "parking zone is full"
}
```

---

### 9. Get My Reservations
Retrieves reservations created by the logged-in user.

- **Endpoint:** `GET {{BASE_URL}}/api/v1/reservations/my-reservations`
- **Headers:** `Authorization: Bearer <driver_jwt_token>`

**cURL Command:**
```bash
curl -X GET http://localhost:8080/api/v1/reservations/my-reservations \
  -H "Authorization: Bearer <driver_jwt_token>"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "My reservations retrieved successfully",
  "data": [
    {
      "id": 1,
      "license_plate": "EV-999-XYZ",
      "status": "active",
      "zone": {
        "id": 1,
        "name": "Terminal 1 EV Charging",
        "type": "ev_charging"
      },
      "created_at": "2026-06-27T12:10:00Z"
    }
  ]
}
```

---

### 10. Cancel Reservation
Sets the reservation status to `cancelled`. Drivers can only cancel their own.

- **Endpoint:** `DELETE {{BASE_URL}}/api/v1/reservations/:id`
- **Headers:** `Authorization: Bearer <driver_jwt_token>`

**cURL Command:**
```bash
curl -X DELETE http://localhost:8080/api/v1/reservations/1 \
  -H "Authorization: Bearer <driver_jwt_token>"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Reservation cancelled successfully"
}
```

**Forbidden Response (403 Forbidden - when a driver tries to cancel another user's reservation):**
```json
{
  "success": false,
  "message": "Failed to cancel reservation",
  "errors": "forbidden: you cannot perform this action"
}
```

---

### 11. Get All Reservations (Admin Only)
Allows administrators to view all bookings in the database.

- **Endpoint:** `GET {{BASE_URL}}/api/v1/reservations`
- **Headers:** `Authorization: Bearer <admin_jwt_token>`

**cURL Command:**
```bash
curl -X GET http://localhost:8080/api/v1/reservations \
  -H "Authorization: Bearer <admin_jwt_token>"
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "All reservations retrieved successfully",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "user": {
        "id": 1,
        "name": "Jane Driver",
        "email": "jane.driver@spotsync.com"
      },
      "zone_id": 1,
      "zone": {
        "id": 1,
        "name": "Terminal 1 EV Charging",
        "type": "ev_charging"
      },
      "license_plate": "EV-999-XYZ",
      "status": "cancelled",
      "created_at": "2026-06-27T12:10:00Z",
      "updated_at": "2026-06-27T12:12:00Z"
    }
  ]
}
```
