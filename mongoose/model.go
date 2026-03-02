package mongoose

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mrhid6/go-mongoose/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ShortWaitTime  time.Duration = 2
	MediumWaitTime time.Duration = 5
	LongWaitTime   time.Duration = 10
)

type Model struct {
	CollectionName string
	SchemaType     reflect.Type
	client         *MongooseClient
}

func (m *Model) FindAll(results interface{}, filter bson.M) error {
	return m.FindAllWithOptions(results, filter, options.Find())
}

func (m *Model) FindAllWithOptions(results interface{}, filter bson.M, opts *options.FindOptionsBuilder) error {
	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()

	col := m.client.db.Collection(m.CollectionName)
	cursor, err := col.Find(ctx, filter, opts)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	return cursor.All(ctx, results)
}

func (m *Model) FindOne(result interface{}, filter bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()

	col := m.client.db.Collection(m.CollectionName)
	return col.FindOne(ctx, filter).Decode(result)
}

func (m *Model) FindOneById(result interface{}, id bson.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()

	col := m.client.db.Collection(m.CollectionName)
	return col.FindOne(ctx, bson.M{"_id": id}).Decode(result)
}

func (m *Model) FindOneAndUpdate(result interface{}, filter bson.M, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	col := m.client.db.Collection(m.CollectionName)
	err := col.FindOneAndUpdate(ctx, filter, update, opts).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) Create(doc interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), ShortWaitTime*time.Second)
	defer cancel()

	col := m.client.db.Collection(m.CollectionName)
	_, err := col.InsertOne(ctx, doc)
	return err
}

func (m *Model) Delete(filter bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()
	col := m.client.db.Collection(m.CollectionName)
	_, err := col.DeleteMany(ctx, filter)
	return err
}

func (m *Model) DeleteById(id bson.ObjectID) error {

	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()

	col := m.client.db.Collection(m.CollectionName)
	_, err := col.DeleteMany(ctx, bson.M{"_id": id})
	return err
}

func (m *Model) UpdateData(doc interface{}, updateData bson.M) error {

	schemaType := m.SchemaType
	validFields := utils.GetValidBsonFields(schemaType)

	// Filter updateData keys
	filteredData := bson.M{}
	for k, v := range updateData {
		if _, ok := validFields[k]; ok {
			filteredData[k] = v
		} else {
			return fmt.Errorf("field %q is not valid for schema %s", k, schemaType.Name())
		}
	}

	col := m.client.db.Collection(m.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()

	_, err := col.UpdateOne(ctx, bson.M{
		"_id": utils.GetID(doc),
	}, bson.M{"$set": updateData})

	if err != nil {
		return err
	}
	return nil
}

func (m *Model) RawUpdateData(doc interface{}, updateData bson.M) error {

	col := m.client.db.Collection(m.CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), MediumWaitTime*time.Second)
	defer cancel()

	_, err := col.UpdateOne(ctx, bson.M{
		"_id": utils.GetID(doc),
	}, updateData)

	if err != nil {
		return err
	}
	return nil
}

func (m *Model) PopulateField(objPtr interface{}, fieldName string) error {
	val := reflect.ValueOf(objPtr)
	if val.Kind() != reflect.Ptr {
		return errors.New("objPtr must be a pointer")
	}

	// Dereference until we hit the struct
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return errors.New("nil pointer passed to PopulateField")
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("PopulateField expects pointer to struct, got %s", val.Kind())
	}

	structType := val.Type()

	// Find the destination field (e.g., "User" or "Users")
	destField, ok := structType.FieldByName(fieldName)
	if !ok {
		return fmt.Errorf("no field with name %s", fieldName)
	}

	// Figure out matching ID field name
	var idFieldName string
	if pluralizer.IsPlural(fieldName) {
		idFieldName = pluralizer.Singular(fieldName) + "Ids"
	} else {
		idFieldName = fieldName + "Id"
	}

	idStructField, ok := structType.FieldByName(idFieldName)
	if !ok {
		return fmt.Errorf("no matching ID field %s", idFieldName)
	}

	// Parse mson tag to get collection
	tagVal := idStructField.Tag.Get("mson")
	var collectionName string
	if tagVal != "" {
		parts := strings.Split(tagVal, ",")
		for _, part := range parts {
			keyVal := strings.SplitN(part, "=", 2)
			if len(keyVal) == 2 && keyVal[0] == "collection" {
				collectionName = keyVal[1]
				break
			}
		}
	}
	if collectionName == "" {
		return fmt.Errorf("no collection tag found on %s", idFieldName)
	}

	// Find model for that collection
	targetModel, err := m.client.GetModelByCollection(collectionName)
	if err != nil {
		return err
	}

	idFieldVal := val.FieldByName(idFieldName)
	idFieldType := idFieldVal.Type()

	// ---- CASE 1: bson.A ----
	if idFieldType == reflect.TypeOf(bson.A{}) {
		var ids []bson.ObjectID
		for _, v := range idFieldVal.Interface().(bson.A) {
			if oid, ok := v.(bson.ObjectID); ok {
				ids = append(ids, oid)
			}
		}

		resultsPtr := reflect.New(destField.Type).Interface()
		if len(ids) > 0 {
			if err := targetModel.FindAll(resultsPtr, bson.M{"_id": bson.M{"$in": ids}}); err != nil {
				return err
			}
		}

		val.FieldByName(fieldName).Set(reflect.ValueOf(resultsPtr).Elem())
		return nil
	}

	// ---- CASE 2: Single ObjectID ----
	if idFieldType == reflect.TypeOf(bson.ObjectID{}) {
		objID, _ := idFieldVal.Interface().(bson.ObjectID)
		targetInstance := reflect.New(targetModel.SchemaType).Interface()
		if err := targetModel.FindOne(targetInstance, bson.M{"_id": objID}); err != nil {
			return err
		}

		val.FieldByName(fieldName).Set(reflect.ValueOf(reflect.ValueOf(targetInstance).Elem().Interface()))
		return nil
	}

	return fmt.Errorf("unsupported ID field type: %s", idFieldType.String())
}
