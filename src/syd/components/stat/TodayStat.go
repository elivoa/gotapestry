package stat

import (
	"got/core"
	"syd/dal/statdao"
	"syd/model"
	"time"
)

type TodayStat struct {
	core.Component
	Stats     []*model.SumStat
	ShowChart bool `default:"true"`
}

func (p *TodayStat) New() *TodayStat {
	return &TodayStat{ShowChart: true}
}

func (p *TodayStat) Setup() {
	stats, err := statdao.TodayStat(7)
	if err != nil {
		panic(err.Error())
	}
	p.Stats = stats
}

func (p *TodayStat) ShowDate(diff int) time.Time {
	t := time.Now()
	t = t.AddDate(0, 0, diff)
	return t
}
