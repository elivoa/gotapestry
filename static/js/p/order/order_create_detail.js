// Generated by CoffeeScript 1.5.0
(function() {
  var OrderCreateDetail;

  window.OrderCreateDetail = OrderCreateDetail = (function() {

    function OrderCreateDetail(customerId) {
      this.initPage();
      this.ops = new OrderProductSelector(customerId);
      this.odf = new OrderDetailsForm;
      this.ops.onSelectProduct = $.proxy(function(productId) {
        var product;
        product = this.odf.data.products[productId];
        if (product) {
          alert("已经添加了这件商品，不能重复添加！如需添加或修改，请点击下面对应商品的编辑按钮！谢谢合作！");
          return this.ops.clear();
        }
      }, this);
      this.ops.onAddToOrder = $.proxy(function(product) {
        console.log(product);
        this.odf.appendProduct(product);
        this.odf.refreshOrderForm();
        return this.ops.clear();
      }, this);
      this.odf.onEdit = $.proxy(function(product) {
        console.log("Edit Product: ", product);
        return this.ops.refresh(product);
      }, this);
      $('.product-trigger').focus();
    }

    OrderCreateDetail.prototype.initPage = function() {
      var ef;
      $(".choose_express label").each(function(idx, obj) {
        return $(obj).click(function(e) {
          $(".choose_express label").removeClass('checked');
          $(obj).addClass("checked");
          return $(obj).find("input:radio").attr("checked", true);
        });
      });
      ef = $(".B_fare");
      return ef.find("input:checkbox").click(function(e) {
        if ($(e.target).prop("checked") === true) {
          return ef.find("input:text").prop('disabled', true);
        } else {
          return ef.find("input:text").prop('disabled', false);
        }
      });
    };

    OrderCreateDetail.prototype.setExpress = function(express, expressFee) {
      return $(".choose_express label").each(function(idx, obj) {
        var ef;
        if ($(obj).find("input:radio").val().toLowerCase() === express.toLowerCase()) {
          $(obj).addClass('checked');
          $(obj).find("input:radio").attr("checked", true);
        }
        ef = $(".B_fare");
        if (expressFee === -1) {
          ef.find(".express-fee").val(0);
          return ef.find("input.daofu").prop('checked', true);
        } else {
          ef.find("express-fee").val(expressFee);
          return ef.find("input.daofu").prop('checked', false);
        }
      });
    };

    return OrderCreateDetail;

  })();

}).call(this);
