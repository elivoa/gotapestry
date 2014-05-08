##################################################
# Order Print
# Author: Elivoa@gmail.com
##################################################

window.OrderPrint =
class OrderPrint
  constructor:(config) ->
    @containerClass = ".order-form-container"
    # all data in product.
    @data = {
      order: []
      products: {}
    }
    @refreshOrderForm()

    ## another ::: printorder button
    # TODO PrintOrder set order status to delivering

  setData:(json) ->
    @data = json
    @refreshOrderForm()

  # generate the whole order form
  refreshOrderForm: ->
    tbody = $("#{@containerClass} tbody")
    tbody.html("") # clear the form

    # calculate part
    sumQuantity = 0
    sumPrice = 0

    # render html
    idx = 1
    for id in @data.order
      product = @data.products[id]
      if product
        # calculate part
        totalQuantity = 0
        totalQuantity += q[2] for q in product.quantity
        totalPrice = totalQuantity * product.price
        sumQuantity += totalQuantity
        sumPrice += totalPrice
        # render html
        tr = $(@generateTR(product, totalQuantity, totalPrice, idx++).join("\n"))
        tbody.append(tr)
        tbody.append("<tr><td colspan=\"8\" class=\"line\">&nbsp;</td></tr>")
        # @bindAction(tr, id)
      else
        console.log "[OrderDetailsForm] Error id in order list #{id}." if console

    # append footer
    footers = @appendFooter(sumQuantity, sumPrice)
    tbody.append(footers.join("\n"))

    # update footer summary quantity and price.
    @updateSummary(sumQuantity, sumPrice)

  updateSummary:(sumQuantity, sumPrice) ->
    tfoot = $("#{@containerClass} .other_info")
    tfoot.find(".sumQuantity").html(sumQuantity)
    tfoot.find(".sumPrice").html(sumPrice)

  ## hide if no stock
  # generate one product. with all it's quantities.
  # parameter is product json
  generateTR:(json, totalQuantity, totalPrice, idx)->
    quantities = []
    quantities.push q for q in json.quantity when q[2] >0
    nquantity = quantities.length
    htmls = []
    htmls.push "<tr>"
    htmls.push "  <td>#{idx}</td>"
    htmls.push "  <td valign='top' rowspan='#{nquantity}'>"
    htmls.push "    <strong>#{json.name}</strong>"
    htmls.push "  </td>"
    htmls.push "  <td valign='top' rowspan='#{nquantity}' class='money'>"
    htmls.push "    <span class='price'>#{json.price}</span>"
    htmls.push "  </td>"
    if quantities[0][0]=="默认颜色"
      htmls.push "  <td>-</td>"
    else
      htmls.push "  <td>#{quantities[0][0]}</td>"

    htmls.push "  <td>#{quantities[0][1]}</td>"
    htmls.push "  <td>#{quantities[0][2]}</td>"
    htmls.push "  <td valign='top' rowspan='#{nquantity}'>"
    htmls.push "      <strong>#{totalQuantity}</strong></td>"
    htmls.push "  <td valign='top' align='right' rowspan='#{nquantity}'>"
    htmls.push "      <strong class='price'>#{totalPrice}</strong></td>"
    if json.note == undefined || json.note==null || json.note.replace(/^\s+|\s+$/g, '') == ""
      htmls.push "  <td valign='top' rowspan='#{nquantity}'>-</td>"
    else
      htmls.push "  <td valign='top' rowspan='#{nquantity}'>#{json.note}</td>"

    htmls.push "</tr>"
    for quantity in quantities.slice(1, nquantity)
      htmls.push "<tr>"
      htmls.push "  <td></td>"
      if quantity[0]=="默认颜色"
        htmls.push "  <td>-</td>"
      else
        htmls.push "  <td>#{quantity[0]}</td>"
      htmls.push "  <td>#{quantity[1]}</td>"
      htmls.push "  <td>#{quantity[2]}</td>"
      htmls.push "</tr>"
    return htmls

  appendFooter:(sumQuantity, sumPrice) ->
    footer = []
    footer.push "<tr class='total'>"
    footer.push "<td valign='top' align='right'><strong>总计</strong></td>"
    footer.push "<td valign='top'>&nbsp;</td>"
    footer.push "<td valign='top'>&nbsp;</td>"
    footer.push "<td valign='top'>&nbsp;</td>"
    footer.push "<td valign='top'>&nbsp;</td>"
    footer.push "<td valign='top'>&nbsp;</td>"
    footer.push "<td valign='top' align='center' rowspan='1'><strong>#{sumQuantity}</strong></td>"
    footer.push "<td valign='top' align='right' rowspan='1'><strong class='price'>#{sumPrice}</strong></td>"
    footer.push "     <td valign='top'>&nbsp;</td>"
    footer.push " </tr>"
    return footer
  #######################################################
  # Test Data || Example prudoct json
  #######################################################
  addTestData:->
    testproduct = {
      id:1,
      name:"绣虎头",
      price:138,
      productPrice : 120,
      note: "no note",
      colors: ["红色", "蓝色"]
      sizes: ["S", "M"]
      quantity:[
        ["红色", "S", 101]
        ["红色", "M", 102]
        ["蓝色", "S", 203]
        ["蓝色", "M", 204]
      ]
    }
    @addProduct(testproduct)

    # product 2
    @addProduct {
      id:2,
      name:"鲸鱼宝宝",
      price:138,
      note: "no note",
      quantity:[
        ["默认颜色","均码", 1098]
      ]
    }

# $ ->
#   new OrderDetailsForm

