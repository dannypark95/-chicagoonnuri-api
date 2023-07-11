package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PDF struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Filename   string             `bson:"filename" json:"filename"`
	URL        string             `bson:"url" json:"url"`
	UploadDate time.Time          `bson:"uploadDate" json:"uploadDate"`
}
