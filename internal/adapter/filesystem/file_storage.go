package filesystem

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/pkg/path"
)

type FileStorageConfig struct {
	MountDir string
}

type FileStorage struct {
	MountDir string
}

func NewFileStorage(c FileStorageConfig) *FileStorage {
	s := &FileStorage{
		MountDir: c.MountDir,
	}
	s.createMountPathIfNotExists()
	return s
}

func (s *FileStorage) createMountPathIfNotExists() {
	s.createPathIfNotExists("")
}

func (s *FileStorage) createPathIfNotExists(path string) {
	err := os.MkdirAll(s.MountDir+path, 0777)
	if err != nil {
		mindiaerr.ExitErrorf("unable to create dir, %v", err)
	}
}

func detectContentType(f *os.File) (string, error) {
	// to sniff the content type only the first
	// 512 bytes are used.
	buf := make([]byte, 512)

	_, err := f.Read(buf)

	if err != nil {
		return "", err
	}

	// the function that actually does the trick
	contentType := http.DetectContentType(buf)

	return contentType, nil
}

func (s *FileStorage) Upload(in media.UploadInput) error {
	s.createPathIfNotExists(in.Path.Dir())

	file, err := os.Create(path.JoinPath(s.MountDir, in.Path.ToString()))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, in.Body)
	return err
}

func (s *FileStorage) Download(p media.Path) (*media.DownloadResult, error) {
	s.createPathIfNotExists(p.Dir())

	file, err := os.Open(path.JoinPath(s.MountDir, p.ToString()))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, &mindiaerr.Error{ErrCode: mindiaerr.ErrCodeMediaNotFound}
		}
		return nil, err
	}
	contentType, err := detectContentType(file)
	if err != nil {
		return nil, err
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &media.DownloadResult{
		Path:          p,
		Body:          io.NopCloser(file),
		ContentType:   contentType,
		ContentLength: stat.Size(),
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
	f, err := os.Open(path.JoinPath(s.MountDir, p.ToString()))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contentType, _ := detectContentType(f)
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return &media.FileInfo{
		Path:          p,
		ContentType:   contentType,
		ContentLength: int(stat.Size()),
	}, nil
}

func (s *FileStorage) GetMultiple(p media.Path) ([]media.FileInfo, error) {
	var (
		files  = []media.FileInfo{}
		dir    = p.Dir()
		prefix = p.Uuid()
	)

	entries, err := os.ReadDir(path.JoinPath(s.MountDir, dir))
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if prefix != "" && !strings.HasPrefix(entry.Name(), prefix) {
			continue
		}
		file, err := s.Get(media.NewPath(path.JoinPath(dir, entry.Name())))
		if err != nil {
			continue
		}
		files = append(files, *file)
	}
	return files, nil
}

func (s *FileStorage) Move(src, dst media.Path) error {
	err := os.MkdirAll(path.JoinPath(s.MountDir, dst.Path), 0777)
	if err != nil {
		return err
	}
	return os.Rename(path.JoinPath(s.MountDir, src.ToString()), path.JoinPath(s.MountDir, dst.ToString()+"/"+src.Filename()))
}

func (s *FileStorage) Copy(src, dst media.Path) error {
	source, err := os.Open(path.JoinPath(s.MountDir, src.ToString()+"/"+src.Filename()))
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(path.JoinPath(s.MountDir, dst.ToString()))
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func (s *FileStorage) Delete(p media.Path) error {
	return os.RemoveAll(path.JoinPath(s.MountDir, p.ToString()))
}

func (s *FileStorage) SpaceUsage() (int64, error) {
	sizes := make(chan int64)
	readSize := func(path string, file os.FileInfo, err error) error {
		if err != nil || file == nil {
			return nil // Ignore errors
		}
		if !file.IsDir() {

			sizes <- file.Size()
		}
		return nil
	}

	go func() {
		_ = filepath.Walk(s.MountDir, readSize)
		close(sizes)
	}()

	size := int64(0)
	for s := range sizes {
		size += s
	}

	return size, nil
}
