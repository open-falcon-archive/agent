package funcs

import (
	"github.com/open-falcon/common/model"
	"github.com/toolkits/nux"
	"log"
	"runtime"
)

func LoadAvgMetrics() []*model.MetricValue {
	load, err := nux.LoadAvg()
	if err != nil {
		log.Println(err)
		return nil
	}
	cpuNum := float64(runtime.NumCPU())
	var load1minPercent, load5minPercent, load15minPercent float64
	if cpuNum != 0{
		load1minPercent = load.Avg1min * 100 / cpuNum
		load5minPercent = load.Avg5min * 100/ cpuNum
		load15minPercent = load.Avg15min * 100/ cpuNum
	}
	return []*model.MetricValue{
		GaugeValue("load.1min.percent", load1minPercent),
		GaugeValue("load.5min.percent", load5minPercent),
		GaugeValue("load.15min.percent", load15minPercent),
		GaugeValue("load.1min", load.Avg1min),
		GaugeValue("load.5min", load.Avg5min),
		GaugeValue("load.15min", load.Avg15min),
	}
}
