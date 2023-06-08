package models

import (
	"errors"
)

type Order struct {
	ID         string `json:"id" db:"id"`
	Number     int    `json:"number" db:"number"`
	Status     string `json:"status" db:"status"`
	UserID     string `json:"user_id" db:"user_id"`
	UploadedAt string `json:"uploaded_at" db:"uploaded_at"`
}
type NullString struct {
	String string
	Valid  bool
}
type OrderResponse struct {
	Number     int    `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}

func (o *Order) Validate() error {
	var err error

	if (o.Number%10+checksum(o.Number/10))%10 != 0 {
		err = errors.New("the Number field is valid luhn")
	}

	return err
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
