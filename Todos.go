package main

import (
	"fmt"
	"net/http"
	"strconv"
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
	todoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed`
	}

	fmtTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
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
func (todoModel) TableName() string {
	return "todos"
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
	db.AutoMigrate(&todoModel{})
}
func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1/todo")
	{
		v1.POST("/", add)                // 添加新条目
		v1.GET("/", all)                 // 查询所有条目
		v1.GET("/:id", take)             // 获取单个条目
		v1.PUT("status/:id", update)     // 更新状态
		v1.DELETE("/:id", del)           // 更新title
		v1.PUT("title/:id", updateTitle) // 更新状态
	}
	r.Run(":9089")
}

func add(c *gin.Context) {
	t_com := c
	fmt.Println(t_com)
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{Title: c.PostForm("title"), Completed: completed}
	db.Save(&todo)
	c.JSON(http.StatusOK, gin.H{
		"status":     JSON_SUCCESS,
		"message":    "创建成功",
		"resourceId": todo.ID,
	})
}

func all(c *gin.Context) {
	var todos []todoModel
	var _todos []fmtTodo
	db.Find(&todos)

	// 没有数据
	if len(todos) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  JSON_ERROR,
			"message": "没有数据",
		})
		return
	}

	// 格式化
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, fmtTodo{
			ID:        item.ID,
			Title:     item.Title,
			Completed: completed,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "ok",
		"data":    _todos,
	})
}
func take(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  JSON_ERROR,
			"message": "条目不存在",
		})
		return
	}
	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}

	_todo := fmtTodo{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: completed,
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "ok",
		"data":    _todo,
	})
}
func update(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  JSON_ERROR,
			"message": "条目不存在",
		})
		return
	}

	// db.Model(&todo).Where("id=?", todo.ID).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Where("id=?", todo.ID).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "更新成功",
	})
}
func del(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  JSON_ERROR,
			"message": "条目不存在",
		})
		return
	}
	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "删除成功!",
	})
}
func updateTitle(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  JSON_ERROR,
			"message": "条目不存在",
		})
		return
	}

	db.Model(&todo).Where("id=?", todo.ID).Update("title", c.PostForm("title"))

	c.JSON(http.StatusOK, gin.H{
		"status":  JSON_SUCCESS,
		"message": "更新成功",
	})
}
