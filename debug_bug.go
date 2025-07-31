package histogram

import (
	"testing"
	"fmt"
)

func TestDebugPercentileTracking(t *testing.T) {
	// Test with simple data to understand the percentile tracking
	data := []float64{10, 20, 30, 40, 50}
	
	fmt.Printf("=== Debugging Percentile Tracking ===\n")
	fmt.Printf("Data: %v\n", data)
	
	hist := NewHistogram(10, 10.0, 1)
	hist.AddPercentilePoint(0.5) // P50
	
	fmt.Printf("\nInitial state:\n")
	fmt.Printf("  RootItem: %v\n", hist.RootItem)
	fmt.Printf("  Percentiles: %v\n", hist.Percentiles)
	
	// Add data one by one and track percentile behavior
	for i, value := range data {
		fmt.Printf("\n--- Adding value %.1f (step %d) ---\n", value, i+1)
		
		// Before adding
		if hist.RootItem != nil {
			fmt.Printf("  Before: RootCount=%d, MinItem=%.1f, MaxItem=%.1f\n", 
				hist.RootItem.Count, hist.MinItem.Value, hist.MaxItem.Value)
		}
		
		// Add the value
		hist.Enqueue(value, 1)
		
		// After adding
		fmt.Printf("  After: RootCount=%d, MinItem=%.1f, MaxItem=%.1f\n", 
			hist.RootItem.Count, hist.MinItem.Value, hist.MaxItem.Value)
		
		// Check percentile
		p50 := hist.GetValueAtPercentile(0.5)
		fmt.Printf("  P50: %.1f\n", p50)
		
		// Check percentile item details
		pItem := hist.GetPercentileItem(0.5)
		if pItem != nil {
			fmt.Printf("  PercentileItem: Item=%.1f, Count=%d, RealPercentage=%.3f\n", 
				pItem.Item.Value, pItem.Count, pItem.RealPercentage)
		}
	}
	
	// Final analysis
	fmt.Printf("\n=== Final Analysis ===\n")
	fmt.Printf("Expected P50 for data %v: %.1f\n", data, data[2]) // Should be 30
	
	// Check what the histogram thinks the P50 should be
	targetCount := int(0.5 * float64(len(data)))
	fmt.Printf("Target count for P50: %d (out of %d)\n", targetCount, len(data))
	
	// Walk through the tree to understand the cumulative counts
	if hist.RootItem != nil {
		fmt.Printf("\nTree structure:\n")
		printTree(hist.RootItem, 0)
	}
}

func printTree(item *HistogramItem, depth int) {
	if item == nil {
		return
	}
	
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}
	
	fmt.Printf("%sValue=%.1f, Count=%d, Duplications=%d\n", 
		indent, item.Value, item.Count, item.Duplications)
	
	if item.Smaller != nil {
		fmt.Printf("%sSmaller:\n", indent)
		printTree(item.Smaller, depth+1)
	}
	
	if item.Larger != nil {
		fmt.Printf("%sLarger:\n", indent)
		printTree(item.Larger, depth+1)
	}
} 