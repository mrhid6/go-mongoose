package test

import (
	"fmt"

	"github.com/mrhid6/go-mongoose/mongoose"
	"github.com/mrhid6/go-mongoose/mutility"
	test "github.com/mrhid6/go-mongoose/test/models"

	"go.mongodb.org/mongo-driver/bson"
)

// RunTest Run Some Minor tests
func RunTest() {
	//mongoose.
	mongoose.InitiateDB(mongoose.DBConnection{
		Database: "test",
		Host:     "localhost",
		Port:     27017,
		User:     "",
		Password: "",
	})

	tM := test.UserModel{
		Name:  "test",
		Email: "Something@gmail.com",
	}
	if _, err := mongoose.InsertOne(&tM); err != nil {
		panic(err)
	}
	fmt.Println("ID ", tM)
	uM := test.UserModel{}
	// allModels := make([]bson.M, 0)
	err := mongoose.FindOne(bson.M{
		"name": "test",
	}, &uM)

	fmt.Printf("%v\n", mutility.GetID(uM).Hex())

	if err != nil {
		fmt.Println("Error 1 ", err)
	}
	// err = mongoose.FindAll(bson.M{}, test.UserModel{}, &allModels)
	// if err != nil {
	// 	fmt.Println("Error 2 ", err)
	// }

	// fmt.Println(allModels)
	// sModel := make([]test.UserModel, len(allModels))
	// mbytes, _ := bson.Marshal(allModels[0])
	// bson.Unmarshal(mbytes, &sModel[0])

	// if err != nil {
	// 	fmt.Println("Error", err)
	// }
	// fmt.Println("Total Length ", len(allModels), sModel[0].ID)
	// fmt.Println(uM.ID, uM.Name)
	// fmt.Println(mutility.GetName(uM))
	// fmt.Println(uM.Test)

	mutility.CreateIndex(test.UserModel{})
	// uM.Teams = append(uM.Teams, "")
	// test := uM.PopulateTest()
	uM.PopulateTeams()
}
