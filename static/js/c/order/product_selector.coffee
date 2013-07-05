##################################################
# Order Product Selector Component script.
# Author: elivoa@gmail.com
# Data Required:
#   data:{...}
##################################################

window.OrderProductSelector =
class OrderProductSelector
  constructor:(customerId) ->
    @customerId = customerId
    @init()
    @containerClass = "product-selector"

  init:->
    ## Suggest on Product
    _=@
    sc = new window.SuggestControl({
      parentClass : ".product-select",
      triggerClass : ".product-trigger",
      hiddenClass : ".product-id"
      category : "product"
      onSelect : (line, suggestion) ->
        console.log "select: ", suggestion
        _.onProductSelect line, suggestion
    })
    sc.init()

  ## on suggest select
  onProductSelect:(line, suggestion) ->
    _=@
    # get customer id. TODO bad design
    productId = suggestion.data
    url = "/api/product/#{productId}"
    $.ajax({
      url: url
      context: document.body
      dataType: 'json'
      success: (data)->
        if data
          # update color-size table
          if data.Colors!= null && data.Sizes != null
            pcstg = new ProductCSTableGenerator(data.Colors, data.Sizes)
            pcstg.replace("cs-container") # TODO make this robust.
          else
            $("#cs-container").html("ERROR Loading Color&Size information. Product Information Has Errors!")
      error: (jqXHR, textStatus, errorThrown)->
        console.log textStatus
    })

    ## update customer price & product price.
    url = "/api/customer_price/#{@customerId}/#{productId}"
    $.ajax({
      url: url
      context: document.body
      dataType: 'json'
      success: (data)->
        if data
          $(".#{_.containerClass} .price").html(data.price)
          if data.productPrice - data.price == 0
            $(".#{_.containerClass} .info").html("")
          else
            $(".#{_.containerClass} .info").html("原价：#{data.productPrice - data.price}"+
            "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; √ 已优惠")
      error: (jqXHR, textStatus, errorThrown)->
        console.log textStatus
    })

  ## TOOLS:
  ## Public: refresh
  refresh: ()->
    # TODO refresh this table.


## starting in page script
# $ ->
#   ops = new OrderProductSelector
