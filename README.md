# basic-redis

ตัวอย่างการใช้งาน Redis เป็น Cache layer ร่วมกับ Go (Fiber v3) และ PostgreSQL

## Tech Stack

- **Go 1.25** — ภาษาหลัก
- **Fiber v3** — HTTP Framework
- **GORM** — ORM สำหรับ PostgreSQL
- **Redis** — Cache layer
- **PostgreSQL 15** — ฐานข้อมูลหลัก
- **Docker / Docker Compose** — สำหรับ infrastructure

## Project Structure

```
.
├── cmd/
│   ├── app/        # main application
│   └── seed/       # seed ข้อมูลเข้า database
├── config/         # โหลด config จาก .env
├── docker/         # Dockerfile และ docker-compose.yml
├── internal/
│   ├── domains/    # interfaces (Repository + Usecase)
│   ├── dto/        # request / response structs
│   ├── entities/   # database models
│   ├── handlers/   # HTTP handlers (Fiber)
│   ├── repositories/ # database layer
│   ├── usecases/   # business logic
│   └── route.go    # ลงทะเบียน routes ทั้งหมด
└── pkg/
    ├── cache/      # Redis cache service
    └── database/   # PostgreSQL connection
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/health` | Health check |
| POST | `/api/v1/items` | สร้าง item |
| GET | `/api/v1/items/:id` | ดึง item ตาม id (มี cache) |
| GET | `/api/v1/items/list-with-redis` | ดึง list items (ผ่าน Redis cache) |
| GET | `/api/v1/items/list-with-out-redis` | ดึง list items (ไม่ผ่าน cache) |
| POST | `/api/v1/ref-item-types` | สร้างประเภทสินค้า |
| GET | `/api/v1/ref-item-types` | ดึง list ประเภทสินค้า |
| GET | `/api/v1/ref-item-types/:id` | ดึงประเภทสินค้าตาม id |

---

## วิธีรัน

### วิธีที่ 1: รันด้วย Docker Compose (แนะนำ)

ไม่ต้องติดตั้ง PostgreSQL หรือ Redis เองบนเครื่อง

```bash
# 1. Clone project
git clone <repo-url>
cd basic-redis

# 2. รัน infrastructure + app พร้อมกัน
docker compose -f docker/docker-compose.yml up --build
```

App จะพร้อมใช้งานที่ `http://localhost:8080`

---

### วิธีที่ 2: รันบนเครื่อง (Local)

ต้องการ: Go 1.25+, PostgreSQL, Redis ติดตั้งบนเครื่องหรือรันผ่าน Docker

#### 2.1 เตรียม Infrastructure (PostgreSQL + Redis)

```bash
# รันแค่ postgres กับ redis ผ่าน docker compose
docker compose -f docker/docker-compose.yml up postgres redis -d
```

#### 2.2 ตั้งค่า Environment

```bash
# คัดลอกไฟล์ config
cp .env.example .env
```

แก้ไข `.env` ตามความเหมาะสม:

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=basic_redis
DB_SSLMODE=disable

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

#### 2.3 รัน App

```bash
go run cmd/app/main.go
```

#### 2.4 Seed ข้อมูลตัวอย่าง (ไม่บังคับ)

```bash
# เพิ่มข้อมูล RefItemType 4 ประเภท และ Item 100 รายการ
go run cmd/seed/main.go
```

> **หมายเหตุ:** seed ใช้ `FirstOrCreate` — รันซ้ำได้ปลอดภัย ข้อมูลไม่ซ้ำ

---

## ทดสอบ API

```bash
# Health check
curl http://localhost:8080/api/v1/health

# ดึง list items ผ่าน Redis cache
curl http://localhost:8080/api/v1/items/list-with-redis

# ดึง list items ไม่ผ่าน cache
curl http://localhost:8080/api/v1/items/list-with-out-redis

# สร้าง ref item type
curl -X POST http://localhost:8080/api/v1/ref-item-types \
  -H "Content-Type: application/json" \
  -d '{"name": "อาหาร"}'

# สร้าง item
curl -X POST http://localhost:8080/api/v1/items \
  -H "Content-Type: application/json" \
  -d '{"name": "ข้าวผัด", "price": 60, "is_active": true, "ref_item_type_id": 1}'
```
