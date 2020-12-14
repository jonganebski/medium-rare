package photo

import (
	"bytes"
	"home/jonganebski/github/medium-rare/config"
	"image"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
)

// Repository interface allows us to access the CRUD Operations in mongo here.
type Repository interface {
	UploadImage(resizedImg *image.NRGBA, filename string) (*s3manager.UploadOutput, error)
	DeleteImage(filename string) (*s3.DeleteObjectOutput, error)
	DeleteImages(objects []*s3.ObjectIdentifier) (*s3.DeleteObjectsOutput, error)
}

type repository struct {
	Session *session.Session
}

//NewRepo is the single instance repo that is being created.
func NewRepo(session *session.Session) Repository {
	return &repository{
		Session: session,
	}
}

var bucketName string = config.Config("BUCKET_NAME")

func (r *repository) DeleteImages(objects []*s3.ObjectIdentifier) (*s3.DeleteObjectsOutput, error) {
	svc := s3.New(r.Session)
	output, err := svc.DeleteObjects(&s3.DeleteObjectsInput{Bucket: aws.String(bucketName), Delete: &s3.Delete{Objects: objects, Quiet: aws.Bool(true)}})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (r *repository) UploadImage(resizedImg *image.NRGBA, filename string) (*s3manager.UploadOutput, error) {
	uploader := s3manager.NewUploader(r.Session)

	buf := new(bytes.Buffer)
	imaging.Encode(buf, resizedImg, imaging.JPEG)
	reader := bytes.NewReader(buf.Bytes())

	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   reader,
	})

	if err != nil {
		return nil, err
	}
	return up, nil
}

func (r *repository) DeleteImage(filename string) (*s3.DeleteObjectOutput, error) {
	svc := s3.New(r.Session)
	deleteOutput, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(filename)})
	if err != nil {
		return nil, err
	}
	return deleteOutput, nil
}
