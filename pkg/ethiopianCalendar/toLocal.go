/* Package ethiopianCalendar for Ethiopians
 */
package ethiopianCalendar

import (
	"fmt"
	"time"

	// "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/jinzhu/gorm"
)

// Date representing the Ethiopian Date time
type Date struct {
	gorm.Model
	Name        string `json:"dayname,omitempty"`
	Day         int    `json:"day,omitempty"`
	Year        int    `json:"year,omitempty"`
	Month       int    `json:"month,omitempty"`
	Hour        int    `json:"hour,omitempty"`
	Minute      int    `json:"minute,omitempty"`
	SubName     string `json:"subname,omitempty"`
	Unix        uint64 `json:"unix,omitempty"`
	IsBreakDate bool   `gorm:"default:false;"`
	BranchID    uint
	OwnerID     uint `json:"ownerid,omitempty"`
	RoundID     uint
	Shift       uint // this repressents the Shift Whether it is 0: 2-4  , 1:4-6  , 2 : 7-9 or 3 : 9-11
}

// DAYS representing the Days name
var DAYS = []string{
	"Hamus",
	"Arb",
	"Kidame",
	"Ehud",
	"Segno",
	"Magsegno",
	"Erob",
}

var (
	// Hamus sonstant
	Hamus = DAYS[0]
	// Arb sonstant
	Arb = DAYS[1]
	// Kidame sonstant
	Kidame = DAYS[2]
	// Ehud sonstant
	Ehud = DAYS[3]
	// Segno sonstant
	Segno = DAYS[4]
	// Magsegno sonstant
	Magsegno = DAYS[5]
	// Erob sonstant
	Erob = DAYS[6]
)

var timo time.Time

// YEAR y
var YEAR = 2012
var MONTH = 8
var DAY = 22
var MINUTE = 40
var HOUR = 3
var SECOND = 55
var UNIX_SECOND = 1588228855
var UNIX_SECOND_PAST = 55
var UNIX_MINUTE_PAST = 40
var HOUR_IN_SECOND = 3600
var MINUTE_IN_SECOND = 60
var DAY_IN_SECOND = 86400
var YEAR_IN_SECOND = 31557842
var MONTH_IN_SECOND = 2592000
var TOtalDay = 365
var NUMBER_OF_DAYS = 232*DAY_IN_SECOND + (HOUR * HOUR_IN_SECOND) + (MINUTE * MINUTE_IN_SECOND) + SECOND
var NUMBER_OF_DAYS_REMAINING = YEAR_IN_SECOND - NUMBER_OF_DAYS
var NUMBER_OF_MONTH_FRACTION_REMAINING_IN_SECOND = MONTH_IN_SECOND - ((DAY - 1) * DAY_IN_SECOND) - (HOUR * HOUR_IN_SECOND) - (MINUTE_IN_SECOND * MINUTE) - (SECOND)
var NUMBER_OF_DAY_REMAINING_FRACTION_IN_SECOND = DAY_IN_SECOND - (HOUR * HOUR_IN_SECOND) - (MINUTE_IN_SECOND * MINUTE) - (SECOND)
var YEAR_REMAINING_FRACTION_FROM_NOW_IN_SEC = YEAR_IN_SECOND - NUMBER_OF_DAYS
var totalSecond int64
var REMAINING_SECONDS_TO_THE_YEAR_EANDING = YEAR_IN_SECOND -
	(MONTH_IN_SECOND*(MONTH-1) +
		DAY_IN_SECOND*(DAY-1) +
		HOUR_IN_SECOND*HOUR +
		MINUTE_IN_SECOND*MINUTE +
		SECOND)

// ETYear return the year In Local Ethiopian
func ETYear() (count int64) {
	count = 0
	val := time.Now().Unix() - int64(UNIX_SECOND)
	if val >= int64(NUMBER_OF_DAYS_REMAINING) {
		count = count + 1
		val -= int64(NUMBER_OF_DAYS_REMAINING)
		count = count + val/int64(YEAR_IN_SECOND)
	}
	count += int64(YEAR)
	return count
}

// ETYearIndexed representing the Year after plus days
func ETYearIndexed(plus int64) (count int64) {
	count = 0
	val := time.Now().Unix() + plus*int64(DAY_IN_SECOND) - int64(UNIX_SECOND)
	if val >= int64(NUMBER_OF_DAYS_REMAINING) {
		count = count + 1
		val -= int64(NUMBER_OF_DAYS_REMAINING)
		count = count + val/int64(YEAR_IN_SECOND)
	}
	count += int64(YEAR)
	return count
}

