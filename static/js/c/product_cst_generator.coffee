################################################################################
# product color-size table generator
################################################################################

window.ProductCSTableGenerator =
class ProductColorSizeTableGenerator
  constructor: (colors, sizes) ->

    # clear parameters
    if colors.length == 0 then @colors=["默认颜色"] else @colors = colors
    if sizes.length == 0 then @sizes=["均码"] else @sizes = sizes

    # generate html
    @generateHtml()

  generateHtml: ->
    #
    htmls = []
    htmls.push "<!-- Generated Table, input quantity here -->"
    htmls.push "<table class=\"prd_tbl\">"
    htmls.push " <tr>"
    htmls.push "  <th align=\"left\">颜色</th>"
    htmls.push "  <th align=\"left\">尺码</th>"
    htmls.push "  <th align=\"left\">数量</th>"
    htmls.push " </tr>"

    nColors = @colors.length
    nSizes = @sizes.length
    for color in @colors
      htmls.push " <tr>"
      htmls.push "  <td rowspan=\"#{nSizes}\">#{color}</td>"
      htmls.push "  <td>#{@sizes[0]}</td>"
      htmls.push "  <td><input type=\"text\" name=\"Stocks\" id=\"csq_#{color}__#{@sizes[0]}\" class=\"stock\" size=\"8\" value=\"\"></td>"
      htmls.push " </tr>"

      for size in @sizes.slice(1,@sizes.length)
        htmls.push " <tr>"
        htmls.push "  <td>#{size}</td>"
        htmls.push "  <td><input type=\"text\" name=\"Stocks\" id=\"csq_#{color}__#{size}\" class=\"stock\" size=\"8\" value=\"\"></td>"
        htmls.push " </tr>"

    htmls.push "</table>"
    @html = htmls.join("\n")

  replace: (divId)->
      $("##{divId}").html(@html)


