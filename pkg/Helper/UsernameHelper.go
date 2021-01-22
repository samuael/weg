package Helper

import (
	"regexp"
	"strconv"
	"strings"
)

// ValidateUsername  function to validate whether the string is a valid Username or not
func ValidateUsername(  username string , minLength uint  ) bool {
	trim := func()bool {
		name := strings.Trim(username  , " ") 
		return (len(name) < int(minLength ))
	}

	numbercheck := func() bool{
		_ , err := strconv.Atoi(username)
		return (err != nil)
		}

	if (len( username) < int(minLength))  || 
		trim() ||
		numbercheck(){
		return false 
	}
	return true
}
// ValidatePassword  function to validate whether the string is a valid Username or not
func ValidatePassword(  password string , minLength uint  ) bool {
	if (len( password) < int(minLength))  || 
		(
	func()bool {
		name := strings.Trim(password  , " ") 
		return (len(name) < int(minLength ))
	}()){
		return false 
	}
	return true
}

// EmailRX represents email address maching pattern
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// MatchesPattern checks if a given input form field matchs a given pattern
func MatchesPattern(email string, pattern *regexp.Regexp) bool  {
	if email == "" {
		return false 
	}
	if !pattern.MatchString(email) {
		return false
	}
	return true
}