// ETMonth returns the Month In Ethiopian Calender
func ETMonth() (count int64) {
	count = 0
	letmeadd := true
	val := time.Now().Unix() - int64(UNIX_SECOND)
	if val >= int64(REMAINING_SECONDS_TO_THE_YEAR_EANDING) {
		val -= int64(REMAINING_SECONDS_TO_THE_YEAR_EANDING)
		letmeadd = false
		val = val % int64(YEAR_IN_SECOND)
	} else if val >= int64(NUMBER_OF_MONTH_FRACTION_REMAINING_IN_SECOND) && letmeadd {
		// this gets called whenever the difference is less than a REMAINING_SECONDS_TO_THE_YEAR_EANDING
		// but greater than  NUMBER_OF_MONTH_FRACTION_REMAINING_IN_SECON
		count = count + 1
		val = val - int64(NUMBER_OF_MONTH_FRACTION_REMAINING_IN_SECOND)
		passed := val / int64(MONTH_IN_SECOND)
		count += passed
	}

	if !letmeadd {
		count = (val / int64(MONTH_IN_SECOND)) + 1
	} else {
		count += int64(MONTH)
	}
	return count
}

// ETMonthIndexed representing the Monthe asfter plus Days
func ETMonthIndexed(plus int64) (count int64) {
	count = 0
	letmeadd := true
	val := (time.Now().Unix() + plus*int64(DAY_IN_SECOND)) - int64(UNIX_SECOND)
	if val >= int64(REMAINING_SECONDS_TO_THE_YEAR_EANDING) {
		val -= int64(REMAINING_SECONDS_TO_THE_YEAR_EANDING)
		letmeadd = false
		val = val % int64(YEAR_IN_SECOND)
	} else if val >= int64(NUMBER_OF_MONTH_FRACTION_REMAINING_IN_SECOND) && letmeadd {
		// this gets called whenever the difference is less than a REMAINING_SECONDS_TO_THE_YEAR_EANDING
		// but greater than  NUMBER_OF_MONTH_FRACTION_REMAINING_IN_SECOND
		count = count + 1
		val = val - int64(NUMBER_OF_MONTH_FRACTION_REMAINING_IN_SECOND)
		passed := val / int64(MONTH_IN_SECOND)
		count += passed
	}

	if !letmeadd {
		count = (val / int64(MONTH_IN_SECOND)) + 1
		// if count == 13 {
		// 	if val/int64(DAY_IN_SECOND) >= 7 {
		// 		count = 0
		// 	}
		// }
	} else {
		count += int64(MONTH)
	}
	return count
}

// ETDayName Returns the String representing the Name of The Day
func ETDayName() string {
	val := time.Now().Unix() - int64(UNIX_SECOND)
	days := val / int64(DAY_IN_SECOND)
	index := days % 7
	return DAYS[index]
}

// ETDayNameIndexed returning the Name of The day after the Number of Days taken as a paramenter
func ETDayNameIndexed(plus int64) string {
	val := time.Now().Unix() + plus*int64(DAY_IN_SECOND) - int64(UNIX_SECOND)
	days := val / int64(DAY_IN_SECOND)
	index := days % 7
	return DAYS[index]
}

// ETDay  Ethiopian Day
func ETDay() (count int64) {
	count = 0
	var Passed = false
	val := time.Now().Unix() - int64(UNIX_SECOND)
	// if the Count is greater than a day
	if val >= int64(NUMBER_OF_DAY_REMAINING_FRACTION_IN_SECOND) {
		// if the count is greater than a year
		if val >= int64(REMAINING_SECONDS_TO_THE_YEAR_EANDING) {
			val -= int64(REMAINING_SECONDS_TO_THE_YEAR_EANDING)
			Passed = true
			count += (val / int64(DAY_IN_SECOND))
		} else {
			count++
			// if the Value is greater than te Number of Fraction of seconds remaining for the Next Day
			// in my Count .
			val -= int64(NUMBER_OF_DAY_REMAINING_FRACTION_IN_SECOND)
			count += (val / int64(DAY_IN_SECOND))
			count += int64(DAY)
		}
	}
	// if the day count is greater than 30 it has to be moduled
	// and if the Number of days passed are greater than Number of Days
	if Passed {
		if ETMonth() == 13 {
			val = val % int64(YEAR_IN_SECOND)
			if val >= int64(YEAR_IN_SECOND-DAY_IN_SECOND) && val <= int64(YEAR_IN_SECOND-DAY_IN_SECOND) {
				return count
			} else if val >= int64(YEAR_IN_SECOND) {
				return 1
			}
		} else {
			if count%30 == 0 {
				count = 30
			} else {
				count = (count % 30)
			}
		}
	} else {
		if count%30 == 0 {
			count = 30
		} else {
			count = (count % 30)
		}
	}
	return count
}

