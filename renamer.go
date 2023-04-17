package main

import (
	"bufio"
	"fmt"
	"github.com/saintfish/chardet"
	"os"
	"path/filepath"
	"strings"
)

func getEncoding(file string) string {
	f, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	detector := chardet.NewTextDetector()
	charset, err := detector.DetectBest(f)
	if err != nil {
		panic(err)
	}
	if charset.Charset == "UTF-8" {
		if isUTF8BOM3(f) {
			return "UTF-8-BOM"
		}
		return "UTF-8"
	}
	return charset.Charset
}

func isUTF8BOM3(buf []byte) bool {
	if len(buf) < 3 {
		return false
	}
	return buf[0] == 0xEF && buf[1] == 0xBB && buf[2] == 0xBF
}

func delUTF8BOM3(buf []byte) []byte {
	if len(buf) < 3 {
		return buf
	}
	if buf[0] == 0xEF && buf[1] == 0xBB && buf[2] == 0xBF {
		return buf[3:]
	}
	return buf
}

func isUTF8Compatible(file string) bool {
	if getEncoding(file) == "UTF-8" || getEncoding(file) == "ISO-8859-1" || getEncoding(file) == "ASCII" || getEncoding(file) == "UTF-8-BOM" {
		return true
	}
	return false
}

func readLines(file string) []string {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func main() {
	fmt.Println(getEncoding("name_list.txt"))
	if isUTF8Compatible("name_list.txt") { // 所有兼容 UTF-8 的编码
		// 判断是否为 UTF-8-BOM 编码并转换
		if getEncoding("name_list.txt") == "UTF-8-BOM" {
			fmt.Println("检测到 UTF-8-BOM 编码，尝试转换为 UTF-8")
			f, err := os.ReadFile("name_list.txt")
			if err != nil {
				panic(err)
			}
			f = delUTF8BOM3(f)
			err = os.WriteFile("name_list.txt", f, 0644)
			if err != nil {
				panic(err)
			}
			if isUTF8Compatible("name_list.txt") {
				fmt.Println("转换成功")
			} else {
				panic("转换失败")
			}
		}
	} else {
		panic("不支持的编码格式")
	}
	var fileList []string = os.Args[1:]
	var names []string = readLines("name_list.txt")
	if len(fileList) != len(names) {
		panic("定义文件名条目数与实际文件数不符！")
	}
	for i := 0; i < len(names); i++ {
		var line = names[i]
		fmt.Println(line)
		var file = fileList[i]
		var path, fileName = filepath.Split(file)
		err := os.Chdir(path)
		if err != nil {
			panic(err)
		}
		var ext = filepath.Ext(fileName)
		var pName = strings.TrimSuffix(fileName, ext)
		if pName == line {
			fmt.Println("文件名与定义条目相同，跳过")
			continue
		}
		var dst = line + ext
		err = os.Rename(file, dst)
		if err != nil {
			panic(err)
		}
		fmt.Println(dst)
	}
}
