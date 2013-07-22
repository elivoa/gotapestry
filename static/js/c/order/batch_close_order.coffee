

window.BatchCloseOrder =
class BatchCloseOrder
  constructor:(param) -> # clientId, customerId, accountBallance
    @param = param
    @t = $("##{@param.clientId}_trigger")
    @m = $("##{@param.clientId}_modal")
    @money = @m.find("input.money")
    @t.click $.proxy @onclick,@
    @money.on "keyup", $.proxy @liveMatch,@
    @m.find(".submit").click $.proxy @submit,@

  onclick:(e) ->
    e.preventDefault()
    @m.on 'shown', $.proxy @onshown,@
    @m.on 'hide', $.proxy ->
      @m.find(".unclosed-orders tbody").html("")
    ,@
    @m.modal("show")

    clearbtn = @m.find(".money-clear")
    @m.find(".money-clear").click $.proxy ->
      if clearbtn.prop('checked') == true
        @m.find("input.money").val(-@param.accountBallance)
        @liveMatch()
      else
        @m.find("input.money").focus()
    ,@

  onshown: ->
    $.ajax {
      type:"GET"
      url:"/order/deliveringunclosedorders/#{@param.customerId}"
      dataType: "json"
      success: $.proxy (data) ->
        @applyJson(data)
      ,@
    }
    @m.find("input.money").focus()

  applyJson:(json)->
    @TotalOrderPrice = json.TotalOrderPrice
    tb = @m.find(".unclosed-orders tbody")
    # first loop all existing tr. mark price as done when done.
    tb.find("tr.order").each (index, tr)->
      tracknumber = $(tr).find(".tn").html()
      order = json.Orders[tracknumber]
      if order == undefined #
        $(tr).removeClass("L_cash_wait")
        $(tr).addClass("L_cash_done")
        priceobj = $(tr).find(".price")
        priceobj.append("<span class=\"order-clear\">Clear!</span>")
      else
        delete json.Orders[tracknumber] # don't add twice
        $(tr).removeClass("L_cash_wait L_cash_done")

    # append others into it.
    for tn in json.Order
      order = json.Orders[tn]
      if order != undefined
        tr = []
        tr.push("<tr class=\"order\" money=\"#{order.price}\">")
        tr.push("  <td class=\"tn\">#{order.tn}</td>")
        tr.push("  <td>#{order.time}</td>")
        tr.push("  <td><span class=\"price\">#{order.price}</span></td>")
        tr.push("</tr>")
        $(tr.join("\n")).appendTo(tb)

  liveMatch:(e) ->
    # money used as total shouldbe: inputmoney + (accountballance - allorder's price)
    totalmoney = parseFloat(@money.val()) + (@param.accountBallance + @TotalOrderPrice)
    @orders_can_clear = 0
    _=@
    @m.find(".unclosed-orders tr.order").each (index, obj)->
      orderMoney = parseFloat($(obj).attr('money'))
      if orderMoney >= 0
        if totalmoney >= orderMoney
          totalmoney -= orderMoney
          $(obj).removeClass("L_cash_wait")
          $(obj).addClass("L_cash_done")
          _.orders_can_clear += 1
        else
          $(obj).removeClass("L_cash_done")
          $(obj).addClass("L_cash_wait")

  submit:(e)->
    totalmoney = parseFloat(@money.val())
    if isNaN(totalmoney) or totalmoney <= 0 || @orders_can_clear == 0
      alert "喜乐说你丫的这么点钱还不够结款一单的呢！！！！！"
      return
    # submit to batch clear
    $.ajax {
      type:"GET"
      url:"/order/deliveringunclosedorders.batchclose/#{totalmoney}/#{@param.customerId}"
      dataType: "json"
      success: $.proxy (data) ->
        @applyJson(data)
        @orders_can_clear = 0
        @money.val("")
      ,@
      error: ()->
        alert 'error occured'
    }
