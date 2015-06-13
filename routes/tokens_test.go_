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

func TestTokensRoute(t *testing.T) {
	Convey("Given a working account", t, func() {
		account := &models.Account{
			Resource: models.MakeResource("", "johnorange5"),
			Status:   "complete",
			AltEmail: "john5@orange.org",
		}
		err := account.SetPassword("fruityloops")
		So(err, ShouldBeNil)

		err = env.Accounts.Insert(account)
		So(err, ShouldBeNil)

		result, err := goreq.Request{
			Method:      "POST",
			Uri:         server.URL + "/tokens",
			ContentType: "application/json",
			Body: `{
				"type": "auth",
				"username": "johnorange5",
				"password": "fruityloops"
			}`,
		}.Do()
		So(err, ShouldBeNil)

		var response routes.TokensCreateResponse
		err = result.Body.FromJsonTo(&response)
		So(err, ShouldBeNil)

		So(response.Success, ShouldBeTrue)
		authToken := response.Token

		Convey("Creating a non-auth token should fail", func() {
			request, err := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/tokens",
				ContentType: "application/json",
				Body: `{
					"type": "not-auth"
				}`,
			}.Do()
			So(err, ShouldBeNil)

			var response routes.TokensCreateResponse
			err = request.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Only auth tokens are implemented")
		})

		Convey("Trying to sign in using wrong username should fail", func() {
			request, err := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/tokens",
				ContentType: "application/json",
				Body: `{
					"type": "auth",
					"username": "not-johnorange5",
					"password": "fruityloops"
				}`,
			}.Do()
			So(err, ShouldBeNil)

			var response routes.TokensCreateResponse
			err = request.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Wrong username or password")
		})

		Convey("Trying to sign in using wrong password should fail", func() {
			request, err := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/tokens",
				ContentType: "application/json",
				Body: `{
					"type": "auth",
					"username": "johnorange5",
					"password": "not-fruityloops"
				}`,
			}.Do()
			So(err, ShouldBeNil)

			var response routes.TokensCreateResponse
			err = request.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Wrong username or password")
		})

		Convey("Trying to sign in using an invalid JSON input should fail", func() {
			request, err := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/tokens",
				ContentType: "application/json",
				Body:        "123123123###434$#$",
			}.Do()
			So(err, ShouldBeNil)

			var response routes.TokensCreateResponse
			err = request.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Invalid input format")
		})

		Convey("Getting the currently used token should succeed", func() {
			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/tokens",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.TokensGetResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeTrue)
			So(response.Token.ExpiryDate.After(time.Now().UTC()), ShouldBeTrue)
		})

		Convey("Deleting the token by ID should succeed", func() {
			request := goreq.Request{
				Method: "DELETE",
				Uri:    server.URL + "/tokens/" + authToken.ID,
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.TokensDeleteResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeTrue)
			So(response.Message, ShouldEqual, "Successfully logged out")
		})

		Convey("Deleting a non-existing token should fail", func() {
			request := goreq.Request{
				Method: "DELETE",
				Uri:    server.URL + "/tokens/123",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.TokensDeleteResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Invalid token ID")
		})

		Convey("Deleting current token should succeed", func() {
			request := goreq.Request{
				Method: "DELETE",
				Uri:    server.URL + "/tokens",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.TokensDeleteResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeTrue)
			So(response.Message, ShouldEqual, "Successfully logged out")
		})
	})
}