// ETDayIndexed Ethiopian Date or Time
// the input index representing Number of Days
func ETDayIndexed(plus int64) (count int64) {
	count = 0
	// Telling if the Year is Passed Or Not
	var Passed = false
	val := time.Now().Unix() + int64(DAY_IN_SECOND)*plus - int64(UNIX_SECOND)
	if val >= int64(NUMBER_OF_DAY_REMAINING_FRACTION_IN_SECOND) {
		if val >= int64(REMAINING_SECONDS_TO_THE_YEAR_EANDING) {
			val -= int64(REMAINING_SECONDS_TO_THE_YEAR_EANDING)
			Passed = true
			count += (val / int64(DAY_IN_SECOND))
		} else {
			count++
			// if the Value is greater than te Number of Fraction of seconds remaining for the Next Day
			// in my Count .
			val -= int64(NUMBER_OF_DAY_REMAINING_FRACTION_IN_SECOND)
			count += (val / int64(DAY_IN_SECOND))
			count += int64(DAY)
		}
	}
	// if the day count is greater than 30 it has to be moduled
	// and if the Number of days passed are greater than Number of Days
	if Passed {
		if count%30 == 0 {
			count = 30
		} else {
			count = (count % 30)
		}

	} else {
		if ETMonthIndexed(plus) == 13 {
			val = val % int64(YEAR_IN_SECOND)
			if val >= int64(YEAR_IN_SECOND-DAY_IN_SECOND) && val <= int64(YEAR_IN_SECOND-DAY_IN_SECOND) {
				return count
			}
		} else if count%30 == 0 {
			count = 30
		} else {
			count = (count % 30)
		}
	}

	return count
}

// FulFill method
func (date *Date) FulFill() *Date {
	// benm yahl ken Ke Zare Beltual . esun kagegnehu  Unixun  , daynamun
	today := NewDate(0)
	passedDays := date.Day - today.Day
	passedMonths := date.Month - today.Month
	passedYears := date.Year - today.Year

	sum := passedDays
	sum += (passedMonths) * 30
	sum += (passedYears * 365)
	date.Unix = today.Unix + uint64(sum*DAY_IN_SECOND)
	index := 0
	for i, element := range DAYS {
		if element == today.Name {
			index = i
		}
	}
	passedFloat := sum % 7
	index += passedFloat
	index = ((index) % 7)
	date.Name = DAYS[index]
	date = date.Modify()
	return date
}

// NextDay unix
func (date *Date) NextDay() *Date {
	date = (date.FulFill())
	date.Unix += uint64(DAY_IN_SECOND)
	date.Day = int(ETDayUnix(int64(date.Unix)))
	date.Year = int(ETYearUnix(int64(date.Unix)))
	date.Month = int(ETMonthUnix(int64(date.Unix)))
	date = date.FulFill()
	date.Modify()
	// Filter if the date is in the reserved dates list
	// if so jump it
	return date
}

// BackShiftTime date function for shifting
func (date Date) BackShiftTime() Date {
	if date.Shift == ShiftB {
		date.Shift = ShiftA
		date.Hour = 2
		date.Minute = 15
	} else if date.Shift == ShiftD {
		date.Shift = ShiftC
		date.Hour = 7
		date.Minute = 15
	} else if date.Shift == ShiftC {
		date.Shift = ShiftA
		date.Hour = 2
		date.Minute = 15
	}
	date = *date.Modify()
	date = *date.FulFill()
	return date
}

// NextMorning function to get The Day Morning
func (date *Date) NextMorning() *Date {
	counto := 26 - date.Hour
	date.Unix += uint64(counto * HOUR_IN_SECOND)
	date.Day = int(ETDayUnix(int64(date.Unix)))
	date.Year = int(ETYearUnix(int64(date.Unix)))
	date.Month = int(ETMonthUnix(int64(date.Unix)))
	date.Hour = 2
	date.Minute = 15
	date = date.FulFill()
	date.Modify()
	return date
}

