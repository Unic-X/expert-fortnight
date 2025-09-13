package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	SecretKey string
}

func NewJWTConfig() *JWTConfig {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "default-secret-key"
	}
	return &JWTConfig{
		SecretKey: secretKey,
	}
}

func (j *JWTConfig) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract the token from the header (format: "Bearer <token>")
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(j.SecretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			c.Set("user_id", claims.UserID)
			c.Set("user_role", claims.UserType)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateToken creates a new JWT token for a user
func (j *JWTConfig) GenerateToken(userID, userType string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		UserType: userType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// GetUserFromContext extracts user information from the context
func GetUserFromContext(c *gin.Context) (userID, userType string, ok bool) {
	userIDVal, exists1 := c.Get("userID")
	userTypeVal, exists2 := c.Get("userType")

	if !exists1 || !exists2 {
		return "", "", false
	}

	userID, ok1 := userIDVal.(string)
	userType, ok2 := userTypeVal.(string)

	if !ok1 || !ok2 {
		return "", "", false
	}

	return userID, userType, true
}
