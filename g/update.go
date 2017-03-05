package g

import (
	"errors"
	"fmt"
	"github.com/toolkits/file"
	"github.com/toolkits/sys"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

const (
	ExecFile    = "falcon-agent"
	TmpExecFile = "falcon-agent-tmp"
)

func RestartDaemon() {
	execSpec := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}
	fork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		log.Printf("Fail to fork new process: %s", err)
		return
	}
	log.Printf("Fork new process id: %d", fork)
	log.Printf("Old server(%d) gracefully shutdown.", os.Getpid())
	pidfile := Config().PidFile
	f, err := os.Create(pidfile)
	if err != nil {
		log.Printf("write %s err: %s", pidfile, err)
		f.Close()
	} else {
		pidstr := fmt.Sprintf("%d\n", fork)
		f.WriteString(pidstr)
		f.Close()
	}
	// stop the old server
	os.Exit(0)
}

func AutoUpdateChk(v string) {
	if !Config().AutoUpdate.Enabled {
		return
	}
	var err error
	file.Remove(TmpExecFile)
	tar_prefix := Config().AutoUpdate.Tar
	agent_new := fmt.Sprintf(tar_prefix, v)
	err = UpdateAgent(agent_new)
	if err != nil {
		log.Println("Get Agent Fail with : ", err)
		return
	}

	RestartDaemon()
}

func UpdateAgent(filename string) error {
	fpath := file.SelfDir()
	server_url := Config().AutoUpdate.Url
	url := fmt.Sprintf("%s/%s", server_url, filename)
	var err error
	debug := Config().Debug
	timeout := time.Duration(Config().AutoUpdate.Timeout) * time.Millisecond
	if debug {
		log.Println("Downloading", url, "to", filename)
	}
	err = file.EnsureDir(Config().AutoUpdate.Dir)
	if err != nil {
		return err
	}
	rel_filename := filepath.Join(Config().AutoUpdate.Dir, filename)

	client := http.Client{
		Timeout: timeout,
	}
	response, err := client.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 && response.StatusCode != 301 && response.StatusCode != 302 {
		errtxt := fmt.Sprintf("response code is %d", response.StatusCode)
		return errors.New(errtxt)
	}

	output, err := os.Create(rel_filename)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return err
	}

	err = UnCompressAgentTest(rel_filename)
	if err != nil {
		return err
	}

	err = UnCompressAgent(rel_filename, fpath)
	if err != nil {
		return err
	}

	return nil
}

func UnCompressAgent(filename, destdir string) error {
	var err error
	err = file.Rename(ExecFile, TmpExecFile)
	if err != nil {
		return err
	}
	err = UnTarGz(filename, destdir)
	if err != nil {
		return err
	}
	if !file.IsExist(ExecFile) {
		file.Rename(TmpExecFile, ExecFile)
		return errors.New("File not found,rollback now.")
	}
	file.Remove(filename)
	return nil
}

func UnCompressAgentTest(filename string) error {
	var err error
	tmpdir := filepath.Join(Config().AutoUpdate.Dir, "tmp")
	if file.IsExist(ExecFile) {
		err = os.RemoveAll(tmpdir)
		if err != nil {
			return err
		}
	}
	file.EnsureDir(tmpdir)
	err = UnTarGz(filename, tmpdir)
	if err != nil {
		return err
	}

	if !file.IsExist(filepath.Join(tmpdir, ExecFile)) {
		return errors.New("File not found.")
	}
	os.RemoveAll(tmpdir)
	return nil
}

func UnTarGz(srcFilePath string, destDirPath string) error {
	_, err := sys.CmdOutBytes("tar", "zxf", srcFilePath, "-C", destDirPath)
	if err != nil {
		return err
	}
	return nil
}
