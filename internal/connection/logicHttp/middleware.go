package logichttp

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("this is jwt secret")

type Claims struct {
	CallerId string `json:caller_id`
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(userid string, username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		userid,
		EncodeMD5(username),
		EncodeMD5(password),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "nbim",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// EncodeMD5 md5 encryption
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value := ctx.Request.Header.Get("Authorization")
		tokenstr := strings.SplitN(value, " ", 2)
		if tokenstr[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "JWT 格式不正确",
			})
			ctx.Abort()
			return
		}
		if tokenstr[1] == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "JWT 为空",
			})
			ctx.Abort()
			return
		}
		cliam, err := ParseToken(tokenstr[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "token 解析失败",
			})
			ctx.Abort()
			return
		} else if cliam.ExpiresAt.Unix() < time.Now().Unix() {
			ctx.JSON(400, gin.H{
				"error": "token 超时",
			})
			ctx.Abort()
			return
		}

		ctx.Set("caller_id", cliam.CallerId)
		ctx.Next()
	}
}
