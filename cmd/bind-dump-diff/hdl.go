// bind dump db diff handler

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/itoolkits/toolkit/collect"
	"github.com/itoolkits/toolkit/dnt"
)

type ZoneDiff struct {
	A []string
	B []string
}

type DiffHandler struct {
	a map[string]map[string][]*dnt.RR
	b map[string]map[string][]*dnt.RR

	as map[string]map[string]*collect.HashSet[string]
	bs map[string]map[string]*collect.HashSet[string]

	zSet *collect.HashSet[string]
	vSet *collect.HashSet[string]

	aNum int
	bNum int

	diffSet *collect.HashSet[string]
	diffMap map[string]*ZoneDiff

	diffRNum int
}

// NewDiffHandler create diff handler
func NewDiffHandler(a, b map[string]map[string][]*dnt.RR) *DiffHandler {
	h := &DiffHandler{
		a:  a,
		b:  b,
		as: make(map[string]map[string]*collect.HashSet[string]),
		bs: make(map[string]map[string]*collect.HashSet[string]),

		diffSet: collect.NewHashSetAllowNilVal[string](),
		diffMap: make(map[string]*ZoneDiff),

		zSet: collect.NewHashSetAllowNilVal[string](),
		vSet: collect.NewHashSetAllowNilVal[string](),

		diffRNum: 0,
	}
	return h
}

// Start diff handler
func (h *DiffHandler) Start() {
	h.aNum = h.convertSet(h.a, h.as)
	h.bNum = h.convertSet(h.b, h.bs)

	h.cut(h.as, h.bs)
	h.cut(h.bs, h.as)

	h.collectDiff()
}

// convertSet convert to set
func (h *DiffHandler) convertSet(
	zvrList map[string]map[string][]*dnt.RR,
	rst ...map[string]map[string]*collect.HashSet[string]) int {
	num := 0
	for z, vrList := range zvrList {
		vrSet := make(map[string]*collect.HashSet[string])
		h.zSet.Add(z)
		for v, rList := range vrList {
			if v == "" {
				v = "None"
			}

			h.vSet.Add(v)

			rSet := collect.NewHashSet[string]()
			for _, rr := range rList {
				rSet.Add(Record2Str(rr))
				num++
			}
			vrSet[v] = rSet
		}
		for _, ele := range rst {
			ele[z] = vrSet
		}
	}
	return num
}

// Cut record each other
func (h *DiffHandler) cut(A, B map[string]map[string]*collect.HashSet[string]) {
	for az, avrSet := range A {
		bvrSet, h := B[az]
		if !h {
			continue
		}
		for av, arSet := range avrSet {
			brSet, h := bvrSet[av]
			if !h {
				continue
			}
			bak := collect.NewHashSetByHashSet(arSet)
			arSet.RemoveHashSet(brSet)
			brSet.RemoveHashSet(bak)
			if arSet.Size() <= 0 {
				delete(avrSet, av)
			}
			if brSet.Size() <= 0 {
				delete(bvrSet, av)
			}
		}
		if len(avrSet) <= 0 {
			delete(A, az)
		}
		if len(bvrSet) <= 0 {
			delete(B, az)
		}
	}
}

