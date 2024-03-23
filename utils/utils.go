/**
* @author: gongquanlin
* @since: 2024/3/23
* @desc:
 */

package utils

import (
	"crypto/rand"
	"github.com/gin-contrib/multitemplate"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes)
}

func getFileList(path string, suffix string) (files []string) {
	// 遍历目录
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		// 将模板后缀的文件放到列表
		if strings.HasSuffix(path, suffix) {
			files = append(files, path)
		}
		return nil
	})
	return
}

// LoadTemplateFiles 加载模板
func LoadTemplateFiles(templateDir, suffix string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	rd, _ := ioutil.ReadDir(templateDir)
	for _, fi := range rd {
		if fi.IsDir() {
			// 如果是目录
			for _, f := range getFileList(path.Join(templateDir, fi.Name()), suffix) {
				// 添加到模板的时候，去掉跟路径
				r.AddFromFiles(f[len(templateDir)+1:], f)
			}
		} else {
			if strings.HasSuffix(fi.Name(), suffix) {
				// 如果再根目录底下的文件直接添加
				r.AddFromFiles(fi.Name(), path.Join(templateDir, fi.Name()))
			}
		}
	}

	return r
}
