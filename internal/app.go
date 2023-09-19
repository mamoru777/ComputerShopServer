package internal

import (
	"ComputerShopServer/internal/DataBaseImplement"
	"ComputerShopServer/internal/DataBaseImplement/Config"
	"ComputerShopServer/internal/Repositories/UserRepository"
	"ComputerShopServer/internal/Services/UserService"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg Config.Config) error {
	db, err := DataBaseImplement.InitDB(cfg)
	if err != nil {
		return err
	}
	serv := UserService.New(UserRepository.New(db), cfg)
	s := &http.Server{
		Addr:    ":13999", //"0.0.0.0:%d", cfg.Port),
		Handler: serv.GetHandler(),
	}
	s.SetKeepAlivesEnabled(true)
	ctx, cancel := context.WithCancel(context.Background())
	err = serv.CreateAdmin()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Пользователь админ создан")
	}

	go func() {
		log.Printf("starting http server at %d", cfg.Port)
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}

	}()

	gracefullyShutdown(ctx, cancel, s)

	return nil

}

func gracefullyShutdown(ctx context.Context, cancel context.CancelFunc, server *http.Server) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)
	<-ch
	if err := server.Shutdown(ctx); err != nil {
		log.Print(err)
	}
	cancel()
}
