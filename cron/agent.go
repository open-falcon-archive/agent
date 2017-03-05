package cron

import (
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/common/model"
	"log"
	"math/rand"
	"time"
)

func SyncMineAgentVersion() {
	if g.Config().AutoUpdate.Enabled && g.Config().Heartbeat.Enabled && g.Config().Heartbeat.Addr != "" {
		go syncMineAgentVersion()
	}
}

func syncMineAgentVersion() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	duration := time.Duration(g.Config().AutoUpdate.Interval) * time.Second
	randInt := g.Config().AutoUpdate.RandInt

	for {
		time.Sleep(duration)
		hostname, err := g.Hostname()
		if err != nil {
			return
		}

		rtime := r.Intn(randInt)
		time.Sleep(time.Duration(rtime) * time.Second)

		req := model.AgentHeartbeatRequest{
			Hostname: hostname,
		}

		var version string
		err = g.HbsClient.Call("Agent.MineAgentVersion", req, &version)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}
		if version == "" {
			continue
		}
		if g.VERSION != version {
			log.Printf("Agent from %s to %s.", g.VERSION, version)
			g.AutoUpdateChk(version)
		}
	}
}
