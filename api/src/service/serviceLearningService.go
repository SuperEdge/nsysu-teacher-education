package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/nsysu/teacher-education/src/assembler"
	"github.com/nsysu/teacher-education/src/errors"
	"github.com/nsysu/teacher-education/src/persistence/gorm"
	"github.com/nsysu/teacher-education/src/specification"
	"github.com/nsysu/teacher-education/src/utils/logger"
	"github.com/nsysu/teacher-education/src/utils/typecast"
)

// CreateServieLearning create service-learning
func CreateServieLearning(serviceType, content, session string, hours uint, start, end time.Time) (result interface{}, e *errors.Error) {
	tx := gorm.DB()

	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
			e = errors.UnexpectedError()
		}
	}()

	serviceLearning := &gorm.ServiceLearning{
		Type:    serviceType,
		Content: content,
		Session: session,
		Hours:   hours,
		Start:   start,
		End:     end,
	}
	gorm.ServiceLearningDao.New(tx, serviceLearning)

	return "success", nil
}

// GetServiceLearningList get service-learning list
func GetServiceLearningList(account, start, length string) (result map[string]interface{}, e *errors.Error) {
	tx := gorm.DB()

	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
			e = errors.UnexpectedError()
		}
	}()

	var serviceLearnings *[]gorm.ServiceLearning
	if operator := gorm.AdminDao.GetByAccount(tx, account); operator != nil {
		serviceLearnings = gorm.ServiceLearningDao.Query(
			tx,
			specification.PaginationSpecification(typecast.StringToInt(start), typecast.StringToInt(length)),
			specification.OrderSpecification("start", specification.OrderDirectionDESC),
			specification.IsNullSpecification("deleted_at"),
		)
	} else {
		serviceLearnings = gorm.ServiceLearningDao.Query(
			tx,
			specification.PaginationSpecification(typecast.StringToInt(start), typecast.StringToInt(length)),
			specification.BiggerSpecification("start", time.Now().String()),
			specification.OrderSpecification("start", specification.OrderDirectionASC),
			specification.IsNullSpecification("deleted_at"),
		)
	}

	total := gorm.ServiceLearningDao.Count(
		tx,
		specification.IsNullSpecification("deleted_at"),
	)

	result = map[string]interface{}{
		"list":            assembler.ServiceLearningDTO(serviceLearnings),
		"recordsTotal":    total,
		"recordsFiltered": total,
	}

	return
}

// SingUpServiceLearning sudent sign up service-learning
func SingUpServiceLearning(account, serviceLearningID string) (result interface{}, e *errors.Error) {
	tx := gorm.DB()

	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
			e = errors.UnexpectedError()
		}
	}()

	student := gorm.StudentDao.GetByAccount(tx, account)
	serviceLearning := gorm.ServiceLearningDao.Query(
		tx,
		specification.IDSpecification(serviceLearningID),
		specification.IsNullSpecification("deleted_at"),
		specification.BiggerSpecification("End", time.Now().String()),
	)

	if len(*serviceLearning) == 0 {
		return nil, errors.NotFoundError("service-learning ID " + serviceLearningID)
	}

	studentServiceLearning := &gorm.StudentServiceLearning{
		StudentID:         student.ID,
		ServiceLearningID: typecast.StringToUint(serviceLearningID),
	}

	gorm.StudentServiceLearningDao.New(tx, studentServiceLearning)

	return "success", nil
}

