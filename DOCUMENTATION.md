gotapestry
==========

A Tapestry5 like WebFramework writen in GoLang

Note: This is a very early version that only contains some of the core features.

!! Don't use this in production environment. !!



Documentation
==============

inject type

## go code
type OrderCreateIndex struct {
	core.Page
	Idinpath int `param:"."`
    Id int `query:"id"`
}

