package config

import (
  "bytes"
  "errors"
  "github.com/SKYBroGardenLush/skyscraper/framework"
  "github.com/SKYBroGardenLush/skyscraper/framework/contract"
  "github.com/fsnotify/fsnotify"
  "github.com/mitchellh/mapstructure"
  "github.com/spf13/cast"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "log"
  "os"
  "path/filepath"
  "strings"
  "sync"
  "time"
)

type HadeConfig struct {
  c        framework.Container // 容器
  folder   string              //本地文件夹
  keyDelim string              // 路径分隔符默认为.

  envMaps  map[string]string      //所有环境变量
  confMaps map[string]interface{} //配置文件结构，key为文件名
  confRaws map[string][]byte      // 配置文件原始信息

  lock sync.RWMutex // 配置文件读写锁
}

//表示使用环境变量 maps替换context中的env(xxx)的环境变量
func replace(content []byte, maps map[string]string) []byte {
  if maps == nil {
    return content
  }
  //使用ReplaceAll 替换。这个性能 可能不是最优的 但是配置文件加载，频率比较低
  for key, val := range maps {
    reKey := "env(" + key + ")"
    content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
  }
  return content
}

//查找某个路径的配置项
func searchMap(source map[string]interface{}, path []string) interface{} {
  if len(path) == 0 {
    return source
  }

  //判断是否有下个路径
  next, ok := source[path[0]]
  if ok {
    if len(path) == 1 {
      return next
    }
  }

  switch next.(type) {
  case map[interface{}]interface{}:
    // 如果interface 的map,使用cast进行value转换
    return searchMap(cast.ToStringMap(next), path[1:])
  case map[string]interface{}:
    return searchMap(next.(map[string]interface{}), path[1:])
  default:
    return nil
  }
}

func (conf *HadeConfig) find(key string) interface{} {
  return searchMap(conf.confMaps, strings.Split(key, conf.keyDelim))
}

func NewHadeConfig(params ...interface{}) (interface{}, error) {
  container := params[0].(framework.Container)
  envFolder := params[1].(string)
  envMaps := params[2].(map[string]string)

  // 检查文件夹是否存在
  if _, err := os.Stat(envFolder); os.IsNotExist(err) {
    return nil, errors.New("folder " + envFolder + " not exist: " + err.Error())
  }

  // 实例化
  hadeConf := &HadeConfig{
    c:        container,
    folder:   envFolder,
    envMaps:  envMaps,
    confMaps: map[string]interface{}{},
    confRaws: map[string][]byte{},
    keyDelim: ".",
    lock:     sync.RWMutex{},
  }

  // 读取每个文件
  files, err := ioutil.ReadDir(envFolder)
  if err != nil {
    return nil, err
  }
  for _, file := range files {
    fileName := file.Name()
    err := hadeConf.loadConfigFile(envFolder, fileName)
    if err != nil {
      log.Println(err)
      continue
    }
  }

  //监控文件夹
  watch, err := fsnotify.NewWatcher()
  if err != nil {
    return nil, err
  }
  err = watch.Add(envFolder)
  if err != nil {
    return nil, err
  }
  go func() {
    defer func() {
      if err := recover(); err != nil {
        log.Println(err)
      }
    }()

    for {
      select {
      case ev := <-watch.Events:
        {
          //判断事件发生类型
          //Create 创建
          //Write 写入
          //Remove 删除
          path, _ := filepath.Abs(ev.Name)
          index := strings.LastIndex(path, string(os.PathSeparator))
          folder := path[:index]
          fileName := path[index+1:]

          if ev.Op&fsnotify.Create == fsnotify.Create {
            log.Println("创键文件:", ev.Name)
            hadeConf.loadConfigFile(folder, fileName)
          }
          if ev.Op&fsnotify.Write == fsnotify.Write {
            log.Println("创键文件:", ev.Name)
            hadeConf.loadConfigFile(folder, fileName)
          }
          if ev.Op&fsnotify.Remove == fsnotify.Remove {
            log.Println("创键文件:", ev.Name)
            hadeConf.removeConfigFile(folder, fileName)
          }

        }
      case err := <-watch.Errors:
        {
          log.Println("error :", err)
          return
        }

      }
    }
  }()

  return hadeConf, nil
}

