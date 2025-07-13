package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) RunAnalysis(c *gin.Context) {
	// Call the analysis service to run
	dolarBid, err := h.analisysService.RunAnalysis(c.Request.Context())
	if err != nil {
		h.log.Printf("erro ao obter cotação do dolar: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "ao obter cotação do dolar"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Valor atual da cotação do dolar", "dolar_bid": dolarBid})
	c.Next()
}
