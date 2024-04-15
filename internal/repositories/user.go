package repositories

// import (
// 	"log"
// 	"pressebo/api/db"
// 	"pressebo/internal/models"
// )

// type UserAttrs struct {
// 	FirstName string `json:"first_name"`
// 	LastName  string `json:"last_name"`
// 	Email     string `json:"email"`
// 	Password  string `json:"password"`
// 	OrgID     int64  `json:"org_id"`
// }

// func CreateUser(db db.DBSession, userAttrs UserAttrs) (*models.User, error) {
// 	userCollection := db.Collection("users")
// 	user := &models.User{
// 		FirstName: userAttrs.FirstName,
// 		LastName:  userAttrs.LastName,
// 		Email:     userAttrs.Email,
// 		Password:  userAttrs.Password,
// 		OrgID:     userAttrs.OrgID,
// 	}

// 	err := userCollection.InsertReturning(user)

// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}

// 	return user, nil
// }
