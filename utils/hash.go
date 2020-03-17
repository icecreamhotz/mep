package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/icecreamhotz/mep-api/configs"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

var AccessKey = []byte(viper.GetString("access_token_secret"))
var RefreshKey = []byte(viper.GetString("refresh_token_secret"))

const AccessTokenMinute = 24 * 60
const RefreshTokenMinute = 24 * 7 * 60

type AccessTokenClaims struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	ID string
	jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(claim *AccessTokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(AccessKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(claim *RefreshTokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(RefreshKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func SetAccessTokenAndRefreshToken(id, name, lastname, role string) (string, string, time.Time, error) {
	expirationTimeAccessToken := time.Now().Add(time.Minute * time.Duration(AccessTokenMinute))
	expirationTimeRefreshToken := time.Now().Add(time.Minute * time.Duration(RefreshTokenMinute))

	accessTokenClaim := &AccessTokenClaims{
		ID:       id,
		Name:     name,
		Lastname: lastname,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeAccessToken.Unix(),
		},
	}

	accessToken, err := GenerateToken(accessTokenClaim)
	if err != nil {
		return "", "", time.Time{}, err
	}

	refreshTokenClaims := &RefreshTokenClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeRefreshToken.Unix(),
		},
	}

	refreshToken, err := GenerateRefreshToken(refreshTokenClaims)
	if err != nil {
		return "", "", time.Time{}, err
	}

	client := configs.ConnectCacheDatabase()

	accessTokenKey := "access_token:" + accessToken
	err = client.HMSet(accessTokenKey, map[string]interface{}{
		"user_id": id,
		"role":    role,
	}).Err()
	if err != nil {
		return "", "", time.Time{}, err
	}

	// set expire token redis
	err = client.ExpireAt(accessTokenKey, expirationTimeAccessToken).Err()
	if err != nil {
		return "", "", time.Time{}, err
	}

	refreshTokenKey := "refresh_token:" + refreshToken
	err = client.HMSet(refreshTokenKey, map[string]interface{}{
		"user_id": id,
		"role":    role,
	}).Err()
	if err != nil {
		return "", "", time.Time{}, err
	}

	// set expire token redis
	err = client.ExpireAt(refreshTokenKey, expirationTimeRefreshToken).Err()
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken, refreshToken, expirationTimeAccessToken, nil
}

func GetUserPayload(token string) (*AccessTokenClaims, *jwt.Token, error) {
	claims := &AccessTokenClaims{}

	parseToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return AccessKey, nil
	})

	if err != nil || !parseToken.Valid {
		return nil, parseToken, err
	}

	return claims, parseToken, nil
}