// NextTime function
func (date *Date) NextTime() *Date {
	if date.Shift == ShiftA {
		date.Hour = 4
		date.Minute = 30
		date.SubName = TEWAT
		date.Shift = ShiftB
	} else if date.Shift == ShiftB {
		date.Shift = ShiftC
		date.Hour = 7
		date.Minute = 15
		date.SubName = KESEAT
	} else if date.Shift == ShiftC {
		date.Shift = ShiftD
		date.Hour = 9
		date.Minute = 30
		date.SubName = KESEAT
	} else if date.Hour >= 9 {
		date = date.NextDay()
		date.SubName = TEWAT
		date.Hour = 2
		date.Minute = 15
		date.Shift = ShiftA
	}
	return date
}

// NextTimeHalf function
func (date *Date) NextTimeHalf() *Date {

	if date.SubName == TEWAT || date.Hour < 7 {
		date.SubName = KESEAT
		date.Hour = 7
		date.Minute = 30
		date.Shift = ShiftC
	} else if date.SubName == KESEAT || (date.Hour >= 7 && date.Hour < 12) {
		date = date.NextDay()
		date.SubName = TEWAT
		date.Shift = ShiftA
		date.Hour = 2
		date.Minute = 15
	}
	return date
}

// NextShift function to find the Next Shift Meaning the Tewat Or the Keseat of the Day
func (date *Date) NextShift() *Date {
	if (date.Hour <= 6) && (date.SubName == TEWAT) {
		date.SubName = KESEAT
		date.Hour = 7
		date.Minute = 0
	} else if date.Hour > 7 {
		date = date.NextDay()
		date.Hour = 2
		date.Minute = 0
		date.SubName = TEWAT
	}
	return date
}

// IsFuture function
func (date *Date) IsFuture(newDate *Date) bool {
	passedDays := date.Day - newDate.Day
	passedMonths := date.Month - newDate.Month
	passedYear := date.Year - newDate.Year
	days := passedDays + (passedMonths * 30) + (passedYear * 365)
	if days < 0 {
		return false
	}
	return true
}

// NewDate returning Ethiopian Date Time Representation
func NewDate(index int) *Date {
	unix := time.Now().Unix()
	day := ETDayIndexed(int64(index))
	mont := ETMonthIndexed(int64(index))
	year := ETYearIndexed(int64(index))
	name := ETDayNameIndexed(int64(index))
	date := &Date{
		Unix:   uint64(unix),
		Name:   name,
		Month:  int(mont),
		Year:   int(year),
		Day:    int(day),
		Hour:   time.Now().Hour() - 6,
		Minute: time.Now().Minute(),
	}
	switch {
	case date.Hour <= 7 && date.Minute <= 30:
		date.SubName = TEWAT
		break
	case date.Hour < 12:
		date.SubName = KESEAT
		break
	case date.Hour < 18:
		date.SubName = MISHIT
		break
	case date.Hour < 24:
		date.SubName = LELIT
		break
	default:
		date.SubName = AYITAWEKM
	}
	return date
}

const (
	// TEWAT const
	TEWAT = "Tewat"
	// KESEAT const
	KESEAT = "Keseat"
	// MISHIT const
	MISHIT = "Mishit"
	// LELIT const
	LELIT = "Lelit"
	// AYITAWEKM const
	AYITAWEKM = "Ayitawekm"
	// ShiftA constant
	ShiftA = 0
	// ShiftB constant
	ShiftB = 1
	// ShiftC constant
	ShiftC = 2
	// ShiftD constant
	ShiftD = 3
)

// ToString method Returns a string representation of the Date Struct
func (date *Date) ToString() string {
	return fmt.Sprintf("%d/%d/%d %s", date.Day, date.Month, date.Year, date.Name)
}

// GetAge method
func (date *Date) GetAge() int {
	now := NewDate(0)
	deltayear := now.Year - date.Year
	deltamonth := now.Month - date.Month
	deltaday := now.Day - date.Day
	return (deltayear*365 + deltamonth*30 + deltaday) / 365
}

