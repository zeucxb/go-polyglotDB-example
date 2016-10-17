package main

import (
	"fmt"
	"log"

	"github.com/jmcvetta/neoism"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	db := connect()
	cq := neoism.CypherQuery{
		Statement: `
			MATCH (u:User)
			WHERE u.ID = {id}
			MATCH (u)-[f:Follow]->(following)
			RETURN following.ID
			`,
		Parameters: neoism.Props{"id": "57eac14aee610931a7732d60"},
		Result:     &result{},
	}
	db.Cypher(&cq)

	r := cq.Result.(*result)

	mongo := connectMongo("users")

	var user mongoUser

	for _, v := range *r {
		mongo.FindId(v.ID).One(&user)

		fmt.Println(user)
	}
}

type result []struct {
	ID bson.ObjectId `json:"following.ID" bson:"_id"`
}

type mongoUser struct {
	Username string `bson:"username"`
	Email    string `bson:"email"`
}

func connect() *neoism.Database {
	db, err := neoism.Connect("http://neo4j:12345678@localhost:7474/db/data")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// connectMongo return a mongo session
func connectMongo(collectionName string) *mgo.Collection {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo runnig?
	if err != nil {
		panic(err)
	}

	return s.DB("project").C(collectionName)
}
