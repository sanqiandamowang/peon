package gui

import (
	// "fmt"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CaoYnag/gocui"
)

var fileList []string

func enterChooseFile(g *gocui.Gui, v *gocui.View) error {
	var err error = nil
	var l string
	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		pageError(g, "No file selected")
		return err
	}
	peonDebug("选择文件：" + l)
	err = pageFileTree(g, l)
	return err
}

// GetJSONFiles 返回指定目录下的所有 JSON 文件名
func GetFilesList(dir string) ([]string, error) {
	var jsonFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查文件扩展名是否为 .json
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			jsonFiles = append(jsonFiles, path)
		}
		return nil
	})

	return jsonFiles, err
}
func pageJsonChoose(g *gocui.Gui) error {
	var err error
	fileList, err = GetFilesList(*cmdConfig.ConfigDir)
	peonDebug("读取文件列表")
	if err != nil {
		return err
	}
	if len(fileList) == 0 {
		pageError(g, "No  file found")
		return err
	}
	updatePreviousView()
	maxX, maxY := g.Size()
	if v, err := g.SetView(fileListView, 0, 3, maxX/2, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = "JSON-List"
		for _, value := range fileList {
			fmt.Fprintln(v, value)
		}
		// cursor_len = len(fileList)-1
		if _, err := g.SetCurrentView(fileListView); err != nil {
			return err
		}
	}
	return nil
}
