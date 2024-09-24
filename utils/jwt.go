package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// setup the jwt secret key / access token expiry / refresh token expiry
var (
	jwtSecret          = []byte(os.Getenv("SECRET_KEY"))
	accessTokenExpiry  = time.Hour * 30
	refreshTokenExpiry = time.Hour * 24 * 7
)

// create a struct to read the JWT claims from
type Claims struct {
	UserID uint
	jwt.StandardClaims
}

// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-New-Hmac
// create a function to generate JWT tokens
func GenerateTokens(userID uint) (string, string, error) {
	// create access token
	// 1. set the expiry time and the user ID (payload)
	accessTokenClaims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenExpiry).Unix(),
		},
	}
	// 2. create the JWT token with the header (HS256)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	// 3. sign and get the complete encoded token using the secret key (Signature)
	signedAccessToken, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}
	// create refresh token
	// 1. set the expiry time and the user ID (payload)
	refreshTokenClaims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshTokenExpiry).Unix(),
		},
	}
	// 2. create the JWT token with the header (HS256)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	// 3. sign and get the complete encoded token using the secret key (Signature)
	signedRefreshToken, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-Parse-Hmac
func ParseToken(tokenString string) (*Claims, error) {
	// parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	// check the token
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, jwt.NewValidationError("expired token", jwt.ValidationErrorExpired)
			}
			if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, jwt.NewValidationError("token not active yet", jwt.ValidationErrorNotValidYet)
			}
			return nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorMalformed)
		}
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorMalformed)
}
