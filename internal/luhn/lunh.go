package luhn

import (
	"errors"
)

func Validate(number int) error {
	var err error

	if (number%10+checksum(number/10))%10 != 0 {
		err = errors.New("the Number field is valid luhn")
	}

	return err
}

func checksum(number int) int {
	var s int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		s += cur
		number = number / 10
	}
	return s % 10
}
