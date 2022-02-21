package mongodb

import "go.mongodb.org/mongo-driver/mongo"

const (
	accountsCollectionName = "Accounts"
)

func (d *Database) Accounts() *mongo.Collection {
	return d.client.Database(d.dbName).Collection(accountsCollectionName)
}
