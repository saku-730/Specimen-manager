package controllers

import (
	"GO-specimen/models"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SpecimenForm はフォームからのデータバインディング用
type SpecimenForm struct {
	Latitude           float64 `form:"Latitude"`
	Longitude          float64 `form:"Longitude"`
	Date               string  `form:"Date"`
	Time               int     `form:"Time"`
	CollectionMethod   string  `form:"CollectionMethod"`
	Collector          string  `form:"Collector"`
	IndividualCount    int     `form:"IndividualCount"`
	Weather            string  `form:"Weather"`
	Temperature        float64 `form:"Temperature"`
	Humidity           float64 `form:"Humidity"`
	Precipitation      float64 `form:"Precipitation"`
	Environment        string  `form:"Environment"`
	CollectionRemarks  string  `form:"CollectionRemarks"`
	SpecimenCreator    string  `form:"SpecimenCreator"`
	SpecimenType       string  `form:"SpecimenType"`
	SpecimenLocation   string  `form:"SpecimenLocation"`
	SpecimenBoxID      string  `form:"SpecimenBoxID"`
	SpecimenCreationDate string  `form:"SpecimenCreationDate"`
	SpecimenRemarks    string  `form:"SpecimenRemarks"`
	Sex                string  `form:"Sex"`
	SpeciesName        string  `form:"SpeciesName"`
	JapaneseName       string  `form:"JapaneseName"`
	Age                string  `form:"Age"`
	Identifier         string  `form:"Identifier"`
	ProjectName        string  `form:"ProjectName"`
	DataInputDate      string  `form:"DataInputDate"`
	
	CommonRemarks      string  `form:"CommonRemarks"`
}

func parseDateString(dateStr string) *time.Time {
	if dateStr == "" || dateStr == "0" {
		return nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}
	return &t
}

func getUsers(db *gorm.DB) []models.User {
	var users []models.User
	db.Find(&users)
	return users
}

// GetUsersList はユーザーリストを返します (main.goから呼び出し可能)
func GetUsersList(db *gorm.DB) []models.User {
	return getUsers(db)
}

func getProjects(db *gorm.DB) []models.Project {
	var projects []models.Project
	db.Find(&projects)
	return projects
}

// GetProjectsList はプロジェクトリストを返します (main.goから呼び出し可能)
func GetProjectsList(db *gorm.DB) []models.Project {
	return getProjects(db)
}

// GetDefaultProject はデフォルトプロジェクトを返します。見つからない場合はnilを返します。
func GetDefaultProject(db *gorm.DB) *models.Project {
	var project models.Project
	if err := db.Where("is_default = ?", true).First(&project).Error; err != nil {
		return nil // デフォルトプロジェクトが見つからない場合
	}
	return &project
}

// GetDefaultUser はデフォルトユーザーを返します。見つからない場合はnilを返します。
func GetDefaultUser(db *gorm.DB) *models.User {
	var user models.User
	if err := db.Where("is_default = ?", true).First(&user).Error; err != nil {
		return nil // デフォルトユーザーが見つからない場合
	}
	return &user
}

func CreateSpecimen(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var form SpecimenForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	specimen := models.Specimen{
		Latitude:           form.Latitude,
		Longitude:          form.Longitude,
		Date:               parseDateString(form.Date),
		Time:               form.Time,
		CollectionMethod:   form.CollectionMethod,
		Collector:          form.Collector,
		IndividualCount:    form.IndividualCount,
		Weather:            form.Weather,
		Temperature:        form.Temperature,
		Humidity:           form.Humidity,
		Precipitation:      form.Precipitation,
		Environment:        form.Environment,
		CollectionRemarks:  form.CollectionRemarks,
		SpecimenCreator:    form.SpecimenCreator,
		SpecimenType:       form.SpecimenType,
		SpecimenLocation:   form.SpecimenLocation,
		SpecimenBoxID:      form.SpecimenBoxID,
		SpecimenCreationDate: parseDateString(form.SpecimenCreationDate),
		SpecimenRemarks:    form.SpecimenRemarks,
		Sex:                form.Sex,
		SpeciesName:        form.SpeciesName,
		JapaneseName:       form.JapaneseName,
		Age:                form.Age,
		Identifier:         form.Identifier,
		ProjectName:        form.ProjectName,
		DataInputDate:      parseDateString(form.DataInputDate),
		
		CommonRemarks:      form.CommonRemarks,
	}

	// Specimenを保存
	if err := db.Create(&specimen).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create specimen"})
		return
	}

	// 写真のアップロード処理
	multipartForm, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	files := multipartForm.File["photos"]

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		uniqueFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)
		dst := filepath.Join("./uploads", uniqueFilename)

		// ファイルを保存
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
			return
		}

		// Photoレコードを作成
		photo := models.Photo{
			SpecimenID: specimen.ID,
			Path:       "/uploads/" + uniqueFilename, // Webアクセス用のパス
		}
		if err := db.Create(&photo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo record"})
			return
		}
	}

	c.Redirect(http.StatusFound, "/specimens/new")
}

