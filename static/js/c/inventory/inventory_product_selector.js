//
// Time-stamp: <[inventory_product_selector.js] Elivoa @ Thursday, 2015-01-15 22:52:35>

//
// TODO need to rewrite.
// $master.Product -- product json.
// $master.Colors  -- Colors [{Value:xxx}, {Value:xxx},...] structure.
// $master.Sizes   -- The same with Colors
//

// app is passed from page's config;
function $InventoryProductSelector(app, $master){
  console.log("init inventory_product_selector.js ...");

  app.controller('InventoryProductSelectorCtrl', function($scope, $rootScope){

    // add global functions.
    $scope.AddToProducts = function(inventory){
      initInv();
      if($scope.InventoryMap[inventory.ProductId]==undefined){
        $scope.InventoryMap[inventory.ProductId] = inventory;
        $scope.Inventories.push(inventory);
      }else{
        $rootScope.errmsg = "重复添加同一个商品, 更新现有数据";
        $scope.InventoryMap[inventory.ProductId].Stocks = inventory.Stocks; // update stock value;
      }
    };

    $scope.GetInventory = function(productId) {
      if($scope.InventoryMap !=undefined){
        inv = $scope.InventoryMap[productId];
        if(inv!=undefined){
          return inv;
        }
      }
      return undefined;
    };

    function initInv(){
      if($scope.Inventories==undefined){
        $scope.Inventories = [];
      }
      if($scope.InventoryMap==undefined){
        $scope.InventoryMap = {};
      }

    }

    $scope.submit = function() {
      fillFormNameWithNGModel(ProductForm);
    };

  });

}
