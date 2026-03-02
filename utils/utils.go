package utils

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func IsPointer(a interface{}) bool {
	t := reflect.TypeOf(a)
	return t.Kind() == reflect.Ptr
}

// GetName Returns the collection Name
func GetName(a interface{}) string {
	t := reflect.TypeOf(a)

	if t.Kind() == reflect.String {
		return fmt.Sprintf("%v", a)
	}

	return getName(t)
}
func getName(t reflect.Type) string {
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr || t.Kind() == reflect.Array || t.Kind() == reflect.Map || t.Kind() == reflect.Chan {
		return getName(t.Elem())
	}

	return strings.ToLower(t.Name())
}

// GetID Returns the Object ID
func GetID(a interface{}) bson.ObjectID {
	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		// Not a struct, return NilObjectID
		return bson.NilObjectID
	}

	// Iterate struct fields to find the one with bson:"_id"
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		bsonTag := field.Tag.Get("bson")
		// bson tag can be e.g. "_id,omitempty", so split by comma and check first part
		bsonName := strings.Split(bsonTag, ",")[0]

		if bsonName == "_id" {
			fieldValue := v.Field(i)
			// Check zero
			if fieldValue.IsZero() {
				return bson.NilObjectID
			}

			// Make sure the type is bson.ObjectID before casting
			if oid, ok := fieldValue.Interface().(bson.ObjectID); ok {
				return oid
			}

			// If not an ObjectID, try to convert from string if needed (optional)
			// or return NilObjectID
			return bson.NilObjectID
		}
	}

	// Not found, return NilObjectID
	return bson.NilObjectID
}

func GetValidBsonFields(t reflect.Type) map[string]struct{} {
	fields := make(map[string]struct{})

	// Iterate through struct fields, including embedded structs
	var walkFields func(reflect.Type)
	walkFields = func(typ reflect.Type) {
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			if f.Anonymous && f.Type.Kind() == reflect.Struct {
				walkFields(f.Type)
				continue
			}
			bsonTag := f.Tag.Get("bson")
			if bsonTag == "" || bsonTag == "-" {
				continue
			}
			// Split tag by comma in case of options
			parts := strings.Split(bsonTag, ",")
			if parts[0] != "" {
				fields[parts[0]] = struct{}{}
			}
		}
	}

	walkFields(t)
	return fields
}
