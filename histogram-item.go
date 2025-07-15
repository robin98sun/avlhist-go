package histogram

import (
    "log"
    "fmt"
)

var DEBUG bool = false


// Histogram tree 
type HistogramItem struct {
    Value float64
    Left *HistogramItem
    Right *HistogramItem
    Parent *HistogramItem
    Smaller *HistogramItem
    Larger *HistogramItem
    Height int64
    Count int64
    Duplications int64
}


func NewHistogramItem(v float64) *HistogramItem {
    return &HistogramItem{
            Value: v,
            Left: nil,
            Right: nil,
            Parent: nil,
            Smaller: nil,
            Larger: nil,
            Height: 1,
            Count: 1,
            Duplications: 1,
    }
}

func (t *HistogramItem) GetRoot() *HistogramItem {
    root := t
    path := fmt.Sprintf("%v", t.Value)
    for c := t.Parent; c != nil; c = c.Parent {
        if DEBUG {
            path = fmt.Sprintf("%v, %v", path, c.Value)
        }
        root = c
    }
    if DEBUG {
        log.Printf("path to root: %v", path)
    }
    return root
}

func (t *HistogramItem) Find(v float64) *HistogramItem {
    if t.Value == v {
        return t
    } else if v < t.Value && t.Left != nil {
        return t.Left.Find(v)
    } else if v > t.Value && t.Right != nil {
        return t.Right.Find(v)
    }
    return nil
}


func (t *HistogramItem) FindSmallestInRight() *HistogramItem {
    if t.Right == nil {
        return nil
    } else {
        c := t.Right
        for ; c.Left != nil;  c = c.Left {
        }
        return c
    }
}

func (t *HistogramItem) FindLargestInLeft() *HistogramItem {
    if t.Left == nil {
        return nil
    } else {
        c := t.Left
        for ; c.Right != nil;  c = c.Right {
        }
        return c
    }
}

func (t *HistogramItem) FindNoLargerThan(v float64) *HistogramItem {

    if t == nil {
        return nil
    }

    p := t
    if p.Value <= v {
        for ; p.Right != nil && p.Value <= v; p = p.Right {}
        if p.Value <= v {
            // reached the largest node but still not large enough
            return p
        } else {
            // p.Value > v, it must has passed v just one step
            candidate := p.Parent
            if p.Left != nil {
                candidate = p.Left.FindNoLargerThan(v)
                if candidate == nil {
                    candidate = p.Parent
                }
            }
            return candidate
        }
    } else {
        // p.Value > v
        for ; p.Left != nil && p.Value > v; p = p.Left {}
        if p.Value > v {
            // reached the smallest node but still larger than v
            return nil
        } else {
            // p.Value <= v
            candidate := p
            if p.Right != nil {
                candidate = p.Right.FindNoLargerThan(v)
                if candidate == nil {
                    candidate = p
                }
            }
            return candidate
        }
    }
}

func (t *HistogramItem) CumulativeCount() int64 {
    if t == nil {return int64(0)}

    cumulative_count := int64(0)
    cumulative_count += t.Duplications
    if t.Left != nil {
        cumulative_count += t.Left.Count
    }

    if t.Parent == nil {
        return cumulative_count
    }

    for pre, cur := t, t.Parent; cur!=nil; pre,cur=cur,cur.Parent {
        if cur.Left != pre {
            cumulative_count += cur.Duplications
            if cur.Left != nil{
                cumulative_count += cur.Left.Count
            }
        }
    }
    return cumulative_count
}

// return the inserted node,
// and if the root could be changed, then return the new root
//     but if the root is not changed, then return nil
func (t *HistogramItem) Insert(v float64, count int64, recursion_level int) (*HistogramItem, *HistogramItem) {
    if recursion_level > 30 {
        log.Printf("[histogram][insert] recursion level: %v, incoming value: %v, count: %v, histogram item value: %v", recursion_level, v, count, t.Value)
    }
    if v == t.Value {
        t.Duplications += count
        t.Count += count
        for c := t.Parent; c!= nil; c = c.Parent {
            c.Count += count
        }
        return t, nil
    } else if (t.Left == nil && v < t.Value) || ( t.Right == nil && v > t.Value ) {
        newItem := NewHistogramItem(v)
        newItem.Duplications = count
        newItem.Parent = t
        var root *HistogramItem = nil
        if v > t.Value {
            t.Right = newItem
            newItem.Larger = t.Larger
            newItem.Smaller = t
            t.Larger = newItem
            if newItem.Larger != nil {
                newItem.Larger.Smaller = newItem
            }
        } else {
            t.Left = newItem
            newItem.Smaller = t.Smaller
            newItem.Larger = t
            t.Smaller = newItem
            if newItem.Smaller != nil {
                newItem.Smaller.Larger = newItem
            }
        }
        // update count before rotation
        for p := t; p != nil; p = p.Parent {
            p.Count += count
        }
        // update height
        if (t.Left == nil && v > t.Value) || (t.Right == nil && v < t.Value) {
            t.Height += 1
            root = t.UpdateHeight(true)
        }
        return newItem, root
    } else if v < t.Value {
        if t.Left == t || t.Left.Value == t.Value {
            log.Printf("[histogram][insert] WARNING: left child is identical, t.Left == t ? %v", t.Left == t)
            t.Left = nil
            return t.Insert(v, count, recursion_level+1)
        } else {
            if recursion_level > 30 {
                log.Printf("[histogram][insert] recursive insert to left child")
            }
            return t.Left.Insert(v, count, recursion_level+1)
        }
    } else {
        if t.Right == t || t.Right.Value == t.Value {
            log.Printf("[histogram][insert] WARNING: right child is identical, t.Right == t ? %v", t.Right == t)
            t.Right = nil
            return t.Insert(v, count, recursion_level+1)
        } else {
            if recursion_level > 30 {
                log.Printf("[histogram][insert] recursive insert to right child")
            }
            return t.Right.Insert(v, count, recursion_level+1)
        }
    }
}

