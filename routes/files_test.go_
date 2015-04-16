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

func TestFilesRoute(t *testing.T) {
	Convey("When uploading a new file", t, func() {
		account := &models.Account{
			Resource: models.MakeResource("", "johnorange"),
			Status:   "complete",
			AltEmail: "john@orange.org",
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
				"username": "johnorange",
				"password": "fruityloops"
			}`,
		}.Do()
		So(err, ShouldBeNil)

		var response routes.TokensCreateResponse
		err = result.Body.FromJsonTo(&response)
		So(err, ShouldBeNil)

		So(response.Success, ShouldBeTrue)
		authToken := response.Token

		Convey("Misformatted body should fail", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/files",
				ContentType: "application/json",
				Body:        "!@#!@#",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.FilesCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Invalid input format")
		})

		Convey("Invalid set of data should fail", func() {
			request := goreq.Request{
				Method: "POST",
				Uri:    server.URL + "/files",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.AccountsCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Invalid request")
		})

		Convey("File upload should succeed", func() {
			request := goreq.Request{
				Method:      "POST",
				Uri:         server.URL + "/files",
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

			var response routes.FilesCreateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "A new file was successfully created")
			So(response.Success, ShouldBeTrue)
			So(response.File.ID, ShouldNotBeEmpty)

			file := response.File

			Convey("Getting that file should succeed", func() {
				request := goreq.Request{
					Method: "GET",
					Uri:    server.URL + "/files/" + file.ID,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.FilesGetResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.File.ID, ShouldNotBeNil)
				So(response.Success, ShouldBeTrue)
			})

			Convey("The file should be visible on the list", func() {
				request := goreq.Request{
					Method: "GET",
					Uri:    server.URL + "/files",
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.FilesListResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(len(*response.Files), ShouldBeGreaterThan, 0)
				So(response.Success, ShouldBeTrue)

				found := false
				for _, a := range *response.Files {
					if a.ID == file.ID {
						found = true
						break
					}
				}

				So(found, ShouldBeTrue)
			})

			Convey("Updating it should succeed", func() {
				request := goreq.Request{
					Method:      "PUT",
					Uri:         server.URL + "/files/" + file.ID,
					ContentType: "application/json",
					Body: `{
		"data": "` + uniuri.NewLen(64) + `",
		"name": "` + uniuri.New() + `",
		"encoding": "xml",
		"version_major": 2,
		"version_minor": 1,
		"pgp_fingerprints": ["` + uniuri.New() + `"]
	}`,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.FilesUpdateResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Success, ShouldBeTrue)
				So(response.File.ID, ShouldEqual, file.ID)
				So(response.File.Encoding, ShouldEqual, "xml")
			})

			Convey("Deleting it should succeed", func() {
				request := goreq.Request{
					Method: "DELETE",
					Uri:    server.URL + "/files/" + file.ID,
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.FilesDeleteResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Message, ShouldEqual, "File successfully removed")
				So(response.Success, ShouldBeTrue)
			})
		})

		Convey("Getting a non-existing file should fail", func() {
			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/files/doesntexist",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.FilesGetResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "File not found")
			So(response.Success, ShouldBeFalse)

			Convey("Updating it should fail too", func() {
				request := goreq.Request{
					Method:      "PUT",
					Uri:         server.URL + "/files/doesntexist",
					ContentType: "application/json",
					Body:        "{}",
				}
				request.AddHeader("Authorization", "Bearer "+authToken.ID)
				result, err := request.Do()
				So(err, ShouldBeNil)

				var response routes.FilesGetResponse
				err = result.Body.FromJsonTo(&response)
				So(err, ShouldBeNil)

				So(response.Message, ShouldEqual, "File not found")
				So(response.Success, ShouldBeFalse)
			})
		})

		Convey("Getting a non-owned file should fail", func() {
			file := &models.File{
				Encrypted: models.Encrypted{
					Encoding:        "json",
					Data:            uniuri.NewLen(64),
					Schema:          "file",
					VersionMajor:    1,
					VersionMinor:    0,
					PGPFingerprints: []string{uniuri.New()},
				},
				Resource: models.MakeResource("nonowned", "photo.jpg"),
			}

			err := env.Files.Insert(file)
			So(err, ShouldBeNil)

			request := goreq.Request{
				Method: "GET",
				Uri:    server.URL + "/files/" + file.ID,
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.FilesGetResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Message, ShouldEqual, "File not found")
			So(response.Success, ShouldBeFalse)
		})

		Convey("Updating using a misformatted body should fail", func() {
			request := goreq.Request{
				Method:      "PUT",
				Uri:         server.URL + "/files/shizzle",
				ContentType: "application/json",
				Body:        "!@#!@#",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.FilesUpdateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "Invalid input format")
		})

		Convey("Updating a non-existing file should fail", func() {
			request := goreq.Request{
				Method:      "PUT",
				Uri:         server.URL + "/files/shizzle",
				ContentType: "application/json",
				Body:        "{}",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.FilesUpdateResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "File not found")
		})

		Convey("Deleting a non-existing file should fail", func() {
			request := goreq.Request{
				Method: "DELETE",
				Uri:    server.URL + "/files/shizzle",
			}
			request.AddHeader("Authorization", "Bearer "+authToken.ID)
			result, err := request.Do()
			So(err, ShouldBeNil)

			var response routes.FilesDeleteResponse
			err = result.Body.FromJsonTo(&response)
			So(err, ShouldBeNil)

			So(response.Success, ShouldBeFalse)
			So(response.Message, ShouldEqual, "File not found")
		})
	})
}
