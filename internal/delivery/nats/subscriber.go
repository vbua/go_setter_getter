package nats

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/vbua/go_setter_getter/config"
	"github.com/vbua/go_setter_getter/internal/entity"
	"log"
	"os"
)

type UserGradeSetterService interface {
	Set(entity.UserGrade, bool)
}

type NatsSubscriber struct {
	userGradeService UserGradeSetterService
}

func NewNatsSubscriber(userGradeService UserGradeSetterService) NatsSubscriber {
	return NatsSubscriber{userGradeService}
}

func (n *NatsSubscriber) SubscribeToNats() {
	sc, err := stan.Connect(config.NatsClusterId, config.NatsSubscriberClientId)
	if err != nil {
		log.Fatalln("Couldn't connect: ", err.Error())
	}

	//var startTime time.Time
	_, err = sc.Subscribe(config.NatsTopic, func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))

		data := struct {
			ReplicaId string
			UserGrade entity.UserGrade
		}{}
		err := json.Unmarshal(m.Data, &data)
		if err != nil {
			return
		}
		// фильтруем свои сообщения с той же реплики, иначе попадем в зацикливание
		if data.ReplicaId != os.Getenv("REPLICA_ID") {
			n.userGradeService.Set(data.UserGrade, false)
		}
	})

	if err != nil {
		log.Fatalln("Couldn't subscribe: ", err.Error())
	}
}