// return the replacing node,
// and if the root could be changed, then return the new root
//     but if the root is not changed, then return nil
func (t *HistogramItem) Delete() (*HistogramItem, *HistogramItem) {
    if t.Duplications > 1 {
        t.Count -= 1
        t.Duplications -= 1

        for c := t.Parent; c!= nil; c = c.Parent {
            c.Count -= 1
        }
        return t, nil
    }

    if t.Parent == nil && t.Left == nil && t.Right == nil {
        return nil, nil
    }

    affectedNode_height := t.Parent
    affectedNode_count := t.Parent
    var replaced_by *HistogramItem = nil

    if t.Left == nil && t.Right == nil {
        if DEBUG {
            log.Printf("deleting a leaf node: %v", t.Value)
        }
        if t.Parent != nil {
            if t.Parent.Left == t {
                t.Parent.Left = nil
            } else if t.Parent.Right == t {
                t.Parent.Right = nil
            }
        }

        if t.Smaller != nil {
            t.Smaller.Larger = t.Larger
        }
        if t.Larger != nil {
            t.Larger.Smaller = t.Smaller
        }

    } else {
        if DEBUG {
            log.Printf("deleting a non-leaf node: %v", t.Value)
        }
        if DEBUG {
                log.Printf("   before searching for replacing node, the node is: %v", t.Describe())
            }
        if t.Left != nil {
            replaced_by = t.FindLargestInLeft()
            if replaced_by.Parent != t {
                if replaced_by.Left != nil {
                    replaced_by.Left.Parent = replaced_by.Parent
                    replaced_by.Parent.Right = replaced_by.Left
                } else {
                    replaced_by.Parent.Right = nil
                }
            }
        } else if t.Right != nil {
            replaced_by = t.FindSmallestInRight()
            if replaced_by.Parent != t {
                if replaced_by.Right != nil {
                    replaced_by.Right.Parent = replaced_by.Parent
                    replaced_by.Parent.Left = replaced_by.Right
                } else {
                    replaced_by.Parent.Left = nil
                }
            }
        }
        if replaced_by != nil && replaced_by.Parent != t {
            for c:= replaced_by.Parent; c != nil && c != t; c = c.Parent {
                c.Count -= replaced_by.Duplications
            }
        }
        if replaced_by != nil {
            affectedNode_height = replaced_by.Parent
            if replaced_by.Parent == t {
                affectedNode_height = replaced_by
            }
            
            if DEBUG {
                log.Printf("   the non-leaf node is going to be replaced by %v", replaced_by.Value)
                log.Printf("   before replacing, the node is: %v", t.Describe())
            }
            // update stats
            replaced_by.Count = t.Count - t.Duplications
            replaced_by.Height = t.Height

            // update pointers
            // parent
            replaced_by.Parent = t.Parent
            if t.Parent != nil && t.Parent.Left == t {
                t.Parent.Left = replaced_by
            } else if t.Parent != nil && t.Parent.Right == t {
                t.Parent.Right = replaced_by
            }

            // left
            if replaced_by != t.Left {
                replaced_by.Left = t.Left
                if t.Left != nil {
                    t.Left.Parent = replaced_by
                }
            }

            // right
            if replaced_by != t.Right {
                replaced_by.Right = t.Right
                if t.Right != nil {
                    t.Right.Parent = replaced_by
                }
            }

            // smaller
            if replaced_by != t.Smaller {
                replaced_by.Smaller = t.Smaller
            }
            if t.Smaller != nil && t.Smaller != replaced_by {
                t.Smaller.Larger = replaced_by
            }

            // larger
            if replaced_by != t.Larger {
                replaced_by.Larger = t.Larger
            }
            if t.Larger != nil && t.Larger != replaced_by {
                t.Larger.Smaller = replaced_by
            }

            if DEBUG {
                log.Printf("   after replacing, the replacing node is: %v", replaced_by.Describe())
            }
        } 
    }

    // nil t's pointers
    t.Right = nil
    t.Left = nil
    t.Smaller = nil
    t.Larger = nil
    t.Parent = nil

    // Update Count
    for p := affectedNode_count; p != nil; p = p.Parent {
        if DEBUG {
            log.Printf("   updating count for node[%v]", p.Value)
        }
        p.Count -= t.Duplications
    }

    root := affectedNode_height.UpdateHeight(false)

    return replaced_by, root
}


