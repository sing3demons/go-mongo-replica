package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func getKeys() (privateKey []byte, publicKey []byte, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Error("can't get working directory err: ", err)
		return nil, nil, err
	}

	privateKey, err = os.ReadFile(pwd + "/cert/id_rsa")
	if err != nil {
		logrus.Error("can't read private key err: ", err)
		return nil, nil, err
	}

	publicKey, err = os.ReadFile(pwd + "/cert/id_rsa.pub")
	if err != nil {
		logrus.Error("can't read public key err: ", err)
		return nil, nil, err
	}

	return privateKey, publicKey, nil

}

func getSecretKeyFromEnv() (privateKey []byte, publicKey []byte, err error) {
	private := os.Getenv("PRIVATE_KEY")
	if private == "" {
		return nil, nil, fmt.Errorf("private key not found")
	}
	privateKey, err = base64.StdEncoding.DecodeString(private)
	if err != nil {
		return nil, nil, err
	}

	public := os.Getenv("PUBLIC_KEY")
	if public == "" {
		return nil, nil, fmt.Errorf("public key not found")
	}
	publicKey, err = base64.StdEncoding.DecodeString(public)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, publicKey, nil
}

func GenerateToken(userId string) (token string, err error) {
	privateKey, _, err := getSecretKeyFromEnv()
	if err != nil {
		return "", err
	}

	// rsa, err := ssh.ParseRawPrivateKey(privateKey)
	rsa, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	claims := jwt.RegisteredClaims{
		Subject:   userId,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(rsa)
}

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.Request.Header.Get("Authorization")
		if s == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "unauthorized"})
			return
		}

		token := strings.TrimPrefix(s, "Bearer ")

		claims, err := validateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "unauthorized"})
			return
		}
		sub, err := claims.GetSubject()
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"message": "unauthorized"})
			return
		}

		c.Set("userId", sub)
		c.Next()
	}

}

func validateToken(token string) (jwt.MapClaims, error) {
	_, publicKey, err := getSecretKeyFromEnv()
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims, nil
}
