package funcs

import (
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/common/model"
	"github.com/toolkits/nux"
	"log"
)

func ProcMetrics() (L []*model.MetricValue) {

	reportProcs := g.ReportProcs()
	sz := len(reportProcs)
	if sz == 0 {
		return
	}

	ps, err := nux.AllProcs()
	if err != nil {
		log.Println(err)
		return
	}

	pslen := len(ps)

	for tags, preProc := range reportProcs {
		cnt := 0
		pids := map[int]struct{}{}
		for i := 0; i < pslen; i++ {
			if is_a(ps[i], preProc) {
				cnt++
				pids[ps[i].Pid] = struct{}{}
			}
		}
		if cnt > 0 && is_restart_proc(pids, preProc) {
			//此进程发生了重启,将进程数置为0
			cnt = 0
		}
		L = append(L, GaugeValue(g.PROC_NUM, cnt, tags))
	}
	return
}

func is_restart_proc(pids map[int]struct{}, preProc *g.CacheProc) bool {
	if len(pids) != len(preProc.Pids) {
		//进程数目前后不一致,则认为发生重启
		return true
	}
	if len(preProc.Pids) == 1 {
		if _, ok := preProc.Pids[-1]; ok {
			//此进程是新加的进程监控,无需对比两次采集的进程号
			preProc.Pids = pids
			return false
		}
	}
	flag := false
	for pid, _ := range pids {
		if _, ok := preProc.Pids[pid]; !ok {
			//由tag name、cmdline标识的同一个进程在两次采集的过程中发生了重启
			flag = true
			break
		}
	}
	preProc.Pids = pids
	return flag
}

func is_a(p *nux.Proc, preProc *g.CacheProc) bool {
	// name
	if len(preProc.Name) > 0 {
		//进程监控设置了name tag
		if p.Name != preProc.Name {
			return false
		}
	}
	// cmdline
	if len(preProc.Cmdline) > 0 {
		//进程监控设置了cmdline tag
		if p.Cmdline != preProc.Cmdline {
			return false
		}
	}
	return true
}
