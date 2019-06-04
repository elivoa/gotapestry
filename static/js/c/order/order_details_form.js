// Generated by CoffeeScript 1.6.3//modified manually.

(function () {
  var OrderDetailsForm;

  var enableSales = true; // enable sales

  window.OrderDetailsForm = OrderDetailsForm = (function () {
    function OrderDetailsForm(hideOperation) {
      this.containerClass = ".order-form-container";
      this.hideOperation = hideOperation !== void 0 ? hideOperation : false;
      this.onDelete = this.defaultOnDelete;
      this.onEdit;
      this.data = {
        order: [],
        products: {}
      };
      this.refreshOrderForm();
    }

    OrderDetailsForm.prototype.addProduct = function (product) {
      if (!product) {
        return;
      }
      if (this.data.products[product.id]) {
        alert("已经添加了这件商品，不能重复添加！如需添加或修改，请点击下面对应商品的编辑按钮！谢谢合作！");
        return false;
      }
      this.data.order.push(product.id);
      this.data.products[product.id] = product;
      return true;
    };

    OrderDetailsForm.prototype.editProduct = function (product) {
      if (!product) {
        return;
      }
      this.data.products[product.id] = product;
      return true;
    };

    OrderDetailsForm.prototype.setData = function (json) {
      this.data = json;
      return this.refreshOrderForm();
    };

    OrderDetailsForm.prototype.refreshOrderForm = function () {
      var id, product, q, sumPrice, osumPrice, sumQuantity, tbody, totalPrice, ototalPrice, totalQuantity, tr, _i, _j, _len, _len1, _ref, _ref1;
      tbody = $("" + this.containerClass + " tbody");
      tbody.html("");
      sumQuantity = 0;
      sumPrice = 0;
      osumPrice = 0;
      _ref = this.data.order;

      for (_i = 0, _len = _ref.length; _i < _len; _i++) {
        id = _ref[_i];
        product = this.data.products[id];
        if (product) {
          totalQuantity = 0;
          _ref1 = product.quantity;
          for (_j = 0, _len1 = _ref1.length; _j < _len1; _j++) {
            q = _ref1[_j];
            totalQuantity += q[2];
          }

          // 没考虑折扣的情况下计算价格
          totalPrice = totalQuantity * product.price;
          ototalPrice = 0;
          // 是否启用折扣功能
          if (enableSales && product.discountPercent > 0) {
            // 如果折扣价格更低，那么更新价格到更低价格.
            var discountTotalPrice = totalQuantity * product.productPrice * (product.discountPercent / 100);
            if (discountTotalPrice < totalPrice) { // 如果打折更优惠
              ototalPrice = totalPrice;
              totalPrice = discountTotalPrice;
            }
          }
          // console.log('sumPrice:', sumPrice, totalPrice, sumPrice + totalPrice);
          // console.log('osumPrice:', osumPrice, ototalPrice, osumPrice + ototalPrice);
          sumQuantity += totalQuantity;
          sumPrice += totalPrice;
          if (enableSales) {
            osumPrice += ototalPrice > 0 ? ototalPrice : totalPrice;
          }

          tr = $(this.generateTR(product, totalQuantity, totalPrice, ototalPrice).join("\n"));
          tbody.append(tr);
          this.bindAction(tr, id);
        } else {
          if (console) {
            console.log("[OrderDetailsForm] Error id in order list " + id + ".");
          }
        }
      }
      return this.updateSummary(sumQuantity, sumPrice, osumPrice);
    };

    OrderDetailsForm.prototype.bindAction = function (tr, id) {
      tr.find(".odf-edit").on("click", $.proxy(this.onODFEdit(id), this));
      return tr.find(".odf-delete").on("click", $.proxy(this.onODFDelete(id), this));
    };

    OrderDetailsForm.prototype.onODFEdit = function (id) {
      return function (e) {
        e.preventDefault();
        if (this.onEdit) {
          return this.onEdit(this.data.products[id]);
        }
      };
    };

    OrderDetailsForm.prototype.onODFDelete = function (id) {
      return function (e) {
        e.preventDefault();
        if (this.onDelete) {
          this.onDelete(this.data.products[id]);
        }
        return console.log('----------------------------------------');
      };
    };

    OrderDetailsForm.prototype.defaultOnDelete = function (product) {
      var idx;
      if (!confirm("真的要删除这条记录么？")) {
        return;
      }
      delete this.data.products[product.id];
      idx = this.data.order.indexOf(product.id);
      if (idx >= 0) {
        if (console) {
          console.log(this.data.order.splice(idx, 1));
        }
      }
      return this.refreshOrderForm();
    };

    OrderDetailsForm.prototype.updateSummary = function (sumQuantity, sumPrice, osumPrice) {
      var tfoot;
      tfoot = $("" + this.containerClass + " tfoot");
      tfoot.find(".sumQuantity").html(sumQuantity);
      tfoot.find(".sumPrice").html(enableSales ? sumPrice.toFixed(2) : osumPrice);
      var additionalLine = $(".additionalSumLine");
      console.log('update summary', sumPrice, osumPrice)

      var diff = (osumPrice - sumPrice).toFixed(2);
      if (diff > 0) {
        additionalLine.html('原价：' + osumPrice.toFixed(2) + '元，已优惠：' + diff + '元。');
      } else {
        additionalLine.html('');
      }
      return
    };

    OrderDetailsForm.prototype.generateTR = function (json, totalQuantity, totalPrice, ototalPrice) {
      // console.log("product --> :",json, totalQuantity, totalPrice)
      var htmls, nquantity, q, quantities, quantity, _i, _j, _len, _len1, _ref, _ref1;
      var downsign = "  ";
      var useDiscount = false;
      quantities = [];
      _ref = json.quantity;
      for (_i = 0, _len = _ref.length; _i < _len; _i++) {
        q = _ref[_i];
        if (q[2] > 0) {
          quantities.push(q);
        }
      }
      if (ototalPrice > 0 && ototalPrice > totalPrice) {
        downsign = "↓ ";
        useDiscount = true;
      }
      nquantity = quantities.length;

      var uniquePrice = enableSales && useDiscount ? json.productPrice : json.price;

      htmls = [];
      htmls.push("<tr>");
      htmls.push("  <td valign='top' rowspan='" + nquantity + "'>");
      htmls.push("    " + json.pid);
      htmls.push("  </td>");
      htmls.push("  <td valign='top' rowspan='" + nquantity + "'>");
      htmls.push("    <strong><a href='/product/detail/" + json.id + "' target='_blank'>" + json.name + "</a></strong>");
      htmls.push("    <input type='hidden' name='Order.Details.ProductId' value='" + json.id + "' />");
      htmls.push("    <input type='hidden' name='Order.Details.SellingPrice' value='" + uniquePrice + "' />");
      htmls.push("    <input type='hidden' name='Order.Details.ProductPrice' value='" + json.productPrice + "' />");
      htmls.push("    <input type='hidden' name='Order.Details.DiscountPercent' value='" + (json.discountPercent || 0) + "' />");
      htmls.push("    <input type='hidden' name='Order.Details.Color' value='" + quantities[0][0] + "' />");
      htmls.push("    <input type='hidden' name='Order.Details.Size' value='" + quantities[0][1] + "' />");
      htmls.push("    <input type='hidden' name='Order.Details.Quantity' value='" + quantities[0][2] + "' />");
      htmls.push("    <input type='hidden' name='Order.Details.Note' value='" + json.note + "' />");
      htmls.push("  </td>");
      htmls.push("  <td valign='top' rowspan='" + nquantity + "'>");
      // 如果折扣价格胜出，那么录入产品原价格和折扣比例，保证计算是准确的。
      htmls.push("    <span class='price'>" + uniquePrice + "</span>");
      htmls.push("  </td>");
      htmls.push("  <td>" + quantities[0][0] + "</td>");
      htmls.push("  <td>" + quantities[0][1] + "</td>");
      htmls.push("  <td>" + quantities[0][2] + "</td>");
      htmls.push("  <td valign='top' align='center' rowspan='" + nquantity + "'>");
      htmls.push("      <strong>" + totalQuantity + "</strong></td>");
      // 输出折扣比例
      if (enableSales) {
        htmls.push("  <td valign='top' align='right' rowspan='" + nquantity + "'>");
        if (useDiscount) {
          htmls.push("      <strong class=''>" + json.discountPercent + "%</strong>");
        }
        htmls.push("</td>");
      }

      htmls.push("  <td valign='top' align='right' rowspan='" + nquantity + "'>");
      htmls.push("      <strong class='price'>" + downsign + totalPrice.toFixed(2) + "</strong></td>");

      htmls.push("  <td valign='top' rowspan='" + nquantity + "'>" + json.note + "</td>");
      if (!this.hideOperation) {
        htmls.push("  <td valign='top' rowspan='" + nquantity + "'>");
        htmls.push("      <a href='#' class='odf-edit'>编辑</a><span class='vline'>|</span>");
        htmls.push("      <a href='#' class='odf-delete'>删除</a>");
        htmls.push("  </td>");
      }
      htmls.push("</tr>");
      _ref1 = quantities.slice(1, nquantity);
      for (_j = 0, _len1 = _ref1.length; _j < _len1; _j++) {
        quantity = _ref1[_j];
        htmls.push("<tr>");
        htmls.push("  <td>" + quantity[0] + "</td>");
        htmls.push("  <td>" + quantity[1] + "</td>");
        htmls.push("  <td>" + quantity[2] + "</td>");
        htmls.push("    <input type='hidden' name='Order.Details.ProductId' value='" + json.id + "' />");
        htmls.push("    <input type='hidden' name='Order.Details.SellingPrice' value='" + json.price + "' />");
        htmls.push("    <input type='hidden' name='Order.Details.Color' value='" + quantity[0] + "' />");
        htmls.push("    <input type='hidden' name='Order.Details.Size' value='" + quantity[1] + "' />");
        htmls.push("    <input type='hidden' name='Order.Details.Quantity' value='" + quantity[2] + "' />");
        htmls.push("    <input type='hidden' name='Order.Details.Note' value='" + json.note + "' />");
        htmls.push("</tr>");
      }
      return htmls;
    };

    OrderDetailsForm.prototype.addTestData = function () {
      var testproduct;
      testproduct = {
        id: 1,
        pid: 2233,
        name: "绣虎头",
        price: 138,
        productPrice: 120,
        note: "no note",
        colors: ["红色", "蓝色"],
        sizes: ["S", "M"],
        quantity: [
          ["红色", "S", 101],
          ["红色", "M", 102],
          ["蓝色", "S", 203],
          ["蓝色", "M", 204]
        ]
      };
      this.addProduct(testproduct);
      return this.addProduct({
        id: 2,
        name: "鲸鱼宝宝",
        price: 138,
        note: "no note",
        quantity: [
          ["默认颜色", "均码", 1098]
        ]
      });
    };

    return OrderDetailsForm;

  })();

  $(function () {
    return new OrderDetailsForm;
  });

}).call(this);
