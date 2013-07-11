##################################################
# Order Product Selector Component script.
# Author: elivoa@gmail.com
# Data Required:
#   data:{...}
##################################################

window.OrderProductSelector =
class OrderProductSelector
  constructor:(customerId) ->

    # fields
    @customerId = customerId
    @containerClass = "product-selector"
    @product = {} # product json of this table.

    # callbacks (on basically means after)
    @onSelectProduct # callback after user select a product
    @onAddToOrder    # callback on add to orderp

    # call init
    @init()


  ## ________________________________________
  ## Initialize
  init: ->
    _=@

    # Suggest on Product
    @sc = new SuggestControl({
      parentClass : ".product-selector",
      triggerClass : ".product-trigger",
      hiddenClass : ".product-id"
      category : "product"
      onSelect : $.proxy (line, suggestion) ->
        @onProductSelect line, suggestion
        @onSelectProduct suggestion.data if @onSelectProduct
      ,@
    })
    @sc.init()

    # bind action on AddToOrder button
    $(".ops-add").bind 'click', $.proxy @onAddToOrderClick,@


  ## ________________________________________
  ## on suggest select
  ## crate the new json
  onProductSelect:(line, suggestion) ->
    _=@
    newproduct = {}
    # get customer id. TODO bad design
    productId = suggestion.data
    url = "/api/product/#{productId}"
    $.ajax({
      url: url
      context: document.body
      dataType: 'json'
      success: (data)->
        if data
          newproduct = {
            id: data.Id
            name: data.Name
            productPrice: data.Price
            colors: data.Colors
            sizes: data.Sizes
          }

        ## 2nd ajax. ajax get price
        ## update customer price & product price.
        urlprice = "/api/customer_price/#{@customerId}/#{productId}"
        $.ajax({
          url: urlprice
          context: document.body
          dataType: 'json'
          success: (data)->
            if data
              newproduct.price = data.price
              newproduct.productPrice = data.productPrice

            ## finally ajax all done, We need to refresh the table
            _.refresh(newproduct)

          error: (jqXHR, textStatus, errorThrown)->
            console.log textStatus
        })
      error: (jqXHR, textStatus, errorThrown)->
        console.log textStatus
    })


  ## ________________________________________
  ## Public: call by others.
  refresh: (product)->
    @product = product
    @sc.select(product.id, product.name)
    @refreshContent()


  ## ________________________________________
  refreshContent: ()->
    console.log 'refresh content, ', @product
    # 1/3 update cstable
    if @product.colors!= null && @product.sizes != null
      pcstg = new ProductCSTableGenerator(@product.colors, @product.sizes, "cs-container")
    else
      $("#cs-container").html("ERROR Loading Color&Size information. Product Information Has Errors!")

    # 2/3 update price
    $(".#{@containerClass} .price").html(@product.price)
    if @product.productPrice - @product.price == 0
      $(".#{@containerClass} .info").html("")
    else
      $(".#{@containerClass} .info").html("原价：#{@product.productPrice}"+
      "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; √ 已优惠")

    # 3/3 refresh stock if has
    @fillQuantities()

  ## ________________________________________
  fillQuantities: ->
    if @product and @product.quantity
      for q in @product.quantity
        o = $("#cs-container #csq_#{q[0]}__#{q[1]}") # TODO hardcode
        if o != undefined
          o.val q[2]




  ## ________________________________________
  ## Data part, extract productjson from web. including stocks.
  onAddToOrderClick:(e) ->
    e.preventDefault()
    if not @sc.selection
      alert "请先输入产品!"
      return
    @onAddToOrder @extractProductJson() if @onAddToOrder


  ## ________________________________________
  extractProductJson: ()->
    # fetch all product values.

    strprice = $(".#{@containerClass} .price").html()
    @product.price = parseInt(strprice)
    @product.note = $(".#{@containerClass} .notes").val()
    @product.quantity = []

    # fetch cst info
    $(".#{@containerClass} .stock").each $.proxy (idx, obj) -> #{
      a = obj.id
      a = a.slice(4,a.length)
      csinfo = a.split("__")
      strValue = obj.value
      value = 0
      if strValue != ""
        value = parseInt strValue
      @product.quantity.push([csinfo[0], csinfo[1], value])
    ,@ #}

    return @product

  ## ________________________________________
  ## clear all values
  clear: ->
    @sc.clearSelect() # clear suggest control
    $("#cs-container").html("Please select product.") # clear cstable
    $(".#{@containerClass} .notes").val("") # clear notes
    # update price
    $(".#{@containerClass} .price").html("")
    $(".#{@containerClass} .info").html("")
