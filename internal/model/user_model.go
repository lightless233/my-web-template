package model

import "my-web-template/internal/entity/vo"

type AppUserModel struct {
	BaseModel `xorm:"extends"`
	Username  string `xorm:"VARCHAR(255) NOT NULL UNIQUE"`
	Password  string `xorm:"VARCHAR(255) NOT NULL"`
	Email     string `xorm:"VARCHAR(255) NOT NULL"`
	State     uint8  `xorm:"TINYINT NOTNULL DEFAULT 0"`
}

func (u *AppUserModel) TableName() string {
	return "app_user"
}

func (u *AppUserModel) ToVO() *vo.UserVO {
	return &vo.UserVO{
		UserId:   u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}