func NewConfig(params ...interface{}) (interface{}, error) {
  if len(params) != 3 {
    return nil, errors.New("params error")
  }
  container := params[0].(framework.Container)
  envFolder := params[1].(string)
  // 检查文件夹是否存在
  if _, err := os.Stat(envFolder); os.IsNotExist(err) {
    return nil, errors.New("folder " + envFolder + " not exist: " + err.Error())
  }

  envMap := params[2].(map[string]string)

  // 实例化
  hadeConf := &HadeConfig{
    c:        container,
    folder:   envFolder,
    keyDelim: ".",
    envMaps:  envMap,
    confMaps: map[string]interface{}{}, //配置文件结构，key为文件名
    confRaws: map[string][]byte{},      // 配置文件原始信息
    lock:     sync.RWMutex{},
  }
  //读取每个文件
  files, err := ioutil.ReadDir(envFolder)
  if err != nil {
    return nil, err
  }
  for _, file := range files {
    fileName := file.Name()
    err := hadeConf.loadConfigFile(envFolder, fileName)
    if err != nil {
      log.Println(err)
      continue
    }
  }
  return hadeConf, nil
}

func (conf *HadeConfig) loadConfigFile(folder string, file string) error {
  conf.lock.Lock()
  defer conf.lock.Unlock()
  //判断文件是否以yaml或yml作为后缀
  s := strings.Split(file, ".")
  if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
    name := s[0]
    //读取文件内容
    bf, err := ioutil.ReadFile(filepath.Join(folder, file))
    if err != nil {
      return err
    }
    //直接针对问本作环境变量的替换
    bf = replace(bf, conf.envMaps)
    //解析对应文件
    confMap := map[string]interface{}{}
    if err := yaml.Unmarshal(bf, confMap); err != nil {
      return err
    }
    conf.confMaps[name] = confMap
    conf.confRaws[name] = bf

    //读取app.path中的信息，更新app对应的folder
    if name == "app" && conf.c.IsBind(contract.AppKey) {
      if p, ok := confMap["path"]; ok {
        appService := conf.c.MustMake(contract.AppKey).(contract.App)
        appService.LoadAppConfig(cast.ToStringMapString(p))
      }

    }
  }
  return nil
}

func (conf *HadeConfig) removeConfigFile(folder string, file string) error {
  conf.lock.Lock()
  defer conf.lock.Unlock()
  s := strings.Split(file, ".")
  if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
    name := s[0]
    //删除内存中对应的key
    delete(conf.confRaws, name)
    delete(conf.confRaws, name)
  }
  return nil

}

func (conf *HadeConfig) Get(key string) interface{} {
  return conf.find(key)
}

// IsExist check setting is exist
func (conf *HadeConfig) IsExist(key string) bool {
  return conf.find(key) != nil
}

// GetBool 获取bool类型配置
func (conf *HadeConfig) GetBool(key string) bool {
  return cast.ToBool(conf.find(key))
}

// GetInt 获取int类型配置
func (conf *HadeConfig) GetInt(key string) int {
  return cast.ToInt(conf.find(key))
}

// GetInt 获取int类型配置
func (conf *HadeConfig) GetInt64(key string) int64 {
  return cast.ToInt64(conf.find(key))
}

// GetFloat64 get float64
func (conf *HadeConfig) GetFloat64(key string) float64 {
  return cast.ToFloat64(conf.find(key))
}

// GetTime get time type
func (conf *HadeConfig) GetTime(key string) time.Time {
  return cast.ToTime(conf.find(key))
}

// GetString get string typen
func (conf *HadeConfig) GetString(key string) string {
  return cast.ToString(conf.find(key))
}

// GetIntSlice get int slice type
func (conf *HadeConfig) GetIntSlice(key string) []int {
  return cast.ToIntSlice(conf.find(key))
}

// GetStringSlice get string slice type
func (conf *HadeConfig) GetStringSlice(key string) []string {
  return cast.ToStringSlice(conf.find(key))
}

// GetStringMap get map which key is string, value is interface
func (conf *HadeConfig) GetStringMap(key string) map[string]interface{} {
  return cast.ToStringMap(conf.find(key))
}

// GetStringMapString get map which key is string, value is string
func (conf *HadeConfig) GetStringMapString(key string) map[string]string {
  return cast.ToStringMapString(conf.find(key))
}

// GetStringMapStringSlice get map which key is string, value is string slice
func (conf *HadeConfig) GetStringMapStringSlice(key string) map[string][]string {
  return cast.ToStringMapStringSlice(conf.find(key))
}

// Load a config to a struct, val should be an pointer
func (conf *HadeConfig) Load(key string, val interface{}) error {
  decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
    TagName: "yaml",
    Result:  val,
  })
  if err != nil {
    return err
  }

  return decoder.Decode(conf.find(key))
}
