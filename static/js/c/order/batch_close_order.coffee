##
## Elivoa @ Time-stamp: <[batch_close_order.coffee] Elivoa @ Thursday, 2016-01-14 23:48:47>
##
window.BatchCloseOrder =
class BatchCloseOrder
  constructor:(param) -> # clientId, customerId, accountBallance
    @param = param
    @clientId = @param.ClientId
    @customerId = @param.CustomerId
    @selectedTrackNumbers = undefined # if not empty, use this instead of customerid
    # events
    @onTriggerClick

    @t = $("##{@clientId}_trigger")
    @m = $("##{@clientId}_modal")
    @money = @m.find("input.money")
    @t.click $.proxy @onclick,@
    @money.on "keyup", $.proxy @liveMatch,@
    @m.find(".submit").click $.proxy @submit,@


  onclick:(e) ->
    e.preventDefault()
    (return false if not @onTriggerClick(e)) if @onTriggerClick # call event

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
    # load person
    $.ajax {
      type:"GET"
      url:"/api/person/#{@customerId}"
      dataType:"json"
      success: $.proxy (json) ->
        @m.find(".customer strong").html(json.Name)
        @m.find(".customer .price").html(json.AccountBallance)
        @param.accountBallance=json.AccountBallance
      ,@
    }

    # load list
    url = "/order/deliveringunclosedorders/#{@customerId}" # customer version
    if @selectedTrackNumbers != undefined && @selectedTrackNumbers.length > 0
      url = "/order/deliveringunclosedorders.byTrackingNumber/" + @selectedTrackNumbers.join(",")
    $.ajax {
      type:"GET"
      url: url
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

    # bad-account display
    abafter = @TotalOrderPrice + @param.accountBallance
    if abafter != 0
      @m.find(".customer .bad-account").html("结后余额:#{abafter}")


  liveMatch:(e) ->
    # money used as total shouldbe: inputmoney + (accountballance - allorder's price)
    totalmoney = parseFloat(@money.val()) + (@param.accountBallance + @TotalOrderPrice)
    @orders_can_clear = 0 # at least one order can clear!
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
    if isNaN(totalmoney) or totalmoney <= 0
      alert "喜乐说: 你丫的不给钱还想结款？"
      return
    if @orders_can_clear == 0
      alert "Warrning! Warrning! 不够结订单，只将钱款加入账户。"

    # submit to batch clear
    $.ajax {
      type:"GET"
      url:"/order/deliveringunclosedorders:batchclose/#{totalmoney}/#{@customerId}"
      dataType: "json"
      success: $.proxy (data) ->
        @applyJson(data)
        @orders_can_clear = 0
        @money.val("")
        @m.find("a.btn_a_s").html("结款完毕！必须刷新！")
        # TODO 延时1秒钟刷新窗口
        # setTimeout(function(){
        #    window.location.reload();
        #  }, 1000);
      ,@
      error: ()->
        alert 'error occured'
    }
