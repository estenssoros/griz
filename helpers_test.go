package griz

import "testing"

func TestDataType(t *testing.T) {
	{
		data := 4.0
		if want, have := FloatType, dataType(data); want != have {
			t.Errorf("have: %d, want: %d", have, want)
		}
	}
	{
		data := "asdf"
		if want, have := StringType, dataType(data); want != have {
			t.Errorf("have: %d, want: %d", have, want)
		}
	}
	{
		data := []float64{1.0}
		if want, have := FloatType, dataType(data); want != have {
			t.Errorf("have: %d, want: %d", have, want)
		}
	}
	{
		data := []string{"asdf"}
		if want, have := StringType, dataType(data); want != have {
			t.Errorf("have: %d, want: %d", have, want)
		}
	}

}
