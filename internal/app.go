package internal

import (
	"ComputerShopServer/internal/ConfigServ"
	"ComputerShopServer/internal/DataBaseImplement"
	"ComputerShopServer/internal/DataBaseImplement/Config"
	"ComputerShopServer/internal/Repositories/UserRepository"
	"ComputerShopServer/internal/Services/UserService"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg Config.Config, config ConfigServ.Config) error {
	db, err := DataBaseImplement.InitDB(cfg)
	if err != nil {
		return err
	}
	serv := UserService.New(UserRepository.New(db))
	s := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: UserService,
	}

	l, err := net.Listen("tcp", ":13999") //config.GRPCAddr)
	if err != nil {
		return fmt.Errorf("failed to listen tcp %s, %v", config.GRPCAddr, err)
	}

	go func() {
		log.Printf("starting listening grpc server at %s", "13999") //config.GRPCAddr)
		if err := s.Serve(l); err != nil {
			log.Fatalf("error service grpc server %v", err)
		}
	}()

	gracefulShutDown(s)

	return nil

}

func gracefulShutDown(s *grpc.Server) {
	const waitTime = 5 * time.Second

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	sig := <-ch
	errorMessage := fmt.Sprintf("%s %v - %s", "Received shutdown signal:", sig, "Graceful shutdown done")
	log.Println(errorMessage)
	s.GracefulStop()

}
