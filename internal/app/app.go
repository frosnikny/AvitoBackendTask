package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project/internal/config"
	"project/internal/dsn"
	"project/internal/repository"
)

type Application struct {
	repo   *repository.Repository
	config *config.Config
}

func (a *Application) StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	// r.Use(ErrorHandler())

	api := r.Group("/api")
	{
		api.GET("/ping", a.GetPing)

		tenders := api.Group("/tenders")
		{
			tenders.GET("", a.GetTenders)
			tenders.POST("/new", a.PostNewTender)
			tenders.GET("/my", a.GetMyTenders)
			tenders.GET("/:tenderId/status", a.GetTenderStatus)
			tenders.PUT("/:tenderId/status", a.PutTenderStatus)
			tenders.PATCH("/:tenderId/edit", a.PatchTender)

			tenders.PUT("/:tenderId/rollback/:version", a.PutTenderRollback)
		}

		bids := api.Group("/bids")
		{
			bids.POST("/new", a.PostNewBid)
			bids.GET("/my", a.GetMyBids)
			// На самом деле здесь tenderId, но из-за ошибки укажем как bidId
			// panic: ':bidId' in new path '/api/bids/:bidId/status' conflicts with existing wildcard ':tenderId' in existing prefix '/api/bids/:tenderId'
			bids.GET("/:bidId/list", a.GetBidsList)
			bids.GET("/:bidId/status", a.GetBidStatus)
			bids.PUT("/:bidId/status", a.PutBidStatus)
			bids.PATCH("/:bidId/edit", a.PatchBid)
			bids.PUT("/:bidId/submit_decision", a.PutBitSubmitDecision)

			bids.PUT("/:bidId/feedback", a.PutBitFeedback)
			bids.GET("/:bidId/reviews", a.GetBidReviews)

			bids.PUT("/:bidId/rollback/:version", a.PutBidRollback)
		}
	}

	s := &http.Server{
		Addr:    a.config.ServerAddress,
		Handler: r,
	}

	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}

	log.Println("Server down")
}

func New() *Application {
	var err error

	a := &Application{}
	a.config, err = config.New()
	if err != nil {
		panic(err)
	}

	a.repo, err = repository.New(dsn.FromCfg(a.config))
	if err != nil {
		panic(err)
	}

	return a
}
