package controllers

import (
	"log"
	"time"
	"backend/database"
	"backend/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func DeleteExpiredFiles() ([]models.File, error) {
	var expiredFiles []models.File
	now := time.Now()

	db := database.GetDB()

	err := db.Where("expire_at IS NOT NULL AND expire_at < ?", now).Find(&expiredFiles).Error
	if err != nil {
		return nil, err
	}

	if len(expiredFiles) == 0 {
		return expiredFiles, nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("your-region"), 
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		return nil, err
	}
	svc := s3.New(sess)

	for _, file := range expiredFiles {
		// Delete the file from S3
		_, err = svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String("your-bucket-name"), 
			Key:    aws.String(file.S3URL),
		})
		if err != nil {
			log.Printf("Failed to delete file %s from S3: %v", file.FileName, err)
			continue
		}

		err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
			Bucket: aws.String("filemger"),
			Key:    aws.String(file.S3URL),
		})
		if err != nil {
			log.Printf("Error while waiting for file %s to be deleted from S3: %v", file.FileName, err)
			continue
		}

		err = db.Delete(&file).Error
		if err != nil {
			log.Printf("Failed to delete file metadata for %s from the database: %v", file.FileName, err)
			continue
		}
	}

	return expiredFiles, nil
}
