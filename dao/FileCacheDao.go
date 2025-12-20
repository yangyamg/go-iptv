package dao

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type FileCache struct {
	Dir          string // 缓存目录
	ExpireAtZero bool   // 是否每天 0 点过期
}

var Cache *FileCache

// 创建缓存目录
func NewFileCache(dir string, expireAtZero bool) (*FileCache, error) {
	os.RemoveAll("/config/cache")
	os.RemoveAll(dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &FileCache{
		Dir:          dir,
		ExpireAtZero: expireAtZero,
	}, nil
}

// 判断是否今天 0 点之后
func expiredAtMidnight(modTime time.Time) bool {
	now := time.Now()
	// 今天 0 点
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return modTime.Before(midnight)
}

// 保存缓存（字节）
func (fc *FileCache) Set(key string, data []byte) error {
	path := filepath.Join(fc.Dir, key)
	return os.WriteFile(path, data, 0644)
}

// 读取缓存（字节）
func (fc *FileCache) Get(key string) ([]byte, error) {
	path := filepath.Join(fc.Dir, key)

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if fc.ExpireAtZero && expiredAtMidnight(info.ModTime()) {
		os.Remove(path)
		return nil, os.ErrNotExist
	}

	return os.ReadFile(path)
}

func (fc *FileCache) GetNotExpired(key string) ([]byte, error) {
	path := filepath.Join(fc.Dir, key)

	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path)
}

// 判断缓存是否存在
func (fc *FileCache) Exists(key string) bool {
	path := filepath.Join(fc.Dir, key)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if fc.ExpireAtZero && expiredAtMidnight(info.ModTime()) {
		os.Remove(path)
		return false
	}
	return true
}

// 判断缓存是否存在
func (fc *FileCache) ChannelExists(key string) bool {
	path := filepath.Join(fc.Dir, key)
	_, err := os.Stat(path)
	return err == nil
}

// 删除缓存
func (fc *FileCache) Delete(pattern string) error {
	fullPattern := filepath.Join(fc.Dir, pattern)

	matches, err := filepath.Glob(fullPattern)
	if err != nil {
		return err
	}

	for _, path := range matches {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

// 清空所有缓存
func (fc *FileCache) Clear() error {
	files, err := os.ReadDir(fc.Dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		_ = os.Remove(filepath.Join(fc.Dir, f.Name()))
	}
	return nil
}

//
// ====== JSON 封装 ======
//

// 保存任意结构为 JSON
func (fc *FileCache) SetJSON(key string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return fc.Set(key, data)
}

// 读取 JSON 并解析到 v
func (fc *FileCache) GetJSON(key string, v interface{}) error {
	data, err := fc.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// 保存任意结构体
func (fc *FileCache) SetStruct(key string, v interface{}) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return err
	}
	return fc.Set(key, buf.Bytes())
}

// 读取结构体
func (fc *FileCache) GetStruct(key string, v interface{}) error {
	data, err := fc.Get(key)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(bytes.NewReader(data))
	return dec.Decode(v)
}
