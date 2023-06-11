package models

import (
	"encoding/json"
	"errors"
)

type Transact struct {
	ID          string `json:"id" db:"id"`
	Amount      int    `json:"current" db:"curren"` // Сумма балов лояльности
	AccountID   string `json:"account_id" db:"account_id"`
	Debit       bool   `json:"debit" db:"debit"`
	OrderNumber int    `json:"order_number" db:"order_number"`
	ProcessedAt string `json:"processed_at" db:"processed_at"`
}

func (t *Transact) Validate() error {
	var err error
	number := t.OrderNumber
	if (number%10+checksum(number/10))%10 != 0 {
		err = errors.New("the Number field is valid luhn")
	}

	return err
}

func (t Transact) MarshalJSON() ([]byte, error) {
	// чтобы избежать рекурсии при json.Marshal, объявляем новый тип
	type TransactAlias Transact

	aliasValue := struct {
		TransactAlias
		// переопределяем поле внутри анонимной структуры
		Amount float64 `json:"amount" db:"amount"`
	}{
		// встраиваем значение всех полей изначального объекта (embedding)
		TransactAlias: TransactAlias(t),
		// задаём значение для переопределённого поля
		Amount: float64(t.Amount) / Rate,
	}

	return json.Marshal(aliasValue) // вызываем стандартный Marshal
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler.
func (t *Transact) UnmarshalJSON(data []byte) (err error) {
	// чтобы избежать рекурсии при json.Unmarshal, объявляем новый тип
	type TransactAlias Transact

	aliasValue := &struct {
		*TransactAlias
		// переопределяем поле внутри анонимной структуры
		Amount float64 `json:"amount" db:"amount"`
	}{
		TransactAlias: (*TransactAlias)(t),
	}
	// вызываем стандартный Unmarshal
	if err = json.Unmarshal(data, aliasValue); err != nil {
		return
	}
	t.Amount = int(aliasValue.Amount * Rate)
	return
}
