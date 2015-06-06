package routes

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"

	"github.com/lavab/api/env"
	"github.com/lavab/api/models"
	"github.com/lavab/api/utils"
)

type AddressesListResponse struct {
	Success   bool              `json:"success"`
	Message   string            `json:"message,omitempty"`
	Addresses []*models.Address `json:"addresses,omitempty"`
}

func AddressesList(c web.C, w http.ResponseWriter, r *http.Request) {
	session := c.Env["token"].(*models.Token)
	addresses, err := env.Addresses.GetOwnedBy(session.Owner)
	if err != nil {
		env.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Unable to fetch addresses")

		utils.JSONResponse(w, 500, &AddressesListResponse{
			Success: false,
			Message: "Internal error (code AD/LI/01)",
		})
		return
	}

	utils.JSONResponse(w, 200, &AddressesListResponse{
		Success:   true,
		Addresses: addresses,
	})
}
