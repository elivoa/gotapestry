// Time-stamp: <[angular.init.js] Elivoa @ Sunday, 2015-01-11 14:32:05>
// Init file of got framework using angularjs framework;

var got = window.got || (window.got = {});

// spec: 初始化funtino的结构
/*
 spec : {name: 'testcomponent', init: function(){...}, other parameters}
 */
function ngRegisterComponent(spec){
  // init ngcomps
  var ngcomps = window['_ng_components_'];
  if(ngcomps == undefined){
    ngcomps = new Array();
    window['_ng_components_'] = ngcomps;
  }

  // TODO check if initfunc is a function.
  ngcomps.push(spec);
}


// Call init function when page has components;
function ngLoadComponent(app){
  var ngcomps = window['_ng_components_'];
  if(ngcomps!=undefined){
    // TODO check if ngcomps is not array;
    for(var i=0;i<ngcomps.length;i++){
      // call components' init method.
      if (ngcomps[i]['init']!=undefined){
        ngcomps[i]['init'](app);
      }
    }
  }

  console.log("load components...");
}

