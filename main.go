package main

import (
	"GO-specimen/controllers"
	"GO-specimen/models"
	"fmt"
	"html/template"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

func main() {
	db, err := gorm.Open(sqlite.Open("specimen.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Specimen{}, &models.Photo{}, &models.User{}, &models.Project{})

		r := gin.Default()

	// Load HTML templates manually
	tmpl := template.New("")
	var errParse error
	tmpl, errParse = tmpl.ParseGlob("templates/*.html")
	if errParse != nil {
		panic(fmt.Sprintf("Failed to parse templates: %v", errParse))
	}
	r.SetHTMLTemplate(tmpl)

	// Static files
	r.Static("/uploads", "./uploads")
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	

	r.GET("/", controllers.GetIndexData)
	r.GET("/specimens", controllers.GetSpecimens)
	r.GET("/specimens/new", func(c *gin.Context) {
		db := c.MustGet("db").(*gorm.DB)
		users := controllers.GetUsersList(db)
		projects := controllers.GetProjectsList(db)
		c.HTML(200, "new.html", gin.H{
			"Today": time.Now(),
			"Users": users,
			"Projects": projects,
		})
	})
	r.POST("/specimens", controllers.CreateSpecimen)
	r.POST("/specimens/delete", controllers.DeleteSpecimens)
	r.GET("/specimens/:id", controllers.GetSpecimen)
	r.GET("/specimens/:id/edit", controllers.EditSpecimen)
	r.POST("/specimens/:id", func(c *gin.Context) {
		method := c.PostForm("_method")
		if method == "PUT" {
			controllers.UpdateSpecimen(c)
			return
		}
		if method == "DELETE" {
			controllers.DeleteSpecimen(c)
			return
		}
		c.AbortWithStatus(http.StatusBadRequest)
	})

	// 設定ページ
	r.GET("/settings", func(c *gin.Context) {
		c.HTML(http.StatusOK, "settings_index.html", nil)
	})

	// ユーザー管理ルート
	r.GET("/settings/users", controllers.GetUsers)
	r.POST("/settings/users", controllers.CreateUser)
	r.POST("/settings/users/:id", func(c *gin.Context) {
		method := c.PostForm("_method")
		if method == "DELETE" {
			controllers.DeleteUser(c)
			return
		}
		c.AbortWithStatus(http.StatusBadRequest)
	})

	// プロジェクト管理ルート
	r.GET("/settings/projects", controllers.GetProjects)
	r.POST("/settings/projects", controllers.CreateProject)
	r.POST("/settings/projects/:id", func(c *gin.Context) {
		method := c.PostForm("_method")
		if method == "DELETE" {
			controllers.DeleteProject(c)
			return
		}
		c.AbortWithStatus(http.StatusBadRequest)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}