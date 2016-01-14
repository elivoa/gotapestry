package stat

import (
	"fmt"
	"github.com/elivoa/got/core"
	"strings"
	"syd/dal/statdao"
	"syd/model"
	"time"
)

type TodayStat struct {
	core.Component
	Stats         []*model.SumStat
	LastYearStats map[string]*model.SumStat
	ShowChart     bool `default:"true"`
	EmptyStats    *model.SumStat

	// const
	Today_Key string
}

func (p *TodayStat) New() *TodayStat {
	return &TodayStat{
		ShowChart: true,
		Today_Key: time.Now().Format("2006-01-02"),
	}
}

var (
	current_year = time.Now().AddDate(0, 0, 0).Format("2006")
	last_year    = time.Now().AddDate(-1, 0, 0).Format("2006")
)

func (p *TodayStat) Setup() {
	var debug_print_time = true

	var see_future_days = 1
	var show_days = 7
	var fetch_days = see_future_days + show_days + 1

	now := time.Now()
	stats, err := statdao.TodayStat(now.AddDate(0, 0, see_future_days), fetch_days)
	if err != nil {
		panic(err.Error())
	}

	// fill days
	p.Stats = []*model.SumStat{}
	for i := -show_days + 1; i <= see_future_days; i++ {
		key := now.AddDate(0, 0, i).Format("2006-01-02")
		var stat *model.SumStat
		for _, k := range stats {
			if k.Id == key {
				stat = k
			}
		}
		if nil == stat {
			stat = &model.SumStat{Id: key}
		}
		p.Stats = append(p.Stats, stat)
	}

	if debug_print_time {
		for _, k := range p.Stats {
			fmt.Println(">> p.Stats: ", k)
		}
	}

	if debug_print_time {
		for _, k := range p.Stats {
			fmt.Println(">> p.Stats: ", k)
		}
	}

	// 去年
	statslastyear, err2 := statdao.TodayStat(now.AddDate(-1, 0, see_future_days), fetch_days)
	p.LastYearStats = map[string]*model.SumStat{}
	if err2 != nil {
		panic(err.Error())
	} else if nil != statslastyear {
		for _, ss := range statslastyear {
			key := strings.Replace(ss.Id, last_year, current_year, -1)
			p.LastYearStats[key] = ss
		}
	}

	// p.Stats = stats
}

func (p *TodayStat) Yestoday(id string) *model.SumStat {

	s := p.LastYearStats[id]
	if s == nil {
		s = model.EmptySumStat
	}
	return s
}

func (p *TodayStat) LineClass(key string) string {
	if key == p.Today_Key {
		return "today"
	}
	return ""
}

func (p *TodayStat) DateLabel(key string) string {
	if key == p.Today_Key {
		return "今天"
	}
	return key
}

func (p *TodayStat) ShowDate(diff int) time.Time {
	t := time.Now()
	t = t.AddDate(0, 0, diff+1)
	return t
}
