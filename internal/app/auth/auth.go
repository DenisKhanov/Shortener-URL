package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

// Claims — claims structure that includes standard claims and UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID uint32
}

const (
	TokenExp  = time.Hour * 3
	SecretKey = "SnJSkf123jlLKNfsNln"
)

// BuildJWTString creates a token with the HS256 signature algorithm and Claims statements and returns it as a string.
func BuildJWTString() (string, error) {
	userID := generateUniqueID()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})
	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return tokenString, nil
}

// generate Unique ID generate a unique UserID from 0 to 999999
func generateUniqueID() uint32 {
	rand.NewSource(time.Now().UnixNano())
	id := uint32(rand.Intn(1000000))
	logrus.Infof("Generated user id is: %v", id)
	return id
}

// GetUserID we check the validity of the token and if it is valid, then we get and return the UserID from it
func GetUserID(tokenString string) (uint32, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signed method: %v", t.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	if !token.Valid {
		err = fmt.Errorf("token is not valid")
		logrus.Error(err)
		return 0, err
	}
	logrus.Infof("Token is valid, userID: %v", claims.UserID)
	return claims.UserID, nil
}

// IsValidToken method to check the token for validity, we return bool
func IsValidToken(tokenString string) bool {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signed method: %v", t.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		logrus.Error(err)
		return false
	}
	if !token.Valid {
		err = fmt.Errorf("token is not valid")
		logrus.Error(err)
		return false
	}
	return true
}
