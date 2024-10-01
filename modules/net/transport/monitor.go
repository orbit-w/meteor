package transport

import "fmt"

type IMonitor interface {
	Fmt()
}

type Monitor struct {
	InboundTraffic      uint64
	OutboundTraffic     uint64
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
	m.incrementTraffic(amount, &m.OutboundTraffic)
}

func (m *Monitor) IncrementRealInboundTraffic(amount uint64) {
	m.incrementTraffic(amount, &m.RealInboundTraffic)
}

func (m *Monitor) IncrementRealOutboundTraffic(amount uint64) {
	m.incrementTraffic(amount, &m.RealOutboundTraffic)
}

func (m *Monitor) Fmt() {
	fmt.Printf("InboundTraffic: %d\n, OutboundTraffic: %d\n, RealInboundTraffic: %d\n, RealOutboundTraffic: %d\n", m.InboundTraffic, m.OutboundTraffic, m.RealInboundTraffic, m.RealOutboundTraffic)
}

func (m *Monitor) incrementTraffic(amount uint64, outboundTraffic *uint64) {
	if *outboundTraffic+amount < *outboundTraffic { // Check for overflow
		*outboundTraffic = ^uint64(0) // Set to max value if overflow
	} else {
		*outboundTraffic += amount
	}
}
