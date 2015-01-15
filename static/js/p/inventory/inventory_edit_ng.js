// ProductEdit
// Time-stamp: <[inventory_edit_ng.js] Elivoa @ Thursday, 2015-01-15 23:09:52>

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

    // 将inventories 初始化到系统内部格式
    $scope.initInventories = function(invs){
      var nvs = []; // new inventories
      var idmap = {};
      if (angular.isArray(invs)){
        for (i=0;i<invs.length;i++){
          var inv = invs[i];
          if(inv!=undefined){
            var currentInv = idmap[inv.ProductId];
            if(currentInv == undefined){
              idmap[inv.ProductId] = currentInv = inv;
              nvs.push(currentInv);
            }
            // add stocks
            setStocks(currentInv, inv.Color, inv.Size, inv.Stock);
          }
        }
      }

      function setStocks(inv, color, size, stock){
        if(inv.Stocks == undefined){
          inv.Stocks = {};
        }
        if(inv.Stocks[color] == undefined){
          inv.Stocks[color] = {};
        }
        inv.Stocks[color][size] = stock;
      }

      // final assign
      $scope.Inventories = nvs;
      $scope.InventoryMap = idmap;
    };


    if($master.InventoryGroup!=undefined ){
      $scope.initInventories($master.InventoryGroup.Inventories);
    }

    // nothing is useful in this page. all things is in components / directives;
    $scope.submit = function() {
      fillFormNameWithNGModel(ProductForm);
    };

  });

}


