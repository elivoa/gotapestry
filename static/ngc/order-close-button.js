// Order/CloseButton
// Time-stamp: <[order-close-button.js] Elivoa @ Thursday, 2015-06-11 23:07:18>

// Development Notes:
// 临时使用侧重调用方式，外部框架有待改进；

//
// Param: todo
// $scope.Inventories    -- inventories equals to InventoryGroup.Inventories
// $scope.InventoryMap   -- id -> Inventory map
//


// console.log("import order-close-button.js");
// console.log( angular.module('app'));

angular.module('app').directive("orderCloseButton", function ($http) {
  return {
    scope:{
      trackNumber     : "@trackNumber",
      customerName    : "@customerName",
      accountBallance : "@accountBallance",
      sumOrderPrice   : "@sumOrderPrice",
      referer         : "@referer"
    },
    link: function (scope, element, attrs) {

      // 结款按钮点击方法：如果上层是ngc，那么这里应该用ngc的参数传递；现在使用原始的传递方法；
      scope.closeclick = function(){
        var $scope = scope;
        return function(e) {
          // replace tempalte contents;

          console.log($scope.sumOrderPrice);
          // e.preventDefault();
          // $scope.accountBallance="100000";// test

          // here use jquery 了，去掉jquery耦合吧；
          var m = angular.element("#order-close-button-modal"); //$("#{{.ClientId}}_modal");
          m.find('.tracking-number').val($scope.trackNumber);
          m.find('.referer').val($scope.referer);
          m.find('.customer-name').html($scope.customerName);
          m.find('.account-ballance').html($scope.accountBallance);
          m.find('strong.sum-order-price').html($scope.sumOrderPrice);
          m.find('input.sum-order-price').val($scope.sumOrderPrice);

          m.on('shown', function(){
            m.find("input.money").focus();
          });
          m.modal("show");

          clearbtn = m.find(".money-clear");
          m.find(".money-clear").click(function(){
            if(clearbtn.prop('checked') == true){
              m.find("input.money").val($scope.sumOrderPrice);
            }else{
              m.find("input.money").val(0);
            }
            m.find("input.money").focus();
          });
        };
      }(); // closure

    },
    restrict: "AEC",
    template: function () {
      return angular.element(
        document.querySelector("#order-close-button-template")).html();
    },
    replace:true
  };

});


function p_I____nventoryEdit($master){

  var sydapp = angular.module('syd', [], function($interpolateProvider){
    // TODO move this into global config;
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
  });

  // if has components, init it first; then init this page;
  ngLoadComponent(sydapp);

  sydapp.controller('InventoryEditCtrl', function($scope){

    // calculate the sum of stocks from Stocks structure.
    $scope.calculateSumStocks = function(stocks){
      var sumStock = 0;
      if (stocks!=undefined){
        var colors = Object.keys(stocks);
        for(j=0;j<colors.length;j++){
          var color = colors[j];
          var sizemap = stocks[color];
          var sizes = Object.keys(sizemap);
          for(l=0;l<sizes.length;l++){
            var size = sizes[l];
            var stock = sizemap[size];
            if(stock>0){
              sumStock+=stock;
            }
          }
        }
      }
      return sumStock;
    };


    // 将inventories 初始化到系统内部格式
    $scope.initInventories = function(invs){
      var nvs = []; // new inventories
      var idmap = {};
      if (angular.isArray(invs)){
        for (i=0;i<invs.length;i++){
          var inv = invs[i];
          if(inv!=undefined){
            var currentInv = idmap[inv.ProductId];
            if(currentInv == undefined){
              idmap[inv.ProductId] = currentInv = inv;
              nvs.push(currentInv);
            }
            // add stocks
            setStocks(currentInv, inv.Color, inv.Size, inv.Stock);
          }
        }
      }

      function setStocks(inv, color, size, stock){
        if(inv.Stocks == undefined){
          inv.Stocks = {};
        }
        if(inv.Stocks[color] == undefined){
          inv.Stocks[color] = {};
        }
        inv.Stocks[color][size] = stock;
        inv.sumStock = $scope.calculateSumStocks(inv.Stocks); // calculate stocks
      }

      // final assign
      $scope.Inventories = nvs;
      $scope.InventoryMap = idmap;
    };

    // set master variables into $scope
    if($master.InventoryGroup!=undefined ){
      $scope.initInventories($master.InventoryGroup.Inventories);
      $scope.InventoryGroup = $master.InventoryGroup;
    }
    $scope.Factories = $master.Factories;
    // console.log("parseGoDate($scope.InventoryGroup.SendTime): ", parseGoDate($scope.InventoryGroup.SendTime))
    // change it.
    $scope.InventoryGroup.SendTime = parseGoDate($scope.InventoryGroup.SendTime);
    $scope.InventoryGroup.ReceiveTime = parseGoDate($scope.InventoryGroup.ReceiveTime);
    // TODO how to add date by one day using javascript;

    // Submit my
    $scope.submit = function(form) {
      fillFormNameWithNGModel(InventoryForm);
      InventoryForm.submit();
    };

  });
}



