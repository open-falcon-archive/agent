package funcs

import (
	"log"

	"github.com/open-falcon/common/model"
	"github.com/toolkits/nux"
)

var USES = map[string]struct{}{
	"PruneCalled":        struct{}{},
	"LockDroppedIcmps":   struct{}{},
	"ArpFilter":          struct{}{},
	"TW":                 struct{}{},
	"DelayedACKLocked":   struct{}{},
	"ListenOverflows":    struct{}{},
	"ListenDrops":        struct{}{},
	"TCPPrequeueDropped": struct{}{},
	"TCPTSReorder":       struct{}{},
	"TCPDSACKUndo":       struct{}{},
	"TCPLoss":            struct{}{},
	"TCPLostRetransmit":  struct{}{},
	"TCPLossFailures":    struct{}{},
	"TCPFastRetrans":     struct{}{},
	"TCPTimeouts":        struct{}{},
	"TCPSchedulerFailed": struct{}{},
	"TCPAbortOnMemory":   struct{}{},
	"TCPAbortOnTimeout":  struct{}{},
	"TCPAbortFailed":     struct{}{},
	"TCPMemoryPressures": struct{}{},
	"TCPSpuriousRTOs":    struct{}{},
	"TCPBacklogDrop":     struct{}{},
	"TCPMinTTLDrop":      struct{}{},
}

func NetstatMetrics() (L []*model.MetricValue) {
	tcpExts, err := nux.Netstat("TcpExt")

	if err != nil {
		log.Println(err)
		return
	}

	cnt := len(tcpExts)
	if cnt == 0 {
		return
	}

	for key, val := range tcpExts {
		if _, ok := USES[key]; !ok {
			continue
		}
		L = append(L, CounterValue("TcpExt."+key, val))
	}

	return
}
