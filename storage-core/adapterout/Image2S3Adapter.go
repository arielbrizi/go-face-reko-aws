package storageadapterout

import (
	"bytes"
	"go-intelligent-monitoring-system/domain"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//Image2S3Adapter ...
type Image2S3Adapter struct {
	bucket   string //TODO: just one?
	uploader s3manager.Uploader
}

//Save ...
func (i2s3 *Image2S3Adapter) Save(image domain.Image) error {

	upParams := &s3manager.UploadInput{
		Bucket: &image.Bucket,
		Key:    &image.Name,
		Body:   bytes.NewReader(image.Bytes),
	}

	// Perform an upload.
	_, err := i2s3.uploader.Upload(upParams)

	if err != nil {
		log.WithFields(log.Fields{"image.Bucket": image.Bucket, "image.Name": image.Name}).WithError(err).Error("Error uploading image to S3")
		return err
	}

	log.WithFields(log.Fields{"image.Bucket": image.Bucket, "image.Name": image.Name}).Info("Image correctly uploaded to S3")

	return nil
}

//GetURL from S3 to get the file
func (i2s3 *Image2S3Adapter) GetURL(image domain.Image) (string, error) {
	url := "https://" + image.Bucket + ".s3.amazonaws.com/" + image.Name
	//Todo: ver si conviene poner una URL de la app que haga de intermediario con S3 (un MS)
	return url, nil
}

//NewImage2S3Adapter initializes an Image2S3Adapter object.
func NewImage2S3Adapter() *Image2S3Adapter {

	bucket := os.Getenv("CAMARA_DOMAIN")

	// The session the S3 Uploader will use
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// S3 service client the Upload manager will use.
	s3Svc := s3.New(sess)

	// Create an uploader with S3 client and default options
	uploader := s3manager.NewUploaderWithClient(s3Svc)

	// Create the S3 Bucket
	_, err := s3Svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		log.WithFields(log.Fields{"bucket": bucket}).WithError(err).Fatal("Unable to create bucket")
	}

	err = s3Svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		log.WithFields(log.Fields{"bucket": bucket}).WithError(err).Fatal("Error occurred while waiting for bucket to be created")
	}

	log.WithFields(log.Fields{"bucket": bucket}).Info("Bucket correctly created")

	createBucketPolicy(bucket, s3Svc)

	return &Image2S3Adapter{
		bucket:   bucket,
		uploader: *uploader,
	}
}

func createBucketPolicy(bucket string, s3Svc *s3.S3) {

	policy := os.Getenv("AWS_S3_BUCKET_POLICY")

	if policy == "" {
		return
	}

	input := &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(policy),
	}

	//For make public policy:
	// - https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.PutBucketPolicy
	// - https://docs.aws.amazon.com/AmazonS3/latest/dev/example-bucket-policies.html
	/*{
		"Version": "2012-10-17", ---> don't change that
		"Statement": [
			{
				"Sid": "PublicRead",
				"Effect": "Allow",
				"Principal": "*",
				"Action": [
					"s3:GetObject",
					"s3:GetObjectVersion"
				],
				"Resource": "arn:aws:s3:::camarasilvia/*"
			}
		]
	}*/

	result, err := s3Svc.PutBucketPolicy(input)
	if err != nil {
		log.WithFields(log.Fields{"bucket": bucket}).WithError(err).Fatal("Unable to create bucket Policy")
	}

	log.WithFields(log.Fields{"bucket": bucket, "result": result}).Info("bucket Policy created")
}

///////////// For Test ////////////

//Image2S3AdapterTest ...
type Image2S3AdapterTest struct {
	bucket   string //TODO: just one?
	uploader s3manager.Uploader
}

//Save ...
func (i2s3 *Image2S3AdapterTest) Save(image domain.Image) error {
	return nil
}

//GetURL from S3 to get the file
func (i2s3 *Image2S3AdapterTest) GetURL(image domain.Image) (string, error) {
	return "", nil
}

//NewImage2S3AdapterTest initializes an NewImage2S3AdapterTest object.
func NewImage2S3AdapterTest() *Image2S3AdapterTest {

	return &Image2S3AdapterTest{}
}
