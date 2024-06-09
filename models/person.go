package models

type PersonID interface {
	int64 | string
}

type Person[T PersonID] struct {
	ID        T            `json:"id"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Data1     string       `json:"data1"`
	Data2     string       `json:"data2"`
	Data3     string       `json:"data3"`
	Data4     string       `json:"data4"`
	Data5     string       `json:"data5"`
	Addresses []Address[T] `json:"addresses"`
}

type Address[T PersonID] struct {
	ID       int64  `json:"id"`
	PersonID T      `json:"person_id"`
	Address  string `json:"address"`
	City     string `json:"city"`
	State    string `json:"state"`
	Zip      string `json:"zip"`
}
