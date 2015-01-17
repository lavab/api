package routes_test

import (
	"testing"
	"time"

	"github.com/franela/goreq"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/routes"
)

func TestMiddleware(t *testing.T) {
	Convey("While querying a secure endpoint", t, func() {
		Convey("No header should fail", func() {
			result, err := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/accounts/me",
			}.Do()
			So(err, ShouldBeNil)

			var response routes.AuthMiddlewareResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Missing auth token")
		})

		Convey("An invalid header should fail", func() {
			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/accounts/me",
			}
			request.AddHeader("Authorization", "123")
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.AuthMiddlewareResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Invalid authorization header")
		})

		Convey("Invalid token should fail", func() {
			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/accounts/me",
			}
			request.AddHeader("Authorization", "Bearer 123")
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.AuthMiddlewareResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Invalid authorization token")
		})

		Convey("Expired token should fail", func() {
			account := &models.Account{
				Resource: models.MakeResource("", "johnorange"),
				Status:   "complete",
				AltEmail: "john@orange.org",
			}
			err := account.SetPassword("fruityloops")
			So(err, ShouldBeNil)

			err = env.Accounts.Insert(account)
			So(err, ShouldBeNil)

			token := models.Token{
				Resource: models.MakeResource(account.ID, "test invite token"),
				Expiring: models.Expiring{
					ExpiryDate: time.Now().UTC().Truncate(time.Hour * 8),
				},
				Type: "auth",
			}

			err = env.Tokens.Insert(token)
			So(err, ShouldBeNil)

			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/accounts/me",
			}
			request.AddHeader("Authorization", "Bearer "+token.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.AuthMiddlewareResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Authorization token has expired")
		})
	})
}
