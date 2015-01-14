//
// Time-stamp: <[inventory_product_input.js] Elivoa @ Wednesday, 2015-01-14 19:36:01>

//
// TODO Rewrite this using Directive.
//

// parameters in scope:
// $scope.query -- query
// $scope.product -- current selecte product model.
// $scope.products -- added products. and it's amount.

// app is passed from page's config;
function $InventoryProductInput(app, $master){

  app.controller('InventoryProductInputCtrl', function($scope,$rootScope,$http){
    console.log("init InventoryProductInput component...");

    // if ($scope.selector == undefined){
    //   $scope.selector = {};
    // };
    // var selector = $scope.selector;

    $scope.init = function() {
      // init values
      $scope.query = $master.query;
      // $scope.Inventories = angular.copy($master.Inventories);
    };
    $scope.init();

    // 1. when change occured in `query` box.
    $scope.$watch('query', function(newValue, oldValue) {
      var trimedValue = $.trim(newValue);
      if (trimedValue == ""){
        // need to clear PKU table;
        return $scope.refreshCST(); // call with empty parameter to clear.
      }
      // ajax send
      var responsePromise = $http.get("/api/suggest:product?query="+ trimedValue);
      responsePromise.success(function(data, status, headers, config) {
        $scope.refreshCandidates(data);
      });
      responsePromise.error(function(data, status, headers, config) {
        alert("AJAX failed!");
      });

      return undefined;
    });

    // 2. should show a list of candidates, go to select one.
    $scope.refreshCandidates  = function(data){
      if (data && data.suggestions && angular.isArray(data.suggestions)){
        if (data.suggestions.length == 0 ){
          console.log("~~~~~~~~~~ no suggestion ~~~~~~~~~~~~~");
          return $scope.refreshCST(); // call with empty parameter to clear.
        }else{
          // TODO! Fake select one. TODO make this real select;
          var first = data.suggestions[0];
          $scope.refreshCST(data.suggestions[0].id);
        }
      }

      // fake here to directly call the first value.
      return false;
    };

    $scope.refreshCST = function(productId){
      if(productId == undefined){
        $scope.product = undefined;
        return;
      }
      // call service to get product details;
      $http.get("/api/product/"+ productId).success(function(data, status, headers, config) {
        $scope.product = data;

        // build stocks
        $scope.stocks = {};
        if (angular.isArray(data.Colors)){
          for (i=0;i<data.Colors.length;i++){
            // console.log("colors: ", data.Colors[i]);
            var color  = data.Colors[i];
            for (j=0;j<data.Sizes.length;j++){
              var size  = data.Sizes[j];
              if ($scope.stocks[color] == undefined){
                $scope.stocks[color] = {};
              }
              $scope.stocks[color][size] = 0;
            }
          }
        }

        // TODO should load it's price;

      }).error(function(data, status, headers, config) {
        alert("AJAX failed!");
      });
    };

    // get stock
    $scope.stock = function(color,size){
      if ($scope != undefined && $scope.stocks[color]!=undefined ){
        return $scope.stocks[color][size];
      }
      return 0;
    };

    // click addToInventory
    $scope.addToInventory = function(){
      if($scope.product == undefined){
        alert("请选择商品先！"); // TODO make alert more humanreadable!
        return;
      }
      // Here need to change product into inventories;
      var p = $scope.product;
      var inventory = {
        Id         :0,
	    GroupId    :0,         // TODO
	    ProductId  : p.Id,
	    Stocks     : $scope.stocks, // This is stocks matrix;
	    ProviderId : 0,        // factory person id.
	    OperatorId : 0,        // TODO
	    Note       : $scope.Note,

        // extended
        Product : p

        // ~ the following items 
	    // Status     : 0,
	    // Type       : 0,
	    // SendTime    time.Time // 发货时间
	    // ReceiveTime time.Time // 收到货的时间
	    // CreateTime  time.Time // 创建收货单的时间
	    // UpdateTime  time.Time

      };
      if ($scope.AddToProducts){
        $scope.AddToProducts(inventory); // call upper function
      }else{
        console.log('WARRNING! No AddToProducts method found!');
      }
    };

    $scope.setStock = function(color, size, $event){
      $scope.stocks[color][size] = $event.target.value; // TODO is this right?
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

