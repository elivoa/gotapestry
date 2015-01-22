//
// Time-stamp: <[inventory_product_input.js] Elivoa @ Wednesday, 2015-01-21 23:36:40>

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

    $scope.query = $master.query;

    // 1. when change occured in `query` box.
    $scope.$watch('query', function(newValue, oldValue) {
      var trimedValue = $.trim(newValue);
      if (trimedValue == ""){
        // need to clear PKU table;
        return $scope.refreshCST(); // call with empty parameter to clear.
      }
      // ajax send
      $http.get("/api/suggest:product?query="+ trimedValue)
        .success(function(data, status, headers, config) {
          $scope.refreshCandidates(data);
        })
        .error(function(data, status, headers, config) {
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

    // 3. refresh PKU Stock table;
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

        // 如果Product重复，将stocks的值还原；
        if($scope.GetInventory){
          var existInv = $scope.GetInventory(productId);
          if (existInv!=undefined){
            var stocks = existInv.Stocks;
            for(color in stocks){
              var sizemap = stocks[color];
              if(sizemap != undefined){
                for(size in sizemap){
                  $scope.stocks[color][size] = sizemap[size];
                }
              }
            }
          }
        }else{
          console.log("Warrning: no GetProduct method found!");
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
      $scope.errmsg = ""; // clear error message.
      if($scope.product == undefined){
        // alert("请选择商品先！"); // TODO make alert more humanreadable!
        $scope.errmsg = "请选择商品先！";
        $scope.focusQuery();
        return;
      }

      // Here need to change product into inventories;
      var p = $scope.product;
      var inventory = {
        Id         :0,
	    GroupId    :0,         // TODO
	    ProductId  : p.Id,
	    Stocks     : angular.copy($scope.stocks), // This is stocks matrix;
	    ProviderId : 0,        // factory person id.
	    OperatorId : 0,        // TODO
	    Note       : $scope.Note,
        Product    : p         // extened
      };

      // calculate sum stock if le 0, don't allow to add.
      var sumStock = 0;
      if (inventory.Stocks!=undefined){
        var colors = Object.keys(inventory.Stocks);
        for(i=0;i<colors.length;i++){
          var color = colors[i];
          var sizemap = inventory.Stocks[color];
          var sizes = Object.keys(sizemap);
          for(j=0;j<sizes.length;j++){
            var size = sizes[j];
            var stock = sizemap[size];
            // TODO convert to number;
            if(stock>0){
              sumStock+=stock;
            }
          }
        }
      }
      if (sumStock<=0){
        // alert("库存必须大于0!");
        $scope.errmsg = "库存必须大于0";
        return;
      }

      // console.log("stocks for display is : ", stocksForDisplay)

      if ($scope.AddToProducts){
        $scope.AddToProducts(inventory); // call upper function
        // add success, remove product
        $scope.product = undefined;
        $scope.query = "";
        $scope.stocks = undefined;
        $scope.focusQuery();

        // $scope.form.query.focus(); // why this does not work?
      }else{
        console.log('WARRNING! No AddToProducts method found!');
      }
    };

    // blur and keyup
    $scope.setStock = function(color, size, $event){
      var intstock = parseInt($event.target.value);
      intstock = isNaN(intstock)? 0 : intstock;
      $scope.stocks[color][size] = intstock;
    };

    // click edit on operator column
    $scope.onEdit = function(invId){
      console.log(">>>>>>>>>>>>> ", invId);
      $scope.refreshCST(invId);
    };

    $scope.submit = function() {
      fillFormNameWithNGModel(ProductForm);
    };


    // focus on query box
    $scope.focusQuery = function(){
      $('._temp_query_box').focus();
    };
  });

}

