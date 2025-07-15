// reference: https://medium.com/rungo/unit-testing-made-easy-in-go-25077669318

package histogram

import (
	// "log"
	"testing"
	"github.com/stretchr/testify/assert"
	"sort"
)


func TestScheduler_InsertHistogramItem(t *testing.T) {
	random_list := SAMPLE_LIST

	root := create_tree(random_list)

	assert.Equal(t, int64(7), root.Height, "height of root should be exactly 7")
	assert.Equal(t, int64(50), root.Count, "count of nodes should be exactly 50")
	assert.Equal(t, float64(1612), root.Value, "value of root should be exactly 1762")

	smallest := root
	for ; smallest.Left != nil; smallest = smallest.Left {}

	largest := root
	for ; largest.Right != nil; largest = largest.Right {}

	assert.Equal(t, smallest.Height, int64(1), "height of smallest should be exactly 1")
	assert.Equal(t, largest.Height, int64(1), "height of largest should be exactly 1")

	assert.Nil(t, smallest.Smaller, "smaller of smallest should be exactly nil")
	assert.Nil(t, largest.Larger, "larger of largest should be exactly nil")


	assert.Equal(t, smallest.Value, float64(134), "value of smallest should be exactly 134")
	assert.Equal(t, largest.Value, float64(18210), "value of largest should be exactly 18210")

	
	sorted_sample_list := sorted_list(SAMPLE_LIST)
	copy(sorted_sample_list, SAMPLE_LIST)
	sort.Float64Slice(sorted_sample_list).Sort()
	p := smallest
	i := 0
	for ; p != nil; p = p.Larger {
		assert.Equal(t, p.Value, sorted_sample_list[i], "list should be sorted")
		i++
	}
}

func TestScheduler_DeleteHistogramItem(t *testing.T) {
	
	delete_in_order(t, "sorted", "asc", 0, true)
	delete_in_order(t, "sorted", "desc", 0, true)
	delete_in_order(t, "original", "asc", 0, true)
	delete_in_order(t, "original", "desc", 0, true)
	delete_in_order(t, "random", "asc", 10, true)
	delete_in_order(t, "random", "desc", 10, true)
	delete_in_order(t, "random", "asc", 500, true)
	delete_in_order(t, "random", "desc", 500, true)
	delete_in_order(t, "random", "asc", 3000, true)
	delete_in_order(t, "random", "desc", 3000, true)
	delete_in_order(t, "random", "desc", 100000, false)
	// delete_in_order(t, "random", "asc", 1000000, false)
	
}

func TestScheduler_RandomInsertDeleteHistogramItem(t *testing.T) {
	randomlyInsertAndDelete(t, 3000, true)
	randomlyInsertAndDelete(t, 10000, true)
	randomlyInsertAndDelete(t, 1000000, true)
}

func TestScheduler_FindNoLargerThan(t *testing.T) {
	is_something_wrong := false
	for round := 0; round < 10; round++ {
		if is_something_wrong {break}
		for size := 10; size <= 100000; size *= 10 {
			if is_something_wrong {break}
			random_sample_list := gen_random_list(size)
			sorted_sample_list := sorted_list(random_sample_list)
			
			root := create_tree(random_sample_list)

			pre := float64(-1)
			for i := 0; i < len(random_sample_list); i++ {
				v := sorted_sample_list[i] 
				if i > 0 && v == pre {continue}
				pre = v

				nlt := root.FindNoLargerThan(v)
				// log.Printf("to find value no larger than %v", v)
				assert.Equal(t, v, nlt.Value, "the no larger value should equal to itself value")
				if nlt == nil || v != nlt.Value {
					is_something_wrong = true
					break
				} else {
					// log.Printf("  ---found value no larger than %v: %v", v, nlt.Value)
				}

				v -= 0.1
				nlt = root.FindNoLargerThan(v)
				// log.Printf("to find value no larger than %v", v)
				if i == 0 || pre < 0 {
					assert.Nil(t, nlt, "it should not have smaller value when it is already the smallest value")
					if nlt != nil {
						is_something_wrong = true
						break
					}
				} else {
					assert.NotNil(t, nlt, "it should have smaller value when it is not the smallest value")
					assert.Equal(t, sorted_sample_list[i-1], nlt.Value, "the smaller value should exist and equal to previous value")
					if nlt == nil || nlt.Value != sorted_sample_list[i-1] {
						is_something_wrong = true
						break
					} else {
						// log.Printf("   ---found value no larger than %v: %v", v, nlt.Value)
					}
				}
			}
		}
	}
	
}

func TestScheduler_CumulativeCount(t *testing.T) {

	is_something_wrong := false
	for round := 0; round < 1; round++ {
		if is_something_wrong {break}
		for size := 10; size <= 1000000; size *= 10 {
			if is_something_wrong {break}
			random_sample_list := gen_random_list(size)
			sorted_sample_list := sorted_list(random_sample_list)
			root := create_tree(random_sample_list)
			pre := float64(-1)
			cc := int64(0)
			for i := 0; i < len(random_sample_list); i++ {
				v := sorted_sample_list[i] 
				if i > 0 && v == pre {continue}
				pre = v

				node := root.Find(v)
				cumulativeCount := node.CumulativeCount()

				cc += node.Duplications

				assert.Equal(t, cc, cumulativeCount, "cumulative count should equal with the sum")	

				if cc != cumulativeCount {
					is_something_wrong = true
					break
				}
				
			}
		}
	}
}



