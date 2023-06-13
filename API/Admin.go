package API

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sunbelife/Prelook_Strobe_Backend/Config"
	"github.com/sunbelife/Prelook_Strobe_Backend/Model"
	"net/http"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func AdminLoginIn(c *gin.Context)  {
	// 建立 Users 结构体实例，用来接受 GetUser 的用户
	var myadmin Model.Admins

	// 接受 POST 请求
	c.ShouldBind(&myadmin)

	VerifyResult := Model.CheckAuth(myadmin.Username, myadmin.Password)

	// 判断密码是否正确
	if VerifyResult == true {

		token, err := GenerateToken(myadmin.Username, myadmin.Password)

		if err != nil {
			// 输出为 JSON
			c.JSON(http.StatusBadRequest, gin.H{
				"Code": http.StatusBadRequest,
				"Message": "fail",
				"Err": err,
			})
		} else {
			// 输出为 JSON
			c.JSON(http.StatusOK, gin.H{
				"Code": http.StatusOK,
				"Message": "success",
				"token": token,
			})
		}
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Incorrect name or password",
		})
		return
	}
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		password,
		jwt.StandardClaims {
			ExpiresAt : expireTime.Unix(),
			Issuer : "Sun-Store",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	AppConfig, err := Config.GetAppConfig()
	var JWT_KEY = []byte(AppConfig.App.JWTKey)

	token, err := tokenClaims.SignedString(JWT_KEY)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	AppConfig, err := Config.GetAppConfig()
	var JWT_KEY = []byte(AppConfig.App.JWTKey)

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWT_KEY, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}