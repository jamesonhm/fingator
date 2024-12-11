package uri_test

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/jamesonhm/fingator/internal/polygon/models"
	"github.com/jamesonhm/fingator/internal/uri"
)

func TestEncodePolyParams(t *testing.T) {
	testPath := "/v1/{num}/{str}"

	type Params struct {
		Num float64 `path:"num"`
		Str string  `path:"str"`

		NumQ *float64 `query:"num"`
		StrQ *string  `query:"str"`
	}

	//num := 1.273
	str := "teststr"
	params := Params{
		Num: 1.273,
		Str: str,
	}

	expected := "/v1/1.273/teststr"
	actual := uri.New("").EncodeParams(testPath, params)
	assert.Equal(t, actual, expected)
}

func TestEncodePolyDate(t *testing.T) {
	testPath := "/v1/{date}"

	type Params struct {
		Date  models.Date  `path:"date"`
		DateQ *models.Date `query:"date"`
	}

	pdate := models.Date(time.Date(2023, 12, 6, 0, 0, 0, 0, time.UTC))
	params := Params{
		Date:  pdate,
		DateQ: &pdate,
	}

	expected := "/v1/2023-12-06?date=2023-12-06"
	actual := uri.New("").EncodeParams(testPath, params)
	assert.Equal(t, actual, expected)
}

func TestEncodePointer(t *testing.T) {
	testPath := "/v2/{date}"

	type Params struct {
		Date models.Date `path:"date"`
		StrQ *string     `query:"str"`
	}

	pdate := models.Date(time.Date(2023, 12, 6, 0, 0, 0, 0, time.UTC))
	str := "teststr"
	params := Params{
		Date: pdate,
		StrQ: &str,
	}
	expected := "/v2/2023-12-06?str=teststr"
	actual := uri.New("").EncodeParams(testPath, &params)
	assert.Equal(t, actual, expected)
}
