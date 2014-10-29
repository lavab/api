package models

import "github.com/lavab/api/models/base"

type File struct {
	base.Resource
	Bytes []byte `json:"bytes" gorethink:"bytes"`
	Size  int    `json:"size" gorethink:"size"`
}

type Picture struct {
	base.Resource
	Data       File `json:"data" gorethink:"data"`
	ResX       int  `json:"res_x" gorethink:"res_x"`
	ResY       int  `json:"res_y" gorethink:"res_y"`
	ThumbSmall File `json:"thumb_small" gorethink:"thumb_small"`
	ThumbLarge File `json:"thumb_large" gorethink:"thumb_large"`
}
