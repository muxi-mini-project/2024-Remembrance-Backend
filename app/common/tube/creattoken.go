package tube

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 用于包含payload信息
type CustomClaims struct {
	UserID uint `json:"userid"`
	jwt.StandardClaims
}

func CreateToken(userid uint) (string, error) {
	extime := time.Now().Add(2 * time.Hour)
	claims := &CustomClaims{
		UserID: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: extime.Unix(),
		},
	}
	// 使用 HMAC SHA256 算法创建一个新的JWT令牌
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("remembrance"))
	return token, err
}
