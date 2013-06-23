###
  SYD Sales System
  @author: Bo Gao, [elivoa@gmail.com]
###

# order_edit page required js.

################################################################################
# onload
################################################################################
class OrderEdit
  constructor: ->
    this.init()

  # fill in sum on start.
  init: ->
    # register onblur function
    fn = $.proxy @onProductBlur, @
    $("body").on "blur", ".product-id",->
    $("body").on "blur", '.product-id', fn
    $("body").on "blur", '.quantity', fn
    $("body").on "blur", '.price', fn
    $("body").on "blur", '.pay', $.proxy @onPaidBlur, @
    # TODO this.优惠2元

    # init call calculate all prices.
    t = this
    $('.product-line').each (idx, obj) ->
      t.calculateSum $ obj
    this.calculateTotal()
    this.calculateTotalQuantity()
    return 1

  # product input's onblur
  onProductBlur:(e) ->
    # calculate the current line
    productLine = $(e.target).parents(".product-line")
    this.calculateSum productLine

    # calculate sum after modify values.
    this.calculateTotal()
    this.calculateTotalQuantity()
    return 1

  onPaidBlur: (e)->
    this.calculatePriceCut()

  # calculate
  calculateSum: (productLine)->
    sumObj = productLine.find('.sum')
    try
      quantity = parseInt productLine.find('.quantity').val()
      price = parseFloat productLine.find('.price').val()
      if isNaN(quantity) or isNaN(price)
        this.markError sumObj
      else
        this.markPass sumObj,quantity*price
    catch error
      this.markError sumObj
      "And the error is: sum can't be calculated..."
    return 1

  calculatePriceCut: ->
    total = parseFloat $(".total-price").html()
    paid = parseInt $(".pay").val()
    priceCut = total-paid
    this.markPass $('.price-cut-display'),priceCut
    $('.price-cut').val(priceCut)
    return 1

  # calculate total-quantity and total-sum of the order
  calculateTotal: ->
    total = 0
    onError = false
    $('.product-line .sum').each (idx, obj)->
      sum = parseFloat $(obj).html()
      if isNaN sum
        onError=true
      else
        total += sum

    totalPriceSpan = $ '.total-price'
    if onError
      this.markError totalPriceSpan
    else
      this.markPass totalPriceSpan,total
      $(".pay").val(total)

    # 优惠的价格 PriceCutDown
    this.calculatePriceCut()
    return 1

  calculateTotalQuantity: ->
    totalQuantity = 0
    onError = false
    $('.product-line .quantity').each (idx, obj)->
      quantity = parseInt $(obj).val()
      if isNaN quantity
        onError=true
      else
        totalQuantity += quantity

    totalQuantitySpan = $ '.total-quantity'
    if onError
      this.markError totalQuantitySpan
    else
      # totalPriceSpan.html(total)
      this.markPass totalQuantitySpan,totalQuantity
    return 1

  markError: (obj)->
    obj.html "-.--"
    obj.css 'color', 'red'

  markPass: (obj, price)->
    # TODO format currency
    obj.html price
    obj.css 'color', 'green'


################################################################################
# SubModule: Add / Remove OrderDetails
################################################################################
class OrderManageDetails
  constructor: ()->
    @container = $ '#product-line' # tr
    @register()
    @onAdd

  register: ->
    $('body').on "click",'.fn_delete_line', (e)->
      e.preventDefault()
      if $('.product-line').length > 1
        $(e.target).parents('.product-line').detach()
      else
        alert("You can't delete the last product.")

    _ = @
    $('.fn_add_line').on "click", (e)->
      e.preventDefault()
      _.AddLine(e)

  # __________________________________________________
  # Create new OrderDetails line and append to table.
  AddLine:(e) ->
    lastline = undefined
    lines = $('.product-line')
    if lines.length > 0
      # clone line
      line = $(lines[lines.length-1])
      newline = line.clone()

      # empty value
      newline.find('input').each (idx, obj) ->
        $(obj).attr "value",""
      newline.find('.sum').html("")

      # attach to table
      newline.insertAfter(line)

      # register suggest control
      @onAdd(newline) if @onAdd


# ____________________________________________________________
# Register
# TODO:
#   This init only support one instance in one page?
$ ->
  window['OrderEdit'] = new OrderEdit

  # --------------------------------
  # Init OrderDetails SuggestControl
  # --------------------------------
  sc = new window.SuggestControl({
    parentClass : ".product-line",
    triggerClass : ".product-trigger",
    hiddenClass : ".product-id"
    category : "product"
    onSelect : (line, suggestion) ->
      # TODO fill in all the following textbox.

      # console.log "--------------------------"
      # console.log line
      # console.log suggestion
  })
  sc.init()

  # Add/remove OrderDetails line
  omd = new OrderManageDetails()
  omd.onAdd = (line) ->
    console.log "----> register line"
    console.log line
    sc.registerLine(line)
