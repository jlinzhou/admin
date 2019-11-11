package models

type User struct {
	Model
	Memo     string `gorm:"column:memo;size:64;" json:"memo" form:"memo"`                                                          //备注
	UserName string `gorm:"column:username;size:32;unique_index:uk_admins_user_name;not null;" json:"username"` // 用户名
	RealName string `gorm:"column:realname;size:32;" json:"realname"`                                           // 真实姓名
	Password string `gorm:"column:password;type:char(32);not null;" json:"password" form:"password"`                               // 密码(sha1(md5(明文))加密)
	Email    string `gorm:"column:email;size:64;" json:"email" form:"email"`                                                       // 邮箱
	Phone    string `gorm:"column:phone;type:char(20);" json:"phone" form:"phone"`                                                 // 手机号
	Status   uint8  `gorm:"column:status;type:tinyint(1);not null;" json:"status" form:"status"`                          // 状态(1:正常 2:未激活 3:暂停使用)
}


// 表名
func (User) TableName() string {
	return "User"
}
func (User) Delete(adminsids []uint64) error {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id in (?)", adminsids).Delete(&User{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("adminsId in (?)", adminsids).Delete(&UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
