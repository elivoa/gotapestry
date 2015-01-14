//
// Time-stamp: <[inventory_product_selector.js] Elivoa @ Wednesday, 2015-01-14 19:14:19>

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
      // $scope.Inventories = angular.copy($master.Inventories);

      $scope.quickinput = "please input capital letter!";

      // bind methods. bind in html.
    };
    $scope.init();

    // add global functions.
    $scope.AddToProducts = function(inventory){
      if($scope.Inventories==undefined){
        $scope.Inventories = [];
      }
      $scope.Inventories.push(inventory);
      console.log(">> success add addtoproducts")
      console.log($scope.Inventories)
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

  });

}
