package v115open

import (
	"Q115-STRM/internal/helpers"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	driver115 "github.com/SheltonZhu/115driver/pkg/driver"
)

// NewCookieClient 创建并校验115 Cookie客户端
func NewCookieClient(cookie string) (*driver115.Pan115Client, *driver115.UserInfo, error) {
	client, _, err := newCookieClient(cookie)
	if err != nil {
		return nil, nil, err
	}
	if err := client.CookieCheck(); err != nil {
		return nil, nil, fmt.Errorf("Cookie已失效: %w", err)
	}
	user, err := client.GetUser()
	if err != nil {
		return nil, nil, fmt.Errorf("获取Cookie账号信息失败: %w", err)
	}
	client.UserID = user.UserID
	return client, user, nil
}

// NewCookieClientUnchecked 创建Cookie客户端，不发起网络请求；用于运行时回退。
func NewCookieClientUnchecked(cookie string) (*driver115.Pan115Client, error) {
	client, _, err := newCookieClient(cookie)
	return client, err
}

func newCookieClient(cookie string) (*driver115.Pan115Client, *driver115.Credential, error) {
	credential := &driver115.Credential{}
	if err := credential.FromCookie(strings.TrimSpace(cookie)); err != nil {
		return nil, nil, fmt.Errorf("Cookie格式无效: %w", err)
	}
	client := driver115.New(driver115.UA("Mozilla/5.0 115Browser/35.6.0.3"))
	client.ImportCredential(credential)
	return client, credential, nil
}

func cookieFileToOpen(file *driver115.File) File {
	fileType := TypeFile
	if file.IsDirectory {
		fileType = TypeDir
	}
	return File{
		FileId:       file.FileID,
		Aid:          "1",
		Pid:          file.ParentID,
		FileCategory: fileType,
		FileName:     file.Name,
		PickCode:     file.PickCode,
		Utime:        file.UpdateTime.Unix(),
		Ptime:        file.CreateTime.Unix(),
		Sha1:         file.Sha1,
		FileSize:     file.Size,
		Fta:          "1",
		Ico:          filepath.Ext(file.Name),
		Thumbnail:    file.ThumbURL,
	}
}

func cookieStatToDetail(fileId string, file *driver115.File, stat *driver115.FileStatInfo) *FileDetail {
	detail := &FileDetail{
		FileSize:     strconv.FormatInt(file.Size, 10),
		FileSizeByte: file.Size,
		Ptime:        strconv.FormatInt(file.CreateTime.Unix(), 10),
		Utime:        strconv.FormatInt(file.UpdateTime.Unix(), 10),
		FileName:     file.Name,
		PickCode:     file.PickCode,
		Sha1:         file.Sha1,
		FileId:       file.FileID,
		FileCategory: TypeFile,
	}
	if file.IsDirectory {
		detail.FileCategory = TypeDir
		detail.Count = json.Number(strconv.Itoa(stat.FileCount))
		detail.FolderCount = json.Number(strconv.Itoa(stat.DirCount))
	}
	pathParts := make([]string, 0, len(stat.Parents))
	for _, parent := range stat.Parents {
		detail.Paths = append(detail.Paths, FileDetailPath{FileId: parent.ID, Name: parent.Name})
		if parent.ID != "0" {
			pathParts = append(pathParts, parent.Name)
		}
	}
	detail.Path = filepath.ToSlash(filepath.Join(pathParts...))
	if detail.FileId == "" {
		detail.FileId = fileId
	}
	return detail
}

func (c *OpenClient) getFsListCookie(ctx context.Context, fileId string, showCur bool, onlyDir bool, showDir bool, offset int, limit int) (*FileListResp, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	client := c.getCookieClient()
	if client == nil {
		return nil, fmt.Errorf("未绑定115 Cookie")
	}
	if fileId == "" {
		fileId = "0"
	}
	req := client.NewRequest()
	if showCur {
		req.SetQueryParam("cur", "1")
	}
	if onlyDir {
		req.SetQueryParam("stdir", "1")
	}
	resp, err := driver115.GetFiles(req, fileId,
		driver115.WithOffset(int64(offset)),
		driver115.WithLimit(int64(limit)),
		driver115.WithShowDirEnable(showDir),
	)
	if err != nil {
		return nil, err
	}
	result := &FileListResp{Count: resp.Count, SysCount: resp.Count}
	result.State = true
	result.Limit = jsonNumber(resp.Limit)
	result.Offset = jsonNumber(resp.Offset)
	result.Data = make([]File, 0, len(resp.Files))
	for i := range resp.Files {
		file := (&driver115.File{}).From(&resp.Files[i])
		result.Data = append(result.Data, cookieFileToOpen(file))
	}
	if fileId != "0" {
		if stat, statErr := client.Stat(fileId); statErr == nil {
			pathParts := make([]string, 0, len(stat.Parents)+1)
			for _, parent := range stat.Parents {
				result.Path = append(result.Path, FileParentPath{Name: parent.Name, FileId: jsonNumberString(parent.ID)})
				if parent.ID != "0" {
					pathParts = append(pathParts, parent.Name)
				}
			}
			result.Path = append(result.Path, FileParentPath{Name: stat.Name, FileId: jsonNumberString(fileId)})
			pathParts = append(pathParts, stat.Name)
			result.PathStr = filepath.ToSlash(filepath.Join(pathParts...))
		}
	}
	helpers.V115Log.Infof("OpenAPI限流，已通过115 Cookie驱动查询目录: %s", fileId)
	return result, nil
}

func jsonNumber(value int) json.Number {
	return json.Number(strconv.Itoa(value))
}

func jsonNumberString(value string) json.Number {
	return json.Number(value)
}

func (c *OpenClient) getFsDetailByCidCookie(ctx context.Context, fileId string) (*FileDetail, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	client := c.getCookieClient()
	if client == nil {
		return nil, fmt.Errorf("未绑定115 Cookie")
	}
	file, err := client.GetFile(fileId)
	if err != nil {
		return nil, err
	}
	stat, err := client.Stat(fileId)
	if err != nil {
		return nil, err
	}
	return cookieStatToDetail(fileId, file, stat), nil
}

func (c *OpenClient) getFsDetailByPathCookie(ctx context.Context, path string) (*FileDetail, error) {
	client := c.getCookieClient()
	if client == nil {
		return nil, fmt.Errorf("未绑定115 Cookie")
	}
	resp, err := client.DirName2CID(path)
	if err != nil {
		return nil, err
	}
	return c.getFsDetailByCidCookie(ctx, fmt.Sprint(resp.CategoryID))
}

func (c *OpenClient) uploadCookie(ctx context.Context, filePath, parentFileId string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	client := c.getCookieClient()
	if client == nil {
		return "", fmt.Errorf("未绑定115 Cookie")
	}
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return "", err
	}
	if err = client.UploadFastOrByMultipart(parentFileId, info.Name(), info.Size(), file); err != nil {
		return "", err
	}
	files, err := client.ListPage(parentFileId, 0, driver115.MaxDirPageLimit)
	if err != nil {
		return "", err
	}
	for _, uploaded := range *files {
		if uploaded.Name == info.Name() && !uploaded.IsDirectory {
			return uploaded.FileID, nil
		}
	}
	return "", fmt.Errorf("上传成功但未查询到文件ID")
}
