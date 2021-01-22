package Helper

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	etc "github.com/samuael/Project/Weg/pkg/ethiopianCalendar"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))


// NUMBERS const numbers
const NUMBERS = "1234567890"

// CHARACTERS const field
const CHARACTERS = "abcdefghijelmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_1234567890"

// GenerateRandomString  function
func GenerateRandomString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// // OrderRooms function
// func OrderRooms(rooms []*entity.Room) []*entity.Room {
// 	sortedArray := []*entity.Room{}
// 	for i := len(rooms) - 1; i >= 0; i-- {
// 		for k := 0; k <= i-1; k++ {
// 			if rooms[k].Capacity > rooms[k+1].Capacity {
// 				temp := rooms[k]
// 				rooms[k] = rooms[k+1]
// 				rooms[k+1] = temp
// 			}
// 		}
// 	}
// 	for j := len(rooms) - 1; j >= 0; j-- {
// 		sortedArray = append(sortedArray, rooms[j])
// 	}
// 	return sortedArray
// }

// ConfigureLectureHour method
func ConfigureLectureHour(time etc.Date, duration uint) (startdate, enddate etc.Date) {
	time = *time.Modify()
	startdate = time
	enddate = time
	enddate.Hour += int(duration)
	return startdate, enddate
}

// ConfigureTrainingaHour method
func ConfigureTrainingaHour(time etc.Date, NumberOfStudents uint) (starttime, endtime etc.Date) {
	time = *time.Modify()
	starttime = time
	endtime = time
	if time.SubName == etc.TEWAT {
		starttime.Hour = 1
		starttime.Minute = 30
		endtime.Hour = starttime.Hour + int(NumberOfStudents)
		endtime.Minute = starttime.Minute
	} else if time.SubName == etc.KESEAT {
		starttime.Hour = 6
		starttime.Minute = 30
		endtime.Hour = starttime.Hour + int(NumberOfStudents)
		endtime.Minute = starttime.Minute
	}
	return starttime, endtime
}

// SelectedRoomsHalfMode function


// IsReserved   function
func IsReserved(date etc.Date, dates []etc.Date) bool {
	date = *date.Modify()
	date = *date.FulFill()
	for _, dato := range dates {
		if date.Year == dato.Year && dato.Month == date.Month && date.Day == dato.Day && date.SubName == dato.SubName && dato.Shift == date.Shift {
			return true
		}
	}
	return false
}

// IsReservedHalf function to check whether the date is reserved at that Half Shift Or Not
func IsReservedHalf(date etc.Date, dates []etc.Date) bool {
	date = *date.Modify()
	date = *date.FulFill()
	for _, dato := range dates {
		if date.Year == dato.Year && dato.Month == date.Month && date.Day == dato.Day && date.SubName == dato.SubName {
			return true
		}
	}
	return false
}

// IsReservedShift function
func IsReservedShift(date etc.Date, dates []etc.Date) bool {
	for _, dato := range dates {
		if date.Year == dato.Year && dato.Month == date.Month && date.Day == dato.Day && date.SubName == dato.SubName {
			return true
		}
	}
	return false
}

// IsReservedFull function
func IsReservedFull(date etc.Date, dates []etc.Date) bool {
	for _, dato := range dates {
		if date.Year == dato.Year && dato.Month == date.Month && date.Day == dato.Day {
			return true
		}
	}
	return false
}

// NextTime function
func NextTime(date etc.Date) etc.Date {
	if date.SubName == etc.TEWAT {
		date.Hour = 7
		date.Minute = 20
		date.SubName = etc.KESEAT
	} else {
		date.Day++
		date.SubName = etc.TEWAT
	}
	return date
}

// GetFirstDate function
func GetFirstDate(dates []etc.Date) *etc.Date {
	if len(dates) > 0 {
		return &dates[0]
	}
	return nil
}

// RemoveFirstDate method representing
func RemoveFirstDate(dates []etc.Date) []etc.Date {
	newDates := []etc.Date{}
	for i := 1; i < len(dates); i++ {
		newDates = append(newDates, dates[i])
	}
	return newDates
}


// IsDateIn function to check whether the date is in the list of dates
func IsDateIn(date etc.Date, dates []etc.Date) bool {
	for _, dato := range dates {
		if date.Month == dato.Month && dato.Year == date.Year && date.Day == dato.Day {
			return true
		}
	}
	return false
}

// MarshalThis function
func MarshalThis(inter interface{}) []byte {
	val, era := json.Marshal(inter)
	if era != nil {
		return nil
	}
	return val
}


// IsImage function checking whether the file is an image or not
func IsImage(filepath string ) bool {
	extension := GetExtension(filepath)
	extension = strings.ToLower(extension)
	for _ , e := range entity.ImageExtensions  {
		if e== extension {
			return true 
		}
	}
	return  false 
}