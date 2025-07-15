package histogram

import (
    "log"
    // "fmt"
    "math"
)

// Histogram bucket
type SubBucketHistogram struct {
    BucketList []*HistogramItem
    BucketSize float64
    LowerBoundary float64
    UpperBoundary float64
}

func NewSubBucketHistogram(unit float64, lower float64, upper float64) *SubBucketHistogram{
    var item *SubBucketHistogram = &SubBucketHistogram{
        BucketSize: unit,
        LowerBoundary: lower,
        UpperBoundary: upper,
    }
    return item
}

func (sb *SubBucketHistogram) CalcPosition(v float64) int64 {
    if sb.BucketSize == 0 {
        return int64(-1)
    }
    idx_value := (v-sb.LowerBoundary)/sb.BucketSize
    if math.Abs(sb.BucketSize) < 1 {
        idx_value = math.Round(idx_value)
    }
    idx := int64(idx_value)

    // log.Printf("     ....v: %v, bucket size: %v, lower boundary: %v, idx_value: %v, idx: %v", v, sb.BucketSize, sb.LowerBoundary, idx_value, idx)
    return idx
}

func (sb *SubBucketHistogram) Insert(n *HistogramItem) {

    if n == nil {return}

    if sb.BucketList == nil {
        sb.BucketList = []*HistogramItem{}
    }

    idx := sb.CalcPosition(n.Value)

    cur_len := int64(len(sb.BucketList))
    for i:=int64(0);i<idx-cur_len+1;i++{
        sb.BucketList = append(sb.BucketList, nil)
    }
    if sb.BucketList[idx] != nil {
        log.Printf("bucket is not empty for %v, there exists %v", n.Value, sb.BucketList[idx].Value)
        log.Printf("    bucket list length: %v, index: %v", len(sb.BucketList), idx)
    }
    sb.BucketList[idx] = n
}

func (sb *SubBucketHistogram) Delete(n *HistogramItem) {
    if n == nil {return}
    idx := sb.CalcPosition(n.Value)
    if idx < int64(len(sb.BucketList)) {
        sb.BucketList[idx] = nil
    }
}

// Top level bucket histogram

type BucketHistogram struct {
    SubBucketHistograms []*SubBucketHistogram
    SubBucketHistogramSize float64
    BucketSize float64
}

func NewBucketHistogram( subhistogramSize float64, bucketSize float64) *BucketHistogram{
    var buckets *BucketHistogram = &BucketHistogram{
        BucketSize: bucketSize,
        SubBucketHistogramSize: subhistogramSize,
    }
    return buckets
}

func (b *BucketHistogram) CalcPosition(v float64) (int64, float64, float64) {
    if b.SubBucketHistogramSize == 0 {
        return int64(-1), float64(-1), float64(-1)
    }
    idx_value := v/b.SubBucketHistogramSize
    if math.Abs(b.SubBucketHistogramSize) < 1 {
        idx_value = math.Round(idx_value)
    }
    idx := int64(idx_value)
    lower, upper := b.GetLowerAndUpperBoundaries(idx)
    return idx, lower, upper
}

func (b *BucketHistogram) GetLowerAndUpperBoundaries(idx int64) (float64, float64) {
    lower := float64(idx)*b.SubBucketHistogramSize
    upper := float64(idx+1)*b.SubBucketHistogramSize
    return lower, upper
}

func (b *BucketHistogram) Insert(n *HistogramItem) {
    if n == nil {return}

    if b.SubBucketHistograms == nil {
        b.SubBucketHistograms = []*SubBucketHistogram{}
    } 

    idx, lower, upper := b.CalcPosition(n.Value)
    cur_len := int64(len(b.SubBucketHistograms))
    for i:=int64(0);i<1+idx-cur_len;i++{
        b.SubBucketHistograms = append(b.SubBucketHistograms, nil)
    }

    if b.SubBucketHistograms[idx] == nil {
        b.SubBucketHistograms[idx] = NewSubBucketHistogram(b.BucketSize, lower, upper)
    } 

    // log.Printf("subhistogram list length: %v, index: %v, value: %v, lower: %v, upper: %v", len(b.SubBucketHistograms), idx, n.Value, lower, upper)


    b.SubBucketHistograms[idx].Insert(n)

}

func (b *BucketHistogram) Delete(n *HistogramItem) {
    if n == nil {return}
    idx, _, _ := b.CalcPosition(n.Value)
    if idx < int64(len(b.SubBucketHistograms)) {
        if b.SubBucketHistograms[idx] != nil {
            b.SubBucketHistograms[idx].Delete(n)
        }
    }
}

