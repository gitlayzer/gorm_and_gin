package model

type User struct {
	ID       int        `gorm:"primary_key" json:"id"`
	Username string     `gorm:"not null" json:"username"`
	Password string     `gorm:"not null" json:"password"`
	Token    string     `gorm:"not null" json:"token"`
	Projects []*Project `gorm:"many2many:project_users"`
}

func (*User) TableName() string {
	return "user"
}
