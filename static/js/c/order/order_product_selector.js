// Generated by CoffeeScript 1.5.0
(function() {
  var OrderProductSelector;

  window.OrderProductSelector = OrderProductSelector = (function() {

    function OrderProductSelector(customerId) {
      this.customerId = customerId;
      this.containerClass = "product-selector";
      this.product = {};
      this.onSelectProduct;
      this.onAddToOrder;
      this.init();
    }

    OrderProductSelector.prototype.init = function() {
      var _;
      _ = this;
      this.sc = new SuggestControl({
        parentClass: ".product-selector",
        triggerClass: ".product-trigger",
        hiddenClass: ".product-id",
        category: "product",
        onSelect: $.proxy(function(line, suggestion) {
          this.onProductSelect(line, suggestion);
          if (this.onSelectProduct) {
            return this.onSelectProduct(suggestion.data);
          }
        }, this)
      });
      this.sc.init();
      return $(".ops-add").bind('click', $.proxy(this.onAddToOrderClick, this));
    };

    OrderProductSelector.prototype.onProductSelect = function(line, suggestion) {
      var newproduct, productId, url, _;
      _ = this;
      newproduct = {};
      productId = suggestion.data;
      url = "/api/product/" + productId;
      return $.ajax({
        url: url,
        context: document.body,
        dataType: 'json',
        success: function(data) {
          var urlprice;
          if (data) {
            newproduct = {
              id: data.Id,
              name: data.Name,
              productPrice: data.Price,
              colors: data.Colors,
              sizes: data.Sizes
            };
          }
          urlprice = "/api/customer_price/" + this.customerId + "/" + productId;
          return $.ajax({
            url: urlprice,
            context: document.body,
            dataType: 'json',
            success: function(data) {
              if (data) {
                newproduct.price = data.price;
                newproduct.productPrice = data.productPrice;
              }
              return _.refresh(newproduct);
            },
            error: function(jqXHR, textStatus, errorThrown) {
              return console.log(textStatus);
            }
          });
        },
        error: function(jqXHR, textStatus, errorThrown) {
          return console.log(textStatus);
        }
      });
    };

    OrderProductSelector.prototype.refresh = function(product) {
      this.product = product;
      this.sc.select(product.id, product.name);
      return this.refreshContent();
    };

    OrderProductSelector.prototype.refreshContent = function() {
      var pcstg;
      console.log('refresh content, ', this.product);
      if (this.product.colors !== null && this.product.sizes !== null) {
        pcstg = new ProductCSTableGenerator(this.product.colors, this.product.sizes);
        pcstg.replace("cs-container");
      } else {
        $("#cs-container").html("ERROR Loading Color&Size information. Product Information Has Errors!");
      }
      $("." + this.containerClass + " .price").html(this.product.price);
      if (this.product.productPrice - this.product.price === 0) {
        $("." + this.containerClass + " .info").html("");
      } else {
        $("." + this.containerClass + " .info").html(("原价：" + this.product.productPrice) + "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; √ 已优惠");
      }
      return this.fillQuantities();
    };

    OrderProductSelector.prototype.fillQuantities = function() {
      var o, q, _i, _len, _ref, _results;
      if (this.product && this.product.quantity) {
        _ref = this.product.quantity;
        _results = [];
        for (_i = 0, _len = _ref.length; _i < _len; _i++) {
          q = _ref[_i];
          o = $("#cs-container #csq_" + q[0] + "__" + q[1]);
          if (o !== void 0) {
            _results.push(o.val(q[2]));
          } else {
            _results.push(void 0);
          }
        }
        return _results;
      }
    };

    OrderProductSelector.prototype.onAddToOrderClick = function(e) {
      e.preventDefault();
      if (!this.sc.selection) {
        alert("请先输入产品!");
        return;
      }
      if (this.onAddToOrder) {
        return this.onAddToOrder(this.extractProductJson());
      }
    };

    OrderProductSelector.prototype.extractProductJson = function() {
      var strprice;
      strprice = $("." + this.containerClass + " .price").html();
      this.product.price = parseInt(strprice);
      this.product.note = $("." + this.containerClass + " .notes").val();
      this.product.quantity = [];
      $("." + this.containerClass + " .stock").each($.proxy(function(idx, obj) {
        var a, csinfo, strValue, value;
        a = obj.id;
        a = a.slice(4, a.length);
        csinfo = a.split("__");
        strValue = obj.value;
        value = 0;
        if (strValue !== "") {
          value = parseInt(strValue);
        }
        return this.product.quantity.push([csinfo[0], csinfo[1], value]);
      }, this));
      return this.product;
    };

    OrderProductSelector.prototype.clear = function() {
      this.sc.clearSelect();
      $("#cs-container").html("Please select product.");
      $("." + this.containerClass + " .notes").val("");
      $("." + this.containerClass + " .price").html("");
      return $("." + this.containerClass + " .info").html("");
    };

    return OrderProductSelector;

  })();

}).call(this);
