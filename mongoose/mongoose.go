package mongoose

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo This is the Mongo struct
type Mongo struct {
	client       *mongo.Client
	Database     *mongo.Database
	dbConnection DBConnection
	Err          error
}

// DBConnection DB Connection Details
type DBConnection struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string

	ConnectionURL string

	SRV bool
}

func (dbConnection *DBConnection) BuildConnectionURL() {
	if dbConnection.ConnectionURL != "" {
		return
	}

	if dbConnection.Port == 0 {
		dbConnection.Port = 27017
	}
	urlHeader := "mongodb://"
	if dbConnection.SRV {
		urlHeader = "mongodb+srv://"
	}

	if dbConnection.User == "" {
		dbConnection.ConnectionURL = urlHeader + dbConnection.Host
	} else {
		dbConnection.ConnectionURL = urlHeader + url.QueryEscape(dbConnection.User) + ":" + url.QueryEscape(dbConnection.Password) + "@" + dbConnection.Host
	}

	if !dbConnection.SRV {
		dbConnection.ConnectionURL += ":" + strconv.Itoa(dbConnection.Port)

	}
}

// ShortWaitTime Small Wait time
// MediumWaitTime Medium Wait Time
// LongWaitTime Long wait time
var (
	_mongo Mongo

	ShortWaitTime  time.Duration = 2
	MediumWaitTime time.Duration = 5
	LongWaitTime   time.Duration = 10
)

// InitiateDB This needs to be called if you are using some other than default DB
func InitiateDB(dbConnection DBConnection) {
	dbConnection.BuildConnectionURL()

	_mongo.dbConnection = dbConnection
}

func TestConnection() error {

	if _, err := Get(); err != nil {
		return err
	}

	return nil
}

// Get This function will recieve the Mongo structure
func Get() (Mongo, error) {
	if _mongo.client == nil {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_mongo.client, _mongo.Err = mongo.Connect(ctx, options.Client().ApplyURI(_mongo.dbConnection.ConnectionURL))
		if _mongo.Err != nil {
			return _mongo, _mongo.Err
		}

		_mongo.Database = _mongo.client.Database(_mongo.dbConnection.Database)
	}
	return _mongo, nil
}
