package nats_repo

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/vbua/go_setter_getter/config"
	"github.com/vbua/go_setter_getter/internal/entity"
	"log"
	"os"
)

type UserGradeNatsRepo struct {
	sc stan.Conn
}

func NewGradeNatsRepo() UserGradeNatsRepo {
	sc, err := stan.Connect(config.NatsClusterId, config.NatsPublisherClientId)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return UserGradeNatsRepo{sc}
}

func (u *UserGradeNatsRepo) Publish(userGrade entity.UserGrade) error {
	data := struct {
		ReplicaId string
		UserGrade entity.UserGrade
	}{os.Getenv("REPLICA_ID"), userGrade}
	dataMarshaled, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := u.sc.Publish(config.NatsTopic, dataMarshaled); err != nil {
		return err
	}
	return nil
}
