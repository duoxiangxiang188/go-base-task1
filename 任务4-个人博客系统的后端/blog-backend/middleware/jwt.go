package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/18850341851/blog-backend/config"
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
)

// 自定义JWT声明（包含用户ID）
type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// 生成JWT令牌
func GenerateToken(userID uint) (string, error) {
	jwtCfg := config.LoadJWTConfig()
	//设置过期时间
	expirationTime := time.Now().Add(jwtCfg.ExpirationTime)

	//创建声明
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	//生成令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtCfg.SecretKey))
}

// 验证JWT令牌并领取用户ID
func validateToken(tokenString string) (uint, error) {
	jwtCfg := config.LoadJWTConfig()
	//解析令牌
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtCfg.SecretKey), nil
		},
	)

	if err != nil {
		return 0, err
	}

	//验证声明
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	}
	return 0, errors.New("invalid token")

}

// JWT认证中间件（用于保护需要登录的接口）
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供Authorization令牌"})
			c.Abort()
			return

		}

		//检查Bearer前缀
		const bearerPrefix = "Bearer "
		if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization格式错误（应为Bearer <token>）"})
			c.Abort()
			return
		}
		//提取并验证token
		tokenString := authHeader[len(bearerPrefix):]
		userID, err := validateToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}

		//将用户ID存入上下文，供后续接口使用
		c.Set("userID", userID)
		c.Next()
	}
}
