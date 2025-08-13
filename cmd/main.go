package main

import (
	"log"

	"github.com/davigomes0/P.S.AfiliadosAPI/internal/api"
	"github.com/davigomes0/P.S.AfiliadosAPI/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("erro ao carregar o arquivo .env: %v", err)
	}

	db := database.NewDB()

	router := gin.Default()

	handlers := api.NewHandlers(db)

	router.POST("/api/v1/conversions", handlers.CreateConversion)

	log.Println("Servidor iniciado na porta 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
