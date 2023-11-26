package s3

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
	"github.com/jeremybastin1207/mindia-core/internal/media"
)

type FileStorageConfig struct {
	S3     *S3
	Bucket string
}

type FileStorage struct {
	s3     *S3
	bucket string
}

func NewFileStorage(c FileStorageConfig) *FileStorage {
	return &FileStorage{
		s3:     c.S3,
		bucket: c.Bucket,
	}
}

func (s *FileStorage) Upload(in media.UploadInput) error {
	return s.s3.PutObject(PutObjectParams{
		Bucket:        s.bucket,
		Key:           in.Path.ToString(),
		Body:          in.Body,
		ContentType:   in.ContentType,
		ContentLength: int64(in.ContentLength),
	})
}

func (s *FileStorage) Download(p media.Path) (*media.DownloadResult, error) {
	body, contentType, contentLength, err := s.s3.DownloadObject(GetObjectParams{
		Bucket: s.bucket,
		Key:    p.ToString(),
	})
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == s3.ErrCodeNoSuchKey {
			return nil, &mindiaerr.Error{ErrCode: mindiaerr.ErrCodeMediaNotFound}
		}
		return nil, err
	}
	return &media.DownloadResult{
		Path:          p,
		Body:          body,
		ContentType:   *contentType,
		ContentLength: *contentLength,
	}, nil
}

func (s *FileStorage) DownloadMultiple(paths []media.Path) ([]*media.DownloadResult, error) {
	var downloadResponses []*media.DownloadResult

	for _, p := range paths {
		downloadResponse, err := s.Download(p)
		if err != nil {
			return nil, err
		}
		downloadResponses = append(downloadResponses, downloadResponse)
	}

	return downloadResponses, nil
}

func (s *FileStorage) Get(p media.Path) (*media.FileInfo, error) {
	obj, err := s.s3.GetObject(GetObjectParams{
		Bucket: s.bucket,
		Key:    p.ToString(),
	})
	if err != nil {
		return nil, err
	}
	return &media.FileInfo{
		Path:          p,
		ContentType:   obj.ContentType,
		ContentLength: int(obj.ContentLength),
	}, nil
}

func (s *FileStorage) GetMultiple(p media.Path) ([]media.FileInfo, error) {
	var (
		medias = []media.FileInfo{}
		dir    = strings.TrimPrefix(p.Dir(), "/")
		prefix = p.Uuid()
	)

	objs, err := s.s3.ListObjects(ListObjectsParams{
		Bucket: s.bucket,
		Prefix: dir + "/" + prefix,
	})
	if err != nil {
		return nil, err
	}

	for _, obj := range objs {
		medias = append(medias, media.FileInfo{
			Path:          media.NewPath("/" + obj.Key),
			ContentType:   obj.ContentType,
			ContentLength: int(obj.ContentLength),
		})
	}
	return medias, nil
}

func (s *FileStorage) Move(src, dst media.Path) error {
	return s.s3.RenameObject(MoveObjectParams{
		Bucket: s.bucket,
		SrcKey: src.ToString()[1:],
		DstKey: dst.ToString()[1:] + "/" + src.Filename(),
	})
}

func (s *FileStorage) Copy(src, dst media.Path) error {
	return s.s3.CopyObject(CopyObjectParams{
		Bucket: s.bucket,
		SrcKey: src.ToString()[1:],
		DstKey: dst.ToString()[1:],
	})
}

func (s *FileStorage) Delete(p media.Path) error {
	return s.s3.DeleteObject(DeleteObjectParams{
		Bucket: s.bucket,
		Key:    p.ToString(),
	})
}

func (s *FileStorage) SpaceUsage() (int64, error) {
	objs, err := s.s3.ListObjects(ListObjectsParams{
		Bucket: s.bucket,
	})
	if err != nil {
		return 0, err
	}
	var totalSpace int64
	for _, obj := range objs {
		totalSpace += obj.ContentLength
	}
	return totalSpace, err
}
