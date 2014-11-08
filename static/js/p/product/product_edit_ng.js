// ProductEdit
// Time-stamp: <[product_edit_ng.js] Elivoa @ Saturday, 2014-11-08 15:51:19>

//
// $master.Product -- product json.
// $master.Colors  -- Colors [{Value:xxx}, {Value:xxx},...] structure.
// $master.Sizes   -- The same with Colors
//
function p_ProductEdit($master){

  var sydapp = angular.module('syd', [], function($interpolateProvider){
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
  });

  sydapp.controller('ProductCtrl', function($scope){
    $scope.init = function() {
      // init values
      $scope.Product = angular.copy($master.Product);
      $scope.Colors = angular.copy($master.Colors);
      $scope.Sizes = angular.copy($master.Sizes);
      $scope.Stocks = transforStocks($scope);
      // bind methods. bind in html.
    };
    $scope.init();

    // functions
    function transforStocks($scope){
      // init struct
      stocks = {};
      for (i=0;i<$scope.Colors.length;i++){
        var color  = $scope.Colors[i].Value;
        for (j=0;j<$scope.Sizes.length;j++){
          var size  = $scope.Sizes[j].Value;
          if (stocks[color] == undefined){
            stocks[color] = {};
          }
          stocks[color][size] = 0;
        }
      }
      if ($scope.Product.Stocks != undefined){
        // fill by stocks;
        for (i=0;i<$scope.Product.Stocks.length;i++){
          item = $scope.Product.Stocks[i];
          if (stocks[item.Color]!=undefined){
            if (stocks[item.Color][item.Size] !=undefined){
              stocks[item.Color][item.Size] = item.Stock;
              continue;
            }
          }
          console.log("NOTE: SKU溢出", item.Color, item.Size, item.Stock);
        }
      }
      return stocks;
    }

    $scope.stock = function(color,size){
      if ($scope != undefined && $scope.Stocks[color]!=undefined ){
        return  $scope.Stocks[color][size];
      }
      return 0;
    };

    // events
    $scope.addColor = function(){
      $scope.Colors.push({Value:""});
    };
    $scope.addSize = function(){
      $scope.Sizes.push({Value:""});
    };
    $scope.removeColor = function(idx){
      $scope.Colors.remove(idx,idx);
    };
    $scope.removeSize = function(idx){
      $scope.Sizes.remove(idx,idx);
    };

    $scope.submit = function() {
      fillFormNameWithNGModel(ProductForm);
    };


    // TEST --------------------------------------------------
    $scope.changeData = function(){
      $scope.data[3].client='我要扯淡扯淡';
    };

    $scope.test = function(){
      console.log($scope.Colors);
    };

  });

}
