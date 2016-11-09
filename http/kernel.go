package http

import (
	"net/http"

	"github.com/open-falcon/agent/g"
	"github.com/toolkits/nux"
	"github.com/toolkits/sys"
)

func configKernelRoutes() {
	http.HandleFunc("/proc/kernel/hostname", func(w http.ResponseWriter, r *http.Request) {
		data, err := g.Hostname()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/maxproc", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.KernelMaxProc()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/maxfiles", func(w http.ResponseWriter, r *http.Request) {
		data, err := nux.KernelMaxFiles()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/version", func(w http.ResponseWriter, r *http.Request) {
		data, err := sys.CmdOutNoLn("uname", "-r")
		AutoRender(w, data, err)
	})

}