// Modify function
func (date *Date) Modify() *Date {
	switch {
	case date.Hour < 6:
		date.SubName = TEWAT
		if date.Hour < 3 {
			date.Shift = ShiftA
		} else {
			date.Shift = ShiftB
		}
		break
	case date.Hour < 12:
		date.SubName = KESEAT
		if date.Hour < 9 || (date.Hour == 9 && date.Minute < 15) {
			date.Shift = ShiftC
		} else {
			date.Shift = ShiftD
		}
		break
	case date.Hour < 18:
		date.SubName = MISHIT
		break
	case date.Hour < 24:
		date.SubName = LELIT
		break
	default:
		date.SubName = AYITAWEKM
	}

	return date
}


// ETDayUnix  Ethiopian Day
func ETDayUnix(unix int64) (count int64) {
	count = 0
	var Passed = false
	val := unix - int64(UNIX_SECOND)
	// if the Count is greater than a day
	if val >= int64(NUMBER_OF_DAY_REMAINING_FRACTION_IN_SECOND) {
		// if the count is greater than a year
		if val >= int64(REMAINING_SECONDS_TO_THE_YEAR_EANDING) {
			val -= int64(REMAINING_SECONDS_TO_THE_YEAR_EANDING)
			Passed = true
			count += (val / int64(DAY_IN_SECOND))
		} else {
			count++
			// if the Value is greater than te Number of Fraction of seconds remaining for the Next Day
			// in my Count .
			val -= int64(NUMBER_OF_DAY_REMAINING_FRACTION_IN_SECOND)
			count += (val / int64(DAY_IN_SECOND))
			count += int64(DAY)
		}
	}
	// if the day count is greater than 30 it has to be moduled
	// and if the Number of days passed are greater than Number of Days
	if Passed {
		if ETMonth() == 13 {
			val = val % int64(YEAR_IN_SECOND)
			if val >= int64(YEAR_IN_SECOND-DAY_IN_SECOND) && val <= int64(YEAR_IN_SECOND-DAY_IN_SECOND) {
				return count
			} else if val >= int64(YEAR_IN_SECOND) {
				return 1
			}
		} else {
			if count%30 == 0 {
				count = 30
			} else {
				count = (count % 30)
			}
		}
	} else {
		if count%30 == 0 {
			count = 30
		} else {
			count = (count % 30)
		}
	}
	return count
}

// ETDayNameUnix method
func ETDayNameUnix(unix int64) string {
	val := unix - int64(UNIX_SECOND)
	days := val / int64(DAY_IN_SECOND)
	index := days % 7
	return DAYS[index]
}

// ETMonthUnix returns the Month In Ethiopian Calender
func ETMonthUnix(unix int64) (count int64) {
	count = 0
	letmeadd := true
	val := unix - int64(UNIX_SECOND)
	if val >= int64(REMAINING_SECONDS_TO_THE_YEAR_EANDING) {
		val -= int64(REMAINING_SECONDS_TO_THE_YEAR_EANDING)
		letmeadd = false
		val = val % int64(YEAR_IN_SECOND)
	} else if val >= int64(NUMBER_OF_MONTH_FRACTION_REMAINING_IN_SECOND) && letmeadd {
		// this gets called whenever the difference is less than a REMAINING_SECONDS_TO_THE_YEAR_EANDING
		// but greater than  NUMBER_OF_MONTH_FRACTION_REMAINING_IN_SECOND

		count = count + 1
		val = val - int64(NUMBER_OF_MONTH_FRACTION_REMAINING_IN_SECOND)
		passed := val / int64(MONTH_IN_SECOND)
		count += passed
	}

	if !letmeadd {
		count = (val / int64(MONTH_IN_SECOND)) + 1
		// if count == 13 {
		// 	if val/int64(DAY_IN_SECOND) >= 7 {
		// 		count = 0
		// 	}
		// }
	} else {
		count += int64(MONTH)
	}
	return count
}

// ETYearUnix method
func ETYearUnix(unix int64) (count int64) {
	count = 0
	val := unix - int64(UNIX_SECOND)
	if val >= int64(NUMBER_OF_DAYS_REMAINING) {
		count = count + 1
		val -= int64(NUMBER_OF_DAYS_REMAINING)
		count = count + val/int64(YEAR_IN_SECOND)
	}
	count += int64(YEAR)
	return count
}

// IsPassed returning true if the Date in the Parameter has passed the Date at consideration
func (date *Date) IsPassed(today *Date) bool {
	passedDates := date.Day - today.Day
	passedMonth := date.Month - today.Month
	passedYear := date.Year - today.Year
	days := passedDates + (passedMonth * 30) + (passedYear * 365)
	if days >= 0 {
		return true
	}
	return false
}
