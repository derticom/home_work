package hw04lrucache

import (
	"testing"

	//nolint:depguard // Применение 'require' необходимо для тестирования.
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("push test", func(t *testing.T) {
		l := NewList()

		l.PushFront(1) // [1]
		l.PushBack(2)  // [1,2]
		l.PushFront(3) // [3,1,2]
		l.PushBack(4)  // [3,1,2,4]
		l.PushFront(5) // [5,3,1,2,4]
		l.PushBack(6)  // [5,3,1,2,4,6]
		l.PushFront(7) // [7,5,3,1,2,4,6]
		require.Equal(t, 7, l.Len())

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{7, 5, 3, 1, 2, 4, 6}, elems)
	})

	t.Run("remove test", func(t *testing.T) {
		l := NewList()

		l.PushFront(1) // [1]
		l.PushBack(2)  // [1,2]
		l.PushFront(3) // [3,1,2]
		l.PushBack(4)  // [3,1,2,4]
		l.PushFront(5) // [5,3,1,2,4]
		l.PushBack(6)  // [5,3,1,2,4,6]
		l.PushFront(7) // [7,5,3,1,2,4,6]

		l.Remove(l.Front().Next) // [7,3,1,2,4,6]
		l.Remove(l.Back())       // [7,3,1,2,4]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{7, 3, 1, 2, 4}, elems)
	})

	t.Run("remove single element test", func(t *testing.T) {
		l := NewList()

		l.PushFront(1) // [1]

		l.Remove(l.Front()) // []

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{}, elems)
	})

	t.Run("move test", func(t *testing.T) {
		l := NewList()

		l.PushFront(1) // [1]
		l.PushBack(2)  // [1,2]
		l.PushFront(3) // [3,1,2]
		l.PushBack(4)  // [3,1,2,4]
		l.PushFront(5) // [5,3,1,2,4]
		l.PushBack(6)  // [5,3,1,2,4,6]
		l.PushFront(7) // [7,5,3,1,2,4,6]

		l.MoveToFront(l.Back().Prev) // [4,7,5,3,1,2,6]
		l.MoveToFront(l.Back())      // [6,4,7,5,3,1,2]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{6, 4, 7, 5, 3, 1, 2}, elems)
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
