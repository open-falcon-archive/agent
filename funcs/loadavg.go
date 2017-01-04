package funcs

import (
	"log"

	"github.com/open-falcon/common/model"
	"github.com/toolkits/nux"
	"runtime"
	"strconv"
)

func LoadAvgMetrics() []*model.MetricValue {
	load, err := nux.LoadAvg()
	if err != nil {
		log.Println(err)
		return nil
	}
	cpuNum := strconv.Itoa(runtime.NumCPU())
	tagCpu := "cpu_num="+cpuNum
	return []*model.MetricValue{
		GaugeValue("load.1min", load.Avg1min, tagCpu),
		GaugeValue("load.5min", load.Avg5min, tagCpu),
		GaugeValue("load.15min", load.Avg15min, tagCpu),
	}

}
