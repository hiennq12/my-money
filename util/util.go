package util

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func GetDateFromFirstCell(day string) time.Time {
	if len(day) < 2 {
		day = fmt.Sprintf("0%v", day)
	}
	month := strconv.Itoa(int(time.Now().Month()))
	if len(month) < 2 {
		month = fmt.Sprintf("0%v", month)
	}
	dateStr := fmt.Sprintf("%v/%v/%v", day, month, time.Now().Year())
	layout := "02/01/2006"
	dateTime, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Panic("Error when parse date time", err.Error())
	}
	dateStr = dateTime.Format("02/01/2006 15:04:05")

	//fmt.Println("Dasdauyduya: ", dateStr)
	return dateTime
}
