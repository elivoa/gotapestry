package layout

import (
	"got/core"
)

type GOTHeader struct {
	core.Component
}

func (c *GOTHeader) Setup()        {}
func (c *GOTHeader) BeforeRender() {}
func (c *GOTHeader) AfterRender()  {}
