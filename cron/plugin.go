package cron

import (
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/agent/plugins"
	"github.com/open-falcon/common/model"
	"log"
	"strings"
	"time"
	"github.com/toolkits/file"
	"os/exec"
	"fmt"
	"errors"
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
	if err := UpdatePlugin(); err != nil {
		log.Fatalln(err.Error())
	}
	go syncMinePlugins()
}

func syncMinePlugins() {

	var (
		timestamp  int64 = -1
		pluginDirs []string
	)

	duration := time.Duration(g.Config().Heartbeat.Interval) * time.Second

	for {
		time.Sleep(duration)

		hostname, err := g.Hostname()
		if err != nil {
			continue
		}

		req := model.AgentHeartbeatRequest{
			Hostname: hostname,
		}

		var resp model.AgentPluginsResponse
		err = g.HbsClient.Call("Agent.MinePlugins", req, &resp)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		if resp.Timestamp <= timestamp {
			continue
		}

		pluginDirs = resp.Plugins
		timestamp = resp.Timestamp

		if g.Config().Debug {
			log.Println(&resp)
		}

		if len(pluginDirs) == 0 {
			plugins.ClearAllPlugins()
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
//如果插件启用,则用git地址下载插件
func UpdatePlugin() error {
	if !g.Config().Plugin.Enabled {
		return nil
	}
	dir := g.Config().Plugin.Dir
	parentDir := file.Dir(dir)
	file.InsureDir(parentDir)

	if file.IsExist(dir) {
		// git pull
		cmd := exec.Command("git", "pull")
		cmd.Dir = dir
		err := cmd.Run()
		if err != nil {
			return errors.New(fmt.Sprintf("git pull in dir:%s fail. error: %s", dir, err))
		}
	} else {
		// git clone
		cmd := exec.Command("git", "clone", g.Config().Plugin.Git, file.Basename(dir))
		cmd.Dir = parentDir
		err := cmd.Run()
		if err != nil {
			return errors.New((fmt.Sprintf("git clone plugin url:%s into dir:%s fail. error: %s", g.Config().Plugin.Git, dir, err)))
		}
	}
	return nil
}
