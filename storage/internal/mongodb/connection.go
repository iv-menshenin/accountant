package mongodb

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/iv-menshenin/accountant/model"
)

type (
	Database struct {
		logger *log.Logger
		client *mongo.Client
		dbName string
	}
)

func New(config model.Config, logger *log.Logger) (*Database, error) {
	var db = Database{logger: logger}
	optClient, err := db.newOptions(config)
	if err != nil {
		return nil, err
	}
	db.client, err = mongo.NewClient(optClient)
	if err != nil {
		return nil, err
	}
	err = db.client.Connect(context.Background())
	if err != nil {
		return nil, err
	}
	err = db.client.Ping(context.Background(), nil)

	return &db, err
}

func (d *Database) newOptions(config model.Config) (*options.ClientOptions, error) {
	var (
		port   = config.IntegerConfig("mongo-port", "mongo-port", "MONGO_PORT", 27017, "mongodb server port")
		host   = config.StringConfig("mongo-host", "mongo-host", "MONGO_HOST", "localhost", "mongodb server DNS name (or IP)")
		dbName = config.StringConfig("mongo-dbname", "mongo-dbname", "MONGO_DBNAME", "mongo", "mongodb database name")
		user   = config.StringConfig("mongo-username", "mongo-username", "MONGO_USERNAME", "mongo", "mongodb username")
		pass   = config.StringConfig("mongo-password", "mongo-password", "MONGO_PASSWORD", "mongo", "mongodb password")
		repl   = config.BooleanConfig("mongo-replica-set", "mongo-replica-set", "MONGO_REPLICA", "use replica set configuration")

		socketTimeout     = config.DurationConfig("mongo-socket-timeout", "mongo-socket-timeout", "MONGO_SOCKET_TIMEOUT", time.Second*5, "specifies how long the driver will wait for a socket read or write to return before returning a network error")
		connectionTimeout = config.DurationConfig("mongo-connection-timeout", "mongo-connection-timeout", "MONGO_CONNECTION_TIMEOUT", time.Second*15, "specifies a timeout that is used for creating connections to the server")
		selectionTimeout  = config.DurationConfig("mongo-server-selection-timeout", "mongo-server-selection-timeout", "MONGO_SERVER_SELECTION_TIMEOUT", time.Second*15, "specifies how long the driver will wait to find an available, suitable server")
		idleTime          = config.DurationConfig("mongo-idle-timeout", "mongo-idle-timeout", "MONGO_IDLE_TIMEOUT", time.Minute, "specifies the maximum amount of time that a connection will remain idle in a connection pool")

		poolSize = config.IntegerConfig("mongo-pool-size", "mongo-pool-size", "MONGO_POOL_SIZE", 32, "specifies that maximum number of connections allowed in the driver's connection pool to each server")
		retryR   = config.BooleanConfig("mongo-retry-reads", "mongo-retry-reads", "MONGO_RETRY_READS", "specifies whether supported read operations should be retried once on certain errors")
		retryW   = config.BooleanConfig("mongo-retry-writes", "mongo-retry-writes", "MONGO_RETRY_WRITES", "specifies whether supported write operations should be retried once on certain errors")
	)
	if err := config.Init(); err != nil {
		return nil, err
	}
	d.dbName = *dbName
	var maker mongoURLMaker
	if *repl {
		maker = makeMongoSrvConnectionURI
	} else {
		maker = makeMongoDbConnectionURI
	}
	URL := maker(fmt.Sprintf("%s:%d", *host, *port), *dbName, *user, *pass)
	optClient := options.Client()
	optClient.ApplyURI(URL.String()).
		SetSocketTimeout(*socketTimeout).
		SetConnectTimeout(*connectionTimeout).
		SetServerSelectionTimeout(*selectionTimeout).
		SetMaxConnIdleTime(*idleTime).
		SetMaxPoolSize(uint64(*poolSize)).
		SetRetryReads(*retryR).
		SetRetryWrites(*retryW)
	return optClient, nil
}

type mongoURLMaker func(hostPort, dbName, user, password string) url.URL

// makeMongoDbConnectionURI creates a connection url.URL for concrete server.
func makeMongoDbConnectionURI(hostPort, dbName, user, password string) url.URL {
	return url.URL{
		Scheme: "mongodb",
		User:   url.UserPassword(user, password),
		Host:   hostPort,
		Path:   "/" + dbName,
	}
}

// makeMongoSrvConnectionURI creates a connection url.URL for replica set.
func makeMongoSrvConnectionURI(hostPort, dbName, user, password string) url.URL {
	return url.URL{
		Scheme: "mongodb+srv",
		User:   url.UserPassword(user, password),
		Host:   hostPort,
		Path:   "/" + dbName,
	}
}
