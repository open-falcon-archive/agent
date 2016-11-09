package funcs

import (
	"log"

	"github.com/open-falcon/common/model"
	"github.com/toolkits/nux"
)

func KernelMetrics() (L []*model.MetricValue) {

	maxFiles, err := nux.KernelMaxFiles()
	if err != nil {
		log.Println(err)
		return
	}

	L = append(L, GaugeValue("kernel.maxfiles", maxFiles))

	maxProc, err := nux.KernelMaxProc()
	if err != nil {
		log.Println(err)
		return
	}

	L = append(L, GaugeValue("kernel.maxproc", maxProc))

	allocateFiles, err := nux.KernelAllocateFiles()
	if err != nil {
		log.Println(err)
		return
	}

	L = append(L, GaugeValue("kernel.files.allocated", allocateFiles))
	L = append(L, GaugeValue("kernel.files.left", maxFiles-allocateFiles))
	return
}
