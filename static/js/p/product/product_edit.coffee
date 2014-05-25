###
  SYD System
  @author: Bo Gao, [elivoa@gmail.com]
###

################################################################################
# onload
################################################################################

class ProductEdit
  constructor: (stockcache) ->
    @colorId = "color-area"
    @sizeId = "size-area"
    @csqTableId = "cs-container"


    # eval init json
    # Stores quantity, used to restore quantity after each CSTable refresh. TODO
    if stockcache != undefined
      # parse json
      @stockcache = JSON && JSON.parse(stockcache) || $.parseJSON(stockcache)
    else
      @stockcache = {}

    this.init()

  init: ->
    ## Enable Color multiline table
    colorMT = new EditableTable(@colorId)
    colorMT.onRemoveLastLine = (line) ->
      # alert "至少需要保留一个颜色！"
      # empty all input values
      $(line).find("input").each (idx, obj) ->
        $(obj).attr "value",""
        $(obj).val ""
    colorMT.afterRemove = $.proxy @onCSQTableRefresh,@

    ## Enable Size multiline table
    sizeMT = new EditableTable(@sizeId)
    sizeMT.onRemoveLastLine = ->
      $(line).find("input").each (idx, obj) ->
        $(obj).attr "value",""
        $(obj).val ""
    sizeMT.afterRemove = $.proxy @onCSQTableRefresh,@

    ## Enable Color-size quantity table generator trigger
    ## class is: csq-trigger
    $('body').on "blur", ".csq-trigger", $.proxy @onCSQTableRefresh,@

    ## First time draw csqTable, use server side stock values.
    @firstTimeDrawCSQTable()

  firstTimeDrawCSQTable: ->
    # TODO init stockcache with server values, ajax?
    # restore cs-table
    [@colors, @sizes] = @readColorSizes()
    pcstg = new ProductCSTableGenerator(@colors, @sizes, @csqTableId)
    # restore stock number.
    @fillProductQuantity()

  onCSQTableRefresh: ->
    [@colors, @sizes] = @readColorSizes()
    @cacheStock()
    pcstg = new ProductCSTableGenerator(@colors, @sizes, @csqTableId)
    [@colors, @sizes] = @readColorSizes()
    @fillProductQuantity()

  readColorSizes:->
    colors = []
    sizes = []
    $("##{@colorId} .mt-line .csq-trigger").each (idx,obj)->
      v = obj.value.replace /^\s+|\s+$/g, ""
      colors.push v if v != ""
    $("##{@sizeId} .mt-line .csq-trigger").each (idx,obj)->
      v = obj.value.replace /^\s+|\s+$/g, ""
      sizes.push v if v != ""
    return [colors, sizes]

#_______________________________________________________________________________
  cacheStock: ->
    for color in @colors
      for size in @sizes
        key = "#{color}__#{size}"
        o = $("##{@csqTableId} #csq_#{key}")
        if o != undefined
          @stockcache[key] = o.val()

  fillProductQuantity: ->
    for k, v of @stockcache
      o = $("##{@csqTableId} #csq_#{k}")
      if o != undefined
        o.val v

##  ____________________________________________________________
## Register
## This should be called in page html.
##
$ ->
  window['ProductEdit'] = ProductEdit
