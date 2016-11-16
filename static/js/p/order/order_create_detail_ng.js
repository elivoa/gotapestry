// OrderCreateDetail
// Time-stamp: <[order_create_detail_ng.js] Elivoa @ Friday, 2016-11-11 19:11:28>

// Development Notes:
// Treat all sub components as one html page, use component just split them.
// Later use directive as really component of angularjs.

//
// xx $scope.Inventories    -- inventories equals to InventoryGroup.Inventories
// xx $scope.InventoryMap   -- id -> Inventory map
//
function p_OrderCreateDetail($master){

  var sydapp = angular.module('syd', [], function($interpolateProvider){
    // TODO move this into global config;
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
  });

  // if has components, init it first; then init this page;
  ngLoadComponent(sydapp);

  sydapp.controller('OrderCreateDetailCtrl', function($scope){


  });
}


