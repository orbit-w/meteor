package transport

import (
	"sync/atomic"
)

type IMonitor interface {
	Fmt()
}

type Monitor struct {
	InboundTraffic      uint64
	OutboundTraffic     atomic.Uint64
	RealInboundTraffic  uint64
	RealOutboundTraffic uint64
}

func NewMonitor() *Monitor {
	return &Monitor{}
}

func (m *Monitor) IncrementInboundTraffic(amount uint64) {
	m.incrementTraffic(amount, &m.InboundTraffic)
}

func (m *Monitor) IncrementOutboundTraffic(amount uint64) {
	m.OutboundTraffic.Add(amount)
}

func (m *Monitor) GetOutboundTraffic() uint64 {
	return m.OutboundTraffic.Load()
}

func (m *Monitor) IncrementRealInboundTraffic(amount uint64) {
	m.incrementTraffic(amount, &m.RealInboundTraffic)
}

func (m *Monitor) IncrementRealOutboundTraffic(amount uint64) {
	m.incrementTraffic(amount, &m.RealOutboundTraffic)
}

func (m *Monitor) incrementTraffic(amount uint64, outboundTraffic *uint64) {
	if *outboundTraffic+amount < *outboundTraffic { // Check for overflow
		*outboundTraffic = ^uint64(0) // Set to max value if overflow
	} else {
		*outboundTraffic += amount
	}
}