func GetSpecimens(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var searchParams struct {
		SpeciesName        string `form:"species_name"`
		IncludeSp          bool   `form:"include_sp"`
		JapaneseName       string `form:"japanese_name"`
		Collector          string `form:"collector"`
		Identifier         string `form:"identifier"`
		SpecimenCreator    string `form:"specimen_creator"`
		Latitude           float64 `form:"latitude"`
		Longitude          float64 `form:"longitude"`
		ProjectName        string `form:"project_name"`
		SpecimenBoxID      string `form:"specimen_box_id"`
		CollectionMethod   string `form:"collection_method"`
		SpecimenType       string `form:"specimen_type"`
		Sex                string `form:"sex"`
		Age                string `form:"age"`
		DateStart          string `form:"date_start"`
		DateEnd            string `form:"date_end"`
		SpecimenCreationDateStart string `form:"specimen_creation_date_start"`
		SpecimenCreationDateEnd string `form:"specimen_creation_date_end"`
		DataInputDateStart string `form:"data_input_date_start"`
		DataInputDateEnd   string `form:"data_input_date_end"`
	}

	if err := c.ShouldBindQuery(&searchParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := db.Model(&models.Specimen{})

	if searchParams.SpeciesName != "" {
		if searchParams.IncludeSp {
			query = query.Where("species_name LIKE ?", "%"+searchParams.SpeciesName+"%")
		} else {
			query = query.Where("species_name = ?", searchParams.SpeciesName)
		}
	}
	if searchParams.JapaneseName != "" {
		query = query.Where("japanese_name LIKE ?", "%"+searchParams.JapaneseName+"%")
	}
	if searchParams.Collector != "" {
		query = query.Where("collector LIKE ?", "%"+searchParams.Collector+"%")
	}
	if searchParams.Identifier != "" {
		query = query.Where("identifier LIKE ?", "%"+searchParams.Identifier+"%")
	}
	if searchParams.SpecimenCreator != "" {
		query = query.Where("specimen_creator LIKE ?", "%"+searchParams.SpecimenCreator+"%")
	}
	if searchParams.Latitude != 0 {
		query = query.Where("latitude = ?", searchParams.Latitude)
	}
	if searchParams.Longitude != 0 {
		query = query.Where("longitude = ?", searchParams.Longitude)
	}
	if searchParams.ProjectName != "" {
		query = query.Where("project_name LIKE ?", "%"+searchParams.ProjectName+"%")
	}
	if searchParams.SpecimenBoxID != "" {
		query = query.Where("specimen_box_id LIKE ?", "%"+searchParams.SpecimenBoxID+"%")
	}
	if searchParams.CollectionMethod != "" {
		query = query.Where("collection_method LIKE ?", "%"+searchParams.CollectionMethod+"%")
	}
	if searchParams.SpecimenType != "" {
		query = query.Where("specimen_type LIKE ?", "%"+searchParams.SpecimenType+"%")
	}
	if searchParams.Sex != "" {
		query = query.Where("sex LIKE ?", "%"+searchParams.Sex+"%")
	}
	if searchParams.Age != "" {
		query = query.Where("age LIKE ?", "%"+searchParams.Age+"%")
	}

	// 日付範囲検索
	if searchParams.DateStart != "" {
		query = query.Where("date >= ?", searchParams.DateStart)
	}
	if searchParams.DateEnd != "" {
		query = query.Where("date <= ?", searchParams.DateEnd)
	}
	if searchParams.SpecimenCreationDateStart != "" {
		query = query.Where("specimen_creation_date >= ?", searchParams.SpecimenCreationDateStart)
	}
	if searchParams.SpecimenCreationDateEnd != "" {
		query = query.Where("specimen_creation_date <= ?", searchParams.SpecimenCreationDateEnd)
	}
	if searchParams.DataInputDateStart != "" {
		query = query.Where("data_input_date >= ?", searchParams.DataInputDateStart)
	}
	if searchParams.DataInputDateEnd != "" {
		query = query.Where("data_input_date <= ?", searchParams.DataInputDateEnd)
	}

	var specimens []models.Specimen
	query.Find(&specimens)

	c.HTML(http.StatusOK, "list.html", gin.H{
		"Specimens": specimens,
		"SearchParams": searchParams,
	})
}

func GetSpecimen(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var specimen models.Specimen
	// PhotosをPreloadする
	if err := db.Preload("Photos").First(&specimen, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Specimen not found"})
		return
	}
	c.HTML(http.StatusOK, "show.html", specimen)
}

func EditSpecimen(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var specimen models.Specimen
	if err := db.Preload("Photos").First(&specimen, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Specimen not found"})
		return
	}
	users := GetUsersList(db)
	c.HTML(http.StatusOK, "edit.html", gin.H{
		"Specimen": specimen,
		"Users":    users,
	})
}

func UpdateSpecimen(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var specimen models.Specimen
	if err := db.First(&specimen, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Specimen not found"})
		return
	}

	var form SpecimenForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	specimen.Latitude = form.Latitude
	specimen.Longitude = form.Longitude
	specimen.Date = parseDateString(form.Date)
	specimen.Time = form.Time
	specimen.CollectionMethod = form.CollectionMethod
	specimen.Collector = form.Collector
	specimen.IndividualCount = form.IndividualCount
	specimen.Weather = form.Weather
	specimen.Temperature = form.Temperature
	specimen.Humidity = form.Humidity
	specimen.Precipitation = form.Precipitation
	specimen.Environment = form.Environment
	specimen.CollectionRemarks = form.CollectionRemarks
	specimen.SpecimenCreator = form.SpecimenCreator
	specimen.SpecimenType = form.SpecimenType
	specimen.SpecimenLocation = form.SpecimenLocation
	specimen.SpecimenBoxID = form.SpecimenBoxID
	specimen.SpecimenCreationDate = parseDateString(form.SpecimenCreationDate)
	specimen.SpecimenRemarks = form.SpecimenRemarks
	specimen.Sex = form.Sex
	specimen.SpeciesName = form.SpeciesName
	specimen.JapaneseName = form.JapaneseName
	specimen.Age = form.Age
	specimen.Identifier = form.Identifier
	specimen.ProjectName = form.ProjectName
	specimen.DataInputDate = parseDateString(form.DataInputDate)
	
	specimen.CommonRemarks = form.CommonRemarks

	db.Save(&specimen)

	// 写真のアップロード処理 (新規追加のみ、既存写真の削除は未実装)
	multipartForm, err := c.MultipartForm()
	if err == nil && multipartForm.File["photos"] != nil {
		files := multipartForm.File["photos"]

		for _, file := range files {
			filename := filepath.Base(file.Filename)
			uniqueFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)
			dst := filepath.Join("./uploads", uniqueFilename)

			if err := c.SaveUploadedFile(file, dst); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
				return
			}

			photo := models.Photo{
				SpecimenID: specimen.ID,
				Path:       "/uploads/" + uniqueFilename,
			}
			if err := db.Create(&photo).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo record"})
				return
			}
		}
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/specimens/%d", specimen.ID))
}

