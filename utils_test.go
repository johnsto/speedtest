package speedtest

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_MaximalSumWindow(t *testing.T) {
	Convey("Calculate window sum, size 1", t, func() {
		So(MaximalSumWindow([]int{1, 2, 3, 4, 5}, 1), ShouldEqual, 5)
		So(MaximalSumWindow([]int{3, 4, 5, 1, 2}, 1), ShouldEqual, 5)
		So(MaximalSumWindow([]int{5, 1, 2, 3, 4}, 1), ShouldEqual, 5)
	})

	Convey("Calculate window sum, size 2", t, func() {
		So(MaximalSumWindow([]int{1, 2, 3, 4, 5}, 2), ShouldEqual, 9)
		So(MaximalSumWindow([]int{3, 4, 5, 1, 2}, 2), ShouldEqual, 9)
		So(MaximalSumWindow([]int{5, 1, 2, 3, 4}, 2), ShouldEqual, 7)
	})

	Convey("Calculate window sum, size 5", t, func() {
		So(MaximalSumWindow([]int{1, 2, 3, 4, 5}, 5), ShouldEqual, 15)
		So(MaximalSumWindow([]int{3, 4, 5, 1, 2}, 5), ShouldEqual, 15)
		So(MaximalSumWindow([]int{5, 1, 2, 3, 4}, 5), ShouldEqual, 15)
	})

	Convey("Should restrict window", t, func() {
		So(MaximalSumWindow([]int{1, 2, 3}, 5), ShouldEqual, 6)
	})
}

func Test_MedianSumWindow(t *testing.T) {
	Convey("Calculate window sum, size 1", t, func() {
		So(MedianSumWindow([]int{1, 2, 3, 4, 5}, 1), ShouldEqual, 3)
		So(MedianSumWindow([]int{3, 4, 5, 1, 2}, 1), ShouldEqual, 3)
		So(MedianSumWindow([]int{5, 1, 2, 3, 4}, 1), ShouldEqual, 3)
	})

	Convey("Calculate window sum, size 3", t, func() {
		So(MedianSumWindow([]int{1, 2, 3, 4, 5}, 3), ShouldEqual, 9)
		So(MedianSumWindow([]int{3, 4, 5, 1, 2}, 3), ShouldEqual, 9)
		So(MedianSumWindow([]int{5, 1, 2, 3, 4}, 3), ShouldEqual, 9)
	})
}
