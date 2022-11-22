package app

import (
	"github.com/vbua/go_setter_getter/internal/delivery/httpserver"
	"github.com/vbua/go_setter_getter/internal/delivery/nats"
	"github.com/vbua/go_setter_getter/internal/respository/nats_repo"
	"github.com/vbua/go_setter_getter/internal/respository/storage"
	"github.com/vbua/go_setter_getter/internal/service"
	"os"
)

func Run() {
	wait := make(chan string)
	userGradeRepo := storage.NewUserGradeRepo()
	userGradeNatsRepo := nats_repo.NewGradeNatsRepo()
	userGradeService := service.NewUserGradeService(&userGradeRepo, &userGradeNatsRepo)
	setterHandler := httpserver.NewSetterHandler(&userGradeService)
	getterHandler := httpserver.NewGetterHandler(&userGradeService)

	go func() {
		setter := httpserver.NewServer(setterHandler.CreateRoutes(), ":"+os.Getenv("SETTER_PORT"))
		setter.Start()
	}()

	go func() {
		getter := httpserver.NewServer(getterHandler.CreateRoutes(), ":"+os.Getenv("GETTER_PORT"))
		getter.Start()
	}()

	natsSubscriber := nats.NewNatsSubscriber(&userGradeService)
	go func() {
		natsSubscriber.SubscribeToNats()
	}()
	<-wait
}
