package authn

import "github.com/gofiber/fiber/v2"

const ClaimsLocalKey = "auth_claims"

func ClaimsFromContext(c *fiber.Ctx) (Claims, bool) {
	claims, ok := c.Locals(ClaimsLocalKey).(Claims)
	return claims, ok
}
