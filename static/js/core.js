// Basic Functions
// Time-stamp: <[core.js] Elivoa @ Thursday, 2014-11-06 17:30:06>

// Array Remove - By John Resig (MIT Licensed)
Array.prototype.remove = function(from, to) {
  var rest = this.slice((to || from) + 1 || this.length);
  this.length = from < 0 ? this.length + from : from;
  return this.push.apply(this, rest);
};


// Angularjs Helper methods

function fillFormNameWithNGModel(form){
  $(form).find("input").each(function(idx, t){
    fillElementWithNG(t);
  });
  $(form).find("textarea").each(function(idx, t){
    fillElementWithNG(t);
  });
}

function fillElementWithNG(target){
  tt = $(target);
  tp = tt.attr("type");
  if (tp == "text" || tp == "hidden" || tp == "textarea" || tp == "number" ){
    if (tt.attr("ng-model") != undefined && tt.attr("name") == undefined){
      tt.attr("name", tt.attr("ng-model"));
    }
  }
}

