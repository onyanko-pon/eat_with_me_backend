package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthUser struct {
	UserID string
}

type JWTClaim struct {
	Sub string `json:"sub"`
	jwt.StandardClaims
}

func (claim JWTClaim) GenAuthUser() (*AuthUser, error) {
	return &AuthUser{
		UserID: claim.Sub,
	}, nil
}

const (
	userIDKey = "sub"
	iatKey    = "iat"
	expKey    = "exp"
	lifetime  = 4 * 365 * 24 * time.Hour
)

func (authUser AuthUser) GenToken() (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		userIDKey: authUser.UserID,
		iatKey:    time.Now().Unix(),
		expKey:    time.Now().Add(lifetime).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SIGNINGKEY")))
}

func ParseToken(signedToken string) (*AuthUser, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			//            return "", err.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SIGNINGKEY")), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("%s is expired", signedToken, err)
			} else {
				return nil, fmt.Errorf("%s is invalid", signedToken, err)
			}
		} else {
			return nil, fmt.Errorf("%s is invalid", signedToken, err)
		}
	}

	if token == nil {
		return nil, fmt.Errorf("not found token in %s:", signedToken)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("not found claims in %s", signedToken)
	}
	userID, ok := claims[userIDKey].(string)
	if !ok {
		return nil, fmt.Errorf("not found %s in %s", userIDKey, signedToken)
	}

	return &AuthUser{
		UserID: userID,
	}, nil
}
