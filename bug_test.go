package histogram

import (
	"testing"
	"fmt"
)

// Bug 1: Negative index in bucket calculation
func TestBug1_NegativeIndex(t *testing.T) {
	fmt.Println("=== Testing Bug 1: Negative Index in Bucket Calculation ===")
	
	// Create histogram with small bucket size that can cause negative indices
	hist := NewHistogram(100, 0.1, 3) // Very small bucket size with high accuracy
	
	// This should cause a negative index when values are very small
	hist.Enqueue(0.01, 1) // This value with bucket size 0.1 and accuracy 3 can cause issues
	
	fmt.Println("Bug 1 test completed - check for negative index panic")
}

// Bug 2: CDF nil pointer dereference
func TestBug2_CDFNilPointer(t *testing.T) {
	fmt.Println("=== Testing Bug 2: CDF Nil Pointer Dereference ===")
	
	cdf := NewCDF(100)
	// The CDF.Points array is initialized but all elements are nil
	// When Histogram() tries to iterate over Points, it will encounter nil pointers
	
	// This should panic with nil pointer dereference
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Bug 2 confirmed: %v\n", r)
		}
	}()
	
	hist := cdf.Histogram() // This will panic
	fmt.Printf("CDF histogram created: %v\n", hist)
}

// Bug 3: Potential division by zero in variance calculation
func TestBug3_DivisionByZero(t *testing.T) {
	fmt.Println("=== Testing Bug 3: Division by Zero in Variance ===")
	
	hist := NewHistogram(10, 1.0, 1)
	
	// Add a single item
	hist.Enqueue(1.0, 1)
	
	// Now dequeue all items
	for i := 0; i < 10; i++ {
		hist.Dequeue()
	}
	
	// At this point, hist.Count should be 0, which could cause division by zero
	// in variance calculations
	
	fmt.Printf("After dequeuing all items - Count: %d, Mean: %f, Variance: %f\n", 
		hist.Count, hist.Mean, hist.Variance)
}

// Bug 4: Potential race condition in concurrent access
func TestBug4_RaceCondition(t *testing.T) {
	fmt.Println("=== Testing Bug 4: Race Condition ===")
	
	hist := NewHistogram(100, 1.0, 1)
	
	// Simulate concurrent access
	done := make(chan bool)
	
	// Goroutine 1: Add items
	go func() {
		for i := 0; i < 50; i++ {
			hist.Enqueue(float64(i), 1)
		}
		done <- true
	}()
	
	// Goroutine 2: Query percentiles
	go func() {
		for i := 0; i < 50; i++ {
			hist.GetValueAtPercentile(0.95)
		}
		done <- true
	}()
	
	// Wait for both goroutines
	<-done
	<-done
	
	fmt.Println("Bug 4 test completed - check for race conditions")
}

// Bug 5: Memory leak in bucket list expansion
func TestBug5_MemoryLeak(t *testing.T) {
	fmt.Println("=== Testing Bug 5: Memory Leak in Bucket List ===")
	
	hist := NewHistogram(1000, 1.0, 1)
	
	// Add items with very large values to force bucket list expansion
	for i := 0; i < 100; i++ {
		hist.Enqueue(float64(i*1000), 1) // Large values
	}
	
	// Now add items with very small values
	for i := 0; i < 100; i++ {
		hist.Enqueue(float64(i), 1) // Small values
	}
	
	// The bucket list may have grown very large due to large values
	// but small values may not use all the allocated space
	
	fmt.Printf("Bucket histogram sub-histograms: %d\n", len(hist.BucketHistogram.SubBucketHistograms))
	
	// Check if there are many empty buckets
	emptyBuckets := 0
	totalBuckets := 0
	for _, sbh := range hist.BucketHistogram.SubBucketHistograms {
		if sbh != nil {
			for _, item := range sbh.BucketList {
				totalBuckets++
				if item == nil {
					emptyBuckets++
				}
			}
		}
	}
	
	fmt.Printf("Empty buckets: %d/%d (%.1f%%)\n", emptyBuckets, totalBuckets, 
		float64(emptyBuckets)/float64(totalBuckets)*100)
}

func TestAllBugs(t *testing.T) {
	TestBug1_NegativeIndex(t)
	TestBug2_CDFNilPointer(t)
	TestBug3_DivisionByZero(t)
	TestBug4_RaceCondition(t)
	TestBug5_MemoryLeak(t)
} 