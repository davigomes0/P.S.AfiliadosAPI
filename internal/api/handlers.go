package api

import (
	"net/http"

	"go-api/internal/database"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	DB *database.DB
}

func NewHandlers(db *database.DB) *Handlers {
	return &Handlers{DB: db}
}

func (h *Handlers) CreateConversion(c *gin.Context) {
	//Pega a chave de API do header
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Chave de API não fornecida"})
		return
	}

	var partner database.Partner
	err := h.DB.QueryRow("SELECT id, name FROM partners WHERE api_key = ?", apiKey).Scan(&partner.ID, &partner.Name)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Chave de API inválida"})
		return
	}

	//Receber Notificação de Conversão
	var req database.CreateConversionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de requisição inválidos"})
		return
	}

	//Evitar Duplicidade (Idempotência)
	var count int
	err = h.DB.QueryRow("SELECT COUNT(*) FROM conversions WHERE transaction_id = ?", req.TransactionID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar duplicidade"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Conversão já registrada", "transaction_id": req.TransactionID})
		return
	}

	//Persistir os Dados
	_, err = h.DB.Exec("INSERT INTO conversions (transaction_id, partner_id, amount) VALUES (?, ?, ?)", req.TransactionID, partner.ID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar a conversão"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Conversão registrada com sucesso!", "transaction_id": req.TransactionID})
}
