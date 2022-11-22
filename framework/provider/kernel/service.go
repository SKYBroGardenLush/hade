package kernel

import (
  "github.com/SKYBroGardenLush/skycraper/framework/gin"
  "net/http"
)

type HadeKernelService struct {
  engine *gin.Engine
}

func (s *HadeKernelService) HttpEngine() http.Handler {
  return s.engine
}

func NewHadeKernelService(params ...interface{}) (interface{}, error) {
  httpEngine := params[0].(*gin.Engine)
  return &HadeKernelService{engine: httpEngine}, nil
}
