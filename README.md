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

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/health` | No | Health check |
| POST | `/api/v1/auth/register` | No | สมัครสมาชิก |
| POST | `/api/v1/auth/login` | No | Login เพื่อรับ JWT token |
| GET | `/api/v1/me` | Yes | ดูข้อมูล user จาก token ปัจจุบัน |
| POST | `/api/v1/items/` | Yes | สร้าง item |
| GET | `/api/v1/items/:id` | Yes | ดึง item ตาม id |
| GET | `/api/v1/items/list-with-redis` | Yes | ดึง list items ผ่าน Redis cache |
| GET | `/api/v1/items/list-with-out-redis` | Yes | ดึง list items ไม่ผ่าน Redis cache |
| PUT | `/api/v1/items/:id` | Yes | แก้ไข item ตาม id |
| DELETE | `/api/v1/items/:id` | Yes | ลบ item ตาม id |
| POST | `/api/v1/ref-item-types/` | Yes | สร้างประเภทสินค้า |
| GET | `/api/v1/ref-item-types/` | Yes | ดึง list ประเภทสินค้า |
| GET | `/api/v1/ref-item-types/:id` | Yes | ดึงประเภทสินค้าตาม id |

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

# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

ใช้ token จาก login:

```bash
TOKEN=<access_token_from_login>

# ดึงข้อมูลผู้ใช้ปัจจุบัน
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer $TOKEN"

# ดึง list items ผ่าน Redis cache
curl http://localhost:8080/api/v1/items/list-with-redis \
  -H "Authorization: Bearer $TOKEN"

# ดึง list items ไม่ผ่าน cache
curl http://localhost:8080/api/v1/items/list-with-out-redis \
  -H "Authorization: Bearer $TOKEN"

# สร้าง ref item type
curl -X POST http://localhost:8080/api/v1/ref-item-types \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "อาหาร"}'

# สร้าง item
curl -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "ข้าวผัด", "price": 60, "is_active": true, "ref_item_type_id": 1}'
```

## ทดสอบ Rate Limit

ระบบปัจจุบันใช้ Redis สำหรับ rate limit แบบ fixed window:

- ทุก route: `60 req/min ต่อ IP`
- `POST /api/v1/auth/login`: `10 req/min ต่อ IP`

ตัวอย่างทดสอบ login rate limit:

```bash
for i in $(seq 1 12); do
  curl -s -o /dev/null -w "request $i => HTTP %{http_code}\n" \
    -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser","password":"wrong-password"}'
done
```

ผลที่คาดหวัง:

- ช่วงแรกจะได้ `401 Unauthorized`
- เมื่อเกิน limit จะได้ `429 Too Many Requests`

สามารถดู response headers ได้ด้วย:

```bash
curl -i -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"wrong-password"}'
```

headers ที่เกี่ยวข้อง:

- `X-RateLimit-Limit`
- `X-RateLimit-Remaining`

## ทดสอบ Distributed Lock

ระบบปัจจุบันใช้ Redis distributed lock เพื่อป้องกัน duplicate create และ concurrent update:

- create: lock ตาม `item name`
- update: lock ตาม `item id`
- update ที่เปลี่ยนชื่อ: lock ตาม `item name` เพิ่มอีกชั้น

### 1. ทดสอบ duplicate create แบบ parallel

ยิง request พร้อมกันหลายครั้งด้วยชื่อเดียวกัน:

```bash
for i in $(seq 1 10); do
  curl -s -o /dev/null -w "request $i => HTTP %{http_code}\n" \
    -X POST http://localhost:8080/api/v1/items \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"race-test-item","price":99,"is_active":true,"ref_item_type_id":1}' &
done
wait
```

ผลที่คาดหวัง:

- 1 request ได้ `201 Created`
- request อื่นจะได้ `409 Conflict`

### 2. ทดสอบ lock ค้างด้วย Redis CLI

เปิด Redis CLI:

```bash
docker exec -it docker-redis-1 redis-cli
```

สร้าง lock ด้วยตัวเอง:

```bash
SET lock:item:name:manual-test fake-token EX 30
```

จากนั้นยิง create:

```bash
curl -i -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"manual-test","price":50,"is_active":true,"ref_item_type_id":1}'
```

ผลที่คาดหวัง:

- ได้ `409 Conflict`
- response body ประมาณ `item is being created, please retry`

### 3. ดู lock key ใน Redis

```bash
docker exec -it docker-redis-1 redis-cli KEYS 'lock:item:*'
```

หรือดู command แบบ realtime:

```bash
docker exec -it docker-redis-1 redis-cli MONITOR
```

สิ่งที่ควรเห็น:

- `SET ... NX EX ...` ตอน acquire lock
- `DEL ...` ผ่าน Lua script ตอน release lock
