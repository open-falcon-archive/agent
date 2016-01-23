package funcs

import (
	"github.com/cepave/agent/g"
	"github.com/cepave/common/model"
	"github.com/toolkits/nux"
	"log"
	"strings"
)

func NetMetrics() []*model.MetricValue {
	return CoreNetMetrics(g.Config().Collector.IfacePrefix)
}

func CoreNetMetrics(ifacePrefix []string) []*model.MetricValue {

	netIfs, err := nux.NetIfs(ifacePrefix)
	if err != nil {
		log.Println(err)
		return []*model.MetricValue{}
	}

	cnt := len(netIfs)
	ret := make([]*model.MetricValue, cnt*20+1+1)

	for idx, netIf := range netIfs {
		iface := "iface=" + netIf.Iface
		ret[idx*20+0] = CounterValue("net.if.in.bits", netIf.InBytes*8, iface)
		ret[idx*20+1] = CounterValue("net.if.in.packets", netIf.InPackages, iface)
		ret[idx*20+2] = CounterValue("net.if.in.errors", netIf.InErrors, iface)
		ret[idx*20+3] = CounterValue("net.if.in.dropped", netIf.InDropped, iface)
		ret[idx*20+4] = CounterValue("net.if.in.fifo.errs", netIf.InFifoErrs, iface)
		ret[idx*20+5] = CounterValue("net.if.in.frame.errs", netIf.InFrameErrs, iface)
		ret[idx*20+6] = CounterValue("net.if.in.compressed", netIf.InCompressed, iface)
		ret[idx*20+7] = CounterValue("net.if.in.multicast", netIf.InMulticast, iface)
		ret[idx*20+8] = CounterValue("net.if.out.bits", netIf.OutBytes*8, iface)
		ret[idx*20+9] = CounterValue("net.if.out.packets", netIf.OutPackages, iface)
		ret[idx*20+10] = CounterValue("net.if.out.errors", netIf.OutErrors, iface)
		ret[idx*20+11] = CounterValue("net.if.out.dropped", netIf.OutDropped, iface)
		ret[idx*20+12] = CounterValue("net.if.out.fifo.errs", netIf.OutFifoErrs, iface)
		ret[idx*20+13] = CounterValue("net.if.out.collisions", netIf.OutCollisions, iface)
		ret[idx*20+14] = CounterValue("net.if.out.carrier.errs", netIf.OutCarrierErrs, iface)
		ret[idx*20+15] = CounterValue("net.if.out.compressed", netIf.OutCompressed, iface)
		ret[idx*20+16] = CounterValue("net.if.total.bits", netIf.TotalBytes*8, iface)
		ret[idx*20+17] = CounterValue("net.if.total.packets", netIf.TotalPackages, iface)
		ret[idx*20+18] = CounterValue("net.if.total.errors", netIf.TotalErrors, iface)
		ret[idx*20+19] = CounterValue("net.if.total.dropped", netIf.TotalDropped, iface)
	}

	inTotalBits := int64(0)
	outTotalBits := int64(0)
	for _, netIf := range netIfs {
		if strings.Contains(netIf.Iface, "eth") {
			inTotalBits += netIf.InBytes * 8
			outTotalBits += netIf.OutBytes * 8
		}
	}
	ret[cnt*20+0] = CounterValue("net.if.in.bits", inTotalBits, "iface=eth_all")
	ret[cnt*20+1] = CounterValue("net.if.out.bits", outTotalBits, "iface=eth_all")

	return ret
}
