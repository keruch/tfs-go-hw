package kraken

import (
	"fmt"
	"testing"

	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestQueryByOperation(t *testing.T) {
	a := assert.New(t)

	testID := 0
	t.Logf("\tTest %d:\tcreate order operation", testID)
	{
		_, err := QueryByOperation(domain.Order{}, CreateOrder)
		a.NoErrorf(err, "Should not be error")
	}
	testID++
	t.Logf("\tTest %d:\topen orders operation", testID)
	{
		_, err := QueryByOperation(domain.Order{}, OpenOrders)
		a.NoErrorf(err, "Should not be error")
	}
	testID++
	t.Logf("\tTest %d:\tcancel order operation", testID)
	{
		_, err := QueryByOperation(domain.Order{}, CreateOrder)
		a.NoErrorf(err, "Should not be error")
	}
	testID++
	t.Logf("\tTest %d:\tcancel order operation", testID)
	{
		_, err := QueryByOperation(domain.Order{}, "TEST_OPERATION")
		a.Errorf(err, "Should not error")
		a.Equalf(ErrOperationNotFound, err, "Errors should be equal")
	}
}

func TestGenerateToken(t *testing.T) {
	a := assert.New(t)

	testID := 0
	t.Logf("\tTest %d:\tcreate order", testID)
	{
		privateKey := "D3xf8GQWr/HylLexjL2055e5Pn5z+vIyu0zsSRfbFA+W07Q4jbpb6qa0H5xzbywHWPG7quEA6V2imZ3iqNEe5Bj/"
		endpoint := "/api/v3/sendorder"
		postData := "chart=&limitPrice=4571.1&orderType=ioc&side=sell&size=100&stopPrice=0.0&symbol=PI_ETHUSD&triggerSignal="
		token, err := GenerateToken(privateKey, endpoint, postData)
		a.NoError(err)
		fmt.Println(token)
	}
}

