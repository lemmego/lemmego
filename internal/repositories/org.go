package repositories

// import (
// 	"log"
// 	"pressebo/api/db"
// 	"pressebo/internal/models"
// )

// type OrgAttrs struct {
// 	Name      string `json:"name"`
// 	Subdomain string `json:"subdomain"`
// 	Email     string `json:"email"`
// }

// func CreateOrg(db db.DBSession, orgAttrs OrgAttrs) (*models.Org, error) {
// 	orgCollection := db.Collection("orgs")

// 	org := &models.Org{
// 		Name:      orgAttrs.Name,
// 		Subdomain: orgAttrs.Subdomain,
// 		Email:     orgAttrs.Email,
// 	}

// 	if err := orgCollection.InsertReturning(org); err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}

// 	return org, nil
// }
