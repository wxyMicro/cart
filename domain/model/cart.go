package model

type Cart struct {
	ID        int64 `gorm:"primaryKey;not null;autoIncrement"`
	ProductID int64 `gorm:"not null" json:"product_id"`
	Num       int64 `gorm:"not null" json:"num"`
	SizeID    int64 `gorm:"not null" json:"size_id"`
	UserID    int64 `gorm:"not null" json:"user_id"`
}
