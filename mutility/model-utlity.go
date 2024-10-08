package mutility

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
func GetID(a interface{}) primitive.ObjectID {
	t := reflect.TypeOf(a)
	tv := reflect.ValueOf(a)

	if t.Kind() == reflect.Ptr {
		tv = reflect.Indirect(tv)
	}

	mVal := tv.FieldByName("ID")
	if reflect.Value.IsZero(mVal) {
		return primitive.NilObjectID
	}
	aV := mVal.Interface().(primitive.ObjectID)
	return aV
}

// CreateIndex This function would be used to create index for the table
func CreateIndex(a interface{}) {
	t := reflect.TypeOf(a)
	t.Name()
	var indexes []interface{} = make([]interface{}, 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagStr := field.Tag.Get("mson")
		tags := strings.Split(tagStr, ",")
		if len(tags) == 0 {
			continue
		}
		index := analyzeAndCreateTagIndex(tags)

		if index != nil {
			indexes = append(indexes, index)
		}
	}
}

func analyzeAndCreateTagIndex(tags []string) *interface{} {
	// TODO Do some stuff to analyze the tags
	// fmt.Println("Tags ", tags)

	return nil
}
