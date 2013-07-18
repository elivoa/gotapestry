

window.OrderDeliverButton =
class OrderDeliverButton
  constructor:(param) ->
    console.log param
    @param = param
    if not param
      @param={
        id: "deliverModal",
      }

    @m = $("##{@param.id}_modal")
    @init()
    @initAction()

  init:->
    $("##{@param.id}_trigger").click $.proxy (e)->
      e.preventDefaults
      @m.on 'show', $.proxy ->
        # TODO ajax init modal
        console.log "onshow"
      ,@
      @m.modal {
        keyboard: true
      }
    ,@

    $("##{@param.id}_modal .confirm").click $.proxy (e)->
      @asyncSubmit()
    ,@


  asyncSubmit: ->
    return 1
    _=@
    _.DeliveryTrackingNumber = $("##{@param.id}_modal .tracking-number").val()

    # submit form
    # TODO this is ajax submit, use form submit first
    $.ajax {
      type: "POST"
      url: "/order/ButtonSubmitHere"
      data : {
        "t:form" : "DeliverForm"
        "Order.DeliveryMethod" : _.DeliveryMethod
        "Order.DeliveryTrackingNumber" : _.DeliveryTrackingNumber
        "Order.ExpressFee" : _.ExpressFee
        "Tab" : "all"
      }
      success: (data)->
        console.log "success"
        console.log data
        @m.modal('hide')
      error: ()->
        console.log "error occured"
    }


  initAction : ->
    # express buttons
    _=@
    $(".choose_express label").each (idx, obj) ->
      $(obj).click (e)->
        $(".choose_express label").removeClass 'checked' # remove checked
        $(obj).addClass "checked"
        $(obj).find("input:radio").attr("checked", true)
        _.DeliveryMethod = $(obj).find("input:radio").val()

    # express fee
    ef = $(".B_fare")
    ef.find("input:checkbox").click (e)->
      if $(e.target).prop("checked") == true
        ef.find("input:text").prop('disabled', true)
        _.ExpressFee = -1
      else
        ef.find("input:text").prop('disabled', false)
        _.ExpressFee = ef.find("input:text").val()

  # set express fee.
  setExpress: (express, expressFee) ->
    $(".choose_express label").each $.proxy (idx, obj) ->
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
    ,@

# $ ->
#   new OrderListComponent
# param_example = {
#   deliverModal
#   }
