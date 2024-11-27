package models

type Supply struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Checklist struct {
	ID        int  `json:"id"`
	UserID    int  `json:"user_id"`
	SupplyID  int  `json:"supply_id"`
	IsChecked bool `json:"is_checked"`
}

type UserSupplies struct {
	Nickname string            `json:"nickname"`
	Supplies map[string]string `json:"supplies"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
