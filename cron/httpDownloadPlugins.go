package cron

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/open-falcon/agent/g"
	"github.com/toolkits/file"
	"github.com/toolkits/sys"
)

func WgetRun(filename string, dst string, Timeout time.Duration) error {
	timeout := Timeout*1000 - 500
	debug := g.Config().Debug
	duri := g.Config().Plugin.Uri + "/" + filename
	shell := "wget " + duri + " -O " + dst
	if debug {
		log.Println(shell, " running...")
	}
	cmd := exec.Command("/bin/bash", "-c", shell)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Start()

	err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Millisecond)

	errStr := stderr.String()
	file_prefix := strings.Split(filename, ".")[0]
	if errStr != "" {
		logFile := filepath.Join(g.Config().Plugin.LogDir, file_prefix+".stderr.log")
		if _, err = file.WriteString(logFile, errStr); err != nil {
			log.Printf("[ERROR] write log to %s fail, error: %s\n", logFile, err)
		}
	}

	if isTimeout {
		// has be killed
		if err == nil && debug {
			log.Println("[INFO] timeout and kill process", cmd, "successfully")
		}

		if err != nil {
			log.Println("[ERROR] kill process", cmd, "occur error:", err)
		}

		return err
	}

	if err != nil {
		log.Println("[ERROR] exec cmd:", cmd, "fail. error:", err)
		return err
	}
	return err
}

func ExeCmd(shell string) error {
	cmdobj := exec.Command("/bin/bash", "-c", shell)
	err := cmdobj.Run()
	if err != nil {
		log.Println("[ERROR] cmd:", shell, "fail.error:", err)
		return err
	}
	return err
}

func CheckMd5Update(filename string, Timeout time.Duration) bool {
	file.InsureDir("./md5file")
	file.InsureDir("./tmp")
	dstfile := filepath.Join("./tmp", filename+".tmp")
	if err := WgetRun(filename, dstfile, Timeout); err != nil {
		return false
	}
	//get local md5 value
	local_md5_path := filepath.Join("./md5file", filename)
	if !file.IsExist(local_md5_path) {
		return true
	}
	lbuff, err := ioutil.ReadFile(local_md5_path)
	if err != nil {
		log.Println("read file from file:", local_md5_path, "fail.error is ", err)
		return true
	}
	local_md5 := string(strings.Trim(string(lbuff), "")[0:32])
	//get remote md5 value
	remote_md5_path := filepath.Join(dstfile)
	rbuff, err := ioutil.ReadFile(remote_md5_path)
	if err != nil {
		log.Println("read file from file:", remote_md5_path, "fail.error is ", err)
		return false
	}
	remote_md5 := string(strings.Trim(string(rbuff), "")[0:32])
	if local_md5 != remote_md5 {
		log.Println("md5 change with remote,update plugin.", filename)
		return true
	}
	return false
}

func ValidateFileMd5(md5file string) bool {
	shell := "cd ./tmp && md5sum -c " + md5file
	err := ExeCmd(shell)
	if err != nil {
		log.Println("the md5 file validate failed:", md5file, "the error:", err)
		return false
	}
	return true
}

func TarPackage(plugin_name string, dstdir string) bool {
	tgzfile := filepath.Join("./tmp", plugin_name+".tar.gz")
	srcplugin := filepath.Join("./plugins", plugin_name)
	RmDir(srcplugin)
	shell := "tar -zxf " + tgzfile + " -C " + dstdir
	err := ExeCmd(shell)
	if err != nil {
		log.Println("extract file ", tgzfile, "faild. error:", err)
		return false
	}
	return true
}

func CpFile(srcdir string, dstdir string) bool {
	shell := "cp " + srcdir + " " + dstdir
	err := ExeCmd(shell)
	if err != nil {
		log.Println("cp src dir", srcdir, "to dst dir", dstdir, "failed. error:", err)
		return false
	}
	return true
}

func RmDir(dir string) bool {
	if file.IsExist(dir) {
		shell := "rm -rf " + dir
		err := ExeCmd(shell)
		if err != nil {
			log.Println("rm dir or file ", dir, "faild. error:", err)
			return false
		}
		return true
	}
	return true
}

func CheckPluginUpdate(plugins []string, cycle time.Duration) {
	for _, plugin := range plugins {
		if ok := CheckMd5Update(plugin+".md5", 60); !ok {
			log.Println("the plugin:", plugin, "has not updated.")
			continue
		}
		log.Println("the plugin:", plugin, "has updated.")
		//Plugin's md5 change with local,update plugin
		pack_name := plugin + ".tar.gz"
		dst_name := filepath.Join("./tmp", plugin+".tar.gz")
		WgetRun(pack_name, dst_name, cycle)
		if ok := ValidateFileMd5(plugin + ".md5.tmp"); !ok {
			RmDir("./tmp/" + plugin + ".md5.tmp")
			RmDir("./tmp/" + plugin + ".tar.gz")
			continue
		}
		//decompress package
		file.InsureDir("./plugins")
		if ok := TarPackage(plugin, "./plugins"); !ok {
			continue
		}
		//update tmp md5file to real md5file
		if ok := CpFile("./tmp/"+plugin+".md5.tmp", "./md5file/"+plugin+".md5"); !ok {
			continue
		}
	}
	return
}
