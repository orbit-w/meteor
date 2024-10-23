package subpub_redis

import (
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/orbit-w/meteor/modules/mlog"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"runtime/debug"
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
	log         *mlog.ZapLogger
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
		log:         mlog.NewLogger("subpub_redis"),
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
	pb := new(PubMessage)
	if err := proto.Unmarshal([]byte(msg.Payload), pb); err != nil {
		ps.log.Error("[PubSub] decode proto failed: ", zap.Error(err))
		return
	}

	ps.handleMessage(pb)
	return
}

func (ps *PubSub) handleMessage(msg *PubMessage) {
	defer func() {
		if r := recover(); r != nil {
			//ps.log.Error("panic", zap.Any("recover", r), zap.String("stack", string(debug.Stack())))
			log.Println(r)
			log.Println("Stack: ", string(debug.Stack()))
		}
	}()
	ps.invoker(msg.Pid, msg.Data)
}
