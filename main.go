package main

import (
	"context"

	"github.com/alexflint/go-arg"
	"github.com/go-pubsub-proto-consumer/config"
	"github.com/go-pubsub-proto-consumer/service"
)

//--table students -k="name" -v="neeraj" -select "id,name,age --data="{}" -delete"
type Args struct {
	Project      string `arg:"required, -p,--project" help:"gcp project name"` //`arg:"required, -p,--project" help:"gcp project name" default:"findaya-staging"`
	Topic        string `arg:"required" help:"--topic topic name"`
	Subscription string `arg:"required, -s,--subscription" help:"subscription name"`
	MessageFile  string `arg:"required, -f,--message-file" help:"proto message file name"`
	MessageName  string `arg:"required, -m,--message-name" help:"proto message name"`
}

func main() {
	config, _ := config.ReadConfiguration()
	var args Args
	arg.MustParse(&args)
	protoRegistryService := service.NewProtoRegistortService(config.ProtoDescFilePath)

	pubSub := service.NewPubSub(args.Project,
		args.Subscription,
		args.MessageFile,
		args.MessageName,
		config,
		protoRegistryService)
	pubSub.PullMsgs(context.Background())
}
