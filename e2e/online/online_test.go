package online

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
)

type testData struct {
	SourcePort    int
	TransformPort int
	SinkPort      int
}

func TestOnlineLinearPipeline(t *testing.T) {
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

	waitFunc, killFunc := runMaestro("run", "-f", cfgPath, "-v")

	<-done
	killFunc()
	waitFunc()

	prev := int64(0)
	for i, msg := range collect {
		if prev >= msg.Val {
			t.Fatalf("wrong value order at %d, %d: values are %d, %d", i-1, i, prev, msg.Val)
		}
		if msg.Val%2 != 0 {
			t.Fatalf("value %d is not pair: %d", i, msg.Val)
		}
		prev = msg.Val
	}
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

func runMaestro(args ...string) (func(), func()) {
	executable := "../../target/maestro"
	cmd := exec.Command(executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintf("start cmd %q: %s", cmd, err))
	}
	waitFunc := func() {
		err := cmd.Wait()
		if err == nil {
			return
		}
		var exitErr *exec.ExitError
		// Exit by signal is ok as we are killing the process
		if errors.As(err, &exitErr) && exitErr.ProcessState.ExitCode() <= 0 {
			return
		}
		panic(fmt.Sprintf("wait cmd %q: %s", cmd, err))
	}
	killFunc := func() {
		if err := cmd.Process.Kill(); err != nil {
			panic(fmt.Sprintf("kill cmd %q: %s", cmd, err))
		}
	}
	return waitFunc, killFunc
}
