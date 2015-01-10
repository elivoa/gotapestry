//
// Time-stamp: <[inventory_product_selector.js] Elivoa @ Saturday, 2015-01-10 15:29:36>

//
// TODO need to rewrite.
// $master.Product -- product json.
// $master.Colors  -- Colors [{Value:xxx}, {Value:xxx},...] structure.
// $master.Sizes   -- The same with Colors
//

// app is passed from page's config;
function $InventoryProductSelector(app, $master){
  console.log("init inventory_product_selector.js ...");

  app.controller('InventoryProductSelectorCtrl', function($scope){

    $scope.init = function() {
      // init values
      $scope.Inventories = angular.copy($master.Inventories);

      $scope.quickinput = "please input capital letter!";

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
