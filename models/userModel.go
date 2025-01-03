package models

import(
	"time"
"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	ID				primitive.ObjectID		`bson:"_id"`
	FirstName		*string					`json:"firstName" validate:"required,min=2,max=100"`
	LastName		*string					`json:"lastName" validate:"required,min=2,max=100"`
	Password		*string					`json:"password" validate:"required,min=6"`
	Email			*string					`json:"email" validate:"email,required"`
	Phone			*string					`json:"phone" validate:"required"`
	Token			*string					`json:"token"`
	UserType		*string					`json:"userType" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken	*string					`json:"refreshToken"`
	CreatedAt		time.Time				`json:"createdAt"`
	UpdatedAt		time.Time				`json:"updatedAt"`
	UserId			string					`json:"userId"`
}