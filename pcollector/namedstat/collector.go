// named stats collector

package namedstat

import (
	"log/slog"
	"math"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/itoolkits/toolkit/dnt"
)

type StatsCollector struct {
	rndc     string
	statFile string
}

// NewStatsCollector create collector
func NewStatsCollector(rndc string, statFile string) *StatsCollector {
	return &StatsCollector{
		rndc:     rndc,
		statFile: statFile,
	}
}

// Describe implements prometheus.Collector.
func (c *StatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	//ch <- bootTime
	ch <- nameServerStatistics
	ch <- incomingQueries
	ch <- incomingRequests
	ch <- incomingRequests
	ch <- socketIO
	ch <- zoneMetricStats
	ch <- cacheRRsetsStats
	ch <- cacheStatistics

	for _, desc := range resolverMetricStatsFile {
		ch <- desc
	}
	for _, desc := range cacheMetricStatsFile {
		ch <- desc
	}
}

// Collect implements prometheus.Collector.
func (c *StatsCollector) Collect(ch chan<- prometheus.Metric) {
	sf := dnt.StatsFile{
		RNDC: c.rndc,
		Path: c.statFile,
	}
	out, err := sf.Build()
	if err != nil {
		slog.Error("bind exporter build stats file error", "error", err)
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
		return
	}
	slog.Info("bind exporter build stats file success", "output", out)

	statsInfo, err := sf.Parse()
	if err != nil {
		slog.Error("bind exporter parse stats file error", "error", err)
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
		return
	}
	//ch <- prometheus.MustNewConstMetric(
	//	bootTime, prometheus.GaugeValue, float64(statsInfo.BootTime),
	//)
	if mds, ok := statsInfo.SubMetric["Incoming Requests"]; ok {
		for _, md := range mds {
			for key, value := range md.Metric {
				ch <- prometheus.MustNewConstMetric(
					incomingRequests, prometheus.GaugeValue, value, key,
				)
			}
		}
	}
	if mds, ok := statsInfo.SubMetric["Incoming Queries"]; ok {
		for _, md := range mds {
			for key, value := range md.Metric {
				ch <- prometheus.MustNewConstMetric(
					incomingQueries, prometheus.CounterValue, value, key,
				)
			}
		}
	}
	// Name Server Statistics
	if mds, ok := statsInfo.SubMetric["Name Server Statistics"]; ok {
		for _, md := range mds {
			for key, value := range md.Metric {
				if tk, kok := nameServerMap[key]; kok {
					ch <- prometheus.MustNewConstMetric(
						nameServerStatistics, prometheus.CounterValue, value, tk,
					)
				}
			}
		}
	}
	// Outgoing Queries
	if mds, ok := statsInfo.SubMetric["Outgoing Queries"]; ok {
		for _, md := range mds {
			for key, value := range md.Metric {
				ch <- prometheus.MustNewConstMetric(
					outgoingQueries, prometheus.CounterValue, value, md.View, key,
				)
			}
		}
	}
	// Resolver Statistics  resolverQueries
	if mds, ok := statsInfo.SubMetric["Resolver Statistics"]; ok {
		for _, md := range mds {
			for key, value := range md.Metric {
				if rk, kok := resolverStatisticsMap[key]; kok {
					if rk == "QryRTTnn" {
						continue
					}
					if pd, pok := resolverMetricStatsFile[rk]; pok {
						ch <- prometheus.MustNewConstMetric(
							pd, prometheus.CounterValue, value, md.View,
						)
					}
				}
			}
			if pd, pok := resolverMetricStatsFile["QryRTTnn"]; pok {
				if buckets, count, err := getHistogram(md); err == nil {
					ch <- prometheus.MustNewConstHistogram(
						pd, count, math.NaN(), buckets, md.View,
					)
				}
			}
		}
	}

	// Socket IO Statistics
	if mds, ok := statsInfo.SubMetric["Socket IO Statistics"]; ok {
		for _, md := range mds {
			for key, value := range md.Metric {
				if tk, kok := socketMap[key]; kok {
					ch <- prometheus.MustNewConstMetric(
						socketIO, prometheus.CounterValue, value, tk,
					)
				}
			}
		}
	}
	// Zone Maintenance Statistics
	if mds, ok := statsInfo.SubMetric["Zone Maintenance Statistics"]; ok {
		for _, md := range mds {
			for key, value := range md.Metric {
				if tk, kok := zoneMap[key]; kok {
					ch <- prometheus.MustNewConstMetric(
						zoneMetricStats, prometheus.CounterValue, value, tk,
					)
				}
			}
		}
	}
	// Cache DB RRsets
	if mds, ok := statsInfo.SubMetric["Cache DB RRsets"]; ok {
		for _, md := range mds {
			for key, value := range md.Metric {
				idx := strings.Index(md.View, "(")
				if idx < 0 {
					idx = len(md.View)
				}
				view := strings.TrimSpace(md.View[0:idx])
				ch <- prometheus.MustNewConstMetric(
					cacheRRsetsStats, prometheus.CounterValue, value, view, key,
				)
			}
		}
	}
	// Cache Statistics
	if mds, ok := statsInfo.SubMetric["Cache Statistics"]; ok {
		for _, md := range mds {
			for key, value := range md.Metric {
				if rk, kok := cacheStatsMap[key]; kok {
					if pd, pok := cacheMetricStatsFile[rk]; pok {
						ch <- prometheus.MustNewConstMetric(
							pd, prometheus.GaugeValue, value, md.View,
						)
					}
				}
			}
		}
	}
	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)
}

type RTTHistog struct {
	key string
	le  float64
}

var (
	rttList = []RTTHistog{
		{"queries with RTT < 10ms", 10},
		{"queries with RTT 10-100ms", 100},
		{"queries with RTT 100-500ms", 500},
		{"queries with RTT 500-800ms", 800},
		{"queries with RTT 800-1600ms", 1600},
		{"queries with RTT > 1600ms", 2000},
	}
)

// getHistogram get rtt histogram
func getHistogram(md *dnt.StatsViewMetric) (map[float64]uint64, uint64, error) {
	buckets := map[float64]uint64{}
	var count uint64
	for _, rtt := range rttList {
		if value, ok := md.Metric[rtt.key]; ok {
			buckets[rtt.le] = count + uint64(value)
			count += uint64(value)
		}
	}
	return buckets, count, nil
}
