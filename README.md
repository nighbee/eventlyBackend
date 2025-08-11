
# Evently — Event Booking Backend

**Evently** is a production-ready backend system for event management and seat booking with role-based access control (user/admin), secure JWT authentication, Redis caching, and concurrency-safe booking.  
It provides public event listings with filtering, an admin CRUD interface, and a robust mechanism to prevent double reservations.

---

## 🚀 Features

### **Authentication & Roles**
- User registration and login
- JWT (HS256) authentication
- Role-based access: `user` and `admin`

### **Event Management**
- Public endpoints for listing events and filtering by date/location
- Admin CRUD operations for events
- Redis caching for:
  - Event list — 15 minutes
  - Event details — 30 minutes

### **Booking System**
- Unique index `(event_id, seat_number, active)` prevents double booking
- Seat capacity validation
- Soft cancellation (sets `active=false` while keeping history)
- Concurrency-safe booking process

### **Infrastructure**
- Docker Compose setup with PostgreSQL, Redis, and pgAdmin
- AutoMigrate for database models
- `.env` configuration with `godotenv`

---

## 📂 Project Structure

internal/
auth/ # Registration, login, JWT, bcrypt, middleware
events/ # Event model, public/admin endpoints, Redis caching
booking/ # Booking model, seat validation, booking endpoints
database/ # PostgreSQL and Redis connections
router/ # API routes setup
cmd/server/ # Entry point (main.go)
migrations/ # SQL migrations


---

## 🔗 API Endpoints

### **Public**
- `POST /register` — Register a new user
- `POST /login` — Login and receive a JWT
- `GET /events` — List events (filters: `location`, `from`, `to`)
- `GET /events/:id` — Get event details

### **User (JWT required)**
- `POST /bookings` — Create a booking
- `GET /bookings` — List my bookings
- `DELETE /bookings/:id` — Cancel a booking

### **Admin (JWT + role=admin)**
- `POST /admin/events` — Create event
- `PUT /admin/events/:id` — Update event
- `DELETE /admin/events/:id` — Delete event

---

## 🛠 Tech Stack

- **Language:** Go 1.23
- **Web Framework:** [Fiber](https://github.com/gofiber/fiber)
- **Database:** PostgreSQL + [GORM](https://gorm.io)
- **Cache:** Redis ([go-redis/v9](https://github.com/redis/go-redis))
- **Auth:** JWT (HS256), bcrypt
- **Config:** godotenv (.env files)
- **Infrastructure:** Docker Compose (PostgreSQL, Redis, pgAdmin)

---
