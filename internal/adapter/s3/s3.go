package s3

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
)

type S3Object struct {
	Key           string
	ContentType   string
	ContentLength int64
	Metadata      map[string]*string
}

type S3Config struct {
	AccessKeyId     string
	SecretAccessKey string
	Endpoint        string
	Region          string
}

type S3 struct {
	S3Config
	s3 *s3.S3
}

func NewS3(config S3Config) *S3 {
	s3 := S3{
		S3Config: config,
	}
	s3.createSession(config)
	return &s3
}

func (s *S3) createSession(config S3Config) {
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(config.AccessKeyId, config.SecretAccessKey, ""),
		Endpoint:    aws.String(config.Endpoint),
		Region:      aws.String(config.Region),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		mindiaerr.ExitErrorf("Unable create a new session, %v", err)
	}
	s.s3 = s3.New(newSession)
}

type ListObjectsParams struct {
	Bucket string
	Prefix string
}

func (s *S3) ListObjects(p ListObjectsParams) ([]S3Object, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(p.Bucket),
		Prefix: aws.String(p.Prefix),
	}
	output, err := s.s3.ListObjectsV2(input)
	if err != nil {
		return nil, err
	}
	var objs []S3Object
	for _, obj := range output.Contents {
		objs = append(objs, S3Object{
			Key:           *obj.Key,
			ContentType:   "",
			ContentLength: *obj.Size,
			Metadata:      nil,
		})
	}
	return objs, nil
}

type GetObjectParams struct {
	Bucket string
	Key    string
}

func (s *S3) DownloadObject(p GetObjectParams) (io.ReadCloser, *string, *int64, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(p.Key),
	}
	output, err := s.s3.GetObject(input)
	if err != nil {
		return nil, nil, nil, err
	}
	return output.Body, output.ContentType, output.ContentLength, nil
}

func (s *S3) GetObject(p GetObjectParams) (*S3Object, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(p.Key),
	}
	output, err := s.s3.GetObject(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && (awsErr.Code() == "Forbidden" || awsErr.Code() == "NotFound") {
			return nil, mindiaerr.New(mindiaerr.ErrCodeMediaNotFound)
		}
		return nil, &mindiaerr.Error{ErrCode: mindiaerr.ErrCodeInternal, Msg: err}
	}
	return &S3Object{
		Key:           p.Key,
		ContentType:   *output.ContentType,
		ContentLength: *output.ContentLength,
		Metadata:      output.Metadata,
	}, nil
}

type PutObjectParams struct {
	Bucket        string
	Key           string
	Body          io.Reader
	ContentType   string
	ContentLength int64
	Metadata      map[string]*string
}

func (s *S3) PutObject(p PutObjectParams) error {
	partSize := int64(5 * 1024 * 1024) // 5MB
	partNumber := int64(1)

	createResp, err := s.s3.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket:      &p.Bucket,
		Key:         &p.Key,
		ContentType: &p.ContentType,
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return err
	}
	uploadID := createResp.UploadId
	completedParts := []*s3.CompletedPart{}

	for {
		buf := make([]byte, partSize)
		n, err := p.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if n < int(partSize) {
			buf = buf[:n]
		}
		partResp, err := s.s3.UploadPart(&s3.UploadPartInput{
			Bucket:     &p.Bucket,
			Key:        &p.Key,
			PartNumber: &partNumber,
			UploadId:   uploadID,
			Body:       bytes.NewReader(buf),
		})
		if err != nil {
			return err
		}
		completedParts = append(completedParts, &s3.CompletedPart{
			ETag:       partResp.ETag,
			PartNumber: &partNumber,
		})
	}
	_, err = s.s3.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   &p.Bucket,
		Key:      &p.Key,
		UploadId: uploadID,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	return err
}

type MoveObjectParams struct {
	Bucket string
	SrcKey string
	DstKey string
}

func (s *S3) RenameObject(p MoveObjectParams) error {
	input := CopyObjectParams{
		Bucket: p.Bucket,
		SrcKey: p.SrcKey,
		DstKey: p.DstKey,
	}
	err := s.CopyObject(input)
	if err != nil {
		return err
	}
	return s.DeleteObject(DeleteObjectParams{
		Bucket: p.Bucket,
		Key:    p.SrcKey,
	})
}

type CopyObjectParams struct {
	Bucket string
	SrcKey string
	DstKey string
}

func (s *S3) CopyObject(p CopyObjectParams) error {
	input := &s3.CopyObjectInput{
		Bucket:     aws.String(p.Bucket),
		CopySource: aws.String(fmt.Sprintf("%v/%v", p.Bucket, p.SrcKey)),
		Key:        aws.String(p.DstKey),
	}
	_, err := s.s3.CopyObject(input)
	return err
}

type DeleteObjectParams struct {
	Bucket string
	Key    string
}

func (s *S3) DeleteObject(p DeleteObjectParams) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(p.Key),
	}
	_, err := s.s3.DeleteObject(input)
	return err
}

type SpaceUsageParams struct {
	Bucket string
	Prefix string
}
