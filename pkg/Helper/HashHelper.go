package Helper

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword password 
func HashPassword( password string ) (string , error) {
	pwd := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
    if err != nil {
        return "" , err
	}    
	// GenerateFromPassword returns a byte slice so we need to
    // convert the bytes to a string and return it
    return string(hash) ,nil
}

// ComparePassword function returning true is the hash and the password are equivalent 
func ComparePassword( hash  , pwd  string    ) bool {
	bhash := []byte(hash)
	bpwd := []byte(pwd)

	era := bcrypt.CompareHashAndPassword(bhash  , bpwd )
	if era != nil {
		return false
	}
	return true
}