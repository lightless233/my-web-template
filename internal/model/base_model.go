package model

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type BaseModel struct {
	ID          uint64 `xorm:"UNSIGNED BIGINT NOTNULL PK AUTOINCR"`
	CreatedTime int64  `xorm:"UNSIGNED BIGINT NOTNULL"`
	UpdatedTime int64  `xorm:"UNSIGNED BIGINT NOTNULL"`
	Deleted     bool   `xorm:"BOOL NOTNULL DEFAULT false"`
}

// BeforeInsert 插入之前
func (b *BaseModel) BeforeInsert() {
	ts := time.Now().UnixMilli()
	b.CreatedTime = ts
	b.UpdatedTime = ts
}

// BeforeUpdate 更新之前
func (b *BaseModel) BeforeUpdate() {
	b.UpdatedTime = time.Now().UnixMilli()
}
