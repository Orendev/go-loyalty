package models

const (
	Rate = 100 // в системе храним копейки 1р = 100 к (1б = 1р)

	//StatusOrderNew        = "NEW"        // заказ загружен в систему, но не попал в обработку;
	//StatusOrderProcessing = "PROCESSING" // вознаграждение за заказ рассчитывается;
	//StatusOrderInvalid    = "INVALID"    // система расчёта вознаграждений отказала в расчёте;
	//StatusOrderProcessed  = "PROCESSED"  // данные по заказу проверены и информация о расчёте успешно
	//
	//StatusAccrualRegistered = "REGISTERED" // заказ зарегистрирован, но вознаграждение не рассчитано;
	//StatusAccrualInvalid    = "INVALID"    // заказ не принят к расчёту, и вознаграждение не будет начислено;
	//StatusAccrualProcessing = "PROCESSING" // расчёт начисления в процессе;
	//StatusAccrualProcessed  = "PROCESSED"  // расчёт начисления окончен;

)

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
