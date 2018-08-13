package ini

import (
	"strings"
	"sync"

	"github.com/go-ini/ini"
	"github.com/gwaylib/errors"
)

// 用于省略前缀长路径写法
// 例如,以下可用于多语言处理：
// ini := NewIni(ini.RootDir()+"/app.default)
// lang := ".zh_cn"
// cfg := ini.Get(lang)
// cfg.String("msg", "1001")
type Ini struct {
	rootPath string
}

type File struct {
	*ini.File
}

func (f *File) String(section, key string) string {
	result := f.Section(section).Key(key).String()
	if len(result) == 0 {
		panic(errors.ErrNoData.As(section, key))
	}
	return result
}
func (f *File) Float64(section, key string) float64 {
	result, err := f.Section(section).Key(key).Float64()
	if err != nil {
		panic(errors.As(err, section, key))
	}
	return result
}
func (f *File) Int64(section, key string) int64 {
	result, err := f.Section(section).Key(key).Int64()
	if err != nil {
		panic(errors.As(err, section, key))
	}
	return result
}
func (f *File) Uint64(section, key string) uint64 {
	result, err := f.Section(section).Key(key).Uint64()
	if err != nil {
		panic(errors.As(err, section, key))
	}
	return result
}
func (f *File) Bool(section, key string) bool {
	result, err := f.Section(section).Key(key).Bool()
	if err != nil {
		panic(errors.As(err, section, key))
	}
	return result
}

var cache = sync.Map{}

func GetFile(fileName string) (*File, error) {
	f, ok := cache.Load(fileName)
	if ok {
		return f.(*File), nil
	}
	file, err := ini.Load(fileName)
	if err != nil {
		if strings.Index(err.Error(), "no such file or directory") > -1 {
			return nil, errors.ErrNoData.As(err, fileName)
		}
		return nil, err
	}
	ff := &File{file}
	cache.Store(fileName, ff)
	return ff, nil
}

func NewIni(rootPath string) *Ini {
	return &Ini{rootPath}
}

func (ini *Ini) Get(fileName string) *File {
	f, err := GetFile(ini.rootPath + fileName)
	if err != nil {
		panic(errors.As(err, ini.rootPath+fileName))
	}
	return f
}

func (ini *Ini) GetDefault(fileName, defFileName string) *File {
	f, err := GetFile(ini.rootPath + fileName)
	if err != nil {
		if !errors.ErrNoData.Equal(err) {
			panic(errors.As(err, ini.rootPath+fileName))
		}
		return ini.Get(defFileName)
	}
	return f
}
