package model

import (
	"time"
)

type ActionType int

const (
	ACTION_LOGIN  ActionType = 100
	ACTION_LOGOUT            = 110

	// inventory
	ACTION_INVENTORY_CREATE      = 210
	ACTION_INVENTORY_EDIT        = 220
	ACTION_MARK_INVENTORY_INUSE  = 230
	ACTION_MARK_INVENTORY_RUNOUT = 240

	// order
	ACTION_ORDER_CREATE = 300
	ACTION_ORDER_EDIT   = 310
	ACTION_ORDER_PAY    = 320
	ACTION_ORDER_REDO   = 325 // inserted.
	ACTION_ORDER_PRINT  = 330
	ACTION_ORDER_RETURN = 340
	ACTION_ORDER_DELETE = 350
	ACTION_ORDER_DROPIN = 360
	ACTION_ORDER_PICKUP = 361

	// customer
	ACTION_CUSTOMER_EDIT   = 400
	ACTION_CUSTOMER_DELETE = 410

	// user
	ACTION_USER_CREATE = 500
	ACTION_USER_EDIT   = 510
	ACTION_USER_DELETE = 520

	// Add Here Only...
)

// for display only.
var DisplayActions = map[ActionType]string{
	ACTION_LOGIN:                 "Login",
	ACTION_LOGOUT:                "Logout",
	ACTION_MARK_INVENTORY_INUSE:  "Mark Inventory as InUse",
	ACTION_MARK_INVENTORY_RUNOUT: "Mark Inventory as Runout",
	ACTION_ORDER_CREATE:          "Create Order",
	ACTION_ORDER_EDIT:            "Edit Order",
	ACTION_ORDER_PAY:             "Pay Order",
	ACTION_ORDER_PRINT:           "Print Order",
	ACTION_ORDER_RETURN:          "Order",
	ACTION_ORDER_DELETE:          "** Delete Order **",
	ACTION_ORDER_DROPIN:          "Drop in Order",
	ACTION_ORDER_PICKUP:          "Pick up Order",
	ACTION_CUSTOMER_EDIT:         "Edit Customer",
	ACTION_CUSTOMER_DELETE:       "Delete Customer",
	ACTION_USER_CREATE:           "Create User",
	ACTION_USER_EDIT:             "Edit User",
	ACTION_USER_DELETE:           "Delete User",
}

func (t ActionType) String() string {
	if r, ok := DisplayActions[t]; ok {
		return r
	} else {
		return "Action Not Defined!"
	}
}

type UserAction struct {
	Id         int64
	UserId     int64
	Action     ActionType
	Context    string
	CreateTime time.Time

	// external
	User *User
}
