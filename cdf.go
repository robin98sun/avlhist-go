package histogram

import (
	"math"
	// "sync"
	// "time"
	// "uta.edu/aces/jade-go/kernel"
	// "strconv"
	// "fmt"
)

type CDFPoint struct {
	Percentile float64
	Value      float64
}

type CDF struct {
	Points 	[]*CDFPoint `json:"points,omitempty"`
	StartPoint float64 `json:"startPoint,omitempty"`
	Increment  float64 `json:"increment,omitempty"`
	Amount int `json:"amount,omitempty"`
	histogram *Histogram
}

func NewCDF(amountOfPoint int) *CDF {
	return &CDF{
		Points: make([]*CDFPoint, amountOfPoint),
		Amount: amountOfPoint,
	}
}

func (c *CDF) Histogram() *Histogram {
	if c == nil {
		return nil
	}
	if c.histogram != nil {
		return c.histogram
	}
	count_zero := int(math.Round((1/(1-c.StartPoint)-1)))*(c.Amount-1)
	total_count := count_zero + c.Amount

	hist := NewHistogram(int64(total_count), float64(0.1), 1)
	hist.AddPercentilePoint(float64(0.99))
	hist.Enqueue(0, count_zero)

	// assume points are sorted
	for _, p := range c.Points {
		hist.Enqueue(p.Value, 1)
	}
	c.histogram = hist

	return hist
}


func SearchCDFProduct(cdf_list []*CDF, percentile float64) float64 {
	if len(cdf_list) == 0 {
		return 0
	}

	hist_list := make([]*Histogram, len(cdf_list))

	for i, cdf := range cdf_list {
		if cdf == nil {continue}
		hist_list[i] = cdf.Histogram()
	}

	return CalcPercentileOfProduct(percentile, hist_list, false)

}
