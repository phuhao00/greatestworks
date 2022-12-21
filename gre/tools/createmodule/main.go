package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var moduleName string
var (
	src string
	dst string
)

func init() {
	flag.StringVar(&moduleName, "module", "task", "please")
	flag.StringVar(&src, "src", "../../template", "please")
	flag.StringVar(&dst, "dst", "../../../business/module", "please")
}

func main() {
	flag.Parse() // 解析参数
	fmt.Printf("%s\n", moduleName)
	src, _ = filepath.Abs(src)
	dst, _ = filepath.Abs(dst)
	dst = dst + "\\" + moduleName
	err := Copy(src, dst)
	fmt.Println(err, src, dst)
}
func Copy(from, to string) error {
	var err error

	f, err := os.Stat(from)
	if err != nil {
		return err
	}

	fn := func(fromFile string) error {
		//复制文件的路径
		rel, err := filepath.Rel(from, fromFile)
		if err != nil {
			return err
		}
		toFile := filepath.Join(to, rel)

		//创建复制文件目录
		if err = os.MkdirAll(filepath.Dir(toFile), 0777); err != nil {
			return err
		}

		//读取源文件
		file, err := os.Open(fromFile)
		if err != nil {
			return err
		}

		defer file.Close()
		bufReader := bufio.NewReader(file)
		// 创建复制文件用于保存
		out, err := os.Create(toFile)
		if err != nil {
			return err
		}

		defer out.Close()
		// 然后将文件流和文件流对接起来
		_, err = io.Copy(out, bufReader)
		return err
	}

	//转绝对路径
	pwd, _ := os.Getwd()
	if !filepath.IsAbs(from) {
		from = filepath.Join(pwd, from)
	}
	if !filepath.IsAbs(to) {
		to = filepath.Join(pwd, to)
	}

	//复制
	if f.IsDir() {
		return filepath.WalkDir(from, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				return fn(path)
			} else {
				if err = os.MkdirAll(path, 0777); err != nil {
					return err
				}
			}
			return err
		})
	} else {
		return fn(from)
	}
}
