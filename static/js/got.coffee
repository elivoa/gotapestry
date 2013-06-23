###
  SYD Sales System
  @author: Bo Gao, [elivoa@gmail.com]
###

# global bindings, values should be element's id.

################################################################################
# onload
################################################################################
class Gotinit
  constructor: ->
    console.log 'got initialized'

  initSelect: ->
    # init select value
    $('select').each (index, selectObj)->
      target = $ selectObj
      value = target.attr('value')
      if value
        target.find("option").each (idx, option)->
          if value == option.value
            option.selected=true
            return true
  initRadio: ->
    # init
    $('.auto-radio').each (idx, obj) ->
      v = $(obj).attr("value")
      $(obj).find('>input[type="radio"]').each (i, r)->
        if $(r).attr('value') == v
          $(r).prop('checked', true)
    return 1

#---------------------------------------------------------------------------------
$ ->
  init = new Gotinit
  init.initSelect()
  init.initRadio()
  return true
