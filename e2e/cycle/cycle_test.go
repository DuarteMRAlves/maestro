package cycle

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	sync "sync"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOfflineCycle(t *testing.T) {
	var mu sync.Mutex
	max := 100
	collect := make([]*SumMessage, 0, max)
	done := make(chan struct{})
	collectFunc := func(msg *SumMessage) {
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

	counterAddr, counterStop := ServeCounter()
	defer counterStop()

	incAddr, incStop := ServeInc()
	defer incStop()

	sumAddr, sumStop := ServeSum(collectFunc)
	defer sumStop()

	tmplData := testData{
		CounterPort: extractPort(counterAddr),
		IncPort:     extractPort(incAddr),
		SumPort:     extractPort(sumAddr),
	}

	tempDir := t.TempDir()
	cfgPath := tempDir + "/config.yaml"
	writeTemplate(cfgPath, "config.yaml", tmplData)

	waitFunc, killFunc := runMaestro("run", "-f", cfgPath, "-v")
	<-done
	killFunc()
	waitFunc()

	incVal := int64(0)
	for i, msg := range collect {
		counterVal := int64(i + 1)
		if diff := cmp.Diff(counterVal, msg.Counter.Val); diff != "" {
			t.Fatalf("mismatch on counter value %d:\n%s", i, diff)
		}
		if diff := cmp.Diff(incVal, msg.Inc.Val); diff != "" {
			t.Fatalf("mismatch on inc value %d:\n%s", i, diff)
		}
		incVal = counterVal + incVal + 1
	}

}

type testData struct {
	CounterPort int
	IncPort     int
	SumPort     int
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
