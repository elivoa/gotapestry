###
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

###
class SuggestControl
  constructor: (param)->
    @param = param
    @default "parentClass", ".parent-class"
    @default "hiddenClass", ".suggest-id"
    @default "triggerClass", ".suggest-display"

  default: (key, defaultValue)->
    if !@param[key]
      @param[key] = defaultValue

  init: ()-> # this params is plugins' params
    console.log 'init suggest'

    _ = @
    @params = {
      serviceUrl: "/ajax/suggest/" + _.param.category
      # lookup: this.countries,

      onSearchStart: (query)->
        parent = $(this).parents(_.param["parentClass"])
        parent.find(_.param["hiddenClass"]).val('')

      onSelect: (suggestion)->
        parent = $(this).parents(_.param["parentClass"])
        parent.find(_.param["hiddenClass"]).val(suggestion.data)
        parent.find(_.param["triggerClass"]).val(_.formatValue(suggestion.value))
        _.param["onSelect"](parent, suggestion) if _.param["onSelect"]

      formatResult:(suggestion, currentValue) ->
        # TODO highlight matches...
        return _.formatValue(suggestion.value)
    }

    # register init items.
    $(@param["triggerClass"]).each $.proxy @eachDisplayClass,@

  eachDisplayClass: (idx, obj)->
    @register obj

  # onClick method on display input box.
  onClick: (e) -> # no use
    console.log "register" + e.target
    @register(e.target)
    $(e.target).off(e)

  # register display input.
  register: (displayInput)->
    console.log "bind" + displayInput
    $(displayInput).autocomplete(@params)

  # register display input.
  registerLine: (parnetLine)->
    @register $(parnetLine).find(@param["triggerClass"])

  # format the composed value
  formatValue:(value)->
    return value.substr(value.indexOf("||") + 2)


## expose
window.SuggestControl = SuggestControl



        # console.log suggestion
        # console.log currentValue
        #reEscape = new RegExp('(\\' + ['/', '.', '*', '+', '?', '|', '(', ')', '[', ']', '{', '}', '\\'].join('|\\') + ')', 'g'),
        #pattern = '(' + currentValue.replace(reEscape, '\\$1') + ')';
        #return ">>> " + suggestion.value.replace(new RegExp(pattern, 'gi'), '<strong>$1<\/strong>')
