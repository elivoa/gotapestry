###
  SYD System
  @author: Bo Gao, [elivoa@gmail.com]
###

################################################################################
# onload
################################################################################

class EditableTable
  constructor: (containerId)->
    @id = containerId
    @config = {
      mtContainer: ".mt-container"
      mtLine: ".mt-line"
      mtAddButton: ".mt-add"
      mtRemoveButton: ".mt-remove"
      mtInsertPlace : ".mt-insert-here" # ~ NotUsed ~
      mtInsertDirection: "above" # [above|after]
    }

    # callbacks
    @onInit
    @onRemove
    @onRemoveLastLine = @defaultRemoveLastLine
    @onNewline # can clear some values, default clear all input's value

    @registerEvent()

  registerEvent: ->
    console.log "register object id is: #{@id}"
    _=@

    $("##{@id} #{@config.mtAddButton}").on "click", $.proxy @addline,@
    $("##{@id}").on "click", @config.mtRemoveButton, $.proxy @removeline,@

  addline: (e)->
    e.preventDefault()
    lines = $("##{@id} #{@config.mtLine}")
    if lines.length > 0
      # clone the last line
      line = $(lines[lines.length-1])
      newline = line.clone()

      # empty values, expose callback functions
      newline.find('input').each (idx, obj) ->
        $(obj).attr "value",""
        $(obj).val ""
      @onNewline(newline) if @onNewline # callback

      # attach to table
      newline.insertAfter(line)

      # register suggest control
      @onAdd(newline) if @onAdd

  removeline: (e)->
    e.preventDefault()
    line = $(e.target).parents(@config.mtLine)
    if $("##{@id} #{@config.mtLine}").length > 1
      @onRemove(line) if @onRemove #callback
      line.detach()
    else
      @onRemoveLastLine(line) if @onRemoveLastLine #callback

  # Default Events
  defaultRemoveLastLine: (line)->
    alert("You can't delete the last item!")

window['EditableTable'] = EditableTable
