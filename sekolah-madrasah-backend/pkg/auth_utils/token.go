package auth_utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"sekolah-madrasah/config"
)

func ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			//nil secret key
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.APP.Rest.JWTSecret), nil
	})
}

func ValidateToken(signedToken string) (claims *AuthClaim, err error) {
	claims = &AuthClaim{}

	parsedData, err := getParsedMapClaims(signedToken)
	if err != nil {
		return
	}

	// Handle UserID
	if userID, ok := parsedData["user_id"].(string); ok {
		claims.UserID, err = uuid.Parse(userID)
		if err != nil {
			return
		}
	}

	// Handle expiration
	if exp, ok := parsedData["exp"].(float64); ok {
		claims.Exp = int64(exp)
		if claims.Exp < time.Now().Local().Unix() {
			err = errors.New("token expired")
			return
		}
	}

	return
}

func getParsedMapClaims(signedToken string) (jwt.MapClaims, error) {
	parsed, err := ParseToken(signedToken)
	if err != nil {
		return nil, err
	}

	parsedData, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || !parsed.Valid {
		return nil, errors.New("couldn't parse claims")
	}
	return parsedData, nil
}

type TokenParams struct {
	UserID uuid.UUID
}

func GenerateToken(params TokenParams, duration time.Duration) (string, error) {
	expiresAt := time.Now().Add(duration).Unix()
	return GenerateTokenWithExpTimestamp(params, expiresAt)
}

func GenerateTokenWithExpTimestamp(params TokenParams, expiresAt int64) (string, error) {
	if params.UserID == uuid.Nil {
		return "", errors.New("user_id is required")
	}

	claims := &AuthClaim{
		UserID: params.UserID,
		Exp:    expiresAt,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.APP.Rest.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
