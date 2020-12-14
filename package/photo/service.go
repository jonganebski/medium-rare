package photo

import (
	"image"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//Service is an interface from which our api module can access our repository of all our models
type Service interface {
	UploadImageToS3(resizedImg *image.NRGBA, filename string) (*s3manager.UploadOutput, error)
	DeleteImageOfS3(filename string) (*s3.DeleteObjectOutput, error)
	DeleteImagesOfS3(objects []*s3.ObjectIdentifier) (*s3.DeleteObjectsOutput, error)
}

type service struct {
	repository Repository
}

//NewService is used to create a single instance of the service
func NewService(r Repository) Service {
	return &service{repository: r}
}

func (s *service) UploadImageToS3(resizedImg *image.NRGBA, filename string) (*s3manager.UploadOutput, error) {
	return s.repository.UploadImage(resizedImg, filename)
}

func (s *service) DeleteImageOfS3(filename string) (*s3.DeleteObjectOutput, error) {
	return s.repository.DeleteImage(filename)
}

func (s *service) DeleteImagesOfS3(objects []*s3.ObjectIdentifier) (*s3.DeleteObjectsOutput, error) {
	return s.repository.DeleteImages(objects)
}
