package histogram

import (
	// "log"
	"testing"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"math"
	"sort"
)


var SAMPLE_LIST []float64 = []float64{
	134, 3693, 1612, 2033, 1762, 669, 296, 567, 547, 935,
	2842, 3135, 645, 4265, 2267, 1170, 399, 635, 2153, 1836,
	348, 2672, 5318, 1662, 6104, 1057, 2900, 2777, 3715, 9208,
	2231, 387, 1181, 1063, 3092, 478, 2039, 781, 11764, 591, 
	271, 1061, 3182, 1470, 4686, 1077, 1997, 2430, 18210, 2618,
}

func gen_random_list(list_size int) []float64 {
	result := make([]float64, list_size)
	for i := 0; i < list_size; i++ {
		result[i] = math.Round(rand.ExpFloat64() * float64(list_size))
	}
	return result
}

func gen_random_list_float(list_size int, sample_mean int, fraction float64) []float64 {
	result := make([]float64, list_size)
	for i := 0; i < list_size; i++ {
		result[i] = math.Round(rand.ExpFloat64() * float64(sample_mean)*fraction) / fraction
	}
	return result
}

func sorted_list(list []float64) []float64{
	sorted_list := make([]float64, len(list))
	copy(sorted_list, list)
	sort.Float64Slice(sorted_list).Sort()
	return sorted_list
}

func create_tree(list []float64) *HistogramItem {
	var root *HistogramItem
	for _, v := range list {
		if root == nil {
			root = NewHistogramItem(v)
		} else {
			_, newRoot := root.Insert(v, 1, 0)
			if newRoot != nil {
				root = newRoot
			}
		}
	}
	return root
}

func delete_in_order(t *testing.T, list string, order string, size int, enforceOrdering bool) {
	random_list := SAMPLE_LIST

	sorted_sample_list := make([]float64, len(SAMPLE_LIST))
	copy(sorted_sample_list, SAMPLE_LIST)
	sort.Float64Slice(sorted_sample_list).Sort()

	if list == "sorted" {
		random_list = sorted_sample_list
	} else if list == "random" {
		random_list = gen_random_list(size)
	}

	var root *HistogramItem
	var smallest *HistogramItem
	var largest *HistogramItem
	min, max := float64(-1), float64(0)
	for _, v := range random_list {
		// log.Printf("inerting: [%v] %v \n", i, v)
		if root == nil {
			root = NewHistogramItem(v)
		} else {
			_, newRoot := root.Insert(v, 1, 0)
			if newRoot != nil {
				root = newRoot
			}
		}
		// log.Printf("after inserting: %v\n\n", root.Describe())
		if enforceOrdering {
			if min < 0 || v < min {
				min = v
			}
			if v > max {
				max = v
			}
			for smallest = root; smallest.Smaller != nil ; smallest = smallest.Smaller {}
			for largest = root; largest.Larger != nil; largest = largest.Larger{}
			assert.Equal(t, min, smallest.Value, "the smallest should be identical with the min in the list")
			assert.Equal(t, max, largest.Value, "the largest should be identical with the max in the list")
		}
	}

	// log.Printf("\n\nStart deleting [%v, %v, %v]\n\n", list, order, size)
	for k := range random_list {
		
		i := len(random_list) - k - 1
		remaining_count := i
		if order == "asc" {
			i = k
		}
		v := random_list[i]
		// log.Printf("\n\ndeleting: %v, %v \n", i, v)
		// log.Printf("before deleting, root: %v", root.Describe())
		n := root.Find(v)
		assert.NotNil(t, n, "node should not be nil")
		count := n.Duplications
		if n == nil {
			break
		}

		if k == len(random_list) - 1 {
			assert.Equal(t, root.Value, n.Value, "when deleted the last element, there should have no other nodes left")
			assert.Equal(t, root.Count, int64(1), "when deleted the last element, there should have no other nodes left")
			if root.Value != n.Value{
				break
			}
		}

		replaced, newRoot := n.Delete()
		if newRoot != nil || (replaced == nil && newRoot == nil) {
			root = newRoot
		}

		if n == root && replaced == nil {
			root = nil
			assert.Equal(t, len(random_list)-1, k, "when root is nil it should be the last element")
			if k != len(random_list) - 1 {
				break
			}
		}

		if root != nil {
			n = root.Find(v)
			if count == 1 {
				assert.Nil(t, n, "the node should not be searchable anymore after being deleted")
			} else {
				assert.NotNil(t, n, "the node should not be delete if it has multiple duplications")
			}
			if (count == 1 && n != nil) || (count != 1 && n == nil) {
				break
			}
			// log.Printf("after deleting, root: %v", root.Describe())
		}

		if k == len(random_list) - 1 {
			assert.Nil(t, replaced, "when deleted the last element, there should have no other nodes left")
			if replaced != nil{
				break
			}
		} else {
			assert.NotNil(t, root, "root should not be nil")
			assert.Equal(t, int64(remaining_count), root.Count, "the remaining nodes after deleting should decreased by 1")
			if root == nil || root.Count != int64(remaining_count) {
				break
			}
		}

		if root != nil && enforceOrdering {
			for smallest = root; smallest.Smaller != nil; smallest = smallest.Smaller{}
			previous := smallest
			for p := smallest.Larger; p != nil ; p = p.Larger {
				assert.Greater(t, p.Value, previous.Value, "the sequence should be in ascending order")
				previous = p
			}
		}
	}	
}

