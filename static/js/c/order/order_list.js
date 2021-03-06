// Generated by CoffeeScript 1.6.3
(function() {
  var OrderList;

  window.OrderList = OrderList = (function() {
    function OrderList(param) {
      this.param = param;
      this.clientId = param.ClientId;
      this.container = $("#" + param.ClientId);
      this.initBatchCloseButton();
      this.initAction();
    }

    OrderList.prototype.initBatchCloseButton = function() {
      if (window['BatchCloseOrder']) {
        this.bco = new BatchCloseOrder({
          ClientId: "" + this.param.ClientId + "_close"
        });
        return this.bco.onTriggerClick = $.proxy(this.closeButtonClick, this);
      }
    };

    OrderList.prototype.initAction = function() {
      this.checkall_btn = $("#" + this.clientId + " .check-all");
      return this.checkall_btn.click($.proxy(function(e) {
        if (this.checkall_btn.prop("checked") === true) {
          return $("#" + this.clientId + " .order-check").prop("checked", true);
        } else {
          return $("#" + this.clientId + " .order-check").prop("checked", false);
        }
      }, this));
    };

    OrderList.prototype.closeButtonClick = function(e) {
      var pass, tmp_customer_id, tns;
      e.preventDefault();
      tns = [];
      tmp_customer_id = void 0;
      pass = true;
      this.container.find(".order-check").each(function(index, obj) {
        var cid;
        if (pass === false) {
          return;
        }
        if ($(obj).prop("checked") === true) {
          tns.push($(obj).val());
          cid = $(obj).attr("CustomerId");
          if (tmp_customer_id !== void 0 && tmp_customer_id !== cid) {
            alert("喜大乐这么小，你怎么能让他一下子处理这么多人的订单呢？只选择一个人的吧！");
            pass = false;
          }
          return tmp_customer_id = cid;
        }
      });
      if (tns.length === 0) {
        alert("喜乐说你至少要选择一个订单才能结款啊亲~");
        return false;
      }
      this.bco.customerId = parseInt(tmp_customer_id);
      this.bco.selectedTrackNumbers = tns;
      return pass;
    };

    return OrderList;

  })();

}).call(this);
