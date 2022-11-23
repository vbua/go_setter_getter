package nats

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/vbua/go_setter_getter/config"
	"github.com/vbua/go_setter_getter/internal/entity"
	"log"
	"os"
	"time"
)

type UserGradeSetterService interface {
	Set(entity.UserGrade, bool)
}

type NatsSubscriber struct {
	sc               stan.Conn
	userGradeService UserGradeSetterService
	sub              stan.Subscription
}

func NewNatsSubscriber(userGradeService UserGradeSetterService) NatsSubscriber {
	sc, err := stan.Connect(config.NatsClusterId, config.NatsSubscriberClientId)
	if err != nil {
		log.Fatalln("Couldn't connect: ", err.Error())
	}
	return NatsSubscriber{sc, userGradeService, nil}
}

func (n *NatsSubscriber) SubscribeToNats(parseTime time.Time) {
	if parseTime.IsZero() {
		parseTime = time.Now()
	}
	sub, err := n.sc.Subscribe(config.NatsTopic, func(m *stan.Msg) {
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
	}, stan.StartAtTime(parseTime))

	if err != nil {
		log.Fatalln("Couldn't subscribe: ", err.Error())
	}

	n.sub = sub
}

func (n *NatsSubscriber) Stop() {
	n.sub.Unsubscribe()
	n.sc.Close()
}
