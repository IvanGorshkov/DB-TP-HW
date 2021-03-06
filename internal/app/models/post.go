package models

//easyjson:json
type Posts []Post


//easyjson:json
type Post struct {
	ID       int    `json:"id"`
	Parent   int    `json:"parent"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	IsEdited bool   `json:"isEdited"`
	Forum    string `json:"forum"`
	Thread   int    `json:"thread"`
	Created  string `json:"created"`
}

//easyjson:json
type PostFull struct {
	Post	Post `json:"post"`
	Author	*User `json:"author, omitempty"`
	Thread	*Thread `json:"thread, omitempty"`
	Forum	*Forum `json:"forum, omitempty"`
}