func randomlyInsertAndDelete(t *testing.T, size int, enforceOrdering bool) {
	var root *HistogramItem = nil

	random_action_list := gen_random_list(size)
	random_value_list := gen_random_list(size)

	action_avg := float64(0)
	for _, v := range random_action_list {
		action_avg += v
	}
	action_avg /= float64(len(random_action_list))

	inserted_value_list := []float64{}

	for i := 0; i<len(random_action_list); i++ {
		action := random_action_list[i]
		valueToInsert := random_value_list[i]
		if action > action_avg {
			if root == nil {
				root = NewHistogramItem(valueToInsert)
			} else {
				_, newRoot := root.Insert(valueToInsert, 1, 0)
				if newRoot != nil {
					root = newRoot
				}
			}
			inserted_value_list = append(inserted_value_list, valueToInsert)
		} else if root != nil && len(inserted_value_list) > 0 {
			idx_delete := rand.Intn(len(inserted_value_list))
			valueToDelete := inserted_value_list[idx_delete]
			n := root.Find(valueToDelete)
			assert.NotNil(t, n, "the node to delete should be guaranteed not nil")
			replaced, newRoot := n.Delete()
			if replaced == nil && newRoot == nil {
				root = nil
			} else if newRoot != nil {
				root = newRoot
			}

			inserted_value_list = append(inserted_value_list[0:idx_delete], inserted_value_list[idx_delete+1:]...)

			if root == nil {
				assert.Equal(t, 0, len(inserted_value_list), "there should not have anymore values left when root is nil")
			}
		}
		if root != nil {
			assert.Equal(t, int64(len(inserted_value_list)), root.Count, "the number of nodes should equal the number of elements in the list")
		}

		if root != nil && enforceOrdering {
			var smallest *HistogramItem
			for smallest = root; smallest.Smaller != nil; smallest = smallest.Smaller{}
			previous := smallest
			for p := smallest.Larger; p != nil ; p = p.Larger {
				assert.Greater(t, p.Value, previous.Value, "the sequence should be in ascending order")
				previous = p
			}
		}
	}
	if root != nil {
		assert.Equal(t, int64(len(inserted_value_list)), root.Count, "if there have some nodes left")
	} else {
		assert.Equal(t, 0, len(inserted_value_list), "there should not have any nodes left")
	}
}

