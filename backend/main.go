package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DivinationRequest struct {
	Subject string `json:"subject" binding:"required"`
	Number1 int    `json:"number1" binding:"required,min=0"`
	Number2 int    `json:"number2" binding:"required,min=0"`
	Number3 int    `json:"number3" binding:"required,min=0"`
}

type DivinationRecord struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Subject      string    `json:"subject"`
	Number1      int       `json:"number1"`
	Number2      int       `json:"number2"`
	Number3      int       `json:"number3"`
	HexagramName string    `json:"hexagramName"`
	ChangingLine int       `json:"changingLine"`
	Summary      string    `json:"summary"`
	CreatedAt    time.Time `json:"createdAt"`
}

var (
	upperTrigrams = []string{"乾", "兑", "离", "震", "巽", "坎", "艮", "坤"}
	lowerTrigrams = []string{"乾", "兑", "离", "震", "巽", "坎", "艮", "坤"}
	hexagramNames = [][]string{
		{"乾为天", "天泽履", "天火同人", "天雷无妄", "天风姤", "天水讼", "天山遯", "天地否"},
		{"泽天夬", "兑为泽", "泽火革", "泽雷随", "泽风大过", "泽水困", "泽山咸", "泽地萃"},
		{"火天大有", "火泽睽", "离为火", "火雷噬嗑", "火风鼎", "火水未济", "火山旅", "火地晋"},
		{"雷天大壮", "雷泽归妹", "雷火丰", "震为雷", "雷风恒", "雷水解", "雷山小过", "雷地豫"},
		{"风天小畜", "风泽中孚", "风火家人", "风雷益", "巽为风", "风水涣", "风山渐", "风地观"},
		{"水天需", "水泽节", "水火既济", "水雷屯", "水风井", "坎为水", "水山蹇", "水地比"},
		{"山天大畜", "山泽损", "山火贲", "山雷颐", "山风蛊", "山水蒙", "艮为山", "山地剥"},
		{"地天泰", "地泽临", "地火明夷", "地雷复", "地风升", "地水师", "地山谦", "坤为地"},
	}
	interpretations = map[string]string{
		"乾为天":  "元亨利贞，代表充满创造力与领导力的状态。",
		"坤为地":  "厚德载物，象征包容与顺从，适合稳健前行。",
		"地天泰":  "阴阳交泰，万事通达，宜抓住机遇发展。",
		"天地否":  "闭塞不通，需要耐心等待转机，避免冒进。",
		"雷风恒":  "持久不息，保持恒心即可成功，切忌急躁。",
		"火水未济": "尚未完成，需要再接再厉，谨慎行事。",
		"水火既济": "大事已成，但仍需谨慎维护成果。",
	}
)

func main() {
	db, err := gorm.Open(sqlite.Open("divinations.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&DivinationRecord{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	router.POST("/api/divination", func(c *gin.Context) {
		var req DivinationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		record := createDivination(req)
		if err := db.Create(&record).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save record"})
			return
		}

		c.JSON(http.StatusOK, record)
	})

	router.GET("/api/divination", func(c *gin.Context) {
		var records []DivinationRecord
		if err := db.Order("created_at desc").Find(&records).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load records"})
			return
		}
		c.JSON(http.StatusOK, records)
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func createDivination(req DivinationRequest) DivinationRecord {
	total := req.Number1 + req.Number2 + req.Number3
	upper := total % len(upperTrigrams)
	lower := (req.Number1 + req.Number2) % len(lowerTrigrams)
	changingLine := total % 6
	if changingLine == 0 {
		changingLine = 6
	}

	hexagram := hexagramNames[upper][lower]
	summary, ok := interpretations[hexagram]
	if !ok {
		summary = "此卦象征局势变化多端，需要依据实际情况灵活应对。"
	}

	return DivinationRecord{
		Subject:      req.Subject,
		Number1:      req.Number1,
		Number2:      req.Number2,
		Number3:      req.Number3,
		HexagramName: hexagram,
		ChangingLine: changingLine,
		Summary:      summary,
		CreatedAt:    time.Now(),
	}
}
