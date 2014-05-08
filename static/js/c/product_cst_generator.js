// Generated by CoffeeScript 1.6.3
(function() {
  var ProductColorSizeTableGenerator;

  window.ProductCSTableGenerator = ProductColorSizeTableGenerator = (function() {
    function ProductColorSizeTableGenerator(colors, sizes, tid, editable, data) {
      if (colors.length === 0) {
        this.colors = ["默认颜色"];
      } else {
        this.colors = colors;
      }
      if (sizes.length === 0) {
        this.sizes = ["均码"];
      } else {
        this.sizes = sizes;
      }
      this.tid = tid === void 0 ? "..." : tid;
      this.editable = editable === void 0 ? true : editable;
      this.generateHtml();
    }

    ProductColorSizeTableGenerator.prototype.generateHtml = function() {
      var color, htmls, nColors, nSizes, size, _i, _j, _len, _len1, _ref, _ref1;
      htmls = [];
      htmls.push('<table class="tbl_s">');
      htmls.push('  <tr>');
      htmls.push('    <th align="left">颜色</th>');
      htmls.push('    <th align="left">尺码</th>');
      htmls.push('    <th align="left">数量</th>');
      htmls.push('  </tr>');
      nColors = this.colors.length;
      nSizes = this.sizes.length;
      _ref = this.colors;
      for (_i = 0, _len = _ref.length; _i < _len; _i++) {
        color = _ref[_i];
        htmls.push("  <tr>");
        htmls.push("    <td rowspan=\"" + nSizes + "\">" + color + "</td>");
        htmls.push("    <td>" + this.sizes[0] + "</td>");
        htmls.push("    <td><input type=\"text\" size=\"8\" name=\"Stocks\" id=\"csq_" + color + "__" + this.sizes[0] + "\" class=\"stock\" autocomplete=\"off\"></td>");
        htmls.push("  </tr>");
        _ref1 = this.sizes.slice(1, this.sizes.length);
        for (_j = 0, _len1 = _ref1.length; _j < _len1; _j++) {
          size = _ref1[_j];
          htmls.push("  <tr>");
          htmls.push("    <td>" + size + "</td>");
          htmls.push("    <td><input type=\"text\" size=\"8\" name=\"Stocks\" id=\"csq_" + color + "__" + size + "\" class=\"stock\" autocomplete=\"off\"></td>");
          htmls.push("  </tr>");
        }
      }
      htmls.push("</table>");
      this.html = htmls.join("\n");
      return $("#" + this.tid).html(this.html);
    };

    return ProductColorSizeTableGenerator;

  })();

}).call(this);
