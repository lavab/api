package routes_test

import (
	"strings"
	"testing"

	"github.com/franela/goreq"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/routes"
)

func TestKeysRoute(t *testing.T) {
	Convey("Given a working account", t, func() {
		account := &models.Account{
			Resource: models.MakeResource("", "johnorange4"),
			Status:   "complete",
			AltEmail: "john4@orange.org",
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
				"username": "johnorange4",
				"password": "fruityloops"
			}`,
		}.Do()
		So(err, ShouldBeNil)

		var response routes.TokensCreateResponse
		err = result.Body.FromJsonTo(&response)
		So(err, ShouldBeNil)

		So(response.Success, ShouldBeTrue)
		authToken := response.Token

		Convey("Uploading a new key using an invalid JSON format should fail", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/keys",
				ContentType: "application/json",
				Body:        "!@#!@!@#",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.KeysCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Invalid input format")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Uploading an invalid key should fail", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/keys",
				ContentType: "application/json",
				Body: `{
					"key": "hbnjmvnbhvm nbhm jhbjmnghnbgjvgbhvf bgvmj gvhnft"
				}`,
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.KeysCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Invalid key format")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Uploading a key should succeed", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/keys",
				ContentType: "application/json",
				Body: `{
			"key": "` + strings.Join(strings.Split(`-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v1

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
=r/yh
-----END PGP PUBLIC KEY BLOCK-----`, "\n"), "\\n") + `"
		}`,
			}

			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.KeysCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "A new key has been successfully inserted")
			So(response.Success, ShouldBeTrue)
			So(response.Key.ID, ShouldNotBeNil)

			key := response.Key

			Convey("Key should be visible on the user's key list", func() {
				request := goreq.Request{
					Method: "GET",
					Uri:    server.URL + "/keys?user=johnorange4",
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.KeysListResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Success, ShouldBeTrue)
				So(len(*response.Keys), ShouldBeGreaterThan, 0)
			})

			Convey("Getting that key should succeed", func() {
				request := goreq.Request{
					Method: "GET",
					Uri:    server.URL + "/keys/" + key.ID,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.KeysGetResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Success, ShouldBeTrue)
				So(response.Key.ID, ShouldEqual, key.ID)
			})
		})

		Convey("Listing keys without passing a username should fail", func() {
			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/keys",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.KeysListResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Invalid username")
		})

		Convey("Getting a non-existing key should fail", func() {
			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/keys/123",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.KeysGetResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Requested key does not exist on our server")
		})
	})
}
