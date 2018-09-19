package core

import (
	"bytes"
	"fmt"
	"math"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Beamer support Redis stats collection
type Beamer struct {
	mutex  sync.RWMutex
	labels string

	scrapeCounter  int64
	scrapeSkiped   int64
	scrapeFailures int64
}

// NewBeamer create a beamer
func NewBeamer(exporters []*Exporter, labels map[string]string) *Beamer {
	delta := viper.GetInt("scanDuration") / len(exporters)
	p := math.Max(float64(delta), 1)
	ticker := time.NewTicker(time.Duration(p) * time.Millisecond)
	running := make(chan struct{}, viper.GetInt("maxConcurrent"))
	i := 0

	b := &Beamer{}

	go func() {
		for {
			select {
			case <-ticker.C:
				select {
				case running <- struct{}{}:
					go func() {
						defer func() {
							<-running
						}()
						e := exporters[i]
						success := e.Scrape()

						b.mutex.Lock()

						if !success {
							b.scrapeFailures++
							log.Errorf("Scrape fail for %v", e.Address)
						}
						b.mutex.Unlock()
					}()

					b.scrapeCounter++
					i++
					if i >= len(exporters) {
						i = 0
					}
				default:
					b.mutex.Lock()
					b.scrapeSkiped++
					b.mutex.Unlock()
				}
			}
		}
	}()

	for k := range labels {
		if len(b.labels) > 0 {
			b.labels += ","
		}
		b.labels += k + "=" + labels[k]
	}

	return b
}

// Metrics delivers beamer stats as warp10 metrics.
func (b *Beamer) Metrics() *bytes.Buffer {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	var buf bytes.Buffer

	head := fmt.Sprintf("%v// redis.exporter.", time.Now().UnixNano()/1000)

	buf.WriteString(fmt.Sprintf("%vscrape{%v} %v\n", head, b.labels, b.scrapeCounter))
	buf.WriteString(fmt.Sprintf("%vscrape_skiped{%v} %v\n", head, b.labels, b.scrapeSkiped))
	buf.WriteString(fmt.Sprintf("%vscrape_failures{%v} %v\n", head, b.labels, b.scrapeFailures))

	return &buf
}
