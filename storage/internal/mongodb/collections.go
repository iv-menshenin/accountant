package mongodb

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	accountsCollectionName = "Accounts"
	paymentsCollectionName = "Payments"
	targetsCollectionName  = "Targets"
	billsCollectionName    = "Bills"
	usersCollectionName    = "Users"
)

type (
	Collection struct {
		Collection *mongo.Collection
		Logger     *log.Logger
	}
)

func (d *Database) Accounts() Collection {
	return Collection{
		Collection: d.client.Database(d.dbName).Collection(accountsCollectionName),
		Logger:     d.logger,
	}
}

func (d *Database) Payments() Collection {
	return Collection{
		Collection: d.client.Database(d.dbName).Collection(paymentsCollectionName),
		Logger:     d.logger,
	}
}

func (d *Database) Targets() Collection {
	return Collection{
		Collection: d.client.Database(d.dbName).Collection(targetsCollectionName),
		Logger:     d.logger,
	}
}

func (d *Database) Bills() Collection {
	return Collection{
		Collection: d.client.Database(d.dbName).Collection(billsCollectionName),
		Logger:     d.logger,
	}
}

func (d *Database) Users() Collection {
	return Collection{
		Collection: d.client.Database(d.dbName).Collection(usersCollectionName),
		Logger:     d.logger,
	}
}
