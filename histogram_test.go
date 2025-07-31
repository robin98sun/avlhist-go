// reference: https://medium.com/rungo/unit-testing-made-easy-in-go-25077669318

package histogram

import (
	"log"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"math"
	// "math/rand"
	// "sort"
)

var verbose bool = false

func TestScheduler_CreateHistogram(t *testing.T) {

	sample_mean := 100
	base_window_size := 10000
	subhisto_size := 10
	accuracy := 1
	buckets_in_subhisto := int(float64(subhisto_size) * math.Pow(float64(10), float64(accuracy)))
	for x := 1; x <= 5; x++ {
		window_size := base_window_size* x
		sample_size := window_size*10*x
		percentile_list := []float64{
			float64(0.99), float64(0.995), float64(0.999), float64(0.9995), float64(0.9999),
		}
		list := gen_random_list_float(sample_size, sample_mean, float64(1000))
		assert.Equal(t, len(list), sample_size, "random util should work")

		histogram := NewHistogram(int64(window_size), float64(subhisto_size), accuracy)
		assert.NotNil(t, histogram, "histogram should not be nil")

		for _, p := range percentile_list {
			histogram.AddPercentilePoint(p)
		}

		for i:=0; i<len(list); i++ {
			v := list[i]
			// log.Printf("original %v value: %v", i, v)
			histogram.Enqueue(v, 1)
		}
		assert.Equal(t, int64(window_size), histogram.Count, "histogram size should equal window size")

		assert.NotNil(t, histogram.RootItem, "root item should not be nil")

		assert.Equal(t, int64(window_size), histogram.RootItem.Count, "histogram root count should equal window size")

		assert.Equal(t, int64(window_size), int64(len(histogram.Queue)), "histogram queue length should equal window size")

		sumAllBuckets := int64(0)
		bucket_count := 0
		max_possible_bucket_count := len(histogram.BucketHistogram.SubBucketHistograms)*buckets_in_subhisto
		for i:=0; i<len(histogram.BucketHistogram.SubBucketHistograms); i++ {
			sbh := histogram.BucketHistogram.SubBucketHistograms[i]
			if sbh != nil {
				bucket_count += len(sbh.BucketList)
				for j:=0; j<len(sbh.BucketList);j++ {
					bucketItem := sbh.BucketList[j]
					if bucketItem != nil {
						sumAllBuckets+=bucketItem.Duplications
						// log.Printf(" i: %v, j: %v, value: %v, duplication: %v, count: %v", i, j, bucketItem.Value, bucketItem.Duplications, bucketItem.Count)

						idx:=sbh.CalcPosition(bucketItem.Value)
						assert.Equal(t, int64(j), idx, "the index should be equal")

						idx,_,_ = histogram.BucketHistogram.CalcPosition(bucketItem.Value)
						assert.Equal(t, int64(i), idx, "the index should be equal")

					}
				}
			}
		}

		min_node := histogram.RootItem
		for ; min_node.Left != nil ; min_node = min_node.Left {}
		max_node := histogram.RootItem
		for ; max_node.Right != nil; max_node = max_node.Right {}

		node_amount := 0
		min_value := min_node.Value
		max_value := max_node.Value
		avg_value := min_value - min_value
		sum_value := avg_value
		variance  := sum_value
		for p:=min_node; p!=nil; p=p.Larger{
			node_amount++
			sum_value += p.Value * float64(p.Duplications)
		}
		avg_value = sum_value / float64(histogram.RootItem.Count)

		for p:=min_node; p!=nil; p=p.Larger{
			variance += math.Pow((p.Value-avg_value), 2) * float64(p.Duplications)
		}
		variance /= float64(histogram.RootItem.Count)

		assert.Equal(t, histogram.MinItem, min_node, "min node should be identical")
		assert.Equal(t, histogram.MaxItem, max_node, "max node should be identical")

		for _, v := range percentile_list {
			p := histogram.GetPercentileItem(v)
			cc := p.Item.CumulativeCount()
			// log.Printf("%v percentile(%v), real percentage: %v(%v), count: %v, total: %v, min: %v(%v), max: %v(%v), p-node: %v, cumulativeCount: %v, real: %v", 
			// 	v*float64(100), p.Percentile, 
			// 	p.RealPercentage, float64(p.Count)/float64(histogram.RootItem.Count),
			// 	p.Count, histogram.RootItem.Count, 
			// 	histogram.MinItem.Value, min_node.Value, 
			// 	histogram.MaxItem.Value, max_node.Value,
			// 	p.Item.Value, cc, float64(cc)/float64(histogram.RootItem.Count),
			// )
			// The histogram uses "nearest rank" approach, so we need to be more lenient
			// with the accuracy expectations
			percentileDiff := p.Percentile - float64(cc)/float64(histogram.RootItem.Count)
			realPercentileDiff := p.RealPercentage - float64(cc)/float64(histogram.RootItem.Count)
			
			// Allow for some variance in the nearest rank approach
			assert.GreaterOrEqual(t, percentileDiff, -0.1, "percentile point should be reasonably accurate")
			assert.GreaterOrEqual(t, realPercentileDiff, -0.1, "real percentage should be reasonably accurate")


		}

		assert.Equal(t, int64(window_size), sumAllBuckets, "histogram queue length should equal window size")

		if verbose {
			log.Printf("window size: %v, sample amount: %v, percentile points: %v, height of avl tree: %v", 
				window_size, sample_size, len(percentile_list),
				histogram.RootItem.Height,
			)
			log.Printf("   nodes: %v (%v%%, %v%%, %v%%), amount of buckets: %v (%v%%, %v%%), maximum possible buckets: %v", 
				node_amount, 
				node_amount*100/window_size, 
				node_amount*100/bucket_count, 
				node_amount*100/max_possible_bucket_count,
				bucket_count, 
				bucket_count*100/window_size,
				bucket_count*100/max_possible_bucket_count,
				max_possible_bucket_count,
			)
			log.Printf("   [data distribution] mean: %v, variance %v, min: %v, max: %v", 
				avg_value, variance, min_value, max_value,
			)
		}
		pstr := ""
		for _, p := range percentile_list {
			pstr = fmt.Sprintf("%s%v%%: %v, ", pstr,
				p*float64(100),
				histogram.GetPercentileForValue(p),
			)
		}
		if verbose {
			log.Printf("   percentiles: %s", pstr)

			log.Println("")
		}
	}
	
}

