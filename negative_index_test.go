package histogram

import (
	"testing"
	"fmt"
)

func TestNegativeIndexBug(t *testing.T) {
	fmt.Println("=== Testing Negative Index Bug ===")
	
	// Create histogram with parameters that can cause negative indices
	hist := NewHistogram(100, 0.1, 3) // Small bucket size with high accuracy
	
	// The issue is in SubBucketHistogram.CalcPosition
	// When v < sb.LowerBoundary, the calculation (v-sb.LowerBoundary)/sb.BucketSize
	// will result in a negative value
	
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Negative index bug confirmed: %v\n", r)
		}
	}()
	
	// Add a value that will cause negative index calculation
	// The bucket boundaries are calculated based on the sub-histogram index
	// If we add a very small value, it might fall into a bucket with lower boundary > value
	
	// First, let's see what happens with a very small value
	hist.Enqueue(0.001, 1)
	
	// Now let's try to force a negative index by adding a value that's smaller than the lower boundary
	// We need to understand the bucket boundaries first
	idx, lower, upper := hist.BucketHistogram.CalcPosition(0.001)
	fmt.Printf("Value 0.001 -> idx: %d, lower: %f, upper: %f\n", idx, lower, upper)
	
	// Try adding a value that's smaller than the calculated lower boundary
	if lower > 0.001 {
		fmt.Printf("Lower boundary (%f) > value (0.001), this should cause issues\n", lower)
		// This should cause a negative index in the sub-bucket calculation
		hist.Enqueue(0.0001, 1)
	}
	
	fmt.Println("No panic occurred - the bug might be intermittent")
} 