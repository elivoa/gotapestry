//
// Time-stamp: <[inventory_list.js] Elivoa @ Friday, 2015-01-02 23:23:27>

//
// TODO need to rewrite.
// $master.Product -- product json.
// $master.Colors  -- Colors [{Value:xxx}, {Value:xxx},...] structure.
// $master.Sizes   -- The same with Colors
//
function $InventoryList($master){
  console.log("init inventory_list.js ...");

  var sydapp = angular.module('syd', [], function($interpolateProvider){
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
  });

  sydapp.controller('InventoryListCtrl', function($scope){
    $scope.init = function() {
      // init values
      $scope.Inventories = angular.copy($master.Inventories);
      // bind methods. bind in html.
    };
    $scope.init();

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
