package http

import (
	"fmt"
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/agent/plugins"
	"github.com/toolkits/file"
	"net/http"
	"os/exec"
	"github.com/open-falcon/agent/cron"
)

func configPluginRoutes() {
	http.HandleFunc("/plugin/update", func(w http.ResponseWriter, r *http.Request) {
		if !g.Config().Plugin.Enabled {
			w.Write([]byte("plugin not enabled"))
			return
		}
		if err := cron.UpdatePlugin(); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("success"))
	})

	http.HandleFunc("/plugin/reset", func(w http.ResponseWriter, r *http.Request) {
		if !g.Config().Plugin.Enabled {
			w.Write([]byte("plugin not enabled"))
			return
		}

		dir := g.Config().Plugin.Dir

		if file.IsExist(dir) {
			cmd := exec.Command("git", "reset", "--hard")
			cmd.Dir = dir
			err := cmd.Run()
			if err != nil {
				w.Write([]byte(fmt.Sprintf("git reset --hard in dir:%s fail. error: %s", dir, err)))
				return
			}
		}
		w.Write([]byte("success"))
	})

	http.HandleFunc("/plugins", func(w http.ResponseWriter, r *http.Request) {
		//TODO: not thread safe
		RenderDataJson(w, plugins.Plugins)
	})
}
