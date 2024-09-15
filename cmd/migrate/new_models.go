package main

import (
	"project/internal/config"
	"project/internal/dsn"
	"project/internal/models"
	"project/internal/repository"
)

func main() {
	var err error

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	db, err := repository.CreateDB(dsn.FromCfg(cfg))
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Bid{},
		&models.Tender{},
		&models.Feedback{},
		&models.TenderVersion{},
		&models.BidVersion{})
	if err != nil {
		panic(err)
	}
}
