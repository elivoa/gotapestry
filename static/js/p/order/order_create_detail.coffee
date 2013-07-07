#######################################################
# OrderEdit related js
# Features:
#   商品选择模块的初始化和屁话，独立成模块
#
#######################################################

window.OrderCreateDetail =
class OrderCreateDetail
  constructor: (customerId)->

    ## Create components
    @ops = new OrderProductSelector(customerId)
    @odf = new OrderDetailsForm

    ## Interactions
    @ops.onSelectProduct = $.proxy (productId)->
      product = @odf.data.products[productId]
      if product
        alert "已经添加了这件商品，不能重复添加！如需添加或修改，请点击下面对应商品的编辑按钮！谢谢合作！"
        @ops.clear()
    ,@

    @ops.onAddToOrder = $.proxy (product) ->
      console.log product
      @odf.appendProduct product
      @odf.refreshOrderForm()
      @ops.clear()
    ,@

    @odf.onEdit = $.proxy (product)->
      console.log "Edit Product: ", product
      @ops.refresh product
    ,@
