package models

import "pressebo/api"

type Org struct {
	ID        int64   `json:"id" db:"id,omitempty"`
	Name      string  `json:"name" db:"name"`
	Subdomain string  `json:"subdomain" db:"subdomain"`
	Email     string  `json:"email" db:"email"`
	Users     []*User `json:"users"`
}

func (o *Org) PullUsers(sess api.DBSession) error {
	return sess.Collection("users").
		Find(api.Cond{"org_id": o.ID}).
		All(&o.Users)
}