// GetSutdentServiceLearningList get the list of student service-learning
func GetSutdentServiceLearningList(account, start, length string) (result map[string]interface{}, e *errors.Error) {
	tx := gorm.DB()

	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
			e = errors.UnexpectedError()
		}
	}()

	var studentServiceLearnings *[]gorm.StudentServiceLearning
	if operator := gorm.AdminDao.GetByAccount(tx, account); operator != nil {
		studentServiceLearnings = gorm.StudentServiceLearningDao.Query(
			tx,
			specification.PaginationSpecification(typecast.StringToInt(start), typecast.StringToInt(length)),
			specification.OrderSpecification("`student_service_learning`."+specification.IDColumn, specification.OrderDirectionDESC),
			specification.IsNullSpecification("`student_service_learning`.deleted_at"),
		)
	} else {
		student := gorm.StudentDao.GetByAccount(tx, account)
		studentServiceLearnings = gorm.StudentServiceLearningDao.Query(
			tx,
			specification.PaginationSpecification(typecast.StringToInt(start), typecast.StringToInt(length)),
			specification.OrderSpecification("`student_service_learning`."+specification.IDColumn, specification.OrderDirectionDESC),
			specification.IsNullSpecification("deleted_at"),
			specification.StudentSpecification(student.ID),
		)
	}

	total := gorm.StudentCourseDao.Count(
		tx,
		specification.IsNullSpecification("deleted_at"),
	)

	result = map[string]interface{}{
		"list":            assembler.StudentServiceLearningsDTO(studentServiceLearnings),
		"recordsTotal":    total,
		"recordsFiltered": total,
	}

	return
}

// UpdateServiceLearning update service-learning
func UpdateServiceLearning(reference, review multipart.File, operator, StudentServiceLearningID, referenceFileName, reviewFileName string) (result string, e *errors.Error) {
	tx := gorm.DB().Begin()

	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
			e = errors.UnexpectedError()
		}
	}()

	StudentServiceLearning := gorm.StudentServiceLearningDao.GetByID(tx, typecast.StringToUint(StudentServiceLearningID))

	if reference != nil {
		file, err := os.OpenFile(
			fmt.Sprintf("./assets/service-learning/%s-%s", operator, referenceFileName),
			os.O_WRONLY|os.O_CREATE,
			0666,
		)
		if err != nil {
			panic(err)
		}
		io.Copy(file, reference)
		defer file.Close()

		StudentServiceLearning.Reference = referenceFileName
		gorm.StudentServiceLearningDao.Update(tx, StudentServiceLearning)
	}

	if review != nil {
		file, err := os.OpenFile(
			fmt.Sprintf("./assets/service-learning/%s-%s", operator, reviewFileName),
			os.O_WRONLY|os.O_CREATE,
			0666,
		)
		if err != nil {
			panic(err)
		}
		io.Copy(file, review)
		defer file.Close()

		StudentServiceLearning.Review = reviewFileName
		gorm.StudentServiceLearningDao.Update(tx, StudentServiceLearning)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	return "success", nil
}

// UpdateStudentServiceLearningStatus update student-service-learning status
func UpdateStudentServiceLearningStatus(StudentServiceLearningID, Status string) (result string, e *errors.Error) {
	tx := gorm.DB()

	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
			e = errors.UnexpectedError()
		}
	}()

	StudentServiceLearning := gorm.StudentServiceLearningDao.GetByID(tx, typecast.StringToUint(StudentServiceLearningID))
	StudentServiceLearning.Status = Status
	gorm.StudentServiceLearningDao.Update(tx, StudentServiceLearning)

	return "success", nil
}

// GetStudentServiceLearningFile get student-service-learning refernce or review file
func GetStudentServiceLearningFile(operator, StudentServiceLearningID, file string) (result map[string]string, e *errors.Error) {
	tx := gorm.DB()

	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
			e = errors.UnexpectedError()
		}
	}()

	StudentServiceLearning := gorm.StudentServiceLearningDao.GetByID(tx, typecast.StringToUint(StudentServiceLearningID))
	if StudentServiceLearning == nil {
		return nil, errors.NotFoundError(file)
	}

	var (
		filePath string
		fileName string
	)

	if file == "reference" {
		fileName = StudentServiceLearning.Reference
		filePath = fmt.Sprintf("./assets/service-learning/%s-%s", operator, StudentServiceLearning.Reference)
	} else {
		fileName = StudentServiceLearning.Review
		filePath = fmt.Sprintf("./assets/service-learning/%s-%s", operator, StudentServiceLearning.Review)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.NotFoundError(fileName)
	}

	return map[string]string{
		"Path": filePath,
		"Name": fileName,
	}, nil
}