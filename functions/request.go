package functions

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Limit  int
	Offset int
}

func GetPagination(c *gin.Context) Pagination {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	return Pagination{Limit: limit, Offset: offset}
}

func (p *Pagination) GetLimit() int {
	return p.Limit
}

func (p *Pagination) GetOffset() int {
	return p.Offset
}

func (p *Pagination) SetLimit(limit int) {
	p.Limit = limit
}

func SendPaginatedResponse(c *gin.Context, data any, pagination Pagination) {
	if c.Query("limit") != "" {
		c.IndentedJSON(200, gin.H{
			"data": data,
			"pagination": gin.H{
				"limit":  pagination.GetLimit(),
				"offset": pagination.GetOffset(),
			},
		})
	} else {
		c.IndentedJSON(200, data)
	}
}
