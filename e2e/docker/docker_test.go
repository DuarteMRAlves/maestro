package docker

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/DuarteMRAlves/maestro/e2e"
	"github.com/google/go-cmp/cmp"
)

func TestDockerMestroImage(t *testing.T) {
	if *e2e.NoDocker {
		t.Skip("Environment does not support docker")
	}
	var mu sync.Mutex
	max := 100
	collect := make([]*Message, 0, max)
	done := make(chan struct{})
	collectFunc := func(msg *Message) {
		mu.Lock()
		defer mu.Unlock()
		if len(collect) < max {
			collect = append(collect, msg)
		}
		if len(collect) == max && done != nil {
			close(done)
			done = nil
		}
	}

	sourceAddr, sourceStop := ServeSource()
	defer sourceStop()

	transfAddr, transformStop := ServeTransform()
	defer transformStop()

	sinkAddr, sinkStop := ServeSink(collectFunc)
	defer sinkStop()

	tmplData := testData{
		SourcePort:    extractPort(sourceAddr),
		TransformPort: extractPort(transfAddr),
		SinkPort:      extractPort(sinkAddr),
	}

	tempDir := t.TempDir()
	cfgPath := tempDir + "/config.yaml"
	writeTemplate(cfgPath, "config.yaml", tmplData)

	waitFunc, killFunc := runMaestroContainer(cfgPath)

	<-done
	killFunc()
	waitFunc()

	for i, msg := range collect {
		if diff := cmp.Diff(int64((i+1)*2), msg.Val); diff != "" {
			t.Fatalf("mismatch at msg %d:\n%s", i, diff)
		}
	}
}

type testData struct {
	SourcePort    int
	TransformPort int
	SinkPort      int
}

func extractPort(addr net.Addr) int {
	tcpAddr, ok := addr.(*net.TCPAddr)
	if !ok {
		panic(fmt.Sprintf("not tcp addr: %s", addr))
	}
	return tcpAddr.Port
}

func writeTemplate(dst, src string, data any) {
	var buf bytes.Buffer

	tmpl := template.Must(template.ParseFiles(src))
	tmpl.Execute(&buf, data)
	err := ioutil.WriteFile(dst, buf.Bytes(), 0777)
	if err != nil {
		panic(fmt.Sprintf("write file from template: %s", err))
	}
}

func runMaestroContainer(cfgPath string) (func(), func()) {
	mount := fmt.Sprintf("type=bind,source=%s,target=/config.yaml", cfgPath)
	name := "maestro-e2e.docker-test"
	img := "duartemralves/maestro:v1-latest"

	startCmd := exec.Command(
		"docker", "run", "--rm",
		"--mount", mount,
		"--name", name,
		"--platform", "linux/amd64",
		img,
	)
	startCmd.Stdout = os.Stdout
	startCmd.Stderr = os.Stderr
	if err := startCmd.Start(); err != nil {
		panic(fmt.Sprintf("start cmd %q: %s", startCmd, err))
	}

	waitFunc := func() {
		err := startCmd.Wait()
		if err == nil {
			return
		}
		var exitErr *exec.ExitError
		// The container was terminated with the SIGKILL by the kill command
		// 137 is the exit code for SIGKILL signal.
		if errors.As(err, &exitErr) && exitErr.ProcessState.ExitCode() == 137 {
			return
		}
		panic(fmt.Sprintf("wait cmd %q: %s", startCmd, err))
	}
	killFunc := func() {
		killCmd := exec.Command("docker", "kill", name)
		killCmd.Stdout = os.Stdout
		killCmd.Stderr = os.Stderr
		if err := killCmd.Run(); err != nil {
			panic(fmt.Sprintf("kill cmd %q: %s", killCmd, err))
		}
	}
	return waitFunc, killFunc
}
