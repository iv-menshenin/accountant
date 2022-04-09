package mongodb

import "go.mongodb.org/mongo-driver/mongo"

const (
	accountsCollectionName = "Accounts"
	paymentsCollectionName = "Payments"
	targetsCollectionName  = "Targets"
	billsCollectionName    = "Bills"
)

func (d *Database) Accounts() *mongo.Collection {
	return d.client.Database(d.dbName).Collection(accountsCollectionName)
}

func (d *Database) Payments() *mongo.Collection {
	return d.client.Database(d.dbName).Collection(paymentsCollectionName)
}

func (d *Database) Targets() *mongo.Collection {
	return d.client.Database(d.dbName).Collection(targetsCollectionName)
}

func (d *Database) Bills() *mongo.Collection {
	return d.client.Database(d.dbName).Collection(billsCollectionName)
}
