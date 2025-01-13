package extensions

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
)

func ConvertToDynamicMessage(m proto.Message) (*dynamicpb.Message, error) {
	dynamicRule := dynamicpb.NewMessage(m.ProtoReflect().Descriptor())
	b, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(b, dynamicRule)
	if err != nil {
		return nil, err
	}
	return dynamicRule, nil
}

func ConvertIntoProtoMessage(payload interface{}, msg proto.Message) error {
	payload_marshalled, err := proto.Marshal(payload.(*dynamicpb.Message))
	if err != nil {
		return err
	}
	return proto.Unmarshal(payload_marshalled, msg)
}
