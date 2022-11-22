package env

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/SKYBroGardenLush/skycraper/framework/contract"
	"io"
	"os"
	"path"
	"strings"
)

type HadeEnv struct {
	folder string
	maps   map[string]string
}

// NewHadeEnv 有一个参数，.env文件所在的目录
// example: NewHadeEnv("/envfolder/") 会读取文件: /envfolder/.env
// .env的文件格式 FOO_ENV=BAR

func NewHadeEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewHadeEnv param error")
	}

	//读取folder文件
	folder := params[0].(string)
	//实例化
	hadeEnv := &HadeEnv{
		folder: folder,
		maps:   map[string]string{"APP_ENV": contract.EnvDevelopment},
	}

	//解析folder/.env文件
	file := path.Join(folder, ".env")
	// 读取.env文件, 不管任意失败，都不影响后续
	//打开.env
	fi, err := os.Open(file)
	if err == nil {
		defer fi.Close()
		//读取文件
		br := bufio.NewReader(fi)
		for {

			//按行进行读取
			line, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			//按照等号解析
			s := bytes.SplitN(line, []byte{'='}, 2)
			if len(s) < 2 {
				continue
			}
			// 保存map
			key := string(s[0])
			val := string(s[1])
			hadeEnv.maps[key] = val

		}
	}

	//获取当前环境变量并覆盖.env文件下的变量
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := pair[0]
		val := pair[1]
		hadeEnv.maps[key] = val
	}

	return hadeEnv, nil

}

func (s *HadeEnv) AppEnv() string {
	return s.Get("APP_ENV")
}

func (s *HadeEnv) Get(key string) string {
	if v, ok := s.maps[key]; ok {
		return v
	}
	return ""
}

func (s *HadeEnv) IsExit(key string) bool {
	if _, ok := s.maps[key]; ok {
		return true
	}
	return false
}

func (s *HadeEnv) All() map[string]string {
	return s.maps
}
