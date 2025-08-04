package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product model
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name" binding:"required"`
	Description string             `bson:"description" json:"description" binding:"required"`
	Price       float64            `bson:"price" json:"price" binding:"required,gt=0"`
	Stock       int                `bson:"stock" json:"stock" binding:"required,gte=0"`
	Category    string             `bson:"category" json:"category" binding:"required"`
	ImageURL    string             `bson:"image_url" json:"image_url"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
