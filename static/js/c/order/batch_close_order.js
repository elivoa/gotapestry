// Generated by CoffeeScript 1.6.3
(function() {
  var BatchCloseOrder;

  window.BatchCloseOrder = BatchCloseOrder = (function() {
    function BatchCloseOrder(param) {
      this.param = param;
      this.clientId = this.param.ClientId;
      this.customerId = this.param.CustomerId;
      this.selectedTrackNumbers = void 0;
      this.onTriggerClick;
      this.t = $("#" + this.clientId + "_trigger");
      this.m = $("#" + this.clientId + "_modal");
      this.money = this.m.find("input.money");
      this.t.click($.proxy(this.onclick, this));
      this.money.on("keyup", $.proxy(this.liveMatch, this));
      this.m.find(".submit").click($.proxy(this.submit, this));
    }

    BatchCloseOrder.prototype.onclick = function(e) {
      var clearbtn;
      e.preventDefault();
      if (this.onTriggerClick) {
        if (!this.onTriggerClick(e)) {
          return false;
        }
      }
      this.m.on('shown', $.proxy(this.onshown, this));
      this.m.on('hide', $.proxy(function() {
        return this.m.find(".unclosed-orders tbody").html("");
      }, this));
      this.m.modal("show");
      clearbtn = this.m.find(".money-clear");
      return this.m.find(".money-clear").click($.proxy(function() {
        if (clearbtn.prop('checked') === true) {
          this.m.find("input.money").val(-this.param.accountBallance);
          return this.liveMatch();
        } else {
          return this.m.find("input.money").focus();
        }
      }, this));
    };

    BatchCloseOrder.prototype.onshown = function() {
      var url;
      $.ajax({
        type: "GET",
        url: "/api/person/" + this.customerId,
        dataType: "json",
        success: $.proxy(function(json) {
          this.m.find(".customer strong").html(json.Name);
          this.m.find(".customer .price").html(json.AccountBallance);
          return this.param.accountBallance = json.AccountBallance;
        }, this)
      });
      url = "/order/deliveringunclosedorders/" + this.customerId;
      if (this.selectedTrackNumbers !== void 0 && this.selectedTrackNumbers.length > 0) {
        url = "/order/deliveringunclosedorders.byTrackingNumber/" + this.selectedTrackNumbers.join(",");
      }
      $.ajax({
        type: "GET",
        url: url,
        dataType: "json",
        success: $.proxy(function(data) {
          return this.applyJson(data);
        }, this)
      });
      return this.m.find("input.money").focus();
    };

    BatchCloseOrder.prototype.applyJson = function(json) {
      var abafter, order, tb, tn, tr, _i, _len, _ref;
      this.TotalOrderPrice = json.TotalOrderPrice;
      tb = this.m.find(".unclosed-orders tbody");
      tb.find("tr.order").each(function(index, tr) {
        var order, priceobj, tracknumber;
        tracknumber = $(tr).find(".tn").html();
        order = json.Orders[tracknumber];
        if (order === void 0) {
          $(tr).removeClass("L_cash_wait");
          $(tr).addClass("L_cash_done");
          priceobj = $(tr).find(".price");
          return priceobj.append("<span class=\"order-clear\">Clear!</span>");
        } else {
          delete json.Orders[tracknumber];
          return $(tr).removeClass("L_cash_wait L_cash_done");
        }
      });
      _ref = json.Order;
      for (_i = 0, _len = _ref.length; _i < _len; _i++) {
        tn = _ref[_i];
        order = json.Orders[tn];
        if (order !== void 0) {
          tr = [];
          tr.push("<tr class=\"order\" money=\"" + order.price + "\">");
          tr.push("  <td class=\"tn\">" + order.tn + "</td>");
          tr.push("  <td>" + order.time + "</td>");
          tr.push("  <td><span class=\"price\">" + order.price + "</span></td>");
          tr.push("</tr>");
          $(tr.join("\n")).appendTo(tb);
        }
      }
      abafter = this.TotalOrderPrice + this.param.accountBallance;
      if (abafter !== 0) {
        return this.m.find(".customer .bad-account").html("结后余额:" + abafter);
      }
    };

    BatchCloseOrder.prototype.liveMatch = function(e) {
      var totalmoney, _;
      totalmoney = parseFloat(this.money.val()) + (this.param.accountBallance + this.TotalOrderPrice);
      this.orders_can_clear = 0;
      _ = this;
      return this.m.find(".unclosed-orders tr.order").each(function(index, obj) {
        var orderMoney;
        orderMoney = parseFloat($(obj).attr('money'));
        if (orderMoney >= 0) {
          if (totalmoney >= orderMoney) {
            totalmoney -= orderMoney;
            $(obj).removeClass("L_cash_wait");
            $(obj).addClass("L_cash_done");
            return _.orders_can_clear += 1;
          } else {
            $(obj).removeClass("L_cash_done");
            return $(obj).addClass("L_cash_wait");
          }
        }
      });
    };

    BatchCloseOrder.prototype.submit = function(e) {
      var totalmoney;
      totalmoney = parseFloat(this.money.val());
      if (isNaN(totalmoney) || totalmoney <= 0) {
        alert("喜乐说: 你丫的不给钱还想结款？");
        return;
      }
      if (this.orders_can_clear === 0) {
        alert("Warrning! Warrning! 不够结订单，只将钱款加入账户。");
      }
      return $.ajax({
        type: "GET",
        url: "/order/deliveringunclosedorders:batchclose/" + totalmoney + "/" + this.customerId,
        dataType: "json",
        success: $.proxy(function(data) {
          this.applyJson(data);
          this.orders_can_clear = 0;
          this.money.val("");
          this.m.find("a.btn_a_s").html("结款完毕！必须刷新！(1秒后刷新)");
          setTimeout(function(){
            window.location.reload();
          }, 1000);
          return "";
        }, this),
        error: function() {
          return alert('error occured');
        }
      });
    };

    return BatchCloseOrder;

  })();

}).call(this);
