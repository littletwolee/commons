package commons

import (
	"errors"
	"fmt"
	"time"

	lib "github.com/dgrijalva/jwt-go"
)

var (
	consJwt *jwt
)

type jwt struct {
	HmacSampleSecret []byte
}

func GetJwt() *jwt {
	if consJwt == nil {
		consJwt = &jwt{}
	}
	jwtKey := GetConfig().GetString("jwt.jwtkey")
	consJwt.HmacSampleSecret = []byte(jwtKey)
	return consJwt
}

func getHmacMethod(method string) lib.SigningMethod {
	switch method {
	case "sha256":
		return lib.SigningMethodHS256
	case "sha384":
		return lib.SigningMethodHS384
	case "sha512":
		return lib.SigningMethodHS512
	default:
		return nil
	}
}

// @Title NewHmac
// @Description create new hmac sign by map
// @Parameters
//            method          string                          method
//            mapClaims       map[string]interface{}          a struct by map
// @Returns jwttoken:string err:error
func (j *jwt) NewHmac(method string, mapClaims map[string]interface{}) (string, error) {
	var (
		token         *lib.Token
		tokenString   string
		err           error
		signingMethod lib.SigningMethod
	)
	signingMethod = getHmacMethod(method)
	if signingMethod == nil {
		err = errors.New(fmt.Sprintf("Signing method : %s is bad", method))
		goto RETURN
	}
	token = lib.NewWithClaims(signingMethod, lib.MapClaims(mapClaims))
	tokenString, err = token.SignedString(j.HmacSampleSecret)
	goto RETURN
RETURN:
	return tokenString, err
}

// @Title ParseHmac
// @Description create new hmac sign by map
// @Parameters
//            tokenString     string                          token string
//            mapClaims       map[string]interface{}          a struct by map
// @Returns mapClaims:map[string]interface{} err:error
func (j *jwt) ParseHmac(tokenString string) (map[string]interface{}, error) {
	var (
		claims map[string]interface{}
		err    error
		ok     bool
		token  *lib.Token
	)
	token, err = lib.Parse(tokenString, func(token *lib.Token) (interface{}, error) {
		if _, ok := token.Method.(*lib.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return j.HmacSampleSecret, nil
	})
	if claims, ok = token.Claims.(lib.MapClaims); ok && token.Valid && beforeExp(claims["exp"].(string)) {
		goto RETURN
	} else {
		err = errors.New("Can't parse token")
		return nil, err
	}
RETURN:
	return claims, err
}
func beforeExp(t string) bool {
	tExp, err := time.Parse(time.RFC3339, t)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return time.Now().Before(tExp)
}
