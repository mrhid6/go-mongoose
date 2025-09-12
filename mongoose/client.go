package mongoose

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	pluralize "github.com/gertd/go-pluralize"
)

var pluralizer = pluralize.NewClient()

type MongooseClient struct {
	db            *mongo.Database
	mu            sync.RWMutex
	modelRegistry map[string]*Model
	debug         bool
}

type DBConnectionOptions struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string

    AuthSource *string

	SRV bool

	Debug bool
}

func GetConnectionOptionsFromEnv() *DBConnectionOptions {
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
    authSource := os.Getenv("DB_AUTHSOURCE")

	return &DBConnectionOptions{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		Database: os.Getenv("DB_DATABASE"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SRV:      os.Getenv("DB_SRV") == "true",
		Debug:    os.Getenv("DB_DEBUG") == "true",
        AuthSource: &authSource,
	}
}

func (dbConnection *DBConnectionOptions) BuildConnectionURI() string {
	if dbConnection.Port == 0 {
		dbConnection.Port = 27017
	}
	urlHeader := "mongodb://"
	if dbConnection.SRV {
		urlHeader = "mongodb+srv://"
	}

	resConnectionURI := ""

	if dbConnection.User == "" {
		resConnectionURI = urlHeader + dbConnection.Host
	} else {
		resConnectionURI = urlHeader + url.QueryEscape(dbConnection.User) + ":" + url.QueryEscape(dbConnection.Password) + "@" + dbConnection.Host
	}

	if !dbConnection.SRV {
		resConnectionURI += ":" + strconv.Itoa(dbConnection.Port)
	}

	if dbConnection.Database != "" {
		resConnectionURI += "/" + dbConnection.Database
	}

    if dbConnection.AuthSource != nil{
        resConnectionURI += fmt.Sprintf("?authSource=%s", *dbConnection.AuthSource)
    }

	return resConnectionURI
}

func NewMongoClient(ctx context.Context, connectionOptions *DBConnectionOptions) (*MongooseClient, error) {

	uri := connectionOptions.BuildConnectionURI()

	tM := reflect.TypeOf(bson.M{})

	reg := bson.NewRegistry()
	reg.RegisterTypeMapEntry(bson.TypeEmbeddedDocument, tM)

	clientOpts := options.Client().ApplyURI(uri).SetRegistry(reg)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// ping
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(connectionOptions.Database)

	if connectionOptions.Debug {
		fmt.Printf("Connected to mongo server: %s on port: %d using database: %s\n",
			connectionOptions.Host,
			connectionOptions.Port,
			connectionOptions.Database,
		)
	}

	return &MongooseClient{
		db:            db,
		modelRegistry: make(map[string]*Model),
		debug:         connectionOptions.Debug,
	}, nil
}

func (c *MongooseClient) Disconnect(ctx context.Context) error {
	if c == nil || c.db == nil {
		return nil
	}
	if err := c.db.Client().Disconnect(ctx); err != nil {
		return err
	}
	if c.debug {
		fmt.Println("Disconnected from database")
	}
	return nil
}

func (c *MongooseClient) RegisterModel(schema interface{}) (*Model, error) {
	t := reflect.TypeOf(schema)

	if t.Kind() != reflect.Ptr {
		return nil, errors.New("schema must be a pointer to struct")
	}

	elemType := t.Elem()
	if elemType.Kind() != reflect.Struct {
		return nil, errors.New("schema must be a pointer to struct")
	}

	name := elemType.Name()
	name = strings.TrimSuffix(name, "Schema")

	collectionName := strings.ToLower(pluralizer.Plural(name))

	c.mu.Lock()
	c.modelRegistry[name] = &Model{
		CollectionName: collectionName,
		SchemaType:     elemType,
		client:         c,
	}
	c.mu.Unlock()

	if c.debug {
		fmt.Printf("Registered model: %s (collection: %s)\n", name, collectionName)
	}
	return c.modelRegistry[name], nil
}

func (c *MongooseClient) GetModel(name string) (*Model, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if m, ok := c.modelRegistry[name]; ok {
		return m, nil
	}
	return nil, fmt.Errorf("model %s not registered", name)
}

func (c *MongooseClient) GetModelByCollection(name string) (*Model, error) {
	for _, model := range c.modelRegistry {
		if model.CollectionName == name {
			return model, nil
		}
	}
	return nil, fmt.Errorf("no model found for collection %s", name)
}

func (c *MongooseClient) GetDatabase() *mongo.Database {
	return c.db
}

func (c *MongooseClient) GetCollection(collectionName string) *mongo.Collection {
	return c.db.Collection(collectionName)
}
