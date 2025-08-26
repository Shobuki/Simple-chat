package auth


import (
"errors"
"time"


"github.com/golang-jwt/jwt/v5"
)


type JWT struct { secret string }


func NewJWT(secret string) *JWT { return &JWT{secret: secret} }


type Claims struct {
UserID int64 `json:"user_id"`
DisplayName string `json:"display_name"`
jwt.RegisteredClaims
}


func (j *JWT) Sign(userID int64, displayName string) (string, error) {
claims := Claims{
UserID: userID,
DisplayName: displayName,
RegisteredClaims: jwt.RegisteredClaims{
ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
IssuedAt: jwt.NewNumericDate(time.Now()),
},
}
return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.secret))
}


func (j *JWT) Parse(tokenStr string) (*Claims, error) {
token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
return []byte(j.secret), nil
})
if err != nil { return nil, err }
if claims, ok := token.Claims.(*Claims); ok && token.Valid { return claims, nil }
return nil, errors.New("invalid token")
}