package model

type Demo struct {
	ID          uint64 `gorm:"column:id"`
	UserID      string `gorm:"column:user_id"`
	CreatedAt   int    `gorm:"column:created_at"`
	UpdatedAt   int    `gorm:"column:updated_at"`
	StatusStage int    `gorm:"column:status_stage"`
}

func (d Demo) TableName() string {
	return "root_demo"
}
