package models

import (
	"encoding/json"
)

type Account struct {
	ID        string `json:"id" db:"id"`
	Current   int    `json:"current" db:"curren"` // Сумма балов лояльности
	Withdrawn int    `json:"withdrawn"`           // Сумма использованных за весь период регистрации баллов.
	UserID    string `json:"user_id" db:"user_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type AccountResponse struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

func (a AccountResponse) MarshalJSON() ([]byte, error) {
	// чтобы избежать рекурсии при json.Marshal, объявляем новый тип
	type AccountResponseAlias AccountResponse

	aliasValue := struct {
		AccountResponseAlias
		// переопределяем поле внутри анонимной структуры
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}{
		// встраиваем значение всех полей изначального объекта (embedding)
		AccountResponseAlias: AccountResponseAlias(a),
		// задаём значение для переопределённого поля
		Current:   float64(a.Current) / Rate,
		Withdrawn: float64(a.Withdrawn) / Rate,
	}

	return json.Marshal(aliasValue) // вызываем стандартный Marshal
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler.
func (a *AccountResponse) UnmarshalJSON(data []byte) (err error) {
	// чтобы избежать рекурсии при json.Unmarshal, объявляем новый тип
	type AccountResponseAlias AccountResponse

	aliasValue := &struct {
		*AccountResponseAlias
		// переопределяем поле внутри анонимной структуры
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}{
		AccountResponseAlias: (*AccountResponseAlias)(a),
	}
	// вызываем стандартный Unmarshal
	if err = json.Unmarshal(data, aliasValue); err != nil {
		return
	}
	a.Current = int(aliasValue.Current * Rate)
	a.Withdrawn = int(aliasValue.Withdrawn * Rate)
	return
}
