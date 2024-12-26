package auth

import (
	"net/http"

	"github.com/abdulazizax/yelp/config"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// AuthzMiddleware is a middleware for role-based access control using Casbin.
func AuthzMiddleware(enforcer *casbin.Enforcer, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the user's role and request details.
		role := getRole(c, config)
		path := c.FullPath()       // Get the full request path.
		method := c.Request.Method // Get the HTTP method.

		// Check permissions using the Casbin enforcer.
		ok := enforcer.Enforce(role, path, method)
		if !ok {
			// If the user is not authorized.
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

// getRole retrieves the user's role from the request context or configuration.
func getRole(c *gin.Context, config *config.Config) string {
	// Get the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "" // No token provided
	}

	// Parse the JWT token
	token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(config.JWT.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return "" // Invalid token
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "" // Claims not in the expected format
	}

	// Get the role from claims
	if role, ok := claims["role"].(string); ok {
		return role
	}

	return "" // Role not found
}
