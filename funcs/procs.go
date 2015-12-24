package funcs

import (
    "github.com/open-falcon/agent/g"
    "github.com/open-falcon/common/model"
    "github.com/shwinpiocess/nux"
    "log"
    "strings"
    "fmt"
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

    for tags, m := range reportProcs {
        cnt := 0
        for i := 0; i < pslen; i++ {
            if is_a(ps[i], m) {
                cnt++
                new_tags := fmt.Sprintf("%s,pid=%d", tags, ps[i].Pid)
                L = append(L, CounterValue("process.cpu.user", ps[i].CpuUser, new_tags))
                L = append(L, CounterValue("process.cpu.sys", ps[i].CpuSys, new_tags))
                L = append(L, CounterValue("process.cpu.all", ps[i].CpuAll, new_tags))
                L = append(L, GaugeValue("process.mem", ps[i].Mem, new_tags))
                L = append(L, GaugeValue("process.swap", ps[i].Swap, new_tags))
                L = append(L, GaugeValue("process.fd", ps[i].Fd, new_tags))
            }
        }

        L = append(L, GaugeValue(g.PROC_NUM, cnt, tags))
    }

    return
}

func is_a(p *nux.Proc, m map[int]string) bool {
    // only one kv pair
    for key, val := range m {
        if key == 1 {
            // name
            if val != p.Name {
                return false
            }
        } else if key == 2 {
            // cmdline
            if !strings.Contains(p.Cmdline, val) {
                return false
            }
        }
    }
    return true
}
