package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnixTime(t *testing.T) {
	a := assert.New(t)

	testID := 0
	t.Logf("\tTest %d:\tUnmarshalJSON unix time no error", testID)
	{
		timeJSON := []byte(`1612269657781`)
		var ts UnixTS
		err := json.Unmarshal(timeJSON, &ts)
		a.NoError(err)
		a.Equal(UnixTS(time.Date(2021, 02, 02, 15, 40, 57, 781000000, time.Local)), ts)
		a.Equal("2021-02-02 15:40:57.781 +0300 MSK", ts.String())
	}

	t.Logf("\tTest %d:\tUnmarshalJSON unix time error", testID)
	{
		timeJSON := []byte(`4161a26487269657781`)
		var ts UnixTS
		err := json.Unmarshal(timeJSON, &ts)
		a.Error(err)
		a.Equal(UnixTS(time.Time{}), ts)
	}
}
