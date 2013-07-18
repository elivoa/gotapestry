#######################################################
# OrderEdit related js
# Features:
#   商品选择模块的初始化和屁话，独立成模块
#
#######################################################

window.OrderCreateDetail =
class OrderCreateDetail
  constructor: (customerId)->

    @initPage()

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
      success = false
      console.log "isedit is", @ops.isEdit
      if @ops.isEdit
        success = @odf.editProduct product
      else
        success = @odf.addProduct product
      if success
        @odf.refreshOrderForm()
        @ops.clear()
    ,@

    @odf.onEdit = $.proxy (product)->
      console.log "Edit Product: ", product
      @ops.refresh product
      console.log "set isedit ", true
      @ops.setEdit true
    ,@

    # focus on input box
    # TODO: auto this?
    $('.product-trigger').focus()


  initPage: ->
    # express buttons
    $(".choose_express label").each (idx, obj) ->
      $(obj).click (e)->
        $(".choose_express label").removeClass 'checked' # remove checked
        $(obj).addClass "checked"
        $(obj).find("input:radio").attr("checked", true)

    # express fee
    ef = $(".B_fare")
    ef.find("input:checkbox").click (e)->
      if $(e.target).prop("checked") == true
        ef.find("input:text").prop('disabled', true)
      else
        ef.find("input:text").prop('disabled', false)

  # express: yto, sf, takeaway
  setExpress: (express, expressFee) ->
    $(".choose_express label").each (idx, obj) ->
      # express
      if $(obj).find("input:radio").val().toLowerCase() == express.toLowerCase()
        $(obj).addClass 'checked'
        $(obj).find("input:radio").attr("checked", true)
      # expressfee
      ef = $(".B_fare")
      if expressFee == -1
        ef.find(".express-fee").val(0)
        ef.find("input.daofu").prop('checked', true)
        ef.find("input:text").prop('disabled', true);
      else
        ef.find("express-fee").val(expressFee)
        ef.find("input.daofu").prop('checked', false)
        ef.find("input:text").prop('disabled', false);
