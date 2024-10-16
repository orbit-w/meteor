package subpub_redis

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/redis/go-redis/v9"
	"log"
	"sync/atomic"
)

type IPubSub interface {
	Publish(pid int32, v any) error
	Subscribe()
	Stop()
}

type PubSub struct {
	state       atomic.Uint32
	encoderEnum int    //JSON｜Proto3: 默认编码协议是JSON
	topic       string //主题名称
	cli         redis.UniversalClient
	sub         *redis.PubSub
	invoker     func(pid int32, body []byte)
}

var (
	ctx = context.Background()
)

type IEncoder interface {
	Marshal(v interface{}) ([]byte, error)
}

func NewPubSub(_cli redis.UniversalClient, ee int, topic string, _invoker func(pid int32, body []byte)) IPubSub {
	return &PubSub{
		topic:       topic,
		invoker:     _invoker,
		cli:         _cli,
		encoderEnum: ee,
	}
}

func (ps *PubSub) Publish(pid int32, v any) error {
	body, err := encode(ps.encoderEnum, pid, v)
	if err != nil {
		return err
	}
	err = ps.cli.Publish(ctx, ps.topic, body).Err()
	if err != nil {
		log.Println("[Publish] failed: ", err.Error())
	}
	return err
}

func (ps *PubSub) Subscribe() {
	ps.subscribe(ps.decodeAndInvoke)
}

func (ps *PubSub) Stop() {
	if ps.state.CompareAndSwap(stateReady, stateStopped) {
		if ps.sub != nil {
			_ = ps.sub.Close()
		}
	}
}

func (ps *PubSub) subscribe(handle func(msg *redis.Message)) {
	pubSub := ps.cli.Subscribe(ctx, ps.topic)
	ps.sub = pubSub
	ch := pubSub.Channel()

	go func() {
		defer func() {
			fmt.Println("close sub, topic: ", ps.topic)
			if ps.sub != nil {
				_ = ps.sub.Close()
			}
		}()
		for msg := range ch {
			handle(msg)
		}
	}()
}

func (ps *PubSub) decodeAndInvoke(msg *redis.Message) {
	pb := new(PubMessage)
	if err := proto.Unmarshal([]byte(msg.Payload), pb); err != nil {
		log.Println("Decode proto failed: ", err.Error())
		return
	}

	ps.handleMessage(pb)
	return
}

func (ps *PubSub) handleMessage(msg *PubMessage) {
	defer func() {
		if r := recover(); r != nil {

		}
	}()
	ps.invoker(msg.Pid, msg.Data)
}
