// ProductEdit
// Time-stamp: <[inventory_edit_ng.js] Elivoa @ Thursday, 2015-01-22 14:30:03>

// Development Notes:
// Treat all sub components as one html page, use component just split them.
// Later use directive as really component of angularjs.

//
// $scope.Inventories    -- inventories equals to InventoryGroup.Inventories
// $scope.InventoryMap   -- id -> Inventory map
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

    // set master variables into $scope
    if($master.InventoryGroup!=undefined ){
      $scope.initInventories($master.InventoryGroup.Inventories);
      $scope.InventoryGroup = $master.InventoryGroup;
    }
    $scope.Factories = $master.Factories;
    $scope.SendTime = Date.now();


    // nothing is useful in this page. all things is in components / directives;
    $scope.submit = function() {
      fillFormNameWithNGModel(ProductForm);
    };

  });

}


