package core

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// Exporter TODO
type Exporter struct {
	Address   string
	metrics   map[string]bool
	mutex     sync.RWMutex
	labels    string // format {key=val,a=b}
	sensision bytes.Buffer
}

// NewExporter todo
func NewExporter(address string, labels map[string]string, metrics []string) (*Exporter, error) {

	e := &Exporter{
		Address: address,
		// default exported metrics
		metrics: map[string]bool{
			"redis_version":                false,
			"redis_git_sha1":               false,
			"redis_git_dirty":              false,
			"redis_build_id":               false,
			"redis_mode":                   false,
			"os":                           false,
			"arch_bits":                    false,
			"multiplexing_api":             false,
			"gcc_version":                  false,
			"process_id":                   false,
			"run_id":                       false,
			"tcp_port":                     false,
			"uptime_in_seconds":            true,
			"uptime_in_days":               false,
			"hz":                           false,
			"lru_clock":                    false,
			"executable":                   false,
			"config_file":                  false,
			"connected_clients":            true,
			"client_longest_output_list":   true,
			"client_biggest_input_buf":     true,
			"blocked_clients":              true,
			"used_memory":                  true,
			"used_memory_human":            false,
			"used_memory_rss":              true,
			"used_memory_rss_human":        false,
			"used_memory_peak":             true,
			"used_memory_peak_human":       false,
			"total_system_memory":          true,
			"total_system_memory_human":    false,
			"used_memory_lua":              true,
			"used_memory_lua_human":        false,
			"maxmemory":                    true,
			"maxmemory_human":              false,
			"maxmemory_policy":             false,
			"mem_fragmentation_ratio":      true,
			"mem_allocator":                false,
			"loading":                      true,
			"rdb_changes_since_last_save":  true,
			"rdb_bgsave_in_progress":       true,
			"rdb_last_save_time":           true,
			"rdb_last_bgsave_status":       true,
			"rdb_last_bgsave_time_sec":     true,
			"rdb_current_bgsave_time_sec":  true,
			"aof_enabled":                  true,
			"aof_rewrite_in_progress":      true,
			"aof_rewrite_scheduled":        true,
			"aof_last_rewrite_time_sec":    true,
			"aof_current_rewrite_time_sec": true,
			"aof_last_bgrewrite_status":    true,
			"aof_last_write_status":        true,
			"total_connections_received":   true,
			"total_commands_processed":     true,
			"instantaneous_ops_per_sec":    true,
			"total_net_input_bytes":        true,
			"total_net_output_bytes":       true,
			"instantaneous_input_kbps":     true,
			"instantaneous_output_kbps":    true,
			"rejected_connections":         true,
			"sync_full":                    true,
			"sync_partial_ok":              true,
			"sync_partial_err":             true,
			"expired_keys":                 true,
			"evicted_keys":                 true,
			"keyspace_hits":                true,
			"keyspace_misses":              true,
			"pubsub_channels":              true,
			"pubsub_patterns":              true,
			"latest_fork_usec":             true,
			"migrate_cached_sockets":       true,
			"role":                           false,
			"connected_slaves":               true,
			"master_repl_offset":             true,
			"repl_backlog_active":            true,
			"repl_backlog_size":              true,
			"repl_backlog_first_byte_offset": true,
			"repl_backlog_histlen":           true,
			"used_cpu_sys":                   true,
			"used_cpu_user":                  true,
			"used_cpu_sys_children":          true,
			"used_cpu_user_children":         true,
			"cluster_enabled":                true,
		},
	}

	// filter
	if len(metrics) > 0 {
		// some metrics are white listed, let's black list all of them
		for metric := range e.metrics {
			e.metrics[metric] = false
		}
		// and allow some of them
		for _, metric := range metrics {
			e.metrics[metric] = true
		}
	}

	// format labels
	labelsArray := make([]string, len(labels))
	i := 0
	for k, v := range labels {
		labelsArray[i] = k + "=" + v
		i++
	}
	e.labels = "{" + strings.Join(labelsArray, ",") + "}"

	return e, nil
}

// Metrics sensision format of metrics
func (e *Exporter) Metrics() *bytes.Buffer {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return bytes.NewBuffer(e.sensision.Bytes())
}

func (e *Exporter) clear() {
	// protect consistency
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.sensision.Reset()
}

// Scrape load sensision buffer with metrics
func (e *Exporter) Scrape() bool {
	buf, err := e.readInfo()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	err = e.parse(buf)
	if err != nil {
		fmt.Println("Error parsing response", err.Error())
	}
	return true
}

func (e *Exporter) parse(infos *bytes.Buffer) error {
	infosReader := bufio.NewReader(infos)
	now := fmt.Sprintf("%v// redis_stats", time.Now().UnixNano()/1000)
	var keyVal []string

	e.clear()
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for {
		line, _, err := infosReader.ReadLine()
		if err != nil {
			break
		}
		if !bytes.HasPrefix(line, []byte("#")) {
			keyVal = strings.Split(string(bytes.TrimSpace(line)), ":")
			if len(keyVal) == 2 {
				// if filtered metric
				if e.metrics[keyVal[0]] {
					if keyVal[1] == "ok" {
						keyVal[1] = "true"
					} else if keyVal[1] == "ko" {
						keyVal[1] = "false"
					}
					e.sensision.WriteString(now + keyVal[0] + e.labels + " " + keyVal[1] + "\n")
				}
			}
		}
	}
	infos.Reset()

	return nil
}

func (e *Exporter) readInfo() (*bytes.Buffer, error) {
	var empty bytes.Buffer
	conn, err := net.Dial("tcp", e.Address)
	if err != nil {
		conn = nil
		return &empty, err
	}
	defer conn.Close()

	conn.Write([]byte("INFO\r\n"))
	b := make([]byte, 4096)
	_, err = conn.Read(b)
	if err != nil {
		return &empty, err
	}
	return bytes.NewBuffer(b), nil
}