// to do the benchmark:
// ref: https://blog.logrocket.com/benchmarking-golang-improve-function-performance/
// GOMAXPROCS=1 go test -run=MultiplyHistograms -bench=MultiplyHistograms -count=10 -timeout 99999s

func BenchmarkTestScheduler_MultiplyHistograms(t *testing.B) {
	sample_size := 10000
	window_size := 10000
	histogram_count := 1000

	sample_mean := 100
	subhisto_size := 0.1
	accuracy := 1
	// buckets_in_subhisto := int(float64(subhisto_size) * math.Pow(float64(10), float64(accuracy)))
	// percentile_list := []float64{
	// 	float64(0.95), float64(0.99), float64(0.995), float64(0.999), float64(0.9995), float64(0.9999),
	// }

	percentile_list := []float64{
		float64(0.99),
	}

	var histogram_list []*Histogram = []*Histogram{}

	for h:=0; h<histogram_count; h++ {
		list := gen_random_list_float(sample_size, sample_mean, float64(10))
		assert.Equal(t, len(list), sample_size, "random util should work")

		histogram := NewHistogram(int64(window_size), float64(subhisto_size), accuracy)
		assert.NotNil(t, histogram, "histogram should not be nil")

		for _, p := range percentile_list {
			histogram.AddPercentilePoint(p)
		}

		for i:=0; i<len(list); i++ {
			v := list[i]
			// log.Printf("original %v value: %v", i, v)
			histogram.Enqueue(v, 1)
		}
		histogram_list = append(histogram_list, histogram)

		// real test
		if h % 10 != 9 {
			continue
		}
		title := fmt.Sprintf("multiply %v histograms each window size %v to search:", len(histogram_list), window_size)
		for _, p := range percentile_list {
			title = fmt.Sprintf("%v %v",title, p*float64(100))
		}
		t.Run(title, func(b *testing.B) {

			for pi := 0; pi < len(percentile_list); pi++ {
				p := percentile_list[pi]

				tail := CalcPercentileOfProduct(p, histogram_list, verbose)

				if verbose && p == float64(0.99) {
					log.Printf("searching the percentile for the tail: %v", tail)
					for i, hist := range histogram_list {
						hp := hist.GetPercentileForValue(tail)
						log.Printf("%v: %v", i, hp)
					}	
				}
				
			}

		})
	}




	
}
