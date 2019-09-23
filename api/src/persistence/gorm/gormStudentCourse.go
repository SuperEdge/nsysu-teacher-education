package gorm

import (
	"github.com/jinzhu/gorm"
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

// StudentCourseDao user data acces object
var StudentCourseDao = &studentCourseDao{
	table:        "student_course",
	Meat:         "meate",
	Vegetable:    "vegetable",
	StatusPass:   "pass",
	StatusFailed: "failed",
}

// New a record
func (dao *studentCourseDao) New(tx *gorm.DB, user *StudentCourse) {
	err := tx.Table(dao.table).
		Create(user).Error

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
		Scan(&result).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return &result
}

// Query custom query
func (dao *studentCourseDao) Query(tx *gorm.DB, funcs ...func(*gorm.DB) *gorm.DB) *[]StudentCourse {
	var result []StudentCourse
	err := tx.Preload("Student").
		Preload("Course").
		Table(dao.table).
		Select("*").
		Scopes(funcs...).
		Find(&result).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return &result
}