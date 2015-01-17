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

func TestLabelsRoute(t *testing.T) {
	Convey("Given a working account", t, func() {
		account := &models.Account{
			Resource: models.MakeResource("", "johnorange2"),
			Status:   "complete",
			AltEmail: "john2@orange.org",
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
				"username": "johnorange2",
				"password": "fruityloops"
			}`,
		}.Do()
		So(err, ShouldBeNil)

		var response routes.TokensCreateResponse
		err = result.Body.FromJsonTo(&response)
		So(err, ShouldBeNil)

		So(response.Success, ShouldBeTrue)
		authToken := response.Token

		Convey("Creating a new label using invalid input should fail", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/labels",
				ContentType: "application/json",
				Body:        "!@#!@!@#",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.LabelsCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Invalid input format")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Updating a label using invalid input should fail", func() {
			request := goreq.Request{
				Method:      "PUT",
				Uri:         server.URL + "/labels/anything",
				ContentType: "application/json",
				Body:        "!@#!@!@#",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.LabelsUpdateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Invalid input format")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Creating a new label without enough information should fail", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/labels",
				ContentType: "application/json",
				Body:        "{}",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.LabelsCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Invalid request")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Getting a non-existing label should fail", func() {
			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/labels/nonexisting",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.LabelsCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Label not found")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Updating a non-existing label should fail", func() {
			request := goreq.Request{
				Method:      "PUT",
				Uri:         server.URL + "/labels/anything",
				ContentType: "application/json",
				Body:        "{}",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.LabelsUpdateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Label not found")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Deleting a non-existing label should fail", func() {
			request := goreq.Request{
				Method: "DELETE",
				Uri:    server.URL + "/labels/anything",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.LabelsDeleteResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Label not found")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Getting a non-owned label should fail", func() {
			label := &models.Label{
				Resource: models.MakeResource("not", uniuri.New()),
				Builtin:  false,
			}

			err := env.Labels.Insert(label)
			So(err, ShouldBeNil)

			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/labels/" + label.ID,
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.LabelsGetResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Label not found")
			So(response.Success, ShouldBeFalse)

			Convey("Updating it should fail", func() {
				request := goreq.Request{
					Method:      "PUT",
					Uri:         server.URL + "/labels/" + label.ID,
					ContentType: "application/json",
					Body:        "{}",
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.LabelsUpdateResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Message, ShouldEqual, "Label not found")
				So(response.Success, ShouldBeFalse)
			})

			Convey("Deleting it should fail", func() {
				request := goreq.Request{
					Method: "DELETE",
					Uri:    server.URL + "/labels/" + label.ID,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.LabelsDeleteResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Message, ShouldEqual, "Label not found")
				So(response.Success, ShouldBeFalse)
			})
		})

		Convey("Creating a label should succeed", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/labels",
				ContentType: "application/json",
				Body: `{
					"name": "` + uniuri.New() + `"
				}`,
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.LabelsCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Label.Name, ShouldNotBeEmpty)
			So(response.Success, ShouldBeTrue)

			label := response.Label

			Convey("That label should be visible on the labels list", func() {
				request := goreq.Request{
					Method: "GET",
					Uri:    server.URL + "/labels",
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.LabelsListResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Message, ShouldBeEmpty)
				So(len(*response.Labels), ShouldBeGreaterThan, 0)
				So(response.Success, ShouldBeTrue)
			})

			Convey("Getting that label should succeed", func() {
				request := goreq.Request{
					Method: "GET",
					Uri:    server.URL + "/labels/" + label.ID,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.LabelsGetResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Label.ID, ShouldEqual, label.ID)
				So(response.Success, ShouldBeTrue)
			})

			Convey("Updating that label should succeed", func() {
				request := goreq.Request{
					Method:      "PUT",
					Uri:         server.URL + "/labels/" + label.ID,
					ContentType: "application/json",
					Body: `{
						"name": "test123"
					}`,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.LabelsUpdateResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Label.Name, ShouldEqual, "test123")
				So(response.Success, ShouldBeTrue)
			})

			Convey("Deleting that label should succeed", func() {
				request := goreq.Request{
					Method: "DELETE",
					Uri:    server.URL + "/labels/" + label.ID,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.LabelsDeleteResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Message, ShouldEqual, "Label successfully removed")
				So(response.Success, ShouldBeTrue)
			})
		})
	})
}
