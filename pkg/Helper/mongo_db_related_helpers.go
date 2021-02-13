package Helper

import (
	"fmt"
	// "net/http"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
	#include <stdlib.h>
*/

// import (
// 	"crypto/cipher"
// 	"net/http"

// 	"C"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )
// import "github.com/oklahomer/go-sarah/examples/simple/plugins/echo"

// // InsertedIDToString function returning Strin representation of Object ID
// // if error returning "" empty string
// func InsertedIDToString( result *mongo.InsertOneResult , echo.Context ) map[string]string {

// }

// ObjectIDFromInsertResult returning an insert ID as a string
// // taking an input as *monogo.InsertOneResult if its not valid
// // it will return an empty string
func ObjectIDFromInsertResult( sires *mongo.InsertOneResult ) string {
	if sires == nil {
		fmt.Println("It Was Nil")
		return ""
	}
	 slices := RemoveObjectIDPrefix(sires.InsertedID.(primitive.ObjectID).String())
	return slices
}