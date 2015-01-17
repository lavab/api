package routes_test

import (
	"crypto/md5"
	"encoding/hex"
	"testing"

	"github.com/dchest/uniuri"
	"github.com/franela/goreq"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAvatarsRoute(t *testing.T) {
	Convey("While querying avatars generator", t, func() {
		Convey("Default settings PNG avatar using a hash", func() {
			email := uniuri.New() + "@lavaboom.io"
			rawHashedEmail := md5.Sum([]byte(email))
			hashedEmail := hex.EncodeToString(rawHashedEmail[:])

			result, err := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/avatars/" + hashedEmail + ".png",
			}.Do()
			So(err, ShouldBeNil)

			avatarFromHash, err := result.Body.ToString()
			So(err, ShouldBeNil)
			So(avatarFromHash, ShouldNotBeEmpty)

			Convey("A non-hashed avatar should be the same", func() {
				result, err := goreq.Request{
					Method: "GET",
					Uri:    server.URL + "/avatars/" + email + ".png",
				}.Do()
				So(err, ShouldBeNil)

				avatarFromEmail, err := result.Body.ToString()
				So(err, ShouldBeNil)
				So(avatarFromEmail, ShouldEqual, avatarFromHash)
			})
		})

		Convey("A 150px-wide SVG avatar should have a proper size", func() {
			email := uniuri.New() + "@lavaboom.io"

			result, err := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/avatars/" + email + ".svg?width=150",
			}.Do()
			So(err, ShouldBeNil)

			avatar, err := result.Body.ToString()
			So(err, ShouldBeNil)
			So(avatar, ShouldNotBeEmpty)
			So(avatar, ShouldContainSubstring, `width="150"`)
		})

		Convey("Invalid custom width should fail", func() {
			email := uniuri.New() + "@lavaboom.io"

			result, err := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/avatars/" + email + ".svg?width=ayylmao",
			}.Do()
			So(err, ShouldBeNil)

			avatar, err := result.Body.ToString()
			So(err, ShouldBeNil)
			So(avatar, ShouldNotBeEmpty)
			So(avatar, ShouldContainSubstring, "Invalid width")
		})
	})
}
