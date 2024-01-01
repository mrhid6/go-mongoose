package mongoose

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/mrhid6/go-mongoose/mutility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PopulateObject an Object
func PopulateObject(objPtr interface{}, fieldName string, modelPtr interface{}) error {

	if mutility.IsPointer(objPtr) {
		return errors.New("objptr should be a Pointer")
	}
	if mutility.IsPointer(modelPtr) {
		return errors.New("modelPtr should be a Pointer")
	}

	t := reflect.TypeOf(objPtr)
	val := reflect.ValueOf(objPtr).Elem().FieldByName(fieldName)

	t = t.Elem()
	f, b := t.FieldByName(fieldName)
	if !b {
		return fmt.Errorf("error populating object no field with %s", fieldName)
	}

	tags := strings.Split(f.Tag.Get("mson"), ",")
	for i := range tags {
		tag := strings.Split(tags[i], "=")
		if len(tag) != 2 {
			return nil
		}
		if tag[0] != "collection" {
			return nil
		}

		t1 := val.Interface().(primitive.ObjectID)

		err := FindByObjectID(t1, modelPtr)
		if err != nil {
			return err
		}
	}

	return nil
}

// PopulateObjectArray Populates the Object Array
func PopulateObjectArray(objPtr interface{}, field string, modelArrPtr interface{}) error {

	if mutility.IsPointer(objPtr) {
		return errors.New("objptr should be a Pointer")
	}
	if mutility.IsPointer(modelArrPtr) {
		return errors.New("modelarrptr should be a Pointer")
	}

	t := reflect.TypeOf(objPtr)

	val := reflect.ValueOf(objPtr).Elem().FieldByName(field)

	t = t.Elem()
	f, b := t.FieldByName(field)
	if !b {
		return fmt.Errorf("error populating array no field with %s", field)
	}

	tags := strings.Split(f.Tag.Get("mson"), ",")
	for i := range tags {
		tag := strings.Split(tags[i], "=")
		if len(tag) != 2 {
			continue
		}
		if tag[0] != "collection" {
			continue
		}

		objIds := val.Interface().(primitive.A)
		err := FindAll(bson.M{
			"_id": bson.M{
				"$in": objIds,
			},
		}, modelArrPtr)

		if err != nil {
			return err
		}
	}

	return nil
}
