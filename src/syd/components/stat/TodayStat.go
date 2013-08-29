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
	p.Stats = statdao.TodayStat(7)
}

func (p *TodayStat) ShowDate(diff int) time.Time {
	t := time.Now()
	t = t.AddDate(0, 0, diff)
	return t
}
