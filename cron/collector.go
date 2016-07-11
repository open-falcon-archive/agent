package cron

import (
	"log"
	"time"

	"github.com/open-falcon/agent/funcs"
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/common/model"
)

func InitDataHistory() {
	for {
		funcs.UpdateCpuStat()
		funcs.UpdateDiskStats()
		time.Sleep(g.COLLECT_INTERVAL)
	}
}

func Collect() {

	if !g.Config().Transfer.Enabled {
		return
	}

	if len(g.Config().Transfer.Addrs) == 0 {
		return
	}

	for _, v := range funcs.Mappers {
		go collect(int64(v.Interval), v.Fs)
	}
}

func collect(sec int64, fns []func() []*model.MetricValue) {
	t := time.NewTicker(time.Second * time.Duration(sec)).C
	for {
		<-t

		hostname, err := g.Hostname()
		if err != nil {
			continue
		}

		mvs := []*model.MetricValue{}
		debug := g.Config().Debug

		for _, fn := range fns {
			items := fn()
			if items == nil {
				continue
			}

			if len(items) == 0 {
				continue
			}

			if debug {
				log.Println(" -> collect ", len(items), " metrics\n")
			}
			mvs = append(mvs, items...)
		}

		filtered := *g.FilterMetrics(&mvs)

		now := time.Now().Unix()
		for j := 0; j < len(filtered); j++ {
			filtered[j].Step = sec
			filtered[j].Endpoint = hostname
			filtered[j].Timestamp = now
		}

		g.SendToTransfer(filtered)
	}
}