func DeleteSpecimen(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// トランザクションを開始
	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	var specimen models.Specimen
	// トランザクション内で標本を検索
	if err := tx.First(&specimen, id).Error; err != nil {
		tx.Rollback() // エラーが発生したらロールバック
		c.JSON(http.StatusNotFound, gin.H{"error": "Specimen not found"})
		return
	}

	// 関連する写真のレコードを削除
	// 注意: ここではDBレコードのみ削除。物理ファイルは削除されません。
	if err := tx.Where("specimen_id = ?", specimen.ID).Delete(&models.Photo{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photos"})
		return
	}

	// 標本レコードを削除
	if err := tx.Delete(&specimen).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete specimen"})
		return
	}

	// トランザクションをコミット
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // コミット失敗時もロールバック
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.Redirect(http.StatusFound, "/specimens")
}

func GetIndexData(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var latestEntries []models.Specimen
	db.Order("created_at desc").Limit(10).Find(&latestEntries)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"LatestEntries": latestEntries,
	})
}

func DeleteSpecimens(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)

    // フォームから "ids" をスライスとして受け取る
    ids := c.PostFormArray("ids")
    if len(ids) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No specimens selected"})
        return
    }

    // トランザクションを開始
    tx := db.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }

    // 選択されたIDに基づいて関連するPhotoレコードを削除
    // 注意: ここではDBレコードのみ削除。物理ファイルは削除されません。
    if err := tx.Where("specimen_id IN ?", ids).Delete(&models.Photo{}).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photos"})
        return
    }

    // 選択されたIDに基づいてSpecimenレコードを削除
    if err := tx.Where("id IN ?", ids).Delete(&models.Specimen{}).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete specimens"})
        return
    }

    // トランザクションをコミット
    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
        return
    }

    // 成功したら一覧ページにリダイレクト
    c.Redirect(http.StatusFound, "/specimens")
}
