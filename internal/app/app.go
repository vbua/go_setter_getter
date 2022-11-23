package app

import (
	"fmt"
	"github.com/vbua/go_setter_getter/internal/delivery/backup_reader"
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

	backupReader := backup_reader.NewBackupReader(&userGradeService)
	fmt.Println("Заполняю хранилище с бэкапа")
	parseTime := backupReader.ReadFromBackup()
	// только после заполнения стораджа идем дальше

	go func() {
		fmt.Println("Запускаю сеттер")
		setter := httpserver.NewServer(setterHandler.CreateRoutes(), ":"+os.Getenv("SETTER_PORT"))
		setter.Start()
	}()

	go func() {
		fmt.Println("Запускаю геттер")
		getter := httpserver.NewServer(getterHandler.CreateRoutes(), ":"+os.Getenv("GETTER_PORT"))
		getter.Start()
	}()

	go func() {
		fmt.Println("Запускаю подписчика к натс стриминг")
		natsSubscriber := nats.NewNatsSubscriber(&userGradeService)
		natsSubscriber.SubscribeToNats(parseTime)
	}()
	<-wait
}
