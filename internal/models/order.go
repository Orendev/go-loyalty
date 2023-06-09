package models

import (
	"encoding/json"

	"github.com/Orendev/go-loyalty/internal/luhn"
)

const (
	Rate = 100 // в системе храним копейки 1р = 100 к (1б = 1р)

	//StatusOrderNew        = "NEW"        // заказ загружен в систему, но не попал в обработку;
	StatusOrderProcessing = "PROCESSING" // вознаграждение за заказ рассчитывается;
	StatusOrderInvalid    = "INVALID"    // система расчёта вознаграждений отказала в расчёте;
	StatusOrderProcessed  = "PROCESSED"  // данные по заказу проверены и информация о расчёте успешно
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
	return luhn.Validate(o.Number)
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
