package middleware

import (
	"mikrotik-api/config"
	"net/http"

	"mikrotik-api/modules/auth"
	"mikrotik-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			c.Abort()
			return
		}

		claims := &utils.Claims{}
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(config.GetEnv("JWT_SECRET")), nil
			},
		)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		user, err := auth.GetActiveUserByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user disabled"})
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("role", user.Role)

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}


// func JWTAuth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
// 			c.Abort()
// 			return
// 		}

// 		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

// 		claims := &utils.Claims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return []byte("super_secret_key"), nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		// Inject into context
// 		c.Set("user_id", claims.UserID)
// 		c.Set("role", claims.Role)

// 		c.Next()
// 	}
// }
