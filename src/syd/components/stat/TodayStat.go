package stat

import (
	"github.com/elivoa/got/core"
	"syd/dal/statdao"
	"syd/model"
	"time"
)

type TodayStat struct {
	core.Component
	Stats         []*model.SumStat
	LastYearStats map[int]*model.SumStat
	ShowChart     bool `default:"true"`
	EmptyStats    *model.SumStat
}

func (p *TodayStat) New() *TodayStat {
	return &TodayStat{ShowChart: true}
}

func (p *TodayStat) Setup() {
	now := time.Now()
	stats, err := statdao.TodayStat(now, 7)
	if err != nil {
		panic(err.Error())
	}
	// for _, ss := range stats {
	// 	fmt.Println("========== ", ss)
	// }

	statslastyear, err2 := statdao.TodayStat(now.AddDate(-1, 0, 0), 7)
	p.LastYearStats = map[int]*model.SumStat{}
	if err2 != nil {
		panic(err.Error())
	} else if nil != statslastyear {
		for _, ss := range statslastyear {
			// fmt.Println("======---- ", ss)
			p.LastYearStats[ss.Id] = ss
		}
	}

	p.Stats = stats
}

func (p *TodayStat) Yestoday(id int) *model.SumStat {
	s := p.LastYearStats[id]
	if s == nil {
		s = model.EmptySumStat
	}
	return s
}

func (p *TodayStat) ShowDate(diff int) time.Time {
	t := time.Now()
	t = t.AddDate(0, 0, diff+1)
	return t
}
