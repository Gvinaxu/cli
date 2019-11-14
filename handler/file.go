package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Gvinaxu/cli/util"
	"github.com/valyala/fasthttp"
)

const (
	treeUrl = "http://127.0.0.1:18080/api/1/file/list"
	mvUrl   = "http://127.0.0.1:18080/api/1/file/mv"
	cpUrl   = "http://127.0.0.1:18080/api/1/file/cp"
	rmUrl   = "http://127.0.0.1:18080/api/1/file/delete"

	mkdirUrl  = "http://127.0.0.1:18080/api/1/file/mkdir"
	searchUrl = "http://127.0.0.1:18080/api/1/search/search"
)

var (
	dirCount  = 0
	fileCount = 0
)

type FileReq struct {
}

func (f *FileReq) Tree(path string) ([]*File, error) {
	args := &fasthttp.Args{}
	args.Add("path", path)
	head := map[string]interface{}{
		"access_token": token.AccessToken,
	}

	code, body, err := util.DoTimeout(args, "POST", treeUrl, head)
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("response code is %d", code))
	}
	type R struct {
		Code int     `json:"code"`
		Data []*File `json:"data"`
	}
	r := &R{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, err
	}
	if path == "" {
		path = "/"
	}
	fmt.Println()
	fmt.Println(path)
	f.printDirectory(r.Data, 0)
	fmt.Println()
	fmt.Printf("%d directories, %d files \n", dirCount, fileCount)
	return r.Data, nil
}

func (f *FileReq) printDirectory(files []*File, depth int) {
	for _, file := range files {
		if file.Dir {
			dirCount++
			f.printListing(file.Name, depth, true)
			f.printDirectory(file.Child, depth+1)
		} else {
			fileCount++
			f.printListing(file.Name, depth+1, false)
		}
	}

}

func (f *FileReq) printListing(entry string, depth int, dir bool) {
	indent := strings.Repeat("|   ", depth)
	fmt.Printf("%s|-- %s\n", indent, entry)
	// output color
}

func (f *FileReq) Stat(path string) (*File, error) {
	args := &fasthttp.Args{}
	args.Add("path", path)
	head := map[string]interface{}{
		"access_token": token.AccessToken,
	}

	code, body, err := util.DoTimeout(args, "POST", treeUrl, head)
	if err != nil {
		return nil, err
	}
	if code == http.StatusNotFound {
		errMsg := fmt.Sprintf("path %s file not found ", path)
		fmt.Println(errMsg)
		return nil, errors.New(errMsg)
	}
	if code != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("response code is %d", code))
	}
	type R struct {
		Code int     `json:"code"`
		Data []*File `json:"data"`
	}
	r := &R{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, err
	}
	file := r.Data[0]

	tm := time.Unix(file.Time, 0)
	fmt.Println()
	fmt.Printf("File: %s \n", file.Name)

	fmt.Printf("Size: %d	Blocks:%d  Time:%s	", file.Size, len(file.Blocks), tm.Format("2006-01-02 15:04:05"))
	if file.Dir {
		fmt.Printf("Regular:dir \n")
	} else {
		fmt.Printf("Regular:file \n")
	}
	fmt.Println()
	return file, nil
}

func (f *FileReq) Mkdir(path string) error {
	args := &fasthttp.Args{}
	args.Add("path", path)
	head := map[string]interface{}{
		"access_token": token.AccessToken,
	}

	code, body, err := util.DoTimeout(args, "POST", mkdirUrl, head)
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return errors.New(fmt.Sprintf("response code is %d", code))
	}
	type R struct {
		Code int `json:"code"`
	}
	r := &R{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return err
	}
	if r.Code != 0 {
		return errors.New("mkdir fail")
	}
	return nil
}

func (f *FileReq) Cp(oldPath, newPath string) error {
	args := &fasthttp.Args{}
	args.Add("path", oldPath)
	args.Add("tagPath", newPath)
	head := map[string]interface{}{
		"access_token": token.AccessToken,
	}

	code, body, err := util.DoTimeout(args, "POST", cpUrl, head)
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return errors.New(fmt.Sprintf("response code is %d", code))
	}
	type R struct {
		Code int `json:"code"`
	}
	r := &R{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return err
	}
	if r.Code != 0 {
		return errors.New("copy fail")
	}
	return nil
}

func (f *FileReq) Mv(oldPath, newPath string) error {
	args := &fasthttp.Args{}
	args.Add("path", oldPath)
	args.Add("tagPath", newPath)
	head := map[string]interface{}{
		"access_token": token.AccessToken,
	}

	code, body, err := util.DoTimeout(args, "POST", mvUrl, head)
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return errors.New(fmt.Sprintf("response code is %d", code))
	}
	type R struct {
		Code int `json:"code"`
	}
	r := &R{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return err
	}
	if r.Code != 0 {
		return errors.New("move fail")
	}
	return nil
}

func (f *FileReq) Rm(path string) error {
	args := &fasthttp.Args{}
	args.Add("path", path)
	head := map[string]interface{}{
		"access_token": token.AccessToken,
	}

	code, body, err := util.DoTimeout(args, "POST", rmUrl, head)
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return errors.New(fmt.Sprintf("response code is %d", code))
	}
	type R struct {
		Code int `json:"code"`
	}
	r := &R{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return err
	}
	if r.Code != 0 {
		return errors.New("remove fail")
	}
	return nil
}

func (f *FileReq) Search(search, path string) ([]*File, error) {
	args := &fasthttp.Args{}
	args.Add("search", search)
	if path != "" {
		args.Add("path", path)
	}
	head := map[string]interface{}{
		"access_token": token.AccessToken,
	}

	code, body, err := util.DoTimeout(args, "POST", searchUrl, head)
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("response code is %d", code))
	}
	type R struct {
		Code int     `json:"code"`
		Data []*File `json:"data"`
	}
	r := &R{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, err
	}
	fmt.Println()
	if len(r.Data) == 0 {
		fmt.Println("result 0 row")
	} else {
		fmt.Println("name			path")
	}

	for _, v := range r.Data {
		fmt.Printf("%s		%s \n", v.Name, v.Path)
	}
	return r.Data, nil

}

type File struct {
	FecEnable bool         `json:"fec_enable"`
	Blocks    []*FileBlock `json:"blocks"`
	Path      string       `json:"path"`
	Name      string       `json:"name"`
	Size      int64        `json:"size"`
	Time      int64        `json:"time"`
	Dir       bool         `json:"dir"`
	Expires   int64        `json:"expires"`
	Child     []*File      `json:"child"`
}

type FileBlock struct {
	Key string `json:"key"`
	Cid string `json:"cid"`
}
