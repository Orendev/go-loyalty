package models

import (
	"encoding/json"
	"time"
)

const (
	StatusAccrualRegistered = "REGISTERED" // заказ зарегистрирован, но вознаграждение не рассчитано;
	StatusAccrualInvalid    = "INVALID"    // заказ не принят к расчёту, и вознаграждение не будет начислено;
	StatusAccrualProcessing = "PROCESSING" // расчёт начисления в процессе;
	StatusAccrualProcessed  = "PROCESSED"  // расчёт начисления окончен;
)

type Accrual struct {
	Order   int
	Status  string
	Accrual int
	Time    time.Duration // Retry-After сек
	UserID  string
}

type AccrualResponse struct {
	Order              string        `json:"order"`
	Status             string        `json:"status"`
	Accrual            int           `json:"accrual"`
	RetryAfterDuration time.Duration `json:"-"`
}

// UnmarshalJSON реализует интерфейс json.Unmarshaler.
func (accrual *AccrualResponse) UnmarshalJSON(data []byte) (err error) {
	// чтобы избежать рекурсии при json.Unmarshal, объявляем новый тип
	type AccrualResponseAlias AccrualResponse

	aliasValue := &struct {
		*AccrualResponseAlias
		// переопределяем поле внутри анонимной структуры
		Accrual float64 `json:"accrual"`
	}{
		AccrualResponseAlias: (*AccrualResponseAlias)(accrual),
	}
	// вызываем стандартный Unmarshal
	if err = json.Unmarshal(data, aliasValue); err != nil {
		return
	}

	accrual.Accrual = int(aliasValue.Accrual * Rate)
	return
}

func (accrual AccrualResponse) MarshalJSON() ([]byte, error) {
	// чтобы избежать рекурсии при json.Marshal, объявляем новый тип
	type AccrualResponseAlias AccrualResponse

	aliasValue := struct {
		AccrualResponseAlias
		// переопределяем поле внутри анонимной структуры
		Accrual float64 `json:"accrual"`
	}{
		// встраиваем значение всех полей изначального объекта (embedding)
		AccrualResponseAlias: AccrualResponseAlias(accrual),
		// задаём значение для переопределённого поля
		Accrual: float64(accrual.Accrual) / Rate,
	}

	return json.Marshal(aliasValue) // вызываем стандартный Marshal
}
