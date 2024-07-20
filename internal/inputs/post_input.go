package inputs

import (
    "lemmego/api/vee"
)

type PostInput struct {
    Title string `json:"title" in:"form=title"`
    PostDescription string `json:"post_description" in:"form=post_description"`
}

func (i *PostInput) Validate() error {
	v := vee.New()
    v.Required("title", i.Title)
    v.Required("post_description", i.PostDescription)
	return v.Errors
}