// collectDiff collect diff record
func (h *DiffHandler) collectDiff() {
	h.zSet.Range(func(z string) bool {
		avrSet, ah := h.as[z]
		bvrSet, bh := h.bs[z]
		if len(avrSet) <= 0 {
			ah = false
		}
		if len(bvrSet) <= 0 {
			bh = false
		}
		if !ah && !bh {
			return true
		}

		if ah && !bh {
			num := 0
			for _, ele := range avrSet {
				num += ele.Size()
			}
			k := fmt.Sprintf("Zone: %s", z)
			h.diffSet.Add(k)
			h.diffMap[k] = &ZoneDiff{
				A: []string{fmt.Sprintf("Has %d Records", num)},
				B: []string{"Zone Not Exist"},
			}

			h.diffRNum += num
			return true
		}
		if !ah && bh {
			num := 0
			for _, ele := range bvrSet {
				num += ele.Size()
			}
			k := fmt.Sprintf("Zone: %s", z)
			h.diffSet.Add(k)
			h.diffMap[k] = &ZoneDiff{
				A: []string{"Zone Not Exist"},
				B: []string{fmt.Sprintf("Has %d Records", num)},
			}
			h.diffRNum += num
			return true
		}

		h.vSet.Range(func(v string) bool {
			arSet, arh := avrSet[v]
			brSet, brh := bvrSet[v]
			if arSet == nil || arSet.Size() <= 0 {
				arh = false
			}
			if brSet == nil || brSet.Size() <= 0 {
				brh = false
			}
			if !arh && !brh {
				return true
			}
			if arh && !brh {
				k := fmt.Sprintf("Zone: %s, View: %s", z, v)
				h.diffSet.Add(k)
				h.diffMap[k] = &ZoneDiff{
					A: []string{fmt.Sprintf("Has %d Records", arSet.Size())},
					B: []string{"View Not Exist"},
				}
				h.diffRNum += arSet.Size()
				return true
			}
			if !arh && brh {
				k := fmt.Sprintf("Zone: %s, View: %s", z, v)
				h.diffSet.Add(k)
				h.diffMap[k] = &ZoneDiff{
					A: []string{"View Not Exist"},
					B: []string{fmt.Sprintf("Has %d Records", brSet.Size())},
				}
				h.diffRNum += brSet.Size()
				return true
			}
			arList := arSet.ToSlice()
			brList := brSet.ToSlice()
			sort.Strings(arList)
			sort.Strings(brList)

			k := fmt.Sprintf("Zone: %s, View: %s", z, v)
			h.diffSet.Add(k)
			h.diffMap[k] = &ZoneDiff{
				A: arList,
				B: brList,
			}
			h.diffRNum += len(arList) + len(brList)
			return true
		})

		return true
	})
}

// PrintDiff print diff
func (h *DiffHandler) PrintDiff(a, b string) {
	rows := make([]table.Row, 0)

	zvList := h.diffSet.ToSlice()
	sort.Strings(zvList)
	for i, zv := range zvList {
		diffList := h.diffMap[zv]
		rows = append(rows, table.Row{i + 1, zv,
			strings.Join(diffList.A, "\n"), strings.Join(diffList.B, "\n")})
		rows = append(rows, table.Row{"", "", "", ""})
	}

	tbl := table.NewWriter()
	tbl.SetOutputMirror(os.Stdout)
	tbl.AppendHeader(table.Row{"#", "Zone & View", a, b})
	tbl.AppendRows(rows)
	tbl.AppendSeparator()
	tbl.SetStyle(table.StyleDefault)
	tbl.SetColumnConfigs([]table.ColumnConfig{
		//{Name: "#", Colors: text.Colors{text.FgHiYellow}, ColorsHeader: text.Colors{text.BgWhite, text.FgHiYellow}},
		{Name: "Zone & View", Colors: text.Colors{text.FgHiMagenta}, ColorsHeader: text.Colors{text.FgHiMagenta}},
		{Name: a, Colors: text.Colors{text.FgHiYellow}, ColorsHeader: text.Colors{text.FgHiYellow}},
		{Name: b, Colors: text.Colors{text.FgHiCyan}, ColorsHeader: text.Colors{text.FgHiCyan}},
	})
	tbl.AppendFooter(table.Row{"SUMMARY",
		fmt.Sprintf("ZONE DIF:%d, RECORD DIFF:%d", h.diffSet.Size(), h.diffRNum),
		fmt.Sprintf("CHECK NU:%d", h.aNum),
		fmt.Sprintf("CHECK NU:%d", h.bNum)})
	tbl.Render()
	fmt.Println()
}

// Record2Str record to string
func Record2Str(r *dnt.RR) string {
	return fmt.Sprintf("%s %s %s", r.Domain, r.RType, r.RData)
}
