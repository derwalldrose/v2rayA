package configure

type ObservatoryType string

func (t ObservatoryType) String() string {
	return string(t)
}

const (
	Random     ObservatoryType = "random"
	RoundRobin ObservatoryType = "roundrobin"
	LeastPing  ObservatoryType = "leastping"
	LeastLoad  ObservatoryType = "leastload"
)

type OutboundSetting struct {
	ProbeURL      string          `json:"probeURL"`
	ProbeInterval string          `json:"probeInterval"`
	Type          ObservatoryType `json:"type"`
	FallbackTag   string          `json:"fallbackTag,omitempty"`
}
