package entity

import (
	"github.com/dgrijalva/jwt-go"
)

var (
	// PathToStudentsFromTemplates constant
	PathToStudentsFromTemplates = "img/Students/"
	// TemplateDirectoryFromMain  varaibe
	TemplateDirectoryFromMain = "../../web/templates/"

	// Languages List of Languages 
	Languages = []string{
		"English","አማርኛ" , "Oromiffa", "ትግርኛ"}
	// Langs map of languages and their representing vals
	Langs = map[string]string{
		"English":"eng"  , "አማርኛ":"amh" , "Oromiffa" :"oro"  , "ትግርኛ":"tig" ,   
	}

	// ImageExtensions list of valid image extensions 
	ImageExtensions= []string{"jpeg", "png", "jpg", "gif", "btmp"}

	// LangCodes list of short form of lang representation
	LangCodes = map[string][]string {
		"eng" : {"en" , "english"  , "eng"  ,  "inglish"  } , 
		"amh" : { "amharic" , "amh"  ,  "amaregna" , "amargna" },
	}

	

)


const (
	
	// MsgSeen message ID 
	MsgSeen =1
	// MsgTyping message ID 
	MsgTyping =2
	// MsgStopTyping message ID 
	MsgStopTyping =3
	// MsgIndividualTxt individual Message
	MsgIndividualTxt =4
	// MsgAlieProfileChange alie profile change message  
	MsgAlieProfileChange=6
	// MsgNewAlie you get new alie 
	MsgNewAlie = 8


	// MsgGroupTxt group Message 
	MsgGroupTxt =5
	// MsgGroupProfileChange group Profile Change  
	MsgGroupProfileChange= 7
	// MsgGroupJoin group join message  
	MsgGroupJoin= 9
	// MsgGroupLeave group Leave Message
	MsgGroupLeave= 10

)


const (

	// USER const
	USER ="user"
	// ADMIN const
	ADMIN ="admin"
	// GROUP const
	GROUP ="grup"
	// ALIE const
	ALIE="alie"
	

	// DBName  constant representing the name of the Database both in Postgres and in mongodb 
	DBName="weg"
	// SessionName constant to hold the absolute sessio name 
	SessionName  = "weg"
	// LanguageSessionName constant to hold the language cookie name 
	LanguageSessionName ="weg-lang"
	// PathToTemplates path to templates folder 
	PathToTemplates = "../../web/templates/"
	// PathToUserImageFromaMain path to UserImage folder from main function 
	PathToUserImageFromaMain = "../../web/templates/img/UserImage/"
	// UserImagesPath to be used as path prefix for images 
	UserImagesPath ="img/UserImage/"
	// PathToResources  constant
	PathToResources = "../../web/templates/Source/Resources/"
	// FileSchema  const
	FileSchema = "file:///"
	// PROTOCOL string
	PROTOCOL = "http://"
	// HOST string
	HOST = "localhost:9900/"
)




// Session representing the Sesstion to Be sent with the request body
// no saving of a session in the database so i Will use this session in place of
type Session struct {
	jwt.StandardClaims
	UserID   string 
	Username string
	Email string 
	Role string 
	// Imageurl and Username are Included in this iggo
}


// 