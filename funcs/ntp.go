/*
 * author zhangkun
 * colletc ntp offset
 */
package funcs

import (
	"github.com/open-falcon/common/model"

	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func NtpMetrics() []*model.MetricValue {
	cmd := exec.Command("ntpq", "-c", "rv 0 offset")

	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput

	err := cmd.Run() // will wait for command to return

	if err != nil {
		log.Println(err)
		return []*model.MetricValue{}
	}

	str := strings.TrimSpace(cmdOutput.String())
	if len(str) == 0 {
		log.Println("cmd out is null")
		return []*model.MetricValue{}
	}

	index := strings.Index(str, "=")
	str_offset := str[index+1:]

	if offset, err := strconv.ParseFloat(str_offset, 64); err == nil {
		return []*model.MetricValue{
			GaugeValue("ntp.offset", offset),
		}
	}
	return []*model.MetricValue{}
}
