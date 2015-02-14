package routes_test

import (
	"testing"

	"github.com/dchest/uniuri"
	"github.com/franela/goreq"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/routes"
)

func TestEmailsRoute(t *testing.T) {
	Convey("Given a working account", t, func() {
		account := &models.Account{
			Resource: models.MakeResource("", "johnorange3"),
			Status:   "complete",
			AltEmail: "john3@orange.org",
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
				"username": "johnorange3",
				"password": "fruityloops"
			}`,
		}.Do()
		So(err, ShouldBeNil)

		err = env.Labels.Insert([]*models.Label{
			&models.Label{
				Resource: models.MakeResource(account.ID, "Inbox"),
				Builtin:  true,
			},
			&models.Label{
				Resource: models.MakeResource(account.ID, "Sent"),
				Builtin:  true,
			},
			&models.Label{
				Resource: models.MakeResource(account.ID, "Trash"),
				Builtin:  true,
			},
			&models.Label{
				Resource: models.MakeResource(account.ID, "Spam"),
				Builtin:  true,
			},
			&models.Label{
				Resource: models.MakeResource(account.ID, "Starred"),
				Builtin:  true,
			},
		})
		So(err, ShouldBeNil)

		var response routes.TokensCreateResponse
		err = result.Body.FromJsonTo(&response)
		So(err, ShouldBeNil)

		So(response.Success, ShouldBeTrue)
		authToken := response.Token

		Convey("Creating a new email using invalid JSON input should fail", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/emails",
				ContentType: "application/json",
				Body:        "!@#!@#!@#",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.EmailsCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Invalid input format")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Creating a new email with missing data should fail", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/emails",
				ContentType: "application/json",
				Body: routes.EmailsCreateRequest{
					To:      []string{"piotr@zduniak.net"},
					Subject: "hello world",
				},
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.EmailsCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Invalid request")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Getting a non-existing email should fail", func() {
			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/emails/nonexisting",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.EmailsGetResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Email not found")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Getting a non-owned email should fail", func() {
			email := &models.Email{
				Resource: models.MakeResource("not", uniuri.New()),
			}

			err := env.Emails.Insert(email)
			So(err, ShouldBeNil)

			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/emails/" + email.ID,
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.EmailsGetResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Email not found")
			So(response.Success, ShouldBeFalse)

			Convey("Deleting it should fail", func() {
				request := goreq.Request{
					Method: "DELETE",
					Uri:    server.URL + "/emails/" + email.ID,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.EmailsDeleteResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Success, ShouldBeFalse)
				So(response.Message, ShouldEqual, "Email not found")
			})
		})

		Convey("Listing emails with invalid offset should fail", func() {
			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/emails?offset=pppp&limit=1&sort=+date",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.EmailsListResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Invalid offset")
		})

		Convey("Listing emails with invalid limit should fail", func() {
			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/emails?offset=0&limit=pppp&sort=+date",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.EmailsListResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Invalid limit")
		})

		Convey("Delete a non-existing email should fail", func() {
			request := goreq.Request{
				Method: "DELETE",
				Uri:    server.URL + "/emails/nonexisting",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.EmailsDeleteResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Email not found")
		})

		Convey("Creating a new email should succeed", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/emails",
				ContentType: "application/json",
				Body: routes.EmailsCreateRequest{
					To:      []string{"test@lavaboom.io"},
					Subject: "hello world",
					Body:    "raw meaty email",
				},
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.EmailsCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldBeEmpty)
			So(len(response.Created), ShouldBeGreaterThan, 0)
			So(response.Success, ShouldBeTrue)

			emailID := response.Created[0]

			Convey("Getting that email should succeed", func() {
				request := goreq.Request{
					Method: "GET",
					Uri:    server.URL + "/emails/" + emailID,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.EmailsGetResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Success, ShouldBeTrue)
				So(response.Email.Name, ShouldEqual, "hello world")
			})

			Convey("That email should be visible on the emails list", func() {
				request := goreq.Request{
					Method: "GET",
					Uri:    server.URL + "/emails?offset=0&limit=1&sort=+date",
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.EmailsListResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Success, ShouldBeTrue)
				So((*response.Emails)[0].Name, ShouldEqual, "hello world")
			})

			Convey("Deleting that email should succeed", func() {
				request := goreq.Request{
					Method: "DELETE",
					Uri:    server.URL + "/emails/" + emailID,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.EmailsDeleteResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Success, ShouldBeTrue)
				So(response.Message, ShouldEqual, "Email successfully removed")
			})
		})
	})
}
