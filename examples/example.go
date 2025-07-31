package main

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/robin98sun/avlhist-go"
)

func main() {
	// Example 1: Basic histogram creation and usage
	fmt.Println("=== Basic Histogram Example ===")
	
	// Create a histogram with:
	// - Queue size: 1000 (sliding window size)
	// - Sub-bucket histogram size: 10.0
	// - Accuracy: 1 (decimal places)
	hist := histogram.NewHistogram(1000, 10.0, 1)
	
	// Add percentile points to track
	hist.AddPercentilePoint(0.95)  // 95th percentile
	hist.AddPercentilePoint(0.99)  // 99th percentile
	hist.AddPercentilePoint(0.999) // 99.9th percentile
	
	// Simulate adding data points
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 500; i++ {
		// Generate random values (e.g., response times)
		value := rand.Float64() * 1000
		hist.Enqueue(value, 1)
	}
	
	// Get statistics
	fmt.Printf("Total count: %d\n", hist.Count)
	fmt.Printf("95th percentile: %.2f\n", hist.GetValueAtPercentile(0.95))
	fmt.Printf("99th percentile: %.2f\n", hist.GetValueAtPercentile(0.99))
	fmt.Printf("99.9th percentile: %.2f\n", hist.GetValueAtPercentile(0.999))
	
	// Example 2: Working with multiple histograms
	fmt.Println("\n=== Multiple Histograms Example ===")
	
	// Create multiple histograms for different metrics
	// Use larger sub-bucket sizes to avoid negative indices
	cpuHist := histogram.NewHistogram(100, 10.0, 1)
	memoryHist := histogram.NewHistogram(100, 50.0, 1)
	latencyHist := histogram.NewHistogram(100, 1.0, 1)
	
	// Add percentile tracking to each
	cpuHist.AddPercentilePoint(0.95)
	memoryHist.AddPercentilePoint(0.95)
	latencyHist.AddPercentilePoint(0.95)
	
	// Simulate monitoring data
	for i := 0; i < 50; i++ {
		cpuHist.Enqueue(rand.Float64()*100, 1)      // CPU usage 0-100%
		memoryHist.Enqueue(rand.Float64()*1000, 1)   // Memory usage 0-1000MB
		latencyHist.Enqueue(rand.Float64()*10, 1)    // Latency 0-10ms
	}
	
	// Calculate product percentiles (useful for SLO calculations)
	histList := []*histogram.Histogram{cpuHist, memoryHist, latencyHist}
	product95th := histogram.CalcPercentileOfProduct(0.95, histList, false)
	fmt.Printf("Product 95th percentile: %.2f\n", product95th)
	
	// Example 3: CDF (Cumulative Distribution Function) usage
	fmt.Println("\n=== CDF Example ===")
	
	// Create a CDF with 100 points
	cdf := histogram.NewCDF(100)
	
	// Add some sample points (in practice, these would come from your data)
	cdf.Points[0] = &histogram.CDFPoint{Percentile: 0.1, Value: 10.0}
	cdf.Points[1] = &histogram.CDFPoint{Percentile: 0.5, Value: 50.0}
	cdf.Points[2] = &histogram.CDFPoint{Percentile: 0.9, Value: 90.0}
	cdf.Points[3] = &histogram.CDFPoint{Percentile: 0.95, Value: 95.0}
	cdf.Points[4] = &histogram.CDFPoint{Percentile: 0.99, Value: 99.0}
	
	// Convert CDF to histogram for analysis
	cdfHist := cdf.Histogram()
	fmt.Printf("CDF histogram count: %d\n", cdfHist.Count)
	
	// Example 4: Working with buckets and sub-histograms
	fmt.Println("\n=== Bucket Analysis Example ===")
	
	// Create histogram with specific bucket configuration
	bucketHist := histogram.NewHistogram(200, 20.0, 1)
	bucketHist.AddPercentilePoint(0.95)
	
	// Add data with different ranges
	for i := 0; i < 100; i++ {
		// Values in different ranges
		if i < 30 {
			bucketHist.Enqueue(rand.Float64()*50, 1)     // 0-50 range
		} else if i < 60 {
			bucketHist.Enqueue(50+rand.Float64()*100, 1) // 50-150 range
		} else {
			bucketHist.Enqueue(150+rand.Float64()*200, 1) // 150-350 range
		}
	}
	
	fmt.Printf("Number of sub-histograms: %d\n", bucketHist.GetLengthOfSubHistograms())
	fmt.Printf("Maximum buckets per sub-histogram: %d\n", bucketHist.GetMaximumSizeOfSubHistograms())
	
	// Example 5: Watermark and queue management
	fmt.Println("\n=== Queue Management Example ===")
	
	// Create a small histogram to demonstrate watermark
	smallHist := histogram.NewHistogram(10, 5.0, 1)
	
	for i := 0; i < 15; i++ {
		smallHist.Enqueue(float64(i*10), 1) // Use larger values
		watermark := smallHist.GetWaterMark()
		fmt.Printf("After adding %d: watermark = %.2f (count=%d, size=%d)\n", 
			i*10, watermark, smallHist.Count, smallHist.QueueSize)
	}
	
	fmt.Println("\n=== Example Complete ===")
} 