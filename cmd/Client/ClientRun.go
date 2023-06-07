package main

import (
	userapi "ComputerShopServer/pkg"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:13999", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connect to grpc server err:", err)
	}
	defer conn.Close()

	client := userapi.NewUserServiceClient(conn)

	if err := createNotification(client); err != nil {
		log.Fatal(err)
	}
}

func createNotification(client userapi.UserServiceClient) error {
	notification := userapi.CreateUserRequest{
		Login:    "test1",
		Password: "test1",
		Email:    "test1@mail.ru",
	}
	if _, err := client.CreateUser(context.Background(), &notification); err != nil {
		return err
	}
	log.Println("User created: ", notification)
	return nil
}
