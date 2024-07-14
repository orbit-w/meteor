package subpub_redis

import (
	"github.com/orbit-w/meteor/modules/database/rdb"
	"log"
	"testing"
	"time"
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
	time.Sleep(time.Second * 5)
	log.Println("count: ", count)
	log.Println("count1: ", count1)
}
