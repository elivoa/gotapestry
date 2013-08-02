GOT Web Framework
==================
  Tapestry on GoLang

> A Tapestry5 like WebFramework writen in GoLang<br>
> Note: This is a very early version that only contains some of the core features.<br>
> !! Don't use this in production environment. !!<br>



Documentation
==============

## Page Lifecircle

### Page value Injection
+ TODO …

### Page events::RenderPage
+ Activate()
  - Called every time this page object is Activated. Page render or call event on this page.<br>
  - ? Take no parameters.
+ Setup() or SetupRender()
+ After() or AfterRender()

### Page event::FormSubmit
+ Activate()
+ OnSubmit() or OnSubmitFromTableTID()
  - before inject submit values, init fields.
+ ?? OnValidate()
  - TODO… 
+ OnSuccess() or OnSuccessFromTalbeTID()
  - after inject values, do submit.
+ ?? OnValidationFailed()
+ ?? OnError()

### Form Submit















~~ laji ~~
-------------------------------------------------------------------------



inject type

## go code
type OrderCreateIndex struct {
	core.Page
	Idinpath int `param:"."`
    Id int `query:"id"`
}


#####