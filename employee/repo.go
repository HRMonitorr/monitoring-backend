package employee

import (
	"context"
	"github.com/HRMonitorr/PasetoprojectBackend"
	"github.com/HRMonitorr/monitoring-backend/structure"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertCommitToDB(mconn *mongo.Database, val structure.Commits) (inserted interface{}) {
	return PasetoprojectBackend.InsertOneDoc(mconn, "commit", val)
}

func InsertCommitsManyToDB(mconn *mongo.Database, commits []structure.Commits) (inserted interface{}, err error) {
	collection := mconn.Collection("commit")

	// Convert []structure.Commits to []interface{}
	var documents []interface{}
	for _, commit := range commits {
		documents = append(documents, commit)
	}

	// Insert multiple documents
	result, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		return nil, err
	}

	return result.InsertedIDs, nil
}
