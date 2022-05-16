package service

import (
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

type ProtoRegistortService struct {
	ProtoDescriptorMap map[string]descriptorpb.FileDescriptorProto
}

func NewProtoRegistortService(filePtah string) ProtoRegistortService {
	protoRegistoryService := ProtoRegistortService{ProtoDescriptorMap: make(map[string]descriptorpb.FileDescriptorProto)}
	protoRegistoryService.registerProtoFile(filePtah)
	return protoRegistoryService
}

func (p *ProtoRegistortService) registerProtoFile(filePtah string) error {

	// Now load that temporary file as a file descriptor set protobuf.
	protoFile, err := ioutil.ReadFile(filePtah)
	if err != nil {
		return err
	}

	pbSet := new(descriptorpb.FileDescriptorSet)
	if err := proto.Unmarshal(protoFile, pbSet); err != nil {
		return err
	}

	// We know protoc was invoked with a multiple .proto file.
	for _, pb := range pbSet.GetFile() {
		protoName := pb.GetName()
		_, ok := p.ProtoDescriptorMap[protoName]

		if ok || strings.Contains(pb.GetName(), "google/protobuf/descriptor.proto") || isAlreadyRegistered(protoName) {
			log.Printf("Ignore ..  --> %v", protoName)
			continue
		}

		log.Printf("loading .. file descriptor --> %v", protoName)
		p.ProtoDescriptorMap[protoName] = *pb
		// Initialize the file descriptor object.
		fd, err := protodesc.NewFile(pb, protoregistry.GlobalFiles)
		if err != nil {
			return err
		}

		// and finally register it.
		err = protoregistry.GlobalFiles.RegisterFile(fd)
		if err != nil {
			return err
		}
	}
	return nil
}

func isAlreadyRegistered(messageName string) bool {

	_, err := protoregistry.GlobalFiles.FindFileByPath(messageName)

	if err == nil {
		log.Printf(" already registered %v \n", messageName)
		return true
	}
	return false
}

func (p *ProtoRegistortService) ToJson(messageFile string, messageName string, messageBytes []byte) string {

	d, _ := protoregistry.GlobalFiles.FindFileByPath(messageFile)
	_, ok := p.ProtoDescriptorMap[messageFile]
	if !ok {
		return string(messageBytes)
	}
	mds := d.Messages()
	md := mds.ByName(protoreflect.Name(messageName))
	m := dynamicpb.NewMessage(md)

	proto.Unmarshal(messageBytes, m)
	marshal := protojson.MarshalOptions{Multiline: true}
	bytes, _ := marshal.Marshal(m)

	return string(bytes)
}
