/*
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
*/


(function() {
  var SuggestControl;

  window.SuggestControl = SuggestControl = (function() {
    function SuggestControl(param) {
      this.param = param;
      this["default"]("parentClass", ".parent-class");
      this["default"]("hiddenClass", ".suggest-id");
      this["default"]("triggerClass", ".suggest-display");
      this.selection;
    }

    SuggestControl.prototype["default"] = function(key, defaultValue) {
      if (!this.param[key]) {
        return this.param[key] = defaultValue;
      }
    };

    SuggestControl.prototype.init = function() {
      var _;
      if (console) {
        console.log('init suggest');
      }
      _ = this;
      this.params = {
        serviceUrl: "/api/suggest/" + _.param.category,
        onSearchStart: function(query) {
          var parent;
          parent = $(this).parents(_.param["parentClass"]);
          return parent.find(_.param["hiddenClass"]).val('');
        },
        onSelect: function(suggestion) {
          var parent;
          parent = $(this).parents(_.param["parentClass"]);
          _.select(suggestion.data, _.formatValue(suggestion), this);
          if (_.param["onSelect"]) {
            return _.param["onSelect"](parent, suggestion);
          }
        },
        formatResult: function(suggestion, currentValue) {
          return _.formatValue(suggestion);
        }
      };
      return $(this.param["triggerClass"]).each($.proxy(this.eachDisplayClass, this));
    };

    SuggestControl.prototype.select = function(id, name, obj) {
      this.selection = {
        id: id,
        name: name
      };
      return this._setSelection(id, name, obj);
    };

    SuggestControl.prototype.clearSelect = function(obj) {
      this.selection = null;
      return this._setSelection('', '', obj);
    };

    SuggestControl.prototype._setSelection = function(id, name, obj) {
      var parent;
      if (obj) {
        parent = $(obj).parents(this.param["parentClass"]);
      } else {
        parent = $(this.param["parentClass"]);
      }
      parent.find(this.param["hiddenClass"]).val(id);
      return parent.find(this.param["triggerClass"]).val(name);
    };

    SuggestControl.prototype.eachDisplayClass = function(idx, obj) {
      return this.register(obj);
    };

    SuggestControl.prototype.onClick = function(e) {
      this.register(e.target);
      return $(e.target).off(e);
    };

    SuggestControl.prototype.register = function(displayInput) {
      return $(displayInput).autocomplete(this.params);
    };

    SuggestControl.prototype.registerLine = function(parnetLine) {
      return this.register($(parnetLine).find(this.param["triggerClass"]));
    };

    SuggestControl.prototype.formatValue = function(suggestionItem) {
      // TODO: move this into style sheet files.
      if (suggestionItem.t ==3){
        var item = [];
        item.push("<span style='width:52px;color:#333;display:block;float:left;white-space:initial;'>");
        item.push(suggestionItem.id);
        if(!suggestionItem.id){
          item.push("----");
        }
        item.push("</span>");

        item.push("<span style='display:block;float:left;'>");
        item.push(suggestionItem.value.substr(suggestionItem.value.indexOf("||") + 2));
        item.push("</span>");
        
        return item.join("");
      }else{
        return suggestionItem.value.substr(suggestionItem.value.indexOf("||") + 2);
      }
    };

    return SuggestControl;

  })();

}).call(this);
