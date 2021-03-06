// Generated by CoffeeScript 1.6.3
/*
  SYD Sales System
  @author: Bo Gao, [elivoa@gmail.com]
*/


(function() {
  var Gotinit;

  Gotinit = (function() {
    function Gotinit() {
      console.log('got initialized');
    }

    Gotinit.prototype.initSelect = function() {
      return $('select').each(function(index, selectObj) {
        var target, value;
        target = $(selectObj);
        value = target.attr('value');
        if (value) {
          return target.find("option").each(function(idx, option) {
            if (value === option.value) {
              option.selected = true;
              return true;
            }
          });
        }
      });
    };

    Gotinit.prototype.initRadio = function() {
      $('.auto-radio').each(function(idx, obj) {
        var v;
        v = $(obj).attr("value");
        return $(obj).find('>input[type="radio"]').each(function(i, r) {
          if ($(r).attr('value') === v) {
            return $(r).prop('checked', true);
          }
        });
      });
      return 1;
    };

    return Gotinit;

  })();

  $(function() {
    var init;
    init = new Gotinit;
    init.initSelect();
    init.initRadio();
    return true;
  });

}).call(this);
