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

	for tags, curProc := range reportProcs {
		cnt := 0
		pids := map[int]struct{}{}
		for i := 0; i < pslen; i++ {
			if is_a(ps[i], curProc) {
				cnt++
				pids[ps[i].Pid] = struct{}{}
			}
		}
		if cnt > 0 && is_restart_proc(pids, curProc) {
			//此进程发生了重启,将进程数置为0
			cnt = 0
		}
		L = append(L, GaugeValue(g.PROC_NUM, cnt, tags))
	}
	return
}

func is_restart_proc(pids map[int]struct{}, curProc *g.CacheProc) bool {
	if len(curProc.Pids) != len(curProc.Pids) {
		return true
	}
	if len(curProc.Pids) == 1 {
		if _, ok := curProc.Pids[-1]; ok {
			//此进程是新加的进程监控,无需对比两次采集的进程号
			curProc.Pids = pids
			return false
		}
	}
	flag := false
	for pid, _ := range pids {
		if _, ok := curProc.Pids[pid]; !ok {
			//由tag name、cmdline标识的同一个进程在两次采集的过程中发生了重启
			flag = true
			break
		}
	}
	curProc.Pids = pids
	return flag
}

func is_a(p *nux.Proc, curProc *g.CacheProc) bool {
	// name
	if len(curProc.Name) > 0 {
		//进程监控设置了name tag
		if p.Name != curProc.Name {
			return false
		}
	}
	// cmdline
	if len(curProc.Cmdline) > 0 {
		//进程监控设置了cmdline tag
		if p.Cmdline != curProc.Cmdline {
			return false
		}
	}
	return true
}
