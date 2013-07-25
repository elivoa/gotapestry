##
## Component
## Elivoa @ Time-stamp: <[order_list.coffee] Elivoa @ Wednesday, 2013-07-24 13:34:46>
##
window.OrderList=
class OrderList
  constructor:(param) ->
    @param = param
    @clientId = param.ClientId
    @container = $("##{param.ClientId}")
    @initBatchCloseButton()
    @initAction()

  initBatchCloseButton:->
    if window['BatchCloseOrder']
      @bco = new BatchCloseOrder {ClientId : "#{@param.ClientId}_close"}
      @bco.onTriggerClick = $.proxy @closeButtonClick,@

  initAction:->
    @checkall_btn = $("##{@clientId} .check-all")
    @checkall_btn.click $.proxy (e)->
      if @checkall_btn.prop("checked") == true
        $("##{@clientId} .order-check").prop("checked", true)
      else
        $("##{@clientId} .order-check").prop("checked", false)
    ,@


  closeButtonClick: (e)->
    e.preventDefault()
    # get selected orders
    tns = []
    tmp_customer_id = undefined
    pass = true
    @container.find(".order-check").each (index, obj)->
      return if pass == false # break
      if $(obj).prop("checked") == true
        tns.push $(obj).val()
        cid = $(obj).attr("CustomerId")
        if tmp_customer_id != undefined && tmp_customer_id != cid
          alert "喜大乐这么小，你怎么能让他一下子处理这么多人的订单呢？只选择一个人的吧！"
          pass = false
        tmp_customer_id = cid
    if tns.length == 0
      alert "喜乐说你至少要选择一个订单才能结款啊亲~"
      return false
    # set to bco
    @bco.customerId = parseInt(tmp_customer_id)
    @bco.selectedTrackNumbers = tns
    return pass

