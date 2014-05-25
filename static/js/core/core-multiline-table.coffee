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

    ## callbacks
    @onInit

    # Call on each line remove
    @onRemove
    @afterRemove

    # Before remove last line; return true to remove last line
    @onRemoveLastLine = @defaultRemoveLastLine

    # can clear some values, default clear all input's value
    @onNewline

    @registerEvent()

  registerEvent: ->
    _=@

    $("##{@id} #{@config.mtAddButton}").on "click", $.proxy @addline,@
    $("##{@id}").on "click", @config.mtRemoveButton, $.proxy @removeline,@

  addline: (e)->
    e.preventDefault()
    lines = $("##{@id} #{@config.mtLine}")
    if lines.length > 0
      # clone the last line
      line = $(lines[0]) # $(lines[lines.length-1]) # last line
      newline = line.clone()
      line.css {display:""}

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
    nItems = $("##{@id} #{@config.mtLine}").length

    # call events
    @onRemove(line) if @onRemove
    shouldRemove = true # always remove, if event returns true
    if nItems == 1
      # callback events
      shouldRemove = @onRemoveLastLine(line) if @onRemoveLastLine
    if shouldRemove == true
      @onRemove(line) if @onRemove #callback
      line.detach()
    @afterRemove(line) if @afterRemove

  # Default Events
  defaultRemoveLastLine: (line)->
    alert("You can't delete the last item!")

window['EditableTable'] = EditableTable
