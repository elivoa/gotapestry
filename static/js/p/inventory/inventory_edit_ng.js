// ProductEdit
// Time-stamp: <[inventory_edit_ng.js] Elivoa @ Wednesday, 2015-01-14 19:34:27>

// Development Notes:
// Treat all sub components as one html page, use component just split them.
// Later use directive as really component of angularjs.

//
// $master.Product -- product json.
// $master.Colors  -- Colors [{Value:xxx}, {Value:xxx},...] structure.
// $master.Sizes   -- The same with Colors
//
function p_InventoryEdit($master){

  var sydapp = angular.module('syd', [], function($interpolateProvider){
    // TODO move this into global config;
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
  });

  // if has components, init it first; then init this page;
  ngLoadComponent(sydapp);

  sydapp.controller('InventoryEditCtrl', function($scope){


    // nothing is useful in this page. all things is in components / directives;
    $scope.submit = function() {
      fillFormNameWithNGModel(ProductForm);
    };

  });

}


