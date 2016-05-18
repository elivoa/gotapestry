// ProductList
// Time-stamp: <[product_list.js] Elivoa @ Wednesday, 2016-05-18 15:28:33>

function p_ProductList($master){

  var sydapp = angular.module('syd', [], function($interpolateProvider){
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
  });

  // if has components, init it first; then init this page;
  ngLoadComponent(sydapp);

  sydapp.controller('ProductListCtrl', function($scope, $http){

    $scope.tabs = ["ALL", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N",
                   "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "#"];

    $scope.firstTabClass = function(index){
      if(index==0){
        return "first";
      }
      return "";
    };

    $scope.tabClass = function(tab, pageTab){
      if(tab.toLowerCase() === pageTab.toLowerCase()){
        return "cur";
      }
      return "";
    };

    // init products.
    // $scope.Products = $master.Products; // TODO change to load something.
    // AjaxLevel 1
    $http.get($master.ProductsLink)
      .success(function (data) {
        $scope.Products = data;

        // AjaxLevel 2: Get Stocks
        $http.get($master.StocksLink)
          .success(function (data) {
            if(data == undefined || data.length != $scope.Products.length){
              console.log("Error!!! with data: ", data);
              return;
            }
            for(var i=0;i <= data.length;i++){
              var p = $scope.Products[i];
              if (p && data[i]){
                p.Stock = data[i].Stock;
                p.Stocks = data[i].Stocks;
              }
            }
          });


        // AjaxLevel 2 in parallel: Get Details.
        $http.get($master.DetailLink)
          .success(function (data) {
            if(data == undefined || data.length != $scope.Products.length){
              console.log("Error!!! with data: ", data);
              return;
            }
            for(var i=0;i <= data.length;i++){
              var p = $scope.Products[i];
              if (p && data[i]){
                p.Colors = data[i].Colors;
                p.Sizes = data[i].Sizes;
                p.Properties = data[i].Properties;
              }
            }
          });

        
      });
    $scope.showall = $master.ShowAll;

    $scope.StockDescription = function(product){
      var str = [];
      if(product.Colors!=undefined){
        for(i=0;i<product.Colors.length;i++){
          if(product.Sizes!=undefined){
            for(j=0;j<product.Sizes.length;j++){
              var color = product.Colors[i];
              var size = product.Sizes[j];
              var has = false;
              // console.log(product);
              if(product.Stocks!=undefined){
                var sizes = product.Stocks[color];
                if(sizes!=undefined){
                  var stock = sizes[size];
                  str.push(stock);
                  has=true;
                }
              }
              if(!has){
                str.push("n/a");
              }
              str.push(" - ");
              str.push(color);
              str.push("/");
              str.push(size);
              str.push("\n");
            }
          }
        }
      }
      return str.join("");
    };

    $scope.SpecDescription = function(product){
      var str = [];
      if(product.Colors!=undefined){
        for(i=0;i<product.Colors.length;i++){
          var color = product.Colors[i];
          if(i>0){
            str.push(" / ");
          }
          str.push(color);
        }
        str.push(' | ');
        if(product.Sizes!=undefined){
          for(j=0;j<product.Sizes.length;j++){
            var size = product.Sizes[j];
            if(j>0){
              str.push(" / ");
            }
            str.push(size);
          }
        }
      }
      return str.join("");
    };

    // // calculate the sum of stocks from Stocks structure.
    // $scope.calculateSumStocks = function(stocks){
    //   var sumStock = 0;
    //   if (stocks!=undefined){
    //     var colors = Object.keys(stocks);
    //     for(j=0;j<colors.length;j++){
    //       var color = colors[j];
    //       var sizemap = stocks[color];
    //       var sizes = Object.keys(sizemap);
    //       for(l=0;l<sizes.length;l++){
    //         var size = sizes[l];
    //         var stock = sizemap[size];
    //         if(stock>0){
    //           sumStock+=stock;
    //         }
    //       }
    //     }
    //   }
    //   return sumStock;
    // };


    // // 将inventories 初始化到系统内部格式
    // $scope.initInventories = function(invs){
    //   var nvs = []; // new inventories
    //   var idmap = {};
    //   if (angular.isArray(invs)){
    //     for (i=0;i<invs.length;i++){
    //       var inv = invs[i];
    //       if(inv!=undefined){
    //         var currentInv = idmap[inv.ProductId];
    //         if(currentInv == undefined){
    //           idmap[inv.ProductId] = currentInv = inv;
    //           nvs.push(currentInv);
    //         }
    //         // add stocks
    //         setStocks(currentInv, inv.Color, inv.Size, inv.Stock);
    //       }
    //     }
    //   }

    //   function setStocks(inv, color, size, stock){
    //     if(inv.Stocks == undefined){
    //       inv.Stocks = {};
    //     }
    //     if(inv.Stocks[color] == undefined){
    //       inv.Stocks[color] = {};
    //     }
    //     inv.Stocks[color][size] = stock;
    //     inv.sumStock = $scope.calculateSumStocks(inv.Stocks); // calculate stocks
    //   }

    //   // final assign
    //   $scope.Inventories = nvs;
    //   $scope.InventoryMap = idmap;
    // };

    // // set master variables into $scope
    // if($master.InventoryGroup!=undefined ){
    //   $scope.initInventories($master.InventoryGroup.Inventories);
    //   $scope.InventoryGroup = $master.InventoryGroup;
    // }
    // $scope.Factories = $master.Factories;

    // $scope.InventoryGroup.SendTime = parseGoDate($scope.InventoryGroup.SendTime);
    // $scope.ReceiveTime =  new Date();

    // // Submit my
    // $scope.submit = function(form) {
    //   fillFormNameWithNGModel(InventoryForm);
    //   InventoryForm.submit();
    // };

  });
}


