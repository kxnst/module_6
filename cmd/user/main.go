package main

import (
	"fmt"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
	"guitar_processor/internal"
	"guitar_processor/internal/entity"
	"guitar_processor/internal/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := fx.New(
		fx.Provide(internal.GetProviders()...),
		fx.Invoke(func(repo *repository.UserRepository) {
			in := os.Args
			if len(in) < 4 {
				log.Fatalf("Expected 3 arguments, received %d", len(in))
			}

			hashedPass, err := bcrypt.GenerateFromPassword([]byte(os.Args[3]), 14)
			if err != nil {
				log.Fatal(err)
			}

			user := entity.User{Name: os.Args[1], Login: os.Args[2], Password: string(hashedPass)}

			err = repo.Put(&user)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("User %s successfully created", user.Name)
		}),
	)

	app.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
