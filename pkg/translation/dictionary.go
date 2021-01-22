// CopyRight all Right are Reserved

// Package translation this package is created for the
// sake of translating the language embedded in the Templates to needed language choise of the User
package translation

import "strings"

// DICTIONARY FOR SAVING DIFFERENT LANGUAGES
var DICTIONARY = map[string]map[string]string{
	"amh": map[string]string{
		"login":                "ግባ",
		"name" : "ስም", 
		"logout":               "ውጣ",
		"full name":            "mulu sim",
		"sex":                  "tsota",
		"age":                  "edme",
		"acadamic status":      "yetmhrt dereja",
		"language":             "quankua",
		"city":                 "ketema",
		"kebele":               "kebele",
		"phone":                "silk kutr",
		"address":              "Adrasha",
		"region":               "kilil",
		"category":             "zerf",
		"id":                   "metawekia",
		"trainers id card":     "ye Seltagn metawekia",
		"disclaimer":           "mastawesha",
		"this id is valid for": "Yih Metawekia yemiageleglew le",
		"months only":          "bicha new ",
		"username" : "መለያ ስም",
		"your medicine infos" : "የመዘገቧቸው ዝርዝር" , 
		"register medicine info" : "የመድሀኒት መረጃ መዝግብ" , 
		"register admin" : "አስተዳዳሪ መዝግብ" , 
		"given date":           "yetesetebet ken",
		"confirm password" : "ይለፍ ቃል ዐረጋግጥ",
		"password":"ይለፍ ቃል",
		"update your profile" : "መለያህን አዘምን" ,
		"home":"መግቢያ" ,
		"register" :"ተመዝገብ" , 
		"profile" :"መለያ" , 
		"update" :"አዘምን",
		"lang" :"ቋንቋ" , 
		"purpose":"አገልግሎት" , 
		"precaution" :"ማስጠንቀቂያዎች" , 
		"utilization":"አጠቃቀም" , 
		"side effect" :"የጎንዮሽ ጉዳት" , 
		"age from" :"መነሻ እድሜ" , 
		"age to":"መጨረሻ እድሜ",
		"time interval" :"የሰአት ልዩነት" , 
		"intake recursion" : "የእወሳሰድ ድግግሞሽ",
		"amount" : "መጠን" , 
		"pieces":"ፍሬዎች",
		"piece" : "ፍሬ" ,
		"initial weight" : "መነሻ ክብደት" , 
		"final weight" : "መጨረሻ ክብደት" ,
		"max weight" : "መጨረሻ ክብደት" , 
		"delete" : "ሰርዝ" ,
		"no medicine information found" : "ምንም የመድሀኒት መረጃ አልተገኘም" , 	
		"tab" :"ንክ" , 
		"overdose"  : "ከተገቢው በላይ ሲወሰድ" ,
		"title" : "ርእስ"  ,	
		"medicine information registration form" : "የመድሀኒት መረጃ መመዝገቢያ ቅፅ",
		"registor admin form" : "ተቆጣጣሪ መመዝገቢያ ቅፅ" ,
		"search" : "ፈልግ" , 
		"search title" : "በርዕስ ፈልግ" , 
		"no record found please try again" : "ምንም የትገኘ መረጃ የለም እባክዎ ደግመው ይሞከሩ" , 
		"search for medicine information results will be listed here" : "የመድሀኒት መረጃ ፍለጋ ውጤቶች እዚህ ጋር ይዘረዘራሉ",
		"update profile picture" : "መለያ ምስል አዘምን" , 
		"profile image" : "መለያ ምስል" ,
		"change profile picture" : "መለያ ምስል ቀይር" ,
		"the" : "የ",
	},
	"oromifa": map[string]string{
		"login":            "giba",
		"logout":           "chufa",
		"register":         "temezgeb",
		"Full name":        "mulu sim",
		"Sex":              "tsota",
		"Age":              "edme",
		"Acadamic Status":  "yetmhrt dereja",
		"Language":         "quankua",
		"City":             "ketema",
		"Kebele":           "kebele",
		"Phone":            "silk kutr",
		"address":          "Adrasha",
		"region":           "kilil",
		"category":         "zerf",
		"id":               "metawekia",
		"trainers id card": "ye Seltagn metawekia",
		"home":"መግቢያ" , 
		
	},
}

// Translate  function to change the word to the needed Language Representation
func Translate(lang string, sentence string) string {
	hold := sentence

	sentence = strings.Trim(sentence , " ")

	switch strings.ToLower(lang) {
	case "en", "eng":
		return sentence
	case "amh", "am", "amharic", "amhara":
		sentence = strings.ToTitle((DICTIONARY["amh"])[strings.ToLower(sentence)])
	case "oro", "or", "oromifa", "oromo":
		sentence = strings.ToTitle((DICTIONARY["oromifa"])[strings.ToLower(sentence)])
	}
	if sentence=="" {
		sentence = hold
	}
	return sentence
}
