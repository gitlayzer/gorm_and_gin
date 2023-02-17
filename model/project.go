package model

type Project struct {
	ID    int     `gorm:"primary_key" json:"id"`
	Name  string  `gorm:"not null" json:"name"`
	Desc  string  `json:"desc"`
	Users []*User `gorm:"many2many:project_users"`
}

func (*Project) TableName() string {
	return "project"
}
