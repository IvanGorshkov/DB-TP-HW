package models

type Thread struct {
	Id		int32  `json:"id"`
	Title	string `json:"title"`
	Author	string `json:"author"`
	Forum	string `json:"forum"`
	Message string `json:"message"`
	Votes	int32  `json:"votes"`
	Created	string `json:"created"`
	Slug	string `json:"slug"`
}


type Vote struct {
	Nickname	string `json:"nickname"`
	Votes		int  `json:"voice"`
	Thread 		int  `json:"id"`
}