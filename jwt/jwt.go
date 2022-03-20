package jwt

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sebrupomo/sebrupomo-backend/db"
	"github.com/sebrupomo/sebrupomo-backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var JwtMiddlware echo.MiddlewareFunc = middleware.JWTWithConfig(middleware.JWTConfig{
	Claims:      &JwtClaim{},
	TokenLookup: "header:Authorization",
	AuthScheme:  "Bearer",
	SigningKey:  []byte("secret"),
})

func Setup() {
	middleware.ErrJWTMissing = echo.NewHTTPError(http.StatusUnauthorized, model.AccessError())
	middleware.ErrJWTInvalid = echo.NewHTTPError(http.StatusUnauthorized, model.AccessError())
}

type JwtClaim struct {
	ID    primitive.ObjectID
	Admin bool
	jwt.StandardClaims
}

func GetJwtUser(c echo.Context) (*model.User, error) {
	id := c.Get("user").(*jwt.Token).Claims.(*JwtClaim).ID
	return db.FindUserByID(id)
}

func GenerateToken(user *model.User) (string, error) {
	claims := &JwtClaim{
		user.ID,
		false,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("secret"))
}
