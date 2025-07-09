// dns stats file build and parse
// version: BIND 9.16.39
// refer: https://github.com/qiangmzsx/bind_stats_exporter

package dnt

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/itoolkits/toolkit/cmdt"
	"github.com/itoolkits/toolkit/filet"
)

var (
	subReg, _  = regexp.Compile(" ?\\+\\+ ?")
	viewReg, _ = regexp.Compile("[\\(|\\)|\\<|/]")
	numReg, _  = regexp.Compile(`[0-9]+`)
)

type StatsViewMetric struct {
	View   string             `json:"view"`
	Metric map[string]float64 `json:"metric"`
}

type StatsMetric struct {
	StatTimestamp int64                         `json:"stat_timestamp"`
	SubMetric     map[string][]*StatsViewMetric `json:"sub_metric"`
}

type StatsFile struct {
	RNDC string
	Path string
}

// Build stats file
func (r *StatsFile) Build() ([]string, error) {
	if err := filet.TruncFile(r.Path); err != nil {
		return nil, err
	}
	cmd := cmdt.NewCommand(time.Second * 10)
	return cmd.Do(r.RNDC, "stats")
}

// Parse named stats file
func (r *StatsFile) Parse() (*StatsMetric, error) {
	file, err := os.Open(r.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sub := ""
	view := ""
	stats := &StatsMetric{
		SubMetric: map[string][]*StatsViewMetric{},
	}

	ts := make([]string, 0)

	sm := &StatsViewMetric{
		Metric: make(map[string]float64),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "+++") {
			ts = numReg.FindAllString(line, -1)
			continue
		}
		if strings.HasPrefix(line, "---") {
			break
		}
		if strings.HasPrefix(line, "++") {
			if len(sm.Metric) > 0 || len(sm.View) > 0 {
				stats.SubMetric[sub] = append(stats.SubMetric[sub], sm)
				sm = &StatsViewMetric{
					Metric: make(map[string]float64),
					View:   "",
				}
			}
			sm.View = ""
			sub = subReg.ReplaceAllString(line, "")
			sub = viewReg.ReplaceAllString(sub, "")
			continue
		}
		if strings.HasPrefix(line, "[") {
			view = strings.ReplaceAll(line, "[", "")
			view = strings.ReplaceAll(view, "View:", "")
			view = strings.ReplaceAll(view, "]", "")
			view = strings.TrimSpace(view)
			if len(sm.Metric) > 0 {
				stats.SubMetric[sub] = append(stats.SubMetric[sub], sm)
			}
			sm = &StatsViewMetric{
				Metric: make(map[string]float64),
				View:   view,
			}
			continue
		}
		seg := strings.Fields(line)
		if len(seg) < 2 {
			continue
		}
		num := seg[0]
		di := strings.Join(seg[1:], " ")

		v, err := strconv.ParseFloat(num, 10)
		if err != nil {
			return nil, fmt.Errorf("parse metric value error, %s, %w", line, err)
		}

		sm.Metric[di] = v
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(ts) > 0 {
		ti, err := strconv.ParseInt(ts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse boot timestamp error, %s, %w", ts, err)
		}
		stats.StatTimestamp = ti
	}
	return stats, nil
}
