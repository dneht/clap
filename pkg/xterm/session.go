package xterm

import (
	"cana.io/clap/pkg/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/websocket/v2"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"strings"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// End of transmission (Ctrl-C Ctrl-C Ctrl-D Ctrl-D).
	endOfTransmission = "\u0003\u0003\u0004\u0004"
)

var toastTypes = [...]string{"alert", "success", "error", "warning", "info"}

// PtyHandler is what remotecommand expects from a pty
type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

// TerminalSession implements PtyHandler
type TerminalSession struct {
	id       string
	ws       *websocket.Conn
	lock     *sync.Mutex
	input    *ExecInput
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

// TerminalMessage is the messaging protocol between ShellController and TerminalSession.
//
// OP      DIRECTION  FIELD(S) USED  DESCRIPTION
// ---------------------------------------------------------------------
// bind    fe->be     SessionID      Id sent back from TerminalResponse
// stdin   fe->be     Data           Keystrokes/paste buffer
// resize  fe->be     Rows, Cols     New terminal size
// stdout  be->fe     Data           Output from the process
// toast   be->fe     Data           OOB message to be shown to the user
//
// ToastType
// alert, success, error, warning, info - ClassName generator uses this value → noty_type__${type}
type TerminalMessage struct {
	Op, Data, SessionID, ToastType string
	Rows, Cols                     uint16
}

// Next handles pty->process resize events
// Called in a loop from remotecommand as long as the process is running
func (t TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

// Read handles pty->process messages (stdin, resize)
// Called in a loop from remotecommand as long as the process is running
func (t TerminalSession) Read(p []byte) (int, error) {
	err := t.ws.SetReadDeadline(time.Now().Add(time.Second * time.Duration(t.input.Timeout)))
	if nil != err {
		return 0, err
	}
	_, m, err := t.ws.ReadMessage()
	if nil != err {
		log.Warnf("ws read message failed", err)
		// Send terminated signal to process to avoid resource leak
		t.Toast("error", err.Error())
		return copy(p, endOfTransmission), err
	}

	var msg TerminalMessage
	if err := json.Unmarshal(m, &msg); nil != err {
		log.Warnf("ws decode message failed", err)
		t.Toast("error", err.Error())
		return copy(p, endOfTransmission), err
	}

	switch msg.Op {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	default:
		t.Toast("error", "Unknown operation.")
		return copy(p, endOfTransmission), fmt.Errorf("unknown message type '%s'", msg.Op)
	}
}

// Write handles process->pty stdout
// Called from remotecommand whenever there is any output
func (t TerminalSession) Write(p []byte) (int, error) {
	// 避免 Write Toast 同时写(Observed a panic: concurrent write to websocket connection)
	t.lock.Lock()
	defer t.lock.Unlock()
	err := t.ws.SetWriteDeadline(time.Now().Add(writeWait))
	if nil != err {
		return 0, err
	}
	data := strings.ReplaceAll(string(p), "\u0000", "")
	data = strings.ReplaceAll(data, "\n", "\r\n")
	if strings.Index(data, "\r\r\n") > 0 {
		data = strings.ReplaceAll(data, "\r\r\n", "\r\n")
	}
	msg, err := json.Marshal(TerminalMessage{
		SessionID: t.id,
		Op:        "stdout",
		Data:      data,
	})
	if nil != err {
		log.Warnf("ws encode message failed", err)
		return 0, err
	}
	if err = t.ws.WriteMessage(websocket.TextMessage, msg); nil != err {
		log.Warnf("ws write message failed", err)
		return 0, err
	}
	return len(p), nil
}

// Read logs
func Logs(t *TerminalSession, clientset *kubernetes.Clientset) error {
	reqLog := clientset.CoreV1().RESTClient().Get().
		Resource("pods").
		Name(t.input.PodName).
		Namespace(t.input.Namespace).
		SubResource("log")
	logOptions := &corev1.PodLogOptions{
		Container:  t.input.ContainerName,
		Follow:     true,
		Previous:   false,
		Timestamps: false,
	}
	if t.input.SinceSecond <= 0 {
		logOptions.TailLines = &t.input.TailLines
	} else {
		logOptions.SinceSeconds = &t.input.SinceSecond
	}
	readCloser, err := reqLog.VersionedParams(logOptions, scheme.ParameterCodec).Stream(context.TODO())
	if nil != err {
		return err
	}
	defer readCloser.Close()
	ticker := time.NewTicker(10 * time.Millisecond)
	for _ = range ticker.C {
		rawLog := make([]byte, 1024)
		_, err := readCloser.Read(rawLog)
		if nil != err {
			ticker.Stop()
			return err
		}
		_, err = t.Write(rawLog)
		if nil != err {
			ticker.Stop()
			break
		}
	}
	return nil
}

// Toast can be used to send the user any OOB messages
func (t TerminalSession) Toast(toastType, p string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if nil == t.ws || nil == t.ws.Conn {
		return
	}
	if err := t.ws.SetWriteDeadline(time.Now().Add(writeWait)); nil != err {
		println(err)
		return
	}
	index := -1
	for i, x := range toastTypes {
		if toastType == x {
			index = i
			break
		}
	}
	if index == -1 {
		toastType = toastTypes[0]
	}
	msg, err := json.Marshal(TerminalMessage{
		Op:        "toast",
		Data:      p,
		ToastType: toastType,
	})
	if nil != err {
		println(err)
		return
	}
	if err := t.ws.WriteMessage(websocket.TextMessage, msg); nil != err {
		println(err)
		return
	}
}

func readFromSteam(t *TerminalSession, capacity int64, readCloser io.ReadCloser) error {
	rawLog := make([]byte, capacity)
	_, err := readCloser.Read(rawLog)
	if nil != err {
		return err
	}
	_, err = t.Write(rawLog)
	return err
}
