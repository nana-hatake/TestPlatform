package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/gin-gonic/gin"
)

const (
	JSON_SUCCESS int = 1
	JSON_ERROR   int = 0
)

type (
	userModel struct {
		gorm.Model
		Username string `json:username`
		Password string `json:password`
		// Title     string `json:"title"`
		// Completed int    `json:"completed`
	}

	fmtuser struct {
		ID       uint   `json:"id"`
		Username string `json:username`
		Password string `json:password`
	}
)

// gorm.Model 定义
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// 指定表名
func (userModel) TableName() string {
	return "user"
}

var db *gorm.DB

// 初始化
func init() {
	var err error
	// var constr string
	var constr = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "Zyr365daysbxsb$", "localhost", 3306, "TestPlatform")
	db, err = gorm.Open("mysql", constr)
	if err != nil {
		panic("数据库连接失败")
	}
	fmt.Print("Connect Success")
	db.AutoMigrate(&userModel{})
}
func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1/user")
	{
		// v1.POST("/", add)                // 添加新条目
		v1.GET("/", all)           // 查询所有条目
		v1.GET("/:username", take) // 获取单个条目
		// v1.PUT("status/:id", update)     // 更新状态
		// v1.DELETE("/:id", del)           // 更新title
		// v1.PUT("title/:id", updateTitle) // 更新状态
	}
	r.Run(":9090")
}

// func add(c *gin.Context) {
// 	t_com := c
// 	fmt.Println(t_com)
// 	completed, _ := strconv.Atoi(c.PostForm("completed"))
// 	todo := todoModel{Title: c.PostForm("title"), Completed: completed}
// 	db.Save(&todo)
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":     JSON_SUCCESS,
// 		"message":    "创建成功",
// 		"resourceId": todo.ID,
// 	})
// }

func all(c *gin.Context) {
	var user []userModel
	var _user []fmtuser
	db.Find(&user)

	// 没有数据
	if len(user) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  JSON_ERROR,
			"message": "没有数据",
			"stage":   "error page",
		})
		return
	}
	// fmt.Print(user)
	// 格式化
	for _, item := range user {

		_user = append(_user, fmtuser{
			ID:       item.ID,
			Username: item.Username,
			Password: item.Password,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "ok",
		"data":    _user,
		"stage":   "next_page",
	})
}

func take(c *gin.Context) {
	var user userModel
	username := c.Param("username")
	// username := c.Params("username")

	db.First(&user, username)
	// if user.Username == "''" {
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"status":  JSON_ERROR,
	// 		"message": "条目不存在",
	// 	})
	// 	return
	// }
	fmt.Print("%s", user)

	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "ok",
		"data":    user,
		"stage":   "next_page",
	})

}

// func update(c *gin.Context) {
// 	var todo todoModel
// 	todoID := c.Param("id")
// 	db.First(&todo, todoID)
// 	if todo.ID == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"status":  JSON_ERROR,
// 			"message": "条目不存在",
// 		})
// 		return
// 	}

// 	// db.Model(&todo).Where("id=?", todo.ID).Update("title", c.PostForm("title"))
// 	completed, _ := strconv.Atoi(c.PostForm("completed"))
// 	db.Model(&todo).Where("id=?", todo.ID).Update("completed", completed)
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  JSON_SUCCESS,
// 		"message": "更新成功",
// 	})
// }
// func del(c *gin.Context) {
// 	var todo todoModel
// 	todoID := c.Param("id")
// 	db.First(&todo, todoID)
// 	if todo.ID == 0 {
// 		c.JSON(http.StatusOK, gin.H{
// 			"status":  JSON_ERROR,
// 			"message": "条目不存在",
// 		})
// 		return
// 	}
// 	db.Delete(&todo)
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  JSON_SUCCESS,
// 		"message": "删除成功!",
// 	})
// }
// func updateTitle(c *gin.Context) {
// 	var todo todoModel
// 	todoID := c.Param("id")
// 	db.First(&todo, todoID)
// 	if todo.ID == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"status":  JSON_ERROR,
// 			"message": "条目不存在",
// 		})
// 		return
// 	}

// 	db.Model(&todo).Where("id=?", todo.ID).Update("title", c.PostForm("title"))

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  JSON_SUCCESS,
// 		"message": "更新成功",
// 	})
// }
