package utils

import (
  "io"
  "io/ioutil"
  "net/http"
  "os"
)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
  _, err := os.Stat(path) //os.Stat获取文件信息
  if err != nil {
    if os.IsExist(err) {
      return true
    }
    return false
  }
  return true
}

func IsHiddenDirectory(path string) bool {
  if path[:1] == "." || path[:2] == ".." {
    return true
  }
  return false
}

func (c *Collection) Contains(value interface{}) bool {
  return false
}

func SubDir(folder string) ([]string, error) {
  subs, err := ioutil.ReadDir(folder)
  if err != nil {
    return nil, err
  }

  ret := []string{}
  for _, sub := range subs {
    if sub.IsDir() {
      ret = append(ret, sub.Name())
    }
  }
  return ret, nil
}

func SubFile(folder string) ([]string, error) {
  subs, err := ioutil.ReadDir(folder)
  if err != nil {
    return nil, err
  }

  ret := []string{}
  for _, sub := range subs {
    if !sub.IsDir() {
      ret = append(ret, sub.Name())
    }
  }
  return ret, nil
}

// DownloadFile 下载url中的内容保存到本地的filepath中
func DownloadFile(filepath string, url string) error {
  //获取
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  //创建文件目标
  out, err := os.Create(filepath)
  if err != nil {
    return err
  }
  defer out.Close()

  //拷贝内容
  _, err = io.Copy(out, resp.Body)
  return err

}
