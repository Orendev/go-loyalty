package models

import (
	"encoding/json"
	"errors"
	"strconv"
)

// WithdrawRequest описывает запрос клиента.
type WithdrawRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

func (withdraw *WithdrawRequest) Validate() error {
	var err error
	number, err := strconv.Atoi(withdraw.Order)
	if err != nil {
		return err
	}

	if (number%10+checksum(number/10))%10 != 0 {
		err = errors.New("the Number field is valid luhn")
	}

	return err
}

func (withdraw WithdrawRequest) MarshalJSON() ([]byte, error) {
	// чтобы избежать рекурсии при json.Marshal, объявляем новый тип
	type AccountResponseAlias WithdrawRequest

	aliasValue := struct {
		AccountResponseAlias
		// переопределяем поле внутри анонимной структуры
		Sum float64 `json:"sum"`
	}{
		// встраиваем значение всех полей изначального объекта (embedding)
		AccountResponseAlias: AccountResponseAlias(withdraw),
		// задаём значение для переопределённого поля
		Sum: float64(withdraw.Sum) / Rate,
	}

	return json.Marshal(aliasValue) // вызываем стандартный Marshal
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler.
func (withdraw *WithdrawRequest) UnmarshalJSON(data []byte) (err error) {
	// чтобы избежать рекурсии при json.Unmarshal, объявляем новый тип
	type WithdrawRequestAlias WithdrawRequest

	aliasValue := &struct {
		*WithdrawRequestAlias
		// переопределяем поле внутри анонимной структуры
		Sum float64 `json:"sum"`
	}{
		WithdrawRequestAlias: (*WithdrawRequestAlias)(withdraw),
	}
	// вызываем стандартный Unmarshal
	if err = json.Unmarshal(data, aliasValue); err != nil {
		return
	}

	withdraw.Sum = int(aliasValue.Sum * Rate)
	return
}
