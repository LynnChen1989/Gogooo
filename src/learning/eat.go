package main

import (
	"fmt"
	"time"
)

const STEP = "24h"
const TimeFormat = "2006-01-02"

//这个有点奇葩，必须要这样写，记忆方法为 ： 612345

func main() {
	k, _ := time.ParseDuration(STEP)
	//start_date := time.Date(2016, 07, 20, 0, 0, 0, 0, time.Local)
	//end_date := time.Date(2016, 10, 01, 0, 0, 0, 0, time.Local)

	start_date, _ := time.Parse(TimeFormat, "2016-07-20")
	end_date, _ := time.Parse(TimeFormat, "2016-10-01")
	diff := end_date.Unix() - start_date.Unix()
	count := int(diff / (24 * 60 * 60) / 3)

	for i := 0; i <= count; i++ {
		noon := start_date.Add(k * 1)
		noon_week := noon.Weekday()
		rest := start_date.Add(k * 2)
		rest_week := rest.Weekday()
		work_week := start_date.Weekday()

		fmt.Printf("上班：%10s[%s] ---- 中午下班：%s[%s] ----  整天在家：%s[%s]\n",
			start_date.Format(TimeFormat), work_week,
			noon.Format(TimeFormat), noon_week,
			rest.Format(TimeFormat), rest_week)

		start_date = start_date.Add(k * 3)
	}
}
