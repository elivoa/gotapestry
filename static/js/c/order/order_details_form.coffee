##################################################
# Order Details Form script
# Author: Elivoa@gmail.com
##################################################

window.OrderDetailsForm =
class OrderDetailsForm
  constructor:(config) ->
    @containerClass = ".order-form-container"

    @onDelete = @defaultOnDelete # delete line callback
    @onEdit # edit line callback

    # all data in product.
    @data = {
      order: []
      products: {}
    }
    @refreshOrderForm()


  # Append product
  # Return true - if success, false - if conflict with existing.
  addProduct:(product) ->
    return if not product
    if @data.products[product.id]
      alert "已经添加了这件商品，不能重复添加！如需添加或修改，请点击下面对应商品的编辑按钮！谢谢合作！"
      return false
    @data.order.push product.id
    @data.products[product.id] = product
    return true

  editProduct:(product) ->
    return if not product
    @data.products[product.id] = product
    return true

  setData:(json) ->
    console.log typeof json
    @data = json
    console.log typeof @data
    @refreshOrderForm()

  # generate the whole order form
  refreshOrderForm: ->
    tbody = $("#{@containerClass} tbody")
    tbody.html("") # clear the form

    # calculate part
    sumQuantity = 0
    sumPrice = 0

    # render html
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
        tr = $(@generateTR(product, totalQuantity, totalPrice).join("\n"))
        tbody.append(tr)
        @bindAction(tr, id)
      else
        console.log "[OrderDetailsForm] Error id in order list #{id}." if console
    # bind callback function on tr


    # update footer summary quantity and price.
    @updateSummary(sumQuantity, sumPrice)

  bindAction:(tr, id)->
    tr.find(".odf-edit").on "click", $.proxy @onODFEdit(id),@
    tr.find(".odf-delete").on "click", $.proxy @onODFDelete(id),@

  onODFEdit:(id) ->
    return (e)->
      e.preventDefault()
      @.onEdit(@data.products[id]) if @onEdit

  onODFDelete:(id) ->
    return (e)->
      e.preventDefault()
      @onDelete(@data.products[id]) if @onDelete

  defaultOnDelete:(product)->
    return if not confirm "真的要删除这条记录么？"
    delete @data.products[product.id]
    idx = @data.order.indexOf(product.id)
    if idx>=0
      console.log @data.order.splice(idx, 1)
    @refreshOrderForm()

  updateSummary:(sumQuantity, sumPrice) ->
    tfoot = $("#{@containerClass} tfoot")
    tfoot.find(".sumQuantity").html(sumQuantity)
    tfoot.find(".sumPrice").html(sumPrice)

  ## hide if no stock
  # generate one product. with all it's quantities.
  # parameter is product json
  generateTR:(json, totalQuantity, totalPrice)->
    quantities = []
    quantities.push q for q in json.quantity when q[2] >0
    nquantity = quantities.length
    htmls = []
    htmls.push "<tr>"
    htmls.push "  <td valign='top' rowspan='#{nquantity}'>"
    htmls.push "    <strong>#{json.name}</strong>"
    htmls.push "    <input type='hidden' name='Order.Details.ProductId' value='#{json.id}' />"
    htmls.push "    <input type='hidden' name='Order.Details.SellingPrice' value='#{json.price}' />"
    htmls.push "    <input type='hidden' name='Order.Details.Color' value='#{quantities[0][0]}' />"
    htmls.push "    <input type='hidden' name='Order.Details.Size' value='#{quantities[0][1]}' />"
    htmls.push "    <input type='hidden' name='Order.Details.Quantity' value='#{quantities[0][2]}' />"
    htmls.push "    <input type='hidden' name='Order.Details.Note' value='#{json.note}' />"
    htmls.push "  </td>"
    htmls.push "  <td valign='top' rowspan='#{nquantity}'>"
    htmls.push "    <span class='price'>#{json.price}</span>"
    htmls.push "  </td>"
    htmls.push "  <td>#{quantities[0][0]}</td>"
    htmls.push "  <td>#{quantities[0][1]}</td>"
    htmls.push "  <td>#{quantities[0][2]}</td>"
    htmls.push "  <td valign='top' align='center' rowspan='#{nquantity}'>"
    htmls.push "      <strong>#{totalQuantity}</strong></td>"
    htmls.push "  <td valign='top' align='right' rowspan='#{nquantity}'>"
    htmls.push "      <strong class='price'>#{totalPrice}</strong></td>"
    htmls.push "  <td valign='top' rowspan='#{nquantity}'>#{json.note}</td>"
    htmls.push "  <td valign='top' rowspan='#{nquantity}'>"
    htmls.push "      <a href='#' class='odf-edit'>编辑</a><span class='vline'>|</span>"
    htmls.push "      <a href='#' class='odf-delete'>删除</a>"
    htmls.push "  </td>"
    htmls.push "</tr>"
    for quantity in quantities.slice(1, nquantity)
      htmls.push "<tr>"
      htmls.push "  <td>#{quantity[0]}</td>"
      htmls.push "  <td>#{quantity[1]}</td>"
      htmls.push "  <td>#{quantity[2]}</td>"
      htmls.push "    <input type='hidden' name='Order.Details.ProductId' value='#{json.id}' />"
      htmls.push "    <input type='hidden' name='Order.Details.SellingPrice' value='#{json.price}' />"
      htmls.push "    <input type='hidden' name='Order.Details.Color' value='#{quantity[0]}' />"
      htmls.push "    <input type='hidden' name='Order.Details.Size' value='#{quantity[1]}' />"
      htmls.push "    <input type='hidden' name='Order.Details.Quantity' value='#{quantity[2]}' />"
      htmls.push "    <input type='hidden' name='Order.Details.Note' value='#{json.note}' />"
      htmls.push "</tr>"
    return htmls


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

$ ->
  new OrderDetailsForm

