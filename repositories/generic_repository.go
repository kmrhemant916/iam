package repositories

import (
	"gorm.io/gorm"
)

type GenericRepository[T any] interface {
	Create(entity *T) error
	FindOne(entity *T, sqlQuery string, sqlQueryParams ...interface{}) (*T, error)
}

type genericRepository[T any] struct {
	db *gorm.DB
}

func NewGenericRepository[T any](db *gorm.DB) *genericRepository[T] {
	return &genericRepository[T]{
		db: db,
	}
}

func (u *genericRepository[T]) Create(entity *T) error {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	if err := tx.Create(entity).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (u *genericRepository[T]) FindOne(entity *T, sqlQuery string, sqlQueryParams ...interface{}) (*T, error) {
	res := u.db.Raw(sqlQuery, sqlQueryParams...).First(entity)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return entity, nil
}

// func GetUsers(c *gin.Context) {
// 	db := c.MustGet("db").(*gorm.DB)
// 	var users []models.User
// 	db.Find(&users)
// 	response := []UserInputResponse{}
// 	for _, user := range users {
// 		response = append(response, UserInputResponse{user.Email, user.ID})
// 	}
// 	c.JSON(http.StatusOK, gin.H{"data": response})
// }

// func UpdateUser(c *gin.Context) {
// 	db := c.MustGet("db").(*gorm.DB)
// 	var user models.User
// 	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
// 		return
// 	}
// 	var input UserInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	update := models.User{ID: user.ID, Email: input.Email, Password: input.Password}
// 	db.Save(&update)
// 	c.JSON(http.StatusOK, gin.H{"data": "Record updated"})
// }

// func DeleteUser(c *gin.Context) {
// 	db := c.MustGet("db").(*gorm.DB)
// 	var user models.User
// 	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
// 		return
// 	}
// 	db.Delete(&user)
// 	c.JSON(http.StatusOK, gin.H{"data": true})
// }