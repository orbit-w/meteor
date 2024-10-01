package transport

import (
	"context"
	"errors"
	"io"
	"log"
	"sync"
	"sync/atomic"
)

func ServePlannedTraffic(protocol, host string, plan int64) (s IServer, wait func(), err error) {
	p := atomic.Int64{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	handle := func(conn IConn) {
		ctx := context.Background()
		for {
			in, err := conn.Recv(ctx)
			if err != nil {
				if IsClosedConnError(err) || IsCancelError(err) || errors.Is(err, io.EOF) {
					break
				}

				log.Println("conn read failed: ", err.Error())
				break
			}
			if p.Add(int64(len(in))) >= plan {
				wg.Done()
				return
			}
		}
	}
	config := DefaultServerConfig()
	config.Stage = DEV
	s, err = ServeByConfig(protocol, host, handle, config)
	if err != nil {
		return nil, nil, err
	}
	wait = func() {
		wg.Wait()
	}
	return
}
