// ProductEdit
// Time-stamp: <[inventory_edit_ng.js] Elivoa @ Tuesday, 2015-03-03 22:09:14>

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

    // calculate the sum of stocks from Stocks structure.
    $scope.calculateSumStocks = function(stocks){
      var sumStock = 0;
      if (stocks!=undefined){
        var colors = Object.keys(stocks);
        for(j=0;j<colors.length;j++){
          var color = colors[j];
          var sizemap = stocks[color];
          var sizes = Object.keys(sizemap);
          for(l=0;l<sizes.length;l++){
            var size = sizes[l];
            var stock = sizemap[size];
            if(stock>0){
              sumStock+=stock;
            }
          }
        }
      }
      return sumStock;
    };


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
        inv.sumStock = $scope.calculateSumStocks(inv.Stocks); // calculate stocks
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
    // console.log("parseGoDate($scope.InventoryGroup.SendTime): ", parseGoDate($scope.InventoryGroup.SendTime))
    // change it.
    $scope.InventoryGroup.SendTime = parseGoDate($scope.InventoryGroup.SendTime);
    $scope.InventoryGroup.ReceiveTime = parseGoDate($scope.InventoryGroup.ReceiveTime);
    // TODO how to add date by one day using javascript;

    // Submit my
    $scope.submit = function(form) {
      fillFormNameWithNGModel(InventoryForm);
      InventoryForm.submit();
    };

  });
}


