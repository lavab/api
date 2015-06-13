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

func TestContactsRoute(t *testing.T) {
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

		Convey("Creating a contact with missing parts should fail", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/contacts",
				ContentType: "application/json",
				Body: `{
					"data": "` + uniuri.NewLen(64) + `",
					"encoding": "json",
					"version_major": 1,
					"version_minor": 0,
					"pgp_fingerprints": ["` + uniuri.New() + `"]
				}`,
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.ContactsCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Invalid request")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Creating a contact with invalid input data should fail", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/contacts",
				ContentType: "application/json",
				Body:        "!@#!@#!@#",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.ContactsCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "Invalid input format")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Getting a non-owned contact should fail", func() {
			contact := &models.Contact{
				Encrypted: models.Encrypted{
					Encoding:     "json",
					Data:         uniuri.NewLen(64),
					Schema:       "contact",
					VersionMajor: 1,
					VersionMinor: 0,
				},
				Resource: models.MakeResource("not", uniuri.New()),
			}

			err := env.Contacts.Insert(contact)
			So(err, ShouldBeNil)

			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/contacts/" + contact.ID,
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.ContactsGetResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Contact not found")

			Convey("Update of a not-owned contact should fail", func() {
				request := goreq.Request{
					Method:      "PUT",
					Uri:         server.URL + "/contacts/" + contact.ID,
					ContentType: "application/json",
					Body: `{
						"data": "` + uniuri.NewLen(64) + `",
						"name": "` + uniuri.New() + `",
						"encoding": "xml",
						"version_major": 8,
						"version_minor": 3
					}`,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.ContactsUpdateResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Success, ShouldBeFalse)
				So(response.Message, ShouldEqual, "Contact not found")
			})

			Convey("Deleting it should fail", func() {
				request := goreq.Request{
					Method: "DELETE",
					Uri:    server.URL + "/contacts/" + contact.ID,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.ContactsDeleteResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Success, ShouldBeFalse)
				So(response.Message, ShouldEqual, "Contact not found")
			})
		})

		Convey("Getting a non-existing contact should fail", func() {
			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/contacts/" + uniuri.New(),
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.ContactsGetResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Contact not found")
		})

		Convey("Creating a contact should succeed", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/contacts",
				ContentType: "application/json",
				Body: `{
					"data": "` + uniuri.NewLen(64) + `",
					"name": "` + uniuri.New() + `",
					"encoding": "json",
					"version_major": 1,
					"version_minor": 0,
					"pgp_fingerprints": ["` + uniuri.New() + `"]
				}`,
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.ContactsCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "A new contact was successfully created")
			So(response.Success, ShouldBeTrue)
			So(response.Contact.ID, ShouldNotBeNil)

			contact := response.Contact

			Convey("The contact should be visible on the list", func() {
				request := goreq.Request{
					Method: "GET",
					Uri:    server.URL + "/contacts",
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.ContactsListResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(len(*response.Contacts), ShouldBeGreaterThan, 0)
				So(response.Success, ShouldBeTrue)

				found := false
				for _, c := range *response.Contacts {
					if c.ID == contact.ID {
						found = true
						break
					}
				}

				So(found, ShouldBeTrue)
			})

			Convey("Getting that contact should succeed", func() {
				request := goreq.Request{
					Method: "GET",
					Uri:    server.URL + "/contacts/" + contact.ID,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.ContactsGetResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Success, ShouldBeTrue)
				So(response.Contact.Name, ShouldEqual, contact.Name)
			})

			Convey("Updating that contact should succeed", func() {
				newName := uniuri.New()

				request := goreq.Request{
					Method:      "PUT",
					Uri:         server.URL + "/contacts/" + contact.ID,
					ContentType: "application/json",
					Body: `{
						"data": "` + uniuri.NewLen(64) + `",
						"name": "` + newName + `",
						"encoding": "xml",
						"version_major": 8,
						"version_minor": 3
					}`,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.ContactsUpdateResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Success, ShouldBeTrue)
				So(response.Contact.Name, ShouldEqual, newName)
			})

			Convey("Deleting that contact should succeed", func() {
				request := goreq.Request{
					Method: "DELETE",
					Uri:    server.URL + "/contacts/" + contact.ID,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.ContactsDeleteResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Success, ShouldBeTrue)
				So(response.Message, ShouldEqual, "Contact successfully removed")
			})
		})

		Convey("Update with invalid input should fail", func() {
			request := goreq.Request{
				Method:      "PUT",
				Uri:         server.URL + "/contacts/" + uniuri.New(),
				ContentType: "application/json",
				Body:        "123123!@#!@#",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.ContactsUpdateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Invalid input format")
		})

		Convey("Update of a non-existing contact should fail", func() {
			request := goreq.Request{
				Method:      "PUT",
				Uri:         server.URL + "/contacts/gibberish",
				ContentType: "application/json",
				Body: `{
						"data": "` + uniuri.NewLen(64) + `",
						"name": "` + uniuri.New() + `",
						"encoding": "xml",
						"version_major": 8,
						"version_minor": 3
					}`,
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.ContactsUpdateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Contact not found")
		})

		Convey("Deleting a non-existing contact should fail", func() {
			request := goreq.Request{
				Method: "DELETE",
				Uri:    server.URL + "/contacts/gibberish",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.ContactsDeleteResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Contact not found")
		})
	})
}
