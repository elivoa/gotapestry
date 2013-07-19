 var LODOP; //声明为全局变量 
 /*打印订单*/
/*function printOrder(){
	LODOP=getLodop(document.getElementById('LODOP_OB'),document.getElementById('LODOP_EM'));  
	var strStyleCSS = "<link href='../../css/pages/order-print.css' type='text/css' rel='stylesheet'>";
	var orderTable=strStyleCSS+"<body>"+document.getElementById("orderTable").innerHTML+"</body>";
	var ordeTotal=strStyleCSS+"<body>"+document.getElementById("ordeTotal").innerHTML+"</body>";
	LODOP.PRINT_INIT("销货清单打印");
	LODOP.SET_PRINT_PAGESIZE(0,2410,1400,"销货清单");	
	LODOP.ADD_PRINT_TABLE("3mm","3mm","RightMargin:3mm","75mm",orderTable);
	LODOP.SET_PRINT_STYLEA(0,"Vorient",3);	
		LODOP.ADD_PRINT_HTM(15,0,"100%","40mm",ordeTotal);
		LODOP.SET_PRINT_STYLEA(0,"LinkedItem",-1);
		LODOP.SET_PRINT_STYLEA(0,"Alignment",2);
		LODOP.SET_PRINT_STYLEA(0,"ItemType",1);
	LODOP.ADD_PRINT_TEXT("100mm","220mm","100%","5mm","第#页/共&页");
		LODOP.SET_PRINT_STYLEA(0,"ItemType",2);
		LODOP.SET_PRINT_STYLEA(0,"Horient",2);	
	LODOP.PREVIEW();
};	
*/
function printOrder(){
	LODOP=getLodop(document.getElementById('LODOP_OB'),document.getElementById('LODOP_EM'));  
	LODOP.PRINT_INIT("销货清单打印");
	LODOP.SET_PRINT_PAGESIZE(0,2410,1400,"销货清单");	 
	var strStyleCSS = "<link href='/static/css/pages/order-print.css' type='text/css' rel='stylesheet'>";
	var orderTop = strStyleCSS+"<body>"+document.getElementById("orderTop").innerHTML+"</body>";
	var orderTable = strStyleCSS+"<body>"+document.getElementById("orderTable").innerHTML+"</body>";
	var orderTotal=strStyleCSS+"<body>"+document.getElementById("orderTotal").innerHTML+"</body>";
	LODOP.ADD_PRINT_TEXT("3mm","2mm","RightMargin:2mm","4mm","销货清单");
	LODOP.SET_PRINT_STYLEA(0,"FontSize",12);
	LODOP.SET_PRINT_STYLEA(0,"Bold",1);
	LODOP.SET_PRINT_STYLEA(0,"Alignment",2);
	LODOP.ADD_PRINT_HTM("8mm","2mm","RightMargin:2mm","18mm",orderTop);
	LODOP.ADD_PRINT_TABLE("20mm","3mm","RightMargin:3mm","98mm",orderTable);
		LODOP.SET_PRINT_STYLEA(0,"Top2Offset","-18mm"); //设置次页偏移把区域向上扩
	LODOP.ADD_PRINT_HTM("2mm","2mm","RightMargin:2mm","20mm",orderTotal);
		LODOP.SET_PRINT_STYLEA(0,"LinkedItem",-1);
	LODOP.ADD_PRINT_TEXT("126mm","2mm","190mm","5mm","第#页/共&页");
		LODOP.SET_PRINT_STYLEA(0,"ItemType",2);
		LODOP.SET_PRINT_STYLEA(0,"Horient",2);	
		LODOP.SET_PRINT_STYLEA(0,"Alignment",3);
	LODOP.PREVIEW();
}
