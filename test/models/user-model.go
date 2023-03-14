package test

import (
	"fmt"

	"go-mongoose/mongoose"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//UserModel This is the model for the Users
type UserModel struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" `
	Name          string             `json:"name,omitempty" bson:"name,omitempty"`
	Password      string             `json:"password,omitempty" bson:"password,omitempty"`
	Email         string             `json:"email,omitempty" bson:"email,omitempty" mson:"cunique"`
	FirebaseToken string             `json:"firebaseToken" bson:"firebaseToken" mson:"unique"`
	Test          primitive.ObjectID `json:"test" bson:"test" mson:"collection=UserModel"`
	Teams         primitive.A        `json:"teams" bson:"teams" mson:"collection=UserModel"`
	MainProfile   string             `json:"mainProfile" bson:"mainProfile"`
	Phone         string             `json:"phone" bson:"phone" mson:"cunique"`
	SocialMedia   string             `json:"socialMedia" bson:"socialMedia"`
	Strengths     string             `json:"strengths" bson:"strengths"`
	Developments  string             `json:"developments" bson:"developments"`
	UserType      int                `json:"userType" bson:"userType"`
	Status        int                `json:"status" bson:"status"`
}

//PopulateTest This will populate the test field and give to you
func (c *UserModel) PopulateTest() *UserModel {
	u := UserModel{}
	mongoose.PopulateObject(c, "Test", &u)
	return &u
}

//PopulateTeams This will populate the teams and give it to you
func (c *UserModel) PopulateTeams() *[]UserModel {
	teams := make([]UserModel, 0)
	err := mongoose.PopulateObjectArray(c, "Teams", &teams)
	if err != nil {
		fmt.Println("Error While Populate ", err)
	}
	return &teams
}
