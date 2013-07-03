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
    # console.log "init ProductEdit page!!!"

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
    pcstg = new ProductColorSizeTableGenerator(@colors, @sizes)
    pcstg.replace(@csqTableId)
    # restore stock number.
    console.log @stockcache
    @fillProductQuantity()

  onCSQTableRefresh: ->
    [@colors, @sizes] = @readColorSizes()
    @cacheStock()
    pcstg = new ProductColorSizeTableGenerator(@colors, @sizes)
    pcstg.replace(@csqTableId)
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

###### backup ########
  cacheStock___backup: ->
    stocks = []
    $("##{@csqTableId} .stock").each (idx, obj)->
      stocks.push obj.value

    i=0
    for color in @colors
      for size in @sizes
        @stockcache["#{color}__#{size}"] = stocks[i++]

  #### ______________________________________
  fillProductQuantity_backup_badversion: ->
    # use map to cache quantity and restore
    stockInputs = $("##{@csqTableId} .stock")
    console.log "--------------------------------------------------------------------------------"
    console.log stockInputs
    console.log ">>> ", @stockcache
    i=0
    for color in @colors
      for size in @sizes
        console.log ">>> Set #{i}'th #{@stockcache["#{color}__#{size}"]}"
        value = @stockcache["#{color}__#{size}"]
        if value == undefined
          stockInputs[i++].value = ""
        else
          stockInputs[i++].value = value




##  ____________________________________________________________
## Register
## This should be called in page html.
##
$ ->
  window['ProductEdit'] = ProductEdit



################################################################################
# product color-size table generator
################################################################################

class ProductColorSizeTableGenerator
  constructor: (colors, sizes) ->

    # clear parameters
    if colors.length == 0 then @colors=["默认颜色"] else @colors = colors
    if sizes.length==0 then @sizes=["均码"] else @sizes = sizes

    # generate html
    @generateHtml()

  generateHtml: ->
    #
    htmls = []
    htmls.push "<!-- Generated Table, input quantity here -->"
    htmls.push "<table class=\"prd_tbl\">"
    htmls.push " <tr>"
    htmls.push "  <th align=\"left\">颜色</th>"
    htmls.push "  <th align=\"left\">尺码</th>"
    htmls.push "  <th align=\"left\">数量</th>"
    htmls.push " </tr>"

    nColors = @colors.length
    nSizes = @sizes.length
    for color in @colors
      htmls.push " <tr>"
      htmls.push "  <td rowspan=\"#{nSizes}\">#{color}</td>"
      htmls.push "  <td>#{@sizes[0]}</td>"
      htmls.push "  <td><input type=\"text\" name=\"Stocks\" id=\"csq_#{color}__#{@sizes[0]}\" class=\"stock\" size=\"8\" value=\"\"></td>"
      htmls.push " </tr>"

      for size in @sizes.slice(1,@sizes.length)
        htmls.push " <tr>"
        htmls.push "  <td>#{size}</td>"
        htmls.push "  <td><input type=\"text\" name=\"Stocks\" id=\"csq_#{color}__#{size}\" class=\"stock\" size=\"8\" value=\"\"></td>"
        htmls.push " </tr>"

    htmls.push "</table>"
    @html = htmls.join("\n")

  replace: (divId)->
      $("##{divId}").html(@html)

