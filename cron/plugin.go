package cron

import (
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/agent/plugins"
	"github.com/open-falcon/common/model"
	"log"
	"strings"
	"time"
)

func SyncMinePlugins() {
	if !g.Config().Plugin.Enabled {
		return
	}

	if !g.Config().Heartbeat.Enabled {
		return
	}

	if g.Config().Heartbeat.Addr == "" {
		return
	}

	go syncMinePlugins()
}

func syncMinePlugins() {

	var (
		timestamp  int64 = -1
		pluginDirs []string
	)

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second
	protocol := g.Config().Plugin.Protocol
	log.Println("the plugin's update protocol is:", protocol)
	for {
	REST:
		time.Sleep(duration)

		hostname, err := g.Hostname()
		if err != nil {
			goto REST
		}

		req := model.AgentHeartbeatRequest{
			Hostname: hostname,
		}

		var resp model.AgentPluginsResponse
		err = g.HbsClient.Call("Agent.MinePlugins", req, &resp)
		if err != nil {
			log.Println("ERROR:", err)
			goto REST
		}

		if resp.Timestamp <= timestamp {
			goto REST
		}

		pluginDirs = resp.Plugins
		timestamp = resp.Timestamp
		if g.Config().Debug {
			log.Println(&resp)
		}

		if len(pluginDirs) == 0 {
			plugins.ClearAllPlugins()
		}
		//download the plugins that has updated with http
		if protocol == "http" {
			CheckPluginUpdate(pluginDirs, duration)
		}
		desiredAll := make(map[string]*plugins.Plugin)

		for _, p := range pluginDirs {
			underOneDir := plugins.ListPlugins(strings.Trim(p, "/"))
			for k, v := range underOneDir {
				desiredAll[k] = v
			}
		}

		plugins.DelNoUsePlugins(desiredAll)
		plugins.AddNewPlugins(desiredAll)

	}
}
