#######################################################
# OrderEdit related js
# Features:
#   商品选择模块的初始化和屁话，独立成模块
#
#######################################################

window.OrderCreateDetail =
class OrderCreateDetail
  constructor: (customerId)->
    _=@

    @ops = new OrderProductSelector(customerId)
    @odf = new OrderDetailsForm

    ## Interactions
    @ops.onAddToOrder = (product) ->
      console.log product
      _.odf.appendProduct product
      _.odf.refreshOrderForm()
      _.ops.clear()
      console.log 'all done'

    @odf.onEdit = $.proxy (product)->
      console.log product
      @ops.refresh product
    ,@
