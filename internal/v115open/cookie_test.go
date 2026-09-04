package v115open

import (
	"errors"
	"testing"
	"time"

	driver115 "github.com/SheltonZhu/115driver/pkg/driver"
)

func TestCookieFileToOpen(t *testing.T) {
	file := &driver115.File{
		IsDirectory: false,
		FileID:      "123",
		ParentID:    "10",
		Name:        "电影.mkv",
		Size:        1024,
		PickCode:    "pick-code",
		Sha1:        "SHA1",
		CreateTime:  time.Unix(100, 0),
		UpdateTime:  time.Unix(200, 0),
		ThumbURL:    "https://example.com/thumb.jpg",
	}

	got := cookieFileToOpen(file)
	if got.FileId != file.FileID || got.Pid != file.ParentID || got.FileCategory != TypeFile {
		t.Fatalf("Cookie文件转换结果错误: %+v", got)
	}
	if got.Ptime != 100 || got.Utime != 200 || got.FileSize != 1024 {
		t.Fatalf("Cookie文件时间或大小转换错误: %+v", got)
	}
}

func TestCookieStatToDetail(t *testing.T) {
	file := &driver115.File{IsDirectory: true, FileID: "20", Name: "电影", CreateTime: time.Unix(100, 0), UpdateTime: time.Unix(200, 0)}
	stat := &driver115.FileStatInfo{
		IsDirectory: true,
		FileCount:   8,
		DirCount:    2,
		Parents: []*driver115.DirInfo{
			{ID: "0", Name: "根目录"},
			{ID: "10", Name: "媒体"},
		},
	}

	got := cookieStatToDetail(file.FileID, file, stat)
	if got.FileCategory != TypeDir || got.Path != "媒体" || len(got.Paths) != 2 {
		t.Fatalf("Cookie目录详情转换结果错误: %+v", got)
	}
	if got.Count.String() != "8" || got.FolderCount.String() != "2" {
		t.Fatalf("Cookie目录统计转换结果错误: %+v", got)
	}
}

func TestShouldFallbackOnlyForThrottle(t *testing.T) {
	client := &OpenClient{cookieClient: driver115.New()}
	if !client.shouldFallback(ErrOpenAPIThrottled) {
		t.Fatal("OpenAPI限流时应启用Cookie回退")
	}
	if client.shouldFallback(errors.New("普通错误")) {
		t.Fatal("非限流错误不应启用Cookie回退")
	}
}
