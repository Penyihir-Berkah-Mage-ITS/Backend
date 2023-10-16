package model

type UserLikeComment struct {
	User      User    `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    string  `gorm:"null" json:"user_id"`
	Comment   Comment `gorm:"ForeignKey:CommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CommentID string  `gorm:"null" json:"comment_id"`
}
