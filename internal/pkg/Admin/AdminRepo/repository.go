package AdminRepo

import (
	"context"
	"errors"

	"github.com/samuael/Project/Weg/internal/pkg/Admin"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/pkg/Helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AdminRepo  struct representing mongodb Admin Repository Class
type AdminRepo struct {
	DB *mongo.Database
}

// NewAdminRepo function returning AdminRepo
func NewAdminRepo(db *mongo.Database) Admin.AdminRepo {
	return &AdminRepo{
		DB: db,
	}
}

// CreateAdmin create admin structure
func (adminrepo *AdminRepo) CreateAdmin(admin *entity.Admin) (*entity.Admin, error) {
	insertOneResult, era := adminrepo.DB.Collection(entity.ADMIN).InsertOne(context.TODO(), admin)
	if era != nil {
		return admin, era
	}
	result := Helper.ObjectIDFromInsertResult(insertOneResult)
	if result != "" {
		admin.ID = result
	}
	return admin, nil
}

// GetAdminByID (id string ) *entity.Admin
func (adminrepo *AdminRepo) GetAdminByID(id string) (*entity.Admin, error) {
	admin := &entity.Admin{}
	oid  , er := primitive.ObjectIDFromHex(id)
	if er != nil {
		return nil  , er 
	}
	er = adminrepo.DB.Collection(entity.ADMIN).FindOne( context.TODO() , bson.D{{ "_id", oid }} ).Decode(admin)
	return admin  , er 
}

// DeleteAdminByID (id string ) error 
func (adminrepo *AdminRepo) DeleteAdminByID(id string ) error {
	oid  , er := primitive.ObjectIDFromHex(id)
	if er != nil {
		return er 
	}
	delres  , er := adminrepo.DB.Collection(entity.ADMIN).DeleteOne(context.TODO(), bson.D{{"_id" , oid }})
	if delres.DeletedCount==0 || er != nil {
		return errors.New(" No Record Deleted ")
	}
	return nil
}

// GetAdminByEmail   method to get the admin by his or her email 
func (adminrepo *AdminRepo)  GetAdminByEmail(email string)   (*entity.Admin  , error){
	filter := bson.D{{"email", email },}
	admin := &entity.Admin{}
	era := adminrepo.DB.Collection(entity.ADMIN).FindOne(context.TODO(), filter).Decode(admin)
	return admin, era
}
// SaveAdmin (admin *entity.Admin)  *entity.Admin
func (adminrepo *AdminRepo) SaveAdmin(admin *entity.Admin) (*entity.Admin , error) {
	oid  , er := primitive.ObjectIDFromHex(admin.ID)
	if er != nil {
		return nil  , er
	}
	_, er = adminrepo.DB.Collection(entity.ADMIN).UpdateOne(context.TODO(), bson.D{{"_id", oid}},
		bson.D{
			{"$set", bson.D{
			 {"email", admin.Email}, 
			 {"password", admin.Password}, 
			 {"username"  , admin.Username}}},
		})
	if er != nil {
		return admin, er
	}
	return admin, nil
} 


// AdminEmailExist (   email string ) error
func (adminrepo *AdminRepo)  AdminEmailExist(   email string ) error{
	filter := bson.D{{"email", email}}
	return adminrepo.DB.Collection( entity.ADMIN ).FindOne(context.TODO(), filter).Err()
}