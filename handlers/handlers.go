package handlers

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pulpfree/gales-fuelsale-report/model"
	"github.com/pulpfree/gales-fuelsale-report/xlsx"
	log "github.com/sirupsen/logrus"
)

// Response struct
// inspired from example at: https://labs.omniti.com/labs/jsend
/*type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
}*/

// DateRange struct used in date validation
type DateRange struct {
	DateStart time.Time `form:"startDate" binding:"required,validatedate" time_format:"2006-01-02"`
	DateEnd   time.Time `form:"endDate" binding:"required,validatedate" time_format:"2006-01-02"`
}

// HeartBeat function
func HeartBeat(c *gin.Context) {
	db := c.MustGet("mongo").(model.DBHandler)
	count, err := db.GetFuelSalesCount()
	if err != nil {
		log.Errorf("Failed to fetch fuel sales count: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Fuel sales doc count: " + strconv.Itoa(count),
	})

}

// FuelSale function
func FuelSale(c *gin.Context) {

	var dr DateRange
	if err := c.ShouldBindWith(&dr, binding.Query); err != nil {
		log.Errorf("Dates failed to validate: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rtype := c.Param("type")

	switch rtype {
	case "json":
		fetchJSON(c)
	case "xlsx":
		fetchXLSX(c)
	default:
		log.Errorf("Invalid report type requested: %s", rtype)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
	}
}

func fetchXLSX(c *gin.Context) {
	db := c.MustGet("mongo").(model.DBHandler)
	sales, err := db.GetFuelSales(c.Query("startDate"), c.Query("endDate"))
	if err != nil {
		log.Errorf("Failed to fetch fuel sales data: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	filename := "FuelSaleReport_" + c.Query("startDate") + "_" + c.Query("endDate") + ".xlsx"

	output := new(bytes.Buffer)
	xlsxFile, err := xlsx.NewFile(sales)
	if err != nil {
		log.Errorf("Failed to create xlsx file: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	err = xlsxFile.Write(output)
	if err != nil {
		log.Errorf("Failed to write xlsx file: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", output.Bytes())
}

func fetchJSON(c *gin.Context) {
	db := c.MustGet("mongo").(model.DBHandler)
	sales, err := db.GetFuelSales(c.Query("startDate"), c.Query("endDate"))
	if err != nil {
		log.Errorf("Failed to fetch fuel sales data: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, sales)
}
