package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	var sdate, name, surname, lastname string
	var val1, val2, val3 float64
	fmt.Scan(&sdate, &name, &surname, &lastname, &val1, &val2, &val3)

	date, _ := time.Parse("02.01.2006", sdate)
	fireDate := date.AddDate(0, 0, 15)

	c1, c2, c3 := val1-math.Floor(val1), val2-math.Floor(val2), val3-math.Floor(val3)
	cents := c1 + c2 + c3 - math.Floor(c1+c2+c3)
	rubles := val1 + val2 + val3 - cents

	message := fmt.Sprintf("Уважаемый, %s %s %s, доводим до вашего сведения, что бухгалтерия сформировала документы по факту выполненной вами работы.\nДата подписания договора: %s. Просим вас подойти в офис в любое удобное для вас время в этот день.\nОбщая сумма выплат составит %d руб. %d коп. \n\nС уважением,\nГл. бух. Иванов А.Е.",
		surname, name, lastname,
		fireDate.Format("02.01.2006"),
		int(rubles), int(100*cents),
	)

	fmt.Println(message)
}
