package types

import (
	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Id         uint32     `json:"id"`
	Permission Permission `json:"permission"`
	jwt.StandardClaims
}

type Permission uint8

const (
	Client Permission = iota + 1
	Admin
)
