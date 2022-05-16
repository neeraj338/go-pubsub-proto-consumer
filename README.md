# Protobuf golang command line pubsub Consumer

## Configuration is read from following environment variables

| Variable               | Description                   | Default Value         | Possible Values |
|------------------------|-------------------------------|-----------------------|-----------------|
| PUBSUB_EMULATOR_HOST   | PUBSUB_EMULATOR_HOST          | localhost:9092        | host post       |
| GOOGLE_CREDENTIAL_JSON | base64  encoded json          | {}                    | base64({})      |
| PROTO_DESC_PATH        | file path for descriptor.desc | ./proto/descriptor.pb | filepath        |

_Note: if PUBSUB_EMULATOR_HOST is set, GOOGLE_CREDENTIAL_JSON will be ignored_

## How to create proto descriptor file

```shell
protoc --include_imports \
--descriptor_set_out=descriptor.desc A.proto B.proto
```

## Build and run

1. build `go build`

2. run

```shell
go run main.go \
-p my-project-id \
--topic pubsub-topic \
--subscription subscription \
--message-file Message.proto \
--message-name MessageName
```

or

```shell
./go-pubsub-proto-consumer \
-p my-project-id \
--topic pubsub-topic \
--subscription subscription \
--message-file Message.proto \
--message-name MessageName
```

## help

```shell
go run main.go -h
```

# Run PubSub emulator

```shell
docker run --rm -it -p 8681:8681 \
-e PUBSUB_PROJECT1=my-project-name,topic-name:subscription \
messagebird/gcloud-pubsub-emulator:latest
```

```text
projectId: my-project-name
topicId: topic-name
subscriptionId: subscription
```

## Publish

```text
curl -d '{"messages": [{"data": "c3Vwc3VwCg=="}]}' \
-H "Content-Type: application/json" \
-X POST localhost:8681/v1/projects/my-project/topics/topic:publish
```

## Pull Message

```Text
curl -d '{"returnImmediately":true, "maxMessages":1}' \
-H "Content-Type: application/json" \
-X POST localhost:8681/v1/projects/my-project/subscriptions/my-subscription:pull
```

## References

1. [gcloud-pubsub-emulator](https://github.com/marcelcorso/gcloud-pubsub-emulator)
