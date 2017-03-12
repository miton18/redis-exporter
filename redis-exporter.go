// REDIS exporter expose REDIS stats.
//
// Usage
//
// 		redis-exporter  [flags]
// Flags:
//       --config string   config file to use
//       --help            display help
//   -v, --verbose         verbose output
package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/miton18/redis-exporter/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Panicf("%v", err)
	}
}
