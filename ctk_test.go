package ctk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/go-curses/cdk"
)

func TestCtk(t *testing.T) {
	Convey("Checking CDK Rigging", t, func() {
		Convey("typical", cdk.WithApp(
			TestingWithCtkWindow,
			func(app cdk.App) {
				So(app.Version(), ShouldEqual, "v0.0.0")
				d := app.Display()
				w := d.ActiveWindow()
				So(w, ShouldNotBeNil)
			},
		))
	})
}
