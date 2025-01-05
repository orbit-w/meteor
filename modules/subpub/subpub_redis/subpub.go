package subpub_redis

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"sync/atomic"

	"github.com/gogo/protobuf/proto"
	"github.com/orbit-w/meteor/modules/mlog"
	"github.com/redis/go-redis/v9"
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
	log         *mlog.Logger
	invoker     func(pid int32, body []byte, err error)
}

var (
	ctx = context.Background()
)

type IEncoder interface {
	Marshal(v interface{}) ([]byte, error)
}

func NewPubSub(_cli redis.UniversalClient, ee int, topic string, _invoker func(pid int32, body []byte, err error)) IPubSub {
	return &PubSub{
		topic:       topic,
		invoker:     _invoker,
		cli:         _cli,
		encoderEnum: ee,
		log:         mlog.WithPrefix("subpub_redis"),
	}
}

func (ps *PubSub) Publish(pid int32, v any) error {
	body, err := encode(ps.encoderEnum, pid, v)
	if err != nil {
		return err
	}
	err = ps.cli.Publish(ctx, ps.topic, body).Err()
	return ErrPublish(err)
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
	var (
		err error
		pb  = new(PubMessage)
	)

	if err = proto.Unmarshal([]byte(msg.Payload), pb); err != nil {
		err = fmt.Errorf("[PubSub] decode proto failed : %w", err)
	}

	ps.invoke(pb.Pid, pb.Data, err)
	return
}

func (ps *PubSub) invoke(pid int32, data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			log.Println("Stack: ", string(debug.Stack()))
		}
	}()
	ps.invoker(pid, data, err)
}
