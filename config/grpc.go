package config

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func connectGPRCServerShop() {
	var err error

	creds, errKey := credentials.NewClientTLSFromFile("keys/server-shop/public.pem", "localhost")
	if errKey != nil {
		log.Fatalln(errKey)
	}

	clientShopGRPC, err = grpc.Dial(host+":20002", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalln(err)
	}
}
