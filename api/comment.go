package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"wedding/config"
	"wedding/models"
	"wedding/utils"
)

type Pagination struct {
	Limit      int         `json:"limit,omitempty;query:limit"`
	Page       int         `json:"page,omitempty;query:page"`
	Sort       string      `json:"sort,omitempty;query:sort"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages uint        `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

const PageLimit = 12

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = PageLimit
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}

func paginate(value interface{}, pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	config.GORMDB.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := uint(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
	fmt.Println(float64(totalRows), float64(pagination.Limit))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func GetDefaultComment(c *gin.Context) {
	pagination := Pagination{
		// Limit: 10,
		Page: 1,
	}

	comments, err := GetComments(c, pagination)
	if err != nil {
		c.HTML(http.StatusFound, "index.html", nil)
	} else {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"admin":    false,
			"comments": comments,
		})
	}
}

func GetAdminComment(c *gin.Context) {
	pagination := Pagination{
		// Limit: 10,
		Page: 1,
	}

	comments, err := GetComments(c, pagination)
	if err != nil {
		c.HTML(http.StatusFound, "index.html", nil)
	} else {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"admin":    true,
			"comments": comments,
		})
	}
}

func GetComment(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	utils.ErrorCheck(err)
	fmt.Println("Testset", page)

	pagination := Pagination{
		Page: page,
	}
	comments, err := GetComments(c, pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"comments": comments,
		})
	}
}

func GetComments(c *gin.Context, pagination Pagination) (string, error) {
	var comments []*models.Comment

	config.GORMDB.Scopes(paginate(comments, &pagination)).Find(&comments)
	pagination.Rows = comments

	commentsJson, err := json.Marshal(pagination)
	utils.ErrorCheck(err)

	return string(commentsJson), err
}

func SetComments(c *gin.Context) {
	var comment models.Comment
	value := requestBody(c)
	fmt.Println(string(value))
	err := json.Unmarshal(value, &comment)
	utils.ErrorCheck(err)
	fmt.Println("😘", comment)
	config.GORMDB.Create(&comment)
	commentJson, err := json.Marshal(comment)
	utils.ErrorCheck(err)
	fmt.Println(commentJson)
	c.JSON(http.StatusOK, gin.H{
		"comment": string(commentJson),
	})
}

//post 형식으로 body에 데이터 왔을 때 추출하는 함수
func requestBody(c *gin.Context) []byte {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	utils.ErrorCheck(err)
	return value
}
