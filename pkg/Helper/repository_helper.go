package Helper

import "strings"

// RemoveObjectIDPrefix function returning the real internal Object ID from
// ObjectID prefixed object ID Result
// Example Input : ObjectID("5fe1b21d88b1deda65a9a507") :
// 		   OutPut : "5fe1b21d88b1deda65a9a507"
func RemoveObjectIDPrefix(objectid string) string {
	objectid = strings.TrimSuffix(strings.TrimPrefix(objectid , "ObjectID(\"") , "\")")
	return objectid
}