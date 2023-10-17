package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	TokenExpiredError = errors.New("Token已過期")
)

var jwtSecret = []byte("ice_moss")

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken 根據用戶的用戶id和密碼 生成JWT token
func GenerateToken(userId uint, iss string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(48 * 30 * time.Hour)

	claims := Claims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //過期時間
			Issuer:    iss,               //token發行人
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret) //簽名
	return token, err
}

func JWY() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.PostForm("token")
		user := c.Query("user")
		userId, err := strconv.Atoi(user)
		zap.S().Infoln("userId:", userId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "user id 不合法",
			})
			c.Abort()
			return
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "token 不能為空",
			})
			c.Abort()
			return
		}

		claims, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "token 失效",
			})
			c.Abort()
			return
		}else if time.Now().Unix() > claims.ExpiresAt{
			err = TokenExpiredError
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg": "token 已過期",
			})
			c.Abort()
			return
		}
		
		if claims.UserID != uint(userId) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  "登入不合法",
			})
			c.Abort()
			return
		}
		fmt.Println("token 驗證成功")
		c.Next()
	}
}

func ParseToken(token string) (*Claims, error) {
	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
