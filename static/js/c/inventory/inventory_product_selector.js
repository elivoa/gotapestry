//
// Time-stamp: <[inventory_product_selector.js] Elivoa @ Thursday, 2015-01-29 16:42:46>

// app is passed from page's config;
function $InventoryProductSelector(app, $master){

  app.controller('InventoryProductSelectorCtrl', function($scope,$rootScope,$http){

    $scope.query = $master.query;

    // register document's keyboard management. ?? How to do bind this method;
    // var d = angular.element("body");
    // d.on('keyDown', $scope.suggestKeycontrol;

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
          $scope.cancelSuggest();
          return $scope.refreshCST(); // call with empty parameter to clear.
        }else{
          $scope.showCandidates(data);
          // TODO! Fake select one. TODO make this real select;
          // var first = data.suggestions[0];
          // $scope.refreshCST(data.suggestions[0].id);
        }
      }

      // fake here to directly call the first value.
      return false;
    };

    // --------------------------------------------------------------------------------
    // Suggestion Area
    // --------------------------------------------------------------------------------
    $scope.showCandidates = function(data){
      $scope.suggestionIndex = 0;
      $scope.suggestionMaxItems = 0;
      if(data.suggestions!=undefined){
        $scope.candidates = data.suggestions;
        $scope.suggestionMaxItems = data.suggestions.length;
      }
    };

    $scope.suggestKeycontrol = function(e){
      console.log("keydown: keyCode is ", e.keyCode, "; Modifier:",e.Modifier,
                  "; index=", $scope.suggestionIndex);
      //     console.log(e);
      if(e.keyCode == 40 || (e.keyCode==78 && (e.ctrlKey==true))){ // arrow-down, ctrl+p
        $scope.suggestionIndex += 1;
        if ($scope.suggestionIndex >= $scope.suggestionMaxItems){
          $scope.suggestionIndex = $scope.suggestionMaxItems;
        }
      }
      if(e.keyCode == 38 || (e.keyCode==80 && (e.ctrlKey==true))){ // arrow-up
        $scope.suggestionIndex-=1;
        if ($scope.suggestionIndex <=0 ){
          $scope.suggestionIndex = 0;
        }
      }

      if(e.keyCode == 13){ // enter
        $scope.selectSuggest($scope.suggestionIndex - 1);
      }

      if(e.keyCode == 39 || (e.keyCode==70 && (e.ctrlKey==true))){ // arrow-right, ctrl+f
        $scope.showTotalInventory($scope.suggestionIndex - 1);
      }
    };

    // clear suggest candidates dropbox.
    $scope.cancelSuggest = function($event){
      $scope.suggestions = undefined;
      $scope.suggestionIndex = 0;
      $scope.suggestionMaxItems = 0;
      $scope.candidates = undefined;
    };

    $scope.selectSuggest = function(idx){
      if($scope.candidates==undefined){
        return;
      }
      var cur = $scope.candidates[idx];
      if (cur!=undefined){
        $scope.cancelSuggest();
        $scope.refreshCST(cur.id);
      }else{
        console.log("error find index ", $scope.suggestionIndex);
      }
    };

    $scope.showTotalInventory = function(idx){
      if($scope.candidates==undefined){
        return;
      }
      var cur = $scope.candidates[idx];
      if (cur!=undefined){
        $http.get("/api/product/"+ cur.id).success(function(data, status, headers, config) {
          if(data!=undefined){
            cur.totalStock = data.Stock;
          }
        }).error(function(data, status, headers, config) {
          alert("AJAX failed!");
        });
      }else{
        console.log("error find index ", $scope.suggestionIndex);
      }
    };

    $scope.suggestSelectedClass = function(idx){
      if (idx+1 == $scope.suggestionIndex){
        return "selected";
      }
      return "";
    };

    $scope.suggestMouseover = function($index) {
      $scope.suggestionIndex = $index+1;
    };

    // 3. refresh PKU Stock table;
    // if invenotryModel is not undefined, assign back it's stock and inventoryNote back into product.
    $scope.refreshCST = function(productId, inventoryModel){
      if(productId == undefined){
        $scope.product = undefined;
        $scope.stocks = undefined;
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

        if(inventoryModel!=undefined && inventoryModel.ProductId == $scope.product.Id){
          $scope.product.InventoryNote = inventoryModel.Note;
          $scope.stocks = angular.copy(inventoryModel.Stocks);
        }

        // TODO should load it's price;

      }).error(function(data, status, headers, config) {
        alert("AJAX failed!");
      });
    };

    // 显示当前选择框的当前库存
    $scope.CurrentLeftStock = function(color, size){
      if($scope.product!=undefined && $scope.product.Stocks!=undefined){
        var stocks = $scope.product.Stocks;
        var sizes = stocks[color];
        if(sizes!=undefined){
          return sizes[size];// stock
        }
      }
      return 0;
    };

    // get stock
    $scope.stock = function(color,size){
      if ($scope.stocks[color]!=undefined ){
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
        Id         : 0,
	    GroupId    : 0,         // TODO
	    ProductId  : p.Id,
	    Stocks     : angular.copy($scope.stocks), // This is stocks matrix;
	    ProviderId : 0,        // factory person id.
	    OperatorId : 0,        // TODO
	    Note       : $scope.Note,
        Product    : p,         // extened
        Note       : p.InventoryNote
      };

      // calculate sum stock if le 0, don't allow to add.
      var sumStock = $scope.calculateSumStocks(inventory.Stocks);
      inventory.sumStock = sumStock;
      if (sumStock<=0){
        $scope.errmsg = "库存必须大于0";
        return;
      }

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

    // add global functions.
    $scope.AddToProducts = function(inventory){
      initInv();
      if($scope.InventoryMap[inventory.ProductId]==undefined){
        $scope.InventoryMap[inventory.ProductId] = inventory;
        $scope.Inventories.push(inventory);
      }else{
        $rootScope.errmsg = "重复添加同一个商品, 更新现有数据";
        $scope.InventoryMap[inventory.ProductId].Stocks = inventory.Stocks; // update stock value;
        $scope.InventoryMap[inventory.ProductId].sumStock = inventory.sumStock;
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

    // click edit on operator column
    $scope.onEdit = function(invId){
      // call with invmodel to retireve Note and stock back;
      $scope.refreshCST(invId, $scope.InventoryMap[invId]);
    };

    // click edit on operator column
    $scope.onDelete = function(invId){
      var inv = $scope.InventoryMap[invId];
      $scope.InventoryMap[invId] = undefined;
      var idx = $scope.Inventories.indexOf(inv);
      $scope.Inventories[idx] = undefined;
    };

    $scope.totalQuantity = function(){
      var total = 0;
      if ($scope.Inventories != undefined){
        for(i=0;i<$scope.Inventories.length;i++){
          var inv = $scope.Inventories[i];
          if (inv!=undefined && inv.sumStock > 0){
            total += inv.sumStock;
          }
        }
      }
      return total;
    };

    $scope.currentSumQuantity = function(){
      if ($scope.stocks != undefined){
        return $scope.calculateSumStocks($scope.stocks);
      }
      return 0;
    };

    // focus on query box
    $scope.focusQuery = function(){
      $('._temp_query_box').focus();
    };

  });

}
