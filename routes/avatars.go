package routes

import (
	"crypto/md5"
	"encoding/hex"
	"image/color"
	"image/png"
	"net/http"
	"strconv"

	"github.com/lavab/api/utils"

	"github.com/cupcake/sigil/gen"
	"github.com/zenazn/goji/web"
)

var avatarConfig = gen.Sigil{
	Rows: 5,
	Foreground: []color.NRGBA{
		color.NRGBA{45, 79, 255, 255},
		color.NRGBA{254, 180, 44, 255},
		color.NRGBA{226, 121, 234, 255},
		color.NRGBA{30, 179, 253, 255},
		color.NRGBA{232, 77, 65, 255},
		color.NRGBA{49, 203, 115, 255},
		color.NRGBA{141, 69, 170, 255},
	},
	Background: color.NRGBA{224, 224, 224, 255},
}

func Avatars(c web.C, w http.ResponseWriter, r *http.Request) {
	// Parse the query params
	query := r.URL.Query()

	// Settings
	var (
		widthString = query.Get("width")
		width       int
	)

	// Read width
	if widthString == "" {
		width = 100
	} else {
		var err error
		width, err = strconv.Atoi(widthString)
		if err != nil {
			utils.JSONResponse(w, 400, map[string]interface{}{
				"succes":  false,
				"message": "Invalid width",
			})
			return
		}
	}

	hash := c.URLParams["hash"]
	ext := c.URLParams["ext"]

	// data to parse
	var data []byte

	// md5 hash
	if len(hash) == 32 {
		data, _ = hex.DecodeString(hash)
	}

	// not md5
	if data == nil {
		hashed := md5.Sum([]byte(hash))
		data = hashed[:]
	}

	// if svg
	if ext == "svg" {
		w.Header().Set("Content-Type", "image/svg+xml")
		avatarConfig.MakeSVG(w, width, false, data)
		return
	}

	// generate the png
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, avatarConfig.Make(width, false, data))
}
