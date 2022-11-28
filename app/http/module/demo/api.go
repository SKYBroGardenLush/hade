package demo

import (
  "fmt"
  demoService "github.com/SKYBroGardenLush/skyscraper/app/provider/demo"
  "github.com/SKYBroGardenLush/skyscraper/framework/contract"
  "github.com/SKYBroGardenLush/skyscraper/framework/gin"
  "github.com/SKYBroGardenLush/skyscraper/framework/provider/orm"
)

type DemoApi struct {
  service *Service
}

func Register(r *gin.Engine) error {
  api := NewDemoApi()
  r.Bind(&demoService.DemoProvider{})

  r.GET("/demo/demo", api.Demo)
  r.GET("/demo/demo2", api.Demo2)
  r.POST("/demo/demo_post", api.DemoPost)
  return nil
}

func NewDemoApi() *DemoApi {
  service := NewService()
  return &DemoApi{service: service}
}

// Demo godoc
// @Summary 获取所有用户
// @Description 获取所有用户
// @Produce  json
// @Tags demo
// @Success 200 array []UserDTO
// @Router /demo/demo [get]
func (api *DemoApi) Demo(c *gin.Context) {
  gormService := c.MustMake(contract.ORMKey).(contract.ORMService)
  db, err := gormService.GetDB(orm.WithConfigPath("database.default"))
  if err != nil {
    fmt.Println(err.Error())
  }
  db.WithContext(c)
  // 将User模型创建到数据库中
  err = db.AutoMigrate(&Product2{})
  var product Product2
  db.First(&product)
  res := map[string]interface{}{
    "name":  product.Name,
    "price": product.Price,
  }

  c.JSON(200, res)
}

// Demo godoc
// @Summary 获取所有学生
// @Description 获取所有学生
// @Produce  json
// @Tags demo
// @Success 200 array []UserDTO
// @Router /demo/demo2 [get]
func (api *DemoApi) Demo2(c *gin.Context) {
  demoProvider := c.MustMake(demoService.DemoKey).(demoService.IService)
  students := demoProvider.GetAllStudent()
  usersDTO := StudentsToUserDTOs(students)
  c.JSON(200, usersDTO)
}

func (api *DemoApi) DemoPost(c *gin.Context) {
  type Foo struct {
    Name string
  }
  foo := &Foo{}
  err := c.BindJSON(&foo)
  if err != nil {
    c.AbortWithError(500, err)
  }
  c.JSON(200, nil)
}
