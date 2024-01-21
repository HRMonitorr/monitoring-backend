package structure

type CommitsTotal struct {
	EmployeeName string    `json:"employee-name" bson:"employee-name"`
	Commit       []Commits `json:"commit" bson:"commit"`
	Total        int       `json:"total" bson:"total"`
}

type Commits struct {
	Author  string `bson:"author" json:"author"`
	Email   string `json:"email" bson:"email"`
	Comment string `bson:"comment" json:"comment" `
}

type Creds struct {
	Status  int         `bson:"status" json:"status"`
	Message string      `json:"message" bson:"message"`
	Data    interface{} `json:"data" bson:"data"`
}

type BodyReq struct {
	OwnerName string `json:"ownerName"`
	RepoName  string `json:"repoName"`
}
