package main

import (
	"context"
	"fmt"
	"time"

	// "fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/samuael/Project/Weg/internal/pkg/db"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	initialize()
}

// var student *entity.Student

func initialize() {
	autoMigrateToMongo()
	// autoMigrate()
	// PopulateDataBase()
}

var dbs *gorm.DB

func autoMigrate() {
	dbs, erro := db.InitializPostgres()
	if erro != nil {
		log.Fatal(erro)
	}
	defer dbs.Close()
	dbs.AutoMigrate(
		entity.Admin{},
		entity.Alie{},
		entity.User{},
		entity.Group{},
		entity.GroupMembers{},

	)
	// PopulateDataBase()
	err := dbs.GetErrors()
	if err != nil {
		log.Println(err)
	}
}
func autoMigrateToMongo() {
	db := db.ConnectMongodb()
	collection := db.Collection("admin")

	collection.InsertOne(context.TODO()  , entity.Admin{
		// ID : "" , 
		Username : "samuael" , 
		 Role:"samuaelfirst"  ,
		Imageurl: "img/sami.jpg",
		Time: time.Now(),
	})
	// admino := []*entity.Admin{}

	// fmt.Println(admin.)
	filter := bson.D{
		// {"$"  , "_id"} ,
	}
	cursor , era :=collection.Find(context.TODO() , filter )
	// if er != nil {
	// 	log.Fatal("Error while Printing the Result ... ")
	// 	return 
	// }
	// era := singleResult.Decode(admino)
	// if era != nil {
	// 	fmt.Println("Decode Error" , era)
	// 	return 
	// }
	// era:=singleResult.Decode(admino)
	if era != nil {
		fmt.Println("Decode Error" , era)
		return 
	}

	for cursor.Next(context.TODO()) {
		ad := &entity.Admin{}
		cursor.Decode(ad)
		fmt.Println(ad.ID , ad.Username  , ad.Imageurl , ad.Time )
	}
	// bson.Unmarshal(a ,admino  )

}

// drop table payments , course_to_durations , address, students , sessions , langs , teachers , field_assistants  , admins ,cources ,categorys , branchs , resources  , rooms ,  sections , lectures , questions  , asked_quetions , field_sessions , rounds , active_rounds

