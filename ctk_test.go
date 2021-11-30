package ctk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCtk(t *testing.T) {
	Convey("Checking CTK Rigging", t, func() {
		Convey("typical", WithApp(
			TestingWithCtkWindow,
			func(app Application) {
				So(app.Version(), ShouldEqual, "v0.0.0")
				d := app.Display()
				w := d.ActiveWindow()
				So(w, ShouldNotBeNil)
			},
		))
	})
}
