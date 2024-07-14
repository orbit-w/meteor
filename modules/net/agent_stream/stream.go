package agent_stream

type IStream interface {
	Send(data []byte) error
	Recv() ([]byte, error)
	Close() error
}
