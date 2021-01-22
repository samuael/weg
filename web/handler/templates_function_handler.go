package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/samuael/Project/Weg/pkg/Helper"
	"github.com/samuael/Project/Weg/pkg/translation"

	etc "github.com/samuael/Project/Weg/pkg/ethiopianCalendar"
)

// FuncMap func Map For The Templat
var FuncMap = template.FuncMap{
	"Minus":              Minus,
	"GetDate":            GetDayString,
	"GetAge":             GetAge,
	"IsAmharic":          IsAmharic,
	"GetDateString":      GetDateString,
	"DateName":           DateName,
	"GetExtension":       Helper.GetExtension,
	"GetJsonResources":   GetJSONResources,
	"ShortTitle":         ShortTitle,
	"Tran":               translation.Translate,
	"Attach":             Attach,
}

// Attach function
func Attach(val uint, vul int) int {
	vals := fmt.Sprintf("%d%d", val, vul)
	num, ers := strconv.Atoi(vals)
	if ers != nil {
		othar, _ := strconv.Atoi(Helper.GenerateRandomString(3, Helper.NUMBERS))
		return othar
	}
	return num
}

// ShortTitle func
func ShortTitle(title string) string {
	if len(title) <= 47 {
		return title
	}
	return fmt.Sprintf(title[:45], "...")
}

// GetJSONResources functio for Converting to json
func GetJSONResources(inter interface{}) []byte {
	val, _ := json.Marshal(inter)
	return val
}

// Minus function For Subtracting float64 params
func Minus(a, b float64) float64 {
	return (a - b)
}

// GetDayString struct
func GetDayString(day uint) string {
	today := etc.NewDate(int(day))
	return fmt.Sprintf("%d/%d/%d ", today.Day, today.Month, today.Year)
}

// GetAge function
func GetAge(date *etc.Date) string {
	fmt.Println("Calles ", date)
	now := etc.NewDate(0)
	deltayear := now.Year - date.Year
	deltamonth := now.Month - date.Month
	deltaday := now.Day - date.Day
	return fmt.Sprintf("%d", ((deltayear*365 + deltamonth*30 + deltaday) / 365))
}

// IsAmharic function
func IsAmharic(value string) bool {
	variables := strings.Split(value, "")
	bytes := []byte(value)
	if 2*len(variables) == len(bytes) {
		return true
	}
	return false
}



// GetDateString functionh
func GetDateString(date etc.Date) string {
	return fmt.Sprintf("%s %d/%d/%d %d:%d", date.Name, date.Day, date.Month, date.Year, date.Hour, date.Minute)
}

// DateName function
func DateName(date etc.Date) string {
	date = *date.Modify()
	date = *date.FulFill()
	return date.Name
}

// OtherSceneOpenable function
// func OtherSceneOpenable(  resource entity.Resource ) bool {
// 	typ := resource.Type
// 	if typ == entity.PDF || typ == entity.IMAGES ||
// }
