###
  SYD System
  @author: Bo Gao, [elivoa@gmail.com]
###

################################################################################
# onload
################################################################################

class ProductEdit
  constructor: ->
    this.init()

  init: ->
    console.log "init ProductEdit page!!!"

    # enable multiline table
    colorMT = new EditableTable("color-area")
    colorMT.onRemoveLastLine = ->
      alert "至少需要保留一个颜色！"

    sizeMT = new EditableTable("size-area")
    sizeMT.onRemoveLastLine = ->
      alert "至少需要保留一个尺码！"

# ____________________________________________________________
# Register
# This should be called in page html.
#
$ ->
  window['ProductEdit'] = new ProductEdit

