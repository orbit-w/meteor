package transport

import (
	"go.uber.org/zap"
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
	if m == nil {
		return
	}
	m.incrementTraffic(amount, &m.InboundTraffic)
}

func (m *Monitor) IncrementOutboundTraffic(amount uint64) {
	if m == nil {
		return
	}
	m.OutboundTraffic.Add(amount)
}

func (m *Monitor) GetOutboundTraffic() uint64 {
	if m == nil {
		return 0
	}
	return m.OutboundTraffic.Load()
}

func (m *Monitor) IncrementRealInboundTraffic(amount uint64) {
	if m == nil {
		return
	}
	m.incrementTraffic(amount, &m.RealInboundTraffic)
}

func (m *Monitor) IncrementRealOutboundTraffic(amount uint64) {
	if m == nil {
		return
	}
	m.incrementTraffic(amount, &m.RealOutboundTraffic)
}

func (m *Monitor) Log() []zap.Field {
	if m == nil {
		return nil
	}
	return []zap.Field{
		zap.Uint64("InboundTraffic", m.InboundTraffic), zap.Uint64("OutboundTraffic", m.GetOutboundTraffic()),
		zap.Uint64("RealOutboundTraffic", m.RealOutboundTraffic), zap.Uint64("RealInboundTraffic", m.RealInboundTraffic),
	}
}

func (m *Monitor) incrementTraffic(amount uint64, outboundTraffic *uint64) {
	if *outboundTraffic+amount < *outboundTraffic { // Check for overflow
		*outboundTraffic = ^uint64(0) // Set to max value if overflow
	} else {
		*outboundTraffic += amount
	}
}
