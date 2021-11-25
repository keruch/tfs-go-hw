package repository

import (
	"context"
	"net/url"
	"testing"

	"github.com/keruch/tfs-go-hw/trading_robot/config"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/log"
	"github.com/stretchr/testify/suite"
)

type DatabaseSuite struct {
	suite.Suite
	repo   Repository
	logger *log.Logger
}

func (db *DatabaseSuite) SetupSuite() {
	logger := log.NewLogger()
	db.logger = logger
	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("Failed to setup log: %v", err)
	}
	u, _ := url.Parse(config.GetDatabaseURL())
	u.Path = "postgres_test"
	repo, err := NewPostgreSQLPool(u.String(), logger)
	if err != nil {
		logger.Fatal(err)
	}
	db.repo = repo
}

func (db *DatabaseSuite) TestStoreToDB() {
	testID := 0
	db.logger.Infof("\tTest %d:\tcreate transaction no error", testID)
	{
		respOrder := domain.CreateOrderResponse{
			Result:       "success",
			Status:       "placed",
			OrderID:      "8dcdbe17-b729-4fef-8b89-36e561535f38",
			ReceivedTime: "2021-11-25T19:05:03.670Z",
		}
		err := db.repo.StoreToDB(context.Background(), respOrder, 5242)
		db.NoError(err)
	}

	testID++
	db.logger.Infof("\tTest %d:\tcreate transaction error", testID)
	{
		respOrder := domain.CreateOrderResponse{
			Result:       "success",
			Status:       "placed",
			OrderID:      "8dcdbe17-b729-4fef-8b89-36e561535f38",
			ReceivedTime: "5435qgs4hv2nq4nugfg",
		}
		err := db.repo.StoreToDB(context.Background(), respOrder, 3425)
		db.Error(err)
	}
}

func TestDatabase(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}
