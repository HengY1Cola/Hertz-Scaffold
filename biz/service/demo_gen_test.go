package service

import (
	"Hertz-Scaffold/biz/utils/env"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestArticleService_GetArticleList(t *testing.T) {
	Convey("test GetArticleList", t, func() {
		ctx := env.GetMockCtx()
		res, err := GetDemoService().GetString(ctx, "1")
		So(err, ShouldBeNil)
		So(res, ShouldEqual, "1")
	})
}
