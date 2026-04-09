package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/Bank-Thanapat-Developer/basic-redis/pkg/cache"
	"github.com/gofiber/fiber/v3"
)

// RateLimit จำกัด request ต่อ IP ภายใน window time
// ใช้ Redis INCR + EXPIRE เพื่อ atomic fixed-window counter
// max    = จำนวน request สูงสุดต่อ window
// window = ช่วงเวลา เช่น 1*time.Minute
func RateLimit(cacheService cache.CacheService, max int64, window time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := c.IP()
		key := fmt.Sprintf("rate_limit:%s", ip)
		ctx := context.Background()

		// INCR คืนค่าหลังบวก 1 เสมอ (atomic)
		count, err := cacheService.Incr(ctx, key)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "rate limiter error",
			})
		}

		// request แรกใน window → ตั้ง TTL
		// request ต่อไปใน window เดิม → TTL ไม่เปลี่ยน
		if count == 1 {
			if err := cacheService.Expire(ctx, key, window); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "rate limiter error",
				})
			}
		}

		remaining := max - count
		if remaining < 0 {
			remaining = 0
		}

		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", max))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		if count > max {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "too many requests",
				"retry_after": window.Seconds(),
			})
		}

		return c.Next()
	}
}
