package models

import (
	"encoding/json"
	"errors"
)

type Order struct {
	ID         string `json:"id" db:"id"`
	Number     int    `json:"number" db:"number"`
	Status     string `json:"status" db:"status"`
	UserID     string `json:"user_id" db:"user_id"`
	Accrual    int    `json:"accrual,omitempty"`
	UploadedAt string `json:"uploaded_at" db:"uploaded_at"`
}

type OrderResponse struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual,omitempty"`
	UploadedAt string `json:"uploaded_at"`
}

func (o *Order) Validate() error {
	var err error

	if (o.Number%10+checksum(o.Number/10))%10 != 0 {
		err = errors.New("the Number field is valid luhn")
	}

	return err
}

func (order OrderResponse) MarshalJSON() ([]byte, error) {
	// чтобы избежать рекурсии при json.Marshal, объявляем новый тип
	type OrderResponseAlias OrderResponse

	aliasValue := struct {
		OrderResponseAlias
		// переопределяем поле внутри анонимной структуры
		Accrual float64 `json:"accrual,omitempty"`
	}{
		// встраиваем значение всех полей изначального объекта (embedding)
		OrderResponseAlias: OrderResponseAlias(order),
		// задаём значение для переопределённого поля
		Accrual: float64(order.Accrual) / Rate,
	}

	return json.Marshal(aliasValue) // вызываем стандартный Marshal
}
