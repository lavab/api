package models

type File struct {
	Resource
	Bytes []byte `json:"bytes" gorethink:"bytes"`
	Size  int    `json:"size" gorethink:"size"`
}

type Picture struct {
	Resource
	Data       File `json:"data" gorethink:"data"`
	ResX       int  `json:"res_x" gorethink:"res_x"`
	ResY       int  `json:"res_y" gorethink:"res_y"`
	ThumbSmall File `json:"thumb_small" gorethink:"thumb_small"`
	ThumbLarge File `json:"thumb_large" gorethink:"thumb_large"`
}
