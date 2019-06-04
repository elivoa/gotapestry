// Basic Functions
// Time-stamp: <[core.js] Elivoa @ Tuesday, 2015-01-27 22:33:49>

// Array Remove - By John Resig (MIT Licensed)
Array.prototype.remove = function (from, to) {
  var rest = this.slice((to || from) + 1 || this.length);
  this.length = from < 0 ? this.length + from : from;
  return this.push.apply(this, rest);
};


// Angularjs Helper methods

function fillFormNameWithNGModel(form) {
  $(form).find("input").each(function (idx, t) {
    fillElementWithNG(t);
  });
  $(form).find("textarea").each(function (idx, t) {
    fillElementWithNGNonInput(t);
  });
}

function fillElementWithNG(target) {
  tt = $(target);
  tp = tt.attr("type");
  if (tp == "text" || tp == "hidden" || tp == "textarea" || tp == "number" || tp == "date") {
    if (tt.attr("ng-model") != undefined && tt.attr("name") == undefined) {
      tt.attr("name", tt.attr("ng-model"));
    }
  }
}

function fillElementWithNGNonInput(target) {
  tt = $(target);
  if (tt.attr("ng-model") != undefined && tt.attr("name") == undefined) {
    tt.attr("name", tt.attr("ng-model"));
  }
}

// TODO make this to Angular filter
function parseGoDate(goDateString) {
  if (goDateString != undefined) {
    try {
      date = goDateString.replace(/[TZ]/g, " ");
      return new Date(Date.parse(date));
    } catch (error) {
      consle.log("Error when parse date", goDateString);
    }
  }
  return undefined;
}


function fix2(value) {
  if (!value) {
    return ""
  }
  if (typeof value === 'string') {
    return value
  }
  return value.toFixed(2);
}
