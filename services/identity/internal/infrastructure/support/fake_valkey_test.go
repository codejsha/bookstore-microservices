package support

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	sharedconfig "github.com/codejsha/shared-library-go/pkg/config"
)

type fakeValkey struct {
	listener net.Listener

	mu       sync.Mutex
	values   map[string]string
	ttls     map[string]time.Duration
	sets     map[string]map[string]struct{}
	commands [][]string
	conns    []net.Conn
}

func newFakeValkey(t *testing.T) *fakeValkey {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	f := &fakeValkey{
		listener: listener,
		values:   map[string]string{},
		ttls:     map[string]time.Duration{},
		sets:     map[string]map[string]struct{}{},
	}
	go f.accept()
	t.Cleanup(f.stop)
	return f
}

func (f *fakeValkey) stop() {
	_ = f.listener.Close()
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, conn := range f.conns {
		_ = conn.Close()
	}
}

func (f *fakeValkey) cacheConfig(t *testing.T) *sharedconfig.CacheConfig {
	t.Helper()
	host, port, err := net.SplitHostPort(f.listener.Addr().String())
	if err != nil {
		t.Fatalf("split addr: %v", err)
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		t.Fatalf("port: %v", err)
	}
	return &sharedconfig.CacheConfig{Host: host, Port: portNumber}
}

func newFakeValkeyClient(t *testing.T) (*CacheClient, *fakeValkey) {
	t.Helper()
	f := newFakeValkey(t)
	c := NewCacheClient(f.cacheConfig(t))
	t.Cleanup(func() { _ = c.Close() })
	return c, f
}

func (f *fakeValkey) value(key string) (string, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	v, ok := f.values[key]
	return v, ok
}

func (f *fakeValkey) ttl(key string) time.Duration {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.ttls[key]
}

func (f *fakeValkey) members(key string) []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return sortedMembers(f.sets[key])
}

func (f *fakeValkey) put(key, value string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.values[key] = value
}

func (f *fakeValkey) commandsNamed(name string) [][]string {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out [][]string
	for _, cmd := range f.commands {
		if strings.EqualFold(cmd[0], name) {
			out = append(out, cmd)
		}
	}
	return out
}

func sortedMembers(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for member := range set {
		out = append(out, member)
	}
	sort.Strings(out)
	return out
}

func (f *fakeValkey) accept() {
	for {
		conn, err := f.listener.Accept()
		if err != nil {
			return
		}
		f.mu.Lock()
		f.conns = append(f.conns, conn)
		f.mu.Unlock()
		go f.serve(conn)
	}
}

func (f *fakeValkey) serve(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		args, err := readRespCommand(reader)
		if err != nil {
			return
		}
		if _, err := conn.Write([]byte(f.handle(args))); err != nil {
			return
		}
	}
}

func readRespCommand(reader *bufio.Reader) ([]string, error) {
	header, err := readRespLine(reader)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(header, "*") {
		return nil, fmt.Errorf("unexpected resp header %q", header)
	}
	count, err := strconv.Atoi(header[1:])
	if err != nil {
		return nil, err
	}
	args := make([]string, count)
	for i := range args {
		sizeLine, err := readRespLine(reader)
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(sizeLine, "$") {
			return nil, fmt.Errorf("unexpected resp bulk header %q", sizeLine)
		}
		size, err := strconv.Atoi(sizeLine[1:])
		if err != nil {
			return nil, err
		}
		payload := make([]byte, size+2)
		if _, err := io.ReadFull(reader, payload); err != nil {
			return nil, err
		}
		args[i] = string(payload[:size])
	}
	return args, nil
}

func readRespLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func respBulk(value string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
}

func respInt(n int) string {
	return fmt.Sprintf(":%d\r\n", n)
}

func respArray(items []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "*%d\r\n", len(items))
	for _, item := range items {
		b.WriteString(item)
	}
	return b.String()
}

const (
	respNil = "$-1\r\n"
	respOK  = "+OK\r\n"
)

func (f *fakeValkey) handle(args []string) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.commands = append(f.commands, args)

	switch strings.ToUpper(args[0]) {
	case "GET":
		if v, ok := f.values[args[1]]; ok {
			return respBulk(v)
		}
		return respNil
	case "MGET":
		items := make([]string, 0, len(args)-1)
		for _, key := range args[1:] {
			if v, ok := f.values[key]; ok {
				items = append(items, respBulk(v))
			} else {
				items = append(items, respNil)
			}
		}
		return respArray(items)
	case "SET":
		return f.handleSet(args)
	case "DEL":
		return respInt(f.deleteKeys(args[1:]))
	case "EXISTS":
		n := 0
		for _, key := range args[1:] {
			if _, ok := f.values[key]; ok {
				n++
			}
		}
		return respInt(n)
	case "SADD":
		set, ok := f.sets[args[1]]
		if !ok {
			set = map[string]struct{}{}
			f.sets[args[1]] = set
		}
		added := 0
		for _, member := range args[2:] {
			if _, exists := set[member]; !exists {
				set[member] = struct{}{}
				added++
			}
		}
		return respInt(added)
	case "SREM":
		removed := 0
		for _, member := range args[2:] {
			if _, exists := f.sets[args[1]][member]; exists {
				delete(f.sets[args[1]], member)
				removed++
			}
		}
		return respInt(removed)
	case "SMEMBERS":
		members := sortedMembers(f.sets[args[1]])
		items := make([]string, len(members))
		for i, member := range members {
			items[i] = respBulk(member)
		}
		return respArray(items)
	case "EVALSHA":
		return "-NOSCRIPT No matching script. Please use EVAL.\r\n"
	case "EVAL":
		return f.handleCompareAndDelete(args)
	}
	return fmt.Sprintf("-ERR unknown command '%s'\r\n", args[0])
}

func (f *fakeValkey) handleSet(args []string) string {
	key, value := args[1], args[2]
	var ttl time.Duration
	onlyIfAbsent := false
	for i := 3; i < len(args); i++ {
		switch strings.ToUpper(args[i]) {
		case "EX":
			seconds, _ := strconv.Atoi(args[i+1])
			ttl = time.Duration(seconds) * time.Second
			i++
		case "PX":
			millis, _ := strconv.Atoi(args[i+1])
			ttl = time.Duration(millis) * time.Millisecond
			i++
		case "NX":
			onlyIfAbsent = true
		}
	}
	if _, exists := f.values[key]; exists && onlyIfAbsent {
		return respNil
	}
	f.values[key] = value
	f.ttls[key] = ttl
	return respOK
}

func (f *fakeValkey) deleteKeys(keys []string) int {
	n := 0
	for _, key := range keys {
		if _, ok := f.values[key]; ok {
			delete(f.values, key)
			delete(f.ttls, key)
			n++
		}
	}
	return n
}

func (f *fakeValkey) handleCompareAndDelete(args []string) string {
	if len(args) != 5 || args[2] != "1" {
		return "-ERR fake valkey only runs the single-key compare-and-delete script\r\n"
	}
	key, owner := args[3], args[4]
	if f.values[key] != owner {
		return respInt(0)
	}
	return respInt(f.deleteKeys([]string{key}))
}
