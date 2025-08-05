package main

import (
	"encoding/csv"
	"log"
	"os"
	"shop/internal/domain"
	"shop/internal/repository/mysql"
	"shop/pkg/database"
	"strconv"
)

func main() {
	db, err := database.NewMySQLConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	productRepo := mysql.NewProductRepository(db)

	file, err := os.Open("/home/aka/Templates/golang/cmd/import_data/products.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Skip header row
	for _, record := range records[1:] {
		price, _ := strconv.ParseFloat(record[2], 64)
		stock, _ := strconv.Atoi(record[3])

		product := &domain.Product{
			Name:        record[0],
			Description: record[1],
			Price:       price,
			Stock:       stock,
			Category:    record[4],
		}

		if err := productRepo.Create(product); err != nil {
			log.Printf("Error creating product %s: %v", product.Name, err)
		} else {
			log.Printf("Created product: %s", product.Name)
		}
	}
}
