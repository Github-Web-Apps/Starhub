package datastore

// Datastore access to database
type Datastore interface {
	Tokenstore
	Execstore
	Userdatastore
}
