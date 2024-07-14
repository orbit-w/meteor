package subpub_redis

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
)

var (
	jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary
)

const (
	CodecJson = iota
	CodecProto3
)

func encode(name int, pid int32, v any) ([]byte, error) {
	var (
		err  error
		body []byte
	)

	switch name {
	case CodecJson:
		body, err = jsonAPI.Marshal(v)
		if err != nil {
			return nil, err
		}
	case CodecProto3:
		pbMsg, ok := v.(proto.Message)
		if !ok {
			return nil, errors.New("")
		}
		body, err = proto.Marshal(pbMsg)
	}

	msg := &PubMessage{
		Pid:  pid,
		Data: body,
	}
	return proto.Marshal(msg)
}