func (t *HistogramItem) CalcHeight() (int64, int64, int64) {
    leftHeight := int64(0)
    rightHeight := int64(0)
    if t.Left != nil {
        leftHeight = t.Left.Height 
    }
    if t.Right != nil {
        rightHeight = t.Right.Height
    }

    t.Height = leftHeight + 1
    if rightHeight > leftHeight {
        t.Height = rightHeight + 1
    }
    return t.Height, leftHeight, rightHeight
}

func (t *HistogramItem) UpdateHeight(isInserting bool) *HistogramItem {
    if DEBUG {
        log.Printf("updating height for node: %v", t.Describe())
    }
    root := t
    for c := t; c != nil; c = c.Parent {
        _, leftHeight, rightHeight := c.CalcHeight()

        if DEBUG {
            log.Printf("   for node[%v], left height: %v, right height: %v",
                         c.Value, leftHeight, rightHeight,
                )
        }
        if leftHeight - rightHeight > 1 {
            if DEBUG {
                log.Printf("      before right rotation, node[%v]: %v", c.Value, c.Describe())
            }

            if isInserting {
               if c.Left.Right != nil {
                    c.Left.LeftRotate()
                    if DEBUG {
                        log.Printf("      after left rotation, node[%v]: %v", c.Value, c.Describe())
                    }
                } 
            }
            
            c = c.RightRotate()
            if DEBUG {
                log.Printf("      after right rotation, node[%v]: %v", c.Value, c.Describe())
            }
        } else if rightHeight - leftHeight > 1 {
            if DEBUG {
                log.Printf("      before left rotation, node[%v]: %v", c.Value, c.Describe())
            }

            if isInserting {
                if c.Right.Left != nil {
                    c.Right.RightRotate()
                    if DEBUG {
                        log.Printf("      after right rotation, node[%v]: %v", c.Value, c.Describe())
                    }
                }
            }

            c = c.LeftRotate()
            if DEBUG {
                log.Printf("      after left rotation, node[%v]: %v", c.Value, c.Describe())
            }
        }

        root = c
    }
    return root
}

func (t *HistogramItem) LeftRotate() *HistogramItem{

    if DEBUG {
        log.Printf("   Left rotate node: %v", t.Value)
    }

    if t.Right == nil {
        return t
    }

    p := t.Right

    // prepare to rotate
    t.Count -= p.Count
    p.Count += t.Count

    // take p's left
    t.Right = p.Left
    if p.Left != nil {
        t.Count += p.Left.Count
        p.Left.Parent = t
    }

    p.Parent = t.Parent
    if t.Parent != nil {
        if t.Parent.Left == t {
            t.Parent.Left = p
        } else if t.Parent.Right == t {
            t.Parent.Right = p
        }
    }
    t.Parent = p
    p.Left = t

    t.CalcHeight()
    p.CalcHeight()

    return p
}

func (t *HistogramItem) RightRotate() *HistogramItem{
    if DEBUG {
        log.Printf("   Right rotate node: %v", t.Value)
    }

    if t.Left == nil {
        return t
    }

    p := t.Left

    // prepare to rotate
    t.Count -= p.Count
    p.Count += t.Count

    // take p's right
    t.Left = p.Right
    if p.Right != nil {
        t.Count += p.Right.Count
        p.Right.Parent = t
    }

    p.Parent = t.Parent
    if t.Parent != nil {
        if t.Parent.Left == t {
            t.Parent.Left = p
        } else if t.Parent.Right == t {
            t.Parent.Right = p
        }
    }
    t.Parent = p
    p.Right = t


    t.CalcHeight()
    p.CalcHeight()

    return p
}

func (t *HistogramItem) Describe() string {
    desc := fmt.Sprintf("value: %v, height: %v, count: %v", t.Value, t.Height, t.Count)
    left_desc := "nil"
    if t.Left != nil {
        left_desc = t.Left.Describe()
    }
    right_desc := "nil"
    if t.Right != nil {
        right_desc = t.Right.Describe()
    }
    desc = fmt.Sprintf("%v, left: [%v], right: [%v]", desc, left_desc, right_desc)
    return desc
}
