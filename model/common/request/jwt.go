package request

import (
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	BaseClaims
	jwt.RegisteredClaims
}

type BaseClaims struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}
