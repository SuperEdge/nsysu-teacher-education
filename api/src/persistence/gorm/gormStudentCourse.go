package gorm

import (
	"errors"

	"gorm.io/gorm"
)

// StudentCourse model
type StudentCourse struct {
	gorm.Model
	StudentID uint    `gorm:"column:student_id;"`
	Student   Student `gorm:"foreignkey:StudentID;"`
	CourseID  uint    `gorm:"column:course_id;"`
	Course    Course  `gorm:"foreignkey:CourseID;"`
	Meal      string  `gorm:"column:meal;"`
	Status    string  `gorm:"column:status"`
	Review    string  `gorm:"column:review"`
	Comment   string  `gorm:"column:comment"`
}

type studentCourseDao struct {
	table        string
	Meat         string
	Vegetable    string
	StatusPass   string
	StatusFailed string
}

// StudentCourseDao student course data access object
var StudentCourseDao = &studentCourseDao{
	table:        "student_course",
	Meat:         "meat",
	Vegetable:    "vegetable",
	StatusPass:   "pass",
	StatusFailed: "failed",
}

// New a record
func (dao *studentCourseDao) New(tx *gorm.DB, studentCourse *StudentCourse) {
	err := tx.Table(dao.table).
		Create(studentCourse).Error

	if err != nil {
		panic(err)
	}
}

// GetByID get a record by id
func (dao *studentCourseDao) GetByID(tx *gorm.DB, id uint) *StudentCourse {
	result := StudentCourse{}
	err := tx.Table(dao.table).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return &result
}

// Update record
func (dao *studentCourseDao) Update(tx *gorm.DB, studentCourse *StudentCourse) {
	err := tx.Model(&studentCourse).
		Updates(map[string]interface{}{
			"StudentID": studentCourse.StudentID,
			"CourseID":  studentCourse.CourseID,
			"Meal":      studentCourse.Meal,
			"Status":    studentCourse.Status,
			"Review":    studentCourse.Review,
			"Comment":   studentCourse.Comment,
		}).Error

	if err != nil {
		panic(err)
	}
}

// Count get total count
func (dao *studentCourseDao) Count(tx *gorm.DB, funcs ...func(*gorm.DB) *gorm.DB) int64 {
	var result []StudentCourse
	count := tx.Joins("Student").
		Joins("Course").
		Table(dao.table).
		Scopes(funcs...).
		Find(&result).RowsAffected

	return count
}

// Query custom query
func (dao *studentCourseDao) Query(tx *gorm.DB, funcs ...func(*gorm.DB) *gorm.DB) *[]StudentCourse {
	var result []StudentCourse
	err := tx.Joins("Student").
		Joins("Course").
		Table(dao.table).
		Scopes(funcs...).
		Find(&result).Error

	if err != nil {
		panic(err)
	}
	return &result
}
