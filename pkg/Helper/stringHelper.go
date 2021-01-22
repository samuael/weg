package Helper

import "github.com/samuael/Project/Weg/pkg/translation"

// SetPlural function returning the singular or plural form of a string accourding to the Count value
// if count > 1 return pulral form of the word else return as it is translation to the lang specified by the lang variable
func SetPlural(lang  , word   string , count int ) string {
	if count > 1 {
		word += "s"
	}
	return translation.Translate(lang  , word)
}