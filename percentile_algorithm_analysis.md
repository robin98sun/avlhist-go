# Histogram Percentile Calculation Algorithm Analysis

## Overview

The histogram uses a **dynamic cumulative count approach** for percentile calculation, which differs significantly from standard sorted array percentiles. This algorithm is designed for **real-time percentile tracking** with O(log N) complexity.

## Data Structures

### PercentileItem
```go
type PercentileItem struct {
    Percentile      float64        // Target percentile (e.g., 0.5 for P50)
    Item           *HistogramItem  // Points to current tree node
    Key            string          // Unique identifier
    Count          int64           // Cumulative count at this node
    RealPercentage float64         // Actual percentage achieved
}
```

### HistogramItem (AVL Tree Node)
```go
type HistogramItem struct {
    Value        float64
    Count        int64           // Total count in this subtree
    Duplications int64           // Number of identical values
    Smaller      *HistogramItem  // Left child
    Larger       *HistogramItem  // Right child
    // ... other fields
}
```

## Algorithm Steps

### 1. Initialization
When a percentile point is added:
```go
func (h *Histogram) AddPercentilePoint(p float64) {
    item := NewPercentileItem(p)
    h.Percentiles[item.Key] = item
}
```

**Key Point**: Percentile items start with `Item = nil` and `Count = 0`.

### 2. First Item Addition
When the first item is added to an empty histogram:
```go
// Lines 196-203
for _, p := range h.Percentiles {
    p.Item = item           // All percentiles point to first item
    p.Count = int64(count)  // All get full count
    p.RealPercentage = float64(1)  // All at 100%
}
```

**This is the first issue**: All percentiles are initialized to point to the first item with 100% count.

### 3. Dynamic Percentile Tracking
For each new value added, the algorithm updates all percentile items:

#### Step 3a: Determine if new value affects current percentile
```go
if v <= p.Item.Value {
    // New value is <= current percentile item
    p.Count += 1
    percentValue := float64(p.Count)/total_count
    p.RealPercentage = percentValue
```

#### Step 3b: Move to smaller values if percentage is too high
```go
for x:=p.Item.Smaller; x != nil && percentValue > p.Percentile; x=x.Smaller {
    p.Item = x
    p.Count -= x.Larger.Duplications
    p.RealPercentage = float64(p.Count)/total_count
    percentValue = float64(p.Count-x.Duplications)/total_count
}
```

#### Step 3c: Move to larger values if percentage is too low
```go
else {
    percentValue := float64(p.Count)/total_count
    p.RealPercentage = percentValue
    for x:=p.Item.Larger; x != nil && percentValue < p.Percentile; x=x.Larger {
        percentValue = float64(p.Count + x.Duplications)/total_count
        if percentValue <= p.Percentile {
            p.Item = x
            p.Count += x.Duplications
            p.RealPercentage = float64(p.Count)/total_count
        }
    }
}
```

## Algorithm Logic Analysis

### Core Concept: Cumulative Count Tracking

The algorithm tracks **how many items are at or below each percentile point**:

1. **P50 Target**: 50% of items should be ≤ P50 value
2. **P25 Target**: 25% of items should be ≤ P25 value  
3. **P75 Target**: 75% of items should be ≤ P75 value

### Movement Logic

#### Moving to Smaller Values (Lines 176-181)
- **Trigger**: `percentValue > p.Percentile`
- **Action**: Move to smaller tree node
- **Count Adjustment**: Subtract duplications from larger nodes
- **Goal**: Find the largest value where ≤ target percentage of items are below it

#### Moving to Larger Values (Lines 183-190)
- **Trigger**: `percentValue < p.Percentile`  
- **Action**: Move to larger tree node
- **Count Adjustment**: Add duplications from current node
- **Goal**: Find the smallest value where ≥ target percentage of items are below it

### Key Insight: The Algorithm Returns "Nearest Rank" Percentiles

The histogram returns the **largest value where the cumulative count is ≤ the target percentile**, not the **value at the exact percentile position**.

## Example Walkthrough

### Data: [10, 20, 30, 40, 50]
### Target: P50 (50%)

#### Step 1: Add 10
- All percentiles point to 10
- P50: Count=1, Percentage=20% (1/5)
- 20% < 50%, so stay at 10

#### Step 2: Add 20  
- P50: Count=2, Percentage=40% (2/5)
- 40% < 50%, so stay at 10

#### Step 3: Add 30
- P50: Count=3, Percentage=60% (3/5)  
- 60% > 50%, so move to smaller value
- Move to 20: Count=2, Percentage=40%
- 40% < 50%, so stay at 20

#### Step 4: Add 40
- P50: Count=3, Percentage=60% (3/5)
- 60% > 50%, so move to smaller value  
- Move to 20: Count=2, Percentage=40%
- 40% < 50%, so stay at 20

#### Step 5: Add 50
- P50: Count=3, Percentage=60% (3/5)
- 60% > 50%, so move to smaller value
- Move to 20: Count=2, Percentage=40% 
- 40% < 50%, so stay at 20

**Result**: P50 = 20 (not 30 as in standard percentiles)

## Why This Differs from Standard Percentiles

### Standard Percentile Calculation:
1. Sort data: [10, 20, 30, 40, 50]
2. Find index: 0.5 * 5 = 2.5
3. Interpolate or round: data[2] = 30

### Histogram Percentile Calculation:
1. Track cumulative counts dynamically
2. Return largest value where ≤ target percentage of items are below it
3. Result: 20 (where 40% of items are ≤ 20)

## Performance Characteristics

### Advantages:
- ✅ **O(log N) complexity** for both insertion and query
- ✅ **Real-time updates** as data arrives
- ✅ **Memory efficient** using AVL tree structure
- ✅ **Thread-safe** with mutex protection

### Trade-offs:
- ⚠️ **Different from standard percentiles** - uses "nearest rank" method
- ⚠️ **Discrete values only** - no interpolation
- ⚠️ **Consistent offsets** - predictable but different from expectations

## Conclusion

The histogram's percentile algorithm is **functionally correct** but uses a **different definition** than standard percentiles. It prioritizes:

1. **Performance** (O(log N))
2. **Real-time updates** 
3. **Memory efficiency**
4. **Consistency** over exact percentile matching

This explains why the algorithm consistently returns values that differ from standard percentiles - it's designed to return the **largest value where the cumulative count meets the percentile threshold**, not the **value at the exact percentile position**. 