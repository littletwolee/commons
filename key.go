package commons

import (
	"crypto/rsa"
	lib "github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

var (
	consKey *Key
)

type Key struct{}

func GetKey() *Key {
	if consKey == nil {
		consKey = &Key{}
	}
	return consKey
}

// @Title LoadRSAPrivateKeyFromDisk
// @Description load private key from disk
// @Parameters
//            path            string            key path
// @Returns key:*rsa.PrivateKey err:error
func (k *Key) LoadRSAPrivateKeyFromDisk(path string) (*rsa.PrivateKey, error) {
	var (
		keyData []byte
		key     *rsa.PrivateKey
		err     error
	)
	keyData, err = ioutil.ReadFile(path)
	if err != nil {
		goto RETURN
	}
	key, err = lib.ParseRSAPrivateKeyFromPEM(keyData)
	goto RETURN
RETURN:
	return key, err
}

// @Title LoadRSAPublicKeyFromDisk
// @Description load pubilc key from disk
// @Parameters
//            path            string            key path
// @Returns key:*rsa.PublicKey err:error
func (k *Key) LoadRSAPublicKeyFromDisk(path string) (*rsa.PublicKey, error) {
	var (
		keyData []byte
		key     *rsa.PublicKey
		err     error
	)
	keyData, err = ioutil.ReadFile(path)
	if err != nil {
		goto RETURN
	}
	key, err = lib.ParseRSAPublicKeyFromPEM(keyData)
	goto RETURN
RETURN:
	return key, err
}
