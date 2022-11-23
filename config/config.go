package config

import "os"

var (
	NatsClusterId          = "user_grades_cluster"
	NatsSubscriberClientId = "user_grades_subscriber_client" + os.Getenv("REPLICA_ID")
	NatsPublisherClientId  = "user_grades_publisher_client" + os.Getenv("REPLICA_ID")
	NatsTopic              = "user_grades"
	UserName               = "test"
	UserPass               = "test"
)
