package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Category struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	Status    int       `bson:"status" json:"status"`
}

type Contact struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	ItemID    string    `bson:"itemId" json:"itemId"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	Message   string    `bson:"message" json:"message"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type LostFoundItem struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Title       string    `bson:"title" json:"title"`
	Description string    `bson:"description" json:"description"`
	Location    string    `bson:"location" json:"location"`
	CategoryID  string    `bson:"categoryId" json:"categoryId"`
	Image       string    `bson:"image" json:"image"`
	StatusType  string    `bson:"statusType" json:"statusType"`
	Type        string    `bson:"type" json:"type"`
	Date        time.Time `bson:"date" json:"date"` // FIXED
	Time        string    `bson:"time" json:"time"`
	Status      int       `bson:"status" json:"status"`
	Email       string    `bson:"email" json:"email"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
}
type Notification struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Title     string    `bson:"title" json:"title"`
	Message   string    `bson:"message" json:"message"`
	FcmToken  string    `bson:"fcmToken" json:"fcmToken"`
	Image     string    `bson:"image" json:"image"`
	Type      string    `bson:"type" json:"type"`
	Date      string    `bson:"date" json:"date"`
	Time      string    `bson:"time" json:"time"`
	Status    int       `bson:"status" json:"status"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
type User struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Email       string    `bson:"email" json:"email"`
	Password    string    `bson:"password" json:"password"`
	Role        string    `bson:"role" json:"role"`
	Status      int       `bson:"status" json:"status"`
	FcmToken    string    `bson:"fcmToken" json:"fcmToken"`
	BearerToken string    `bson:"bearerToken" json:"bearerToken"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
}

func NewUser(email, password, role string) User {
	return User{
		ID:        primitive.NewObjectID().String(),
		Email:     email,
		Password:  password,
		Role:      role,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
