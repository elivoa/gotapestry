#######################################################
# OrderEdit related js
# Features:
#   商品选择模块的初始化和屁话，独立成模块
#
#######################################################

class OrderEdit
  constructor: ->
    #pass

## starting...
$ ->
  ## Suggest on Product
  sc = new window.SuggestControl({
    parentClass : ".product-select",
    triggerClass : ".product-trigger",
    hiddenClass : ".product-id"
    category : "product"
    onSelect : (line, suggestion) ->
      console.log suggestion
      productId = suggestion.data

      # get customer id. TODO bad design
      customerId = $(".suggest-id").val()
      url = "/api/customer_price/" + customerId + "/" + productId
      $.ajax({
        url: url
        context: document.body
        dataType: 'json'
        success: (data)->
          if data
            $(line).find(".price").val(data.price)
        error: (jqXHR, textStatus, errorThrown)->
          console.log textStatus
      }).done(->
        console.log 'done...'
      )
  })

