package upbit

import (
	"crypto/sha512"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func RequestToken(query string) (string, error) {
	var token *jwt.Token

	accesskey := os.Getenv("UPBIT_ACCESS_KEY")
	secretkey := os.Getenv("UPBIT_SECRET_KEY")

	if query == "" {
		token = jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"access_key": accesskey,
				"nonce":      uuid.New().String(),
			})
	} else {
		sha_512 := sha512.New()
		sha_512.Write([]byte(query))
		query_hash := fmt.Sprintf("%x", sha_512.Sum(nil))
		token = jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"access_key":     accesskey,
				"nonce":          uuid.New().String(),
				"query_hash":     query_hash,
				"query_hash_alg": "SHA512",
			})
	}

	return token.SignedString([]byte(secretkey))
}
