package graph

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)


func isNewUserID(idNum int) bool {
	if idNum < 0 {
		return false
	}
	client , ctx:= getDbClient()
	q := client.Collection(USER_COLLECTION_NAME).Where(USERID_FIELD_NAME, "==", idNum).Limit(1)
	docArr, err := q.Documents(ctx).GetAll()
	//fmt.Println(len(docArr))
	if len(docArr) == 0  {
		return false
	}
	if err != nil{
		fmt.Println(err.Error())
	}
	return true
	//if there are 0 users with the same ID number then we can put a new user with this ID
}
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

