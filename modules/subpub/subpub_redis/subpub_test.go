package subpub_redis

import (
	"log"
	"testing"
	"time"

	"github.com/orbit-w/meteor/modules/database/rdb"
	"github.com/orbit-w/meteor/modules/mlog"
)

type Config struct {
	ID    int64  `bson:"id" json:"id,omitempty"` // 用时间戳作为ID
	AppID string `bson:"app_id" json:"app_id"`
}

func TestPubSub_Publish(t *testing.T) {
	if err := rdb.NewClient(rdb.RedisClientOps{
		Username: "root",
		Addr:     []string{"127.0.0.1:6379"},
		Cluster:  false,
	}); err != nil {
		panic(err)
	}

	var (
		count  = 0
		count1 = 0
	)

	ps := NewPubSub(rdb.UniversalClient(), CodecJson, "test", func(pid int32, body []byte) {
		count++
	})

	ps2 := NewPubSub(rdb.UniversalClient(), CodecJson, "test", func(pid int32, body []byte) {
		count1++
	})

	ps.Subscribe()
	ps2.Subscribe()

	c := &Config{
		AppID: "abc_reading",
	}
	time.Sleep(time.Second)

	for i := int64(0); i < 10000; i++ {
		c.ID = i
		if err := ps.Publish(1002, c); err != nil {
			log.Fatalln("decode config failed: ", err.Error())
		}
	}
	time.Sleep(time.Second * 30)
	log.Println("count: ", count)
	log.Println("count1: ", count1)
}

func TestPubSub_Close(t *testing.T) {
	if err := rdb.NewClient(rdb.RedisClientOps{
		Username: "root",
		Addr:     []string{"127.0.0.1:6379"},
		Cluster:  false,
	}); err != nil {
		panic(err)
	}

	var (
		count  = 0
		count1 = 0
	)

	var (
		h1 = func(pid int32, body []byte) {
			count++
		}

		h2 = func(pid int32, body []byte) {
			count1++
		}
	)

	ps := NewPubSub(rdb.UniversalClient(), CodecJson, "test", h1)

	ps2 := NewPubSub(rdb.UniversalClient(), CodecJson, "test", h2)

	ps.Subscribe()
	ps2.Subscribe()

	c := &Config{
		AppID: "abc_reading",
	}

	go func() {
		for i := int64(0); i < 10000; i++ {
			c.ID = i
			if err := ps.Publish(1002, c); err != nil {
				log.Fatalln("decode config failed: ", err.Error())
			}
		}
	}()

	ps2.Stop()
	time.Sleep(time.Second * 10)
	log.Println("count: ", count)
	log.Println("count1: ", count1)
}

func TestPubSub_Create(t *testing.T) {
	if err := rdb.NewClient(rdb.RedisClientOps{
		Username: "root",
		Addr:     []string{"127.0.0.1:6379"},
		Cluster:  false,
	}); err != nil {
		panic(err)
	}

	psList := make([]IPubSub, 0)

	defer func() {
		for i := range psList {
			ps := psList[i]
			ps.Stop()
		}
	}()

	for i := 0; i < 500; i++ {
		ps := NewPubSub(rdb.UniversalClient(), CodecJson, "test", func(pid int32, body []byte) {

		})
		ps.Subscribe()
		psList = append(psList, ps)
	}

	time.Sleep(time.Second * 30)
}

func TestPubSub_Recovery(t *testing.T) {
	ps := new(PubSub)
	var s []byte
	ps.log = mlog.WithPrefix("subpub_redis")
	ps.invoker = func(pid int32, body []byte) {
		s[2] = 1
	}
	ps.handleMessage(&PubMessage{
		Pid:  10001,
		Data: []byte("hello"),
	})
	time.Sleep(time.Second * 5)
}
