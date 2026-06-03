package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Id       int
	Username string
	Password string
	jwt.RegisteredClaims
}
