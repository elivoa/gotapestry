// Time-stamp: <[angular.init.js] Elivoa @ Monday, 2015-06-08 13:44:46>
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


// Deprecated!! Call init function when page has components;
// Usages: ngcLoadTempalte($http,"order-close-button.html");
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

// ngc Load tempalte
function ngcLoadTempalte($http,template){
  $http.get('/static/ngc/'+template).
    success(function(data, status, headers, config) {
      angular.element("body").append(data);
      // this callback will be called asynchronously
      // when the response is available
    }).
    error(function(data, status, headers, config) {
      console.log("load something error!");
      // called asynchronously if an error occurs
      // or server returns response with an error status.
    });
}


// Page loading information;
// 这个js只能被引用一次
// 这个js会初始化页面app，只在允许ng的页面使用；

// app
var app = angular.module('app', [], function($interpolateProvider){
  // TODO move this into global config;
  $interpolateProvider.startSymbol('[[');
  $interpolateProvider.endSymbol(']]');
});
