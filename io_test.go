package speedtest

import (
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"testing"
)

func Test_JunkReader(t *testing.T) {
	Convey("JunkReader should return bytes read", t, func() {
		jr := NewJunkReader(10)
		buf := make([]byte, 8)

		num, err := jr.Read(buf)
		So(num, ShouldEqual, 8)
		So(err, ShouldEqual, nil)

		num, err = jr.Read(buf)
		So(num, ShouldEqual, 2)
		So(err, ShouldEqual, io.EOF)
	})
}

func Test_CallbackWriter(t *testing.T) {
	Convey("CallbackWriter should call the callback", t, func() {
		callbackNum := -1
		callbackErr := io.EOF
		w := NewCallbackWriter(func(num int) error {
			callbackNum = num
			return callbackErr
		})

		buf := make([]byte, 8)
		num, err := w.Write(buf)
		So(callbackNum, ShouldEqual, 8)
		So(num, ShouldEqual, 8)
		So(err, ShouldEqual, io.EOF)
	})
}
