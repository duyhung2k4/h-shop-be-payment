package config

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func GetDB() *gorm.DB {
	return db
}

func GetRDB() *redis.Client {
	return rdb
}

func GetAppPort() string {
	return appPort
}

func GetJWT() *jwtauth.JWTAuth {
	return jwt
}

func GetHost() string {
	return host
}

func GetClientGRPCShop() *grpc.ClientConn {
	return clientShopGRPC
}

func GetVnpTmncode() string {
	return vnpTmncode
}

func GetVnpHashsecret() string {
	return vnpHashsecret
}

func GetVnpUrl() string {
	return vnpUrl
}

func GetVnpReturnUrl() string {
	return vnpReturnUrl
}
