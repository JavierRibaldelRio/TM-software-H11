package sensors

import(

	"math/rand"
	"time"
)

type Sensor struct {
    Name     	string
    Unit     	string
    Period   	time.Duration
	Avg			float64	
	Sdev		float64
}


type Data struct{
	Measure 		float64
	Unit		string
	Timestamp	time.Time
}

// It executes when the package is loaded
func init() {
	rand.Seed(time.Now().UnixNano())
}

func (s Sensor) Read () Data {
	v := s.Avg + s.Sdev*rand.NormFloat64() // N(Avg, Sdev)
	return Data{
		Timestamp: time.Now(),
		Measure:   v,
		Unit:      s.Unit,
	}
}