###
---------------------------------------------
 SYD Project
 Suggest Control

  param:
    parnetClass    - line container
    id             -
    category       - suggest category
    onSelect       - on select callback

  hard-coded param:
    suggestId = suggest-id
    suggestDisplay = suggets-display

  TODO:
    reset value if canceled.
----------------------------------------------
###

window.SuggestControl =
class SuggestControl
  constructor: (param)->
    @param = param
    @default "parentClass", ".parent-class"
    @default "hiddenClass", ".suggest-id"
    @default "triggerClass", ".suggest-display"

    # PUBLIC:
    @selection # id{id:, name:}


  default: (key, defaultValue)->
    if !@param[key]
      @param[key] = defaultValue

  init: ()-> # this params is plugins' params
    console.log 'init suggest' if console

    _ = @
    @params = {
      serviceUrl: "/api/suggest/" + _.param.category
      # lookup: this.countries,

      onSearchStart: (query)->
        parent = $(this).parents(_.param["parentClass"])
        parent.find(_.param["hiddenClass"]).val('')

      onSelect: (suggestion)->
        parent = $(this).parents(_.param["parentClass"])
        _.select(suggestion.data, _.formatValue(suggestion.value), this)
        _.param["onSelect"](parent, suggestion) if _.param["onSelect"]

      formatResult:(suggestion, currentValue) ->
        # TODO highlight matc3hes...
        return _.formatValue(suggestion.value)
    }

    # register init items.
    $(@param["triggerClass"]).each $.proxy @eachDisplayClass,@

  # PUBLIC
  select:(id, name, obj) -> # optional obj
    @selection = {id: id, name: name}
    @_setSelection(id, name, obj)

  # PUBLIC
  clearSelect: (obj)-> # optional objN
    @selection = null
    @_setSelection('', '', obj)

  # set selection
  _setSelection: (id, name, obj) ->
    if obj
      parent = $(obj).parents(@param["parentClass"])
    else
      parent = $(@param["parentClass"])
    parent.find(@param["hiddenClass"]).val(id)
    parent.find(@param["triggerClass"]).val(name)


  eachDisplayClass: (idx, obj)->
    @register obj

  # onClick method on display input box.
  onClick: (e) -> # no use
    @register(e.target)
    $(e.target).off(e)

  # register display input.
  register: (displayInput)->
    $(displayInput).autocomplete(@params)

  # register display input.
  registerLine: (parnetLine)->
    @register $(parnetLine).find(@param["triggerClass"])

  # format the composed value
  formatValue:(value)->
    return value.substr(value.indexOf("||") + 2)


#reEscape = new RegExp('(\\' + ['/', '.', '*', '+', '?', '|', '(', ')', '[', ']', '{', '}', '\\'].join('|\\') + ')', 'g'),
#pattern = '(' + currentValue.replace(reEscape, '\\$1') + ')';
#return ">>> " + suggestion.value.replace(new RegExp(pattern, 'gi'), '<strong>$1<\/strong>')
