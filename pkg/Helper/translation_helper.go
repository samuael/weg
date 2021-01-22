package Helper

import (
	"strings"

	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

// GetAppropriateLangLabel function representing language label
func GetAppropriateLangLabel(  langString string    ) string {
	
	langString = strings.ToLower( langString)

	for key , group := range entity.LangCodes {
		for _ , el := range group {
			if el == langString {
				return key
			}
		}
	}
	return "en"
}