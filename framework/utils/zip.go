package utils

import (
  "archive/zip"
  "fmt"
  "io"
  "os"
  "path/filepath"
  "strings"
)

//Unzip 解压缩zip文件，复制文件和目录到目标目录中
func Unzip(src string, dest string) ([]string, error) {
  var filenames []string

  //使用archive/zip读取
  r, err := zip.OpenReader(src)
  if err != nil {
    return filenames, err
  }
  defer r.Close()

  //所有内部文件读取
  for _, f := range r.File {
    //目标路径
    fpath := filepath.Join(dest, f.Name)

    if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
      return filenames, fmt.Errorf("%s : ilegal file path", fpath)
    }

    filenames = append(filenames, fpath)

    if f.FileInfo().IsDir() {
      //如果是目录，则创建目录
      os.MkdirAll(fpath, os.ModePerm)
      continue
    }

    outfile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
    if err != nil {
      return filenames, err
    }

    rc, err := f.Open()
    if err != nil {
      return filenames, err
    }

    //复制内容
    _, err = io.Copy(outfile, rc)
    outfile.Close()
    rc.Close()

    if err != nil {
      return filenames, err
    }
  }
  return filenames, nil
}
