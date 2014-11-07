// ProductEdit
// Time-stamp: <[product_edit_ng.js] Elivoa @ Thursday, 2014-11-06 16:47:55>

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

      // bind methods. bind in html.
    };
    $scope.init();

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
