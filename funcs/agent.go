package funcs

import (
	"github.com/cepave/common/model"
)

func AgentMetrics() []*model.MetricValue {
	return []*model.MetricValue{GaugeValue("agent.alive", 1)}
}
