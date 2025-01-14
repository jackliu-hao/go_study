package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

func TestJwtHandler_setJwt(t *testing.T) {
	var claim UserClaims

	claim.Uid = 1
	claim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(30 * time.Minute))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedString, err := token.SignedString([]byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(signedString)
	// 解密
	_, err = jwt.ParseWithClaims(signedString, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK"), nil
	})
	_, err = jwt.ParseWithClaims(signedString, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte("0776f450dd575004ba7c69930c579cae"), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claim)

}
