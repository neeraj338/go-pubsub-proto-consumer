package service

import (
	"context"

	log "github.com/sirupsen/logrus"

	"cloud.google.com/go/pubsub"
	"github.com/go-pubsub-proto-consumer/config"
	"google.golang.org/api/option"
)

type PubSub struct {
	GcpProjectName        string
	Subscription          string
	Config                config.Configuration
	ProtoFile             string
	ProtoMessageName      string
	ProtoRegistortService ProtoRegistortService
}

func NewPubSub(gcpProject string,
	subscription string,
	protoFile string,
	protoMessageName string,
	config config.Configuration,
	protoRegistortService ProtoRegistortService) PubSub {

	return PubSub{GcpProjectName: gcpProject,
		Subscription:          subscription,
		ProtoFile:             protoFile,
		ProtoMessageName:      protoMessageName,
		Config:                config,
		ProtoRegistortService: protoRegistortService}
}

func (ps *PubSub) PullMsgs(ctx context.Context) error {

	client, err := pubsub.NewClient(ctx, ps.GcpProjectName, option.WithCredentialsJSON([]byte(ps.Config.CredentialJson)))
	if err != nil {
		log.Errorf("failed to create pub/sub client %v", err)
		return err
	}
	defer client.Close()

	sub := client.Subscription(ps.Subscription)
	// Receive blocks until the context is cancelled or an error occurs.
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf("Got message: \n %v\n", ps.ProtoRegistortService.ToJson(ps.ProtoFile, ps.ProtoMessageName, msg.Data))
		log.Printf("Attributes:")
		for key, value := range msg.Attributes {
			log.Printf("%s = %s\n", key, value)
		}
		msg.Ack()
	})
	if err != nil {
		log.Errorf("sub.Receive: %v", err)
	}

	return nil
}
