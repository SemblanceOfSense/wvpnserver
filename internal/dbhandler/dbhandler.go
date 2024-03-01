package dbhandler

import (
	"encoding/json"
	"vpnserver/internal/requesthandler"

	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddPublicKey(body requesthandler.PublicKeyRequestStruct) error {
    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil { return err }

    collection := client.Database("wvpn").Collection("PublicKeys")

    json, err := json.Marshal(body)
    if err != nil { return err }

    var doc interface{}
    err = bson.UnmarshalExtJSON(json, true, &doc)
    if err != nil { return err }

    _, err = collection.InsertOne(context.TODO(), doc)
    if err != nil { return err }

    return nil
}
