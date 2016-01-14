package base

import ()

var (
	Role_Admin = "admin"
	// Role_Manager          = "manager"
	// Role_Sales            = "sales"
	// Role_InventoryManager = "inventorymanager"
	// Role_Accountant       = "accountant"
	// Role_Engineer         = "engineer"

	Roles = []string{
		Role_Admin, // Role_Manager, Role_Sales, Role_InventoryManager,
		/*Role_Accountant, Role_Engineer*/
	}
	RoleDisplay = []SelectOption{
		SelectOption{Role_Admin, "Admin"},
		// SelectOption{Role_Manager, "Manager"},
		// SelectOption{Role_Sales, "Sale"},
		// SelectOption{Role_InventoryManager, "Inventory Manager"},
		// SelectOption{Role_Accountant, "Accountant"},
		// SelectOption{Role_Engineer, "Engineer"},
	}

	// RoleSet_Customer   = []string{Role_Admin, Role_Manager, Role_Sales}
	// RoleSet_Orders     = []string{Role_Admin, Role_Manager, Role_Sales}
	// RoleSet_Inventory  = []string{Role_Admin, Role_Manager, Role_Sales, Role_InventoryManager}
	// RoleSet_Product    = []string{Role_Admin, Role_Manager}
	// RoleSet_Finance    = []string{Role_Admin, Role_Accountant}
	// RoleSet_Management = []string{Role_Admin, Role_Manager}
	// RoleSet_Mgr_Staff  = []string{Role_Admin, Role_Manager}

	// RoleSet_Announcement  = []string{Role_Admin, Role_Accountant}
	// RoleSet_Preference  = []string{Role_Admin, Role_Accountant}
	// RoleSet_Staff  = []string{Role_Admin, Role_Accountant}

	CONST_DB_DEFAULT_MAX_ITEMS = 50

	// Application specific settings.
	Order_create_tax_enable = false
)

// const key, as a list.
const (
// CONSTKEY_STORE                = "store"
// CONSTKEY_CAR_MADE             = "car_made"
// CONSTKEY_WINDOW_TINT          = "window_tint"
// CONSTKEY_WINDOW_TINT_Optional = "window_tint_optional"
// CONSTKEY_CLEAR_BRA            = "clear_bra"
// CONSTKEY_CLEAR_BRA_Optional   = "clear_bra_optional"
// CONSTKEY_CLEAR_BRAS           = "clear_bras" // no preference/consts
)

type SelectOption struct {
	Value string
	Name  string
}

type BasicConfig struct {
	RoleOptions          [][]string
	RoleOptionsCanCreate [][]string
}

var Basic = &BasicConfig{
	// 其他的隐藏起来，只显示现在的3种。
	RoleOptions: [][]string{
		[]string{Role_Admin, "Admin"},
		// []string{Role_Manager, "Manager"},
		// []string{Role_Sales, "Sales"},
		// []string{Role_InventoryManager, "Inventory Manager"},
		// []string{Role_Accountant, "Accountant"},
		// []string{Role_Engineer, "Engineer"},
	},
	RoleOptionsCanCreate: [][]string{
	// []string{Role_Admin, "Admin"},
	// []string{Role_Manager, "Manager"},
	// []string{Role_Sales, "Sales"},
	},
}

// -- Application Configurations.
var (
	STAT_EXCLUDED_PRODUCT int = 69 // 去掉叫[样衣]的商品
)
