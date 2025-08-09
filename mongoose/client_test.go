package mongoose

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	mongooseClient *MongooseClient
)

type UserSchema struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
}

type AccountSchema struct {
	ID           primitive.ObjectID `bson:"_id"`
	OwningUserId primitive.ObjectID `bson:"owningUser" mson:"collection=users"`
	OwningUser   UserSchema         `bson:"-"`

	Users   []UserSchema `bson:"-"`
	UserIds primitive.A  `bson:"users" mson:"collection=users"`
}

func TestClient(t *testing.T) {
	dbConnectionOptions := GetConnectionOptionsFromEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := NewMongoClient(ctx, dbConnectionOptions)
	if err != nil {
		t.Errorf(`mongooseClient = %v, error = %v`, mongooseClient, err)
	}

	mongooseClient = c
}

func TestClientModels_RegisterModels(t *testing.T) {
	mongooseClient.RegisterModel(&UserSchema{})
	mongooseClient.RegisterModel(&AccountSchema{})

	// Check getting the Account model
	AccountModel, err := mongooseClient.GetModel("Account")
	if err != nil {
		t.Errorf(`AccountModel = %v, error = %v`, AccountModel, err)
		return
	}

	// Check getting the User model
	UserModel, err := mongooseClient.GetModel("User")
	if err != nil {
		t.Errorf(`UserModel = %v, error = %v`, UserModel, err)
		return
	}
}

func TestClientModels(t *testing.T) {

	// Check getting the Account model
	AccountModel, err := mongooseClient.GetModel("Account")
	if err != nil {
		t.Errorf(`AccountModel = %v, error = %v`, AccountModel, err)
		return
	}

	// Check getting the User model
	UserModel, err := mongooseClient.GetModel("User")
	if err != nil {
		t.Errorf(`UserModel = %v, error = %v`, UserModel, err)
		return
	}

	NewUser := &UserSchema{
		ID:       primitive.NewObjectID(),
		Username: "test",
		Password: "Hello",
	}

	// Check insert of one User document
	if err := UserModel.Create(NewUser); err != nil {
		t.Errorf(`error = %v`, err)
		return
	}

	userIds := make(primitive.A, 0)
	userIds = append(userIds, NewUser.ID)

	NewAccount := &AccountSchema{
		ID:           primitive.NewObjectID(),
		OwningUserId: NewUser.ID,
		UserIds:      userIds,
	}

	// Check insert of one Account document
	if err := AccountModel.Create(NewAccount); err != nil {
		t.Errorf(`error = %v`, err)
		return
	}

	theAccount := &AccountSchema{}

	// Check finding one Account document
	if err := AccountModel.FindOne(theAccount, bson.M{"_id": NewAccount.ID}); err != nil {
		t.Errorf(`theAccount = %v, error = %v`, theAccount, err)
		return
	}

	// Check single populate of user
	if err := AccountModel.PopulateField(&theAccount, "OwningUser"); err != nil {
		t.Errorf(`theAccount = %v, error = %v`, theAccount, err)
		return
	}

	// Check populate array of users
	if err := AccountModel.PopulateField(&theAccount, "Users"); err != nil {
		t.Errorf(`theAccount = %v, error = %v`, theAccount, err)
		return
	}
	t.Logf(`theAccount = %v`, theAccount)

	// Check document update
	if err := UserModel.UpdateData(NewUser, bson.M{"username": "test1"}); err != nil {
		t.Errorf(`NewUser = %v, error = %v`, NewUser, err)
		return
	}

	// Check invalid schema document update
	if err := UserModel.UpdateData(NewUser, bson.M{"username1": "test1"}); err == nil {
		t.Errorf(`NewUser = %v, error = %v`, NewUser, err)
		return
	}

	// Check Delete User Document
	if err := UserModel.Delete(bson.M{"_id": NewUser.ID}); err != nil {
		t.Errorf(`error = %v`, err)
		return
	}

	// Check Delete Account Document
	if err := AccountModel.Delete(bson.M{"_id": NewAccount.ID}); err != nil {
		t.Errorf(`error = %v`, err)
		return
	}
}

func TestDisconnect(t *testing.T) {

	if mongooseClient != nil {
		if err := mongooseClient.Disconnect(context.Background()); err != nil {
			t.Errorf(`mongooseClient = %v, error = %v`, mongooseClient, err)
			return
		}
	}
}
