package main

import (
	"context"
	"gRPC_Example/internal/config"
	"gRPC_Example/internal/repository/mongo"
	"gRPC_Example/internal/service"
	"gRPC_Example/proto"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	port = ":5000"
	
)

func main() {
	ctx := context.Background()

	err := setupViper()

	if err != nil {
		log.Fatalf("error reading uml file: %v", err)
	}

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("error starting tcp listener: %v", err)
	}

	mongoDataBase, err := config.SetupMongoDataBase(ctx)

	if err != nil {
		log.Fatalf("error starting mongo: %v", err)
	}

	userRepository := mongo.NewUserRepository(mongoDataBase.Collection("users"))
	userService := service.NewUserService(userRepository)

	grpcServer := grpc.NewServer()

	proto.RegisterUserServiceServer(grpcServer, userService)

	log.Printf("gRPC stated at %v\n", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error statring gRPC: %v", err)
	}

}

func setupViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
