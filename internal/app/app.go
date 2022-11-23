package app

import (
	"context"
	"fmt"
	"github.com/vbua/go_setter_getter/internal/delivery/backup_reader"
	"github.com/vbua/go_setter_getter/internal/delivery/httpserver"
	"github.com/vbua/go_setter_getter/internal/delivery/nats"
	"github.com/vbua/go_setter_getter/internal/respository/nats_repo"
	"github.com/vbua/go_setter_getter/internal/respository/storage"
	"github.com/vbua/go_setter_getter/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	//wait := make(chan string)
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userGradeRepo := storage.NewUserGradeRepo()
	userGradeNatsRepo := nats_repo.NewGradeNatsRepo()
	userGradeService := service.NewUserGradeService(&userGradeRepo, &userGradeNatsRepo)
	setterHandler := httpserver.NewSetterHandler(&userGradeService)
	getterHandler := httpserver.NewGetterHandler(&userGradeService)

	backupReader := backup_reader.NewBackupReader(&userGradeService)
	fmt.Println("Заполняю хранилище с бэкапа")
	parseTime := backupReader.ReadFromBackup()
	// только после заполнения стораджа идем дальше

	setter := httpserver.NewServer(setterHandler.CreateRoutes(), ":"+os.Getenv("SETTER_PORT"))
	go func() {
		fmt.Println("Запускаю сеттер")
		setter.Start()
	}()

	getter := httpserver.NewServer(getterHandler.CreateRoutes(), ":"+os.Getenv("GETTER_PORT"))
	go func() {
		fmt.Println("Запускаю геттер")
		getter.Start()
	}()

	natsSubscriber := nats.NewNatsSubscriber(&userGradeService)
	go func() {
		fmt.Println("Запускаю подписчика к натс стриминг")
		natsSubscriber.SubscribeToNats(parseTime)
	}()
	<-stopCh
	natsSubscriber.Stop()
	getter.Stop(ctx)
	setter.Stop(ctx)
}
