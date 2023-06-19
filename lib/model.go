package lib

// Repository is a collection of images, each of which may have multiple tags.
type Repo struct {
	Name       string `json:"name"`        // name
	Desc       string `json:"desc"`        // desc
	LastUpdate int64  `json:"last_update"` // last update
	Tags       []*Tag `json:"tags"`        // tags
}

// tag may have multiple images for different os/arch
type Tag struct {
	Name      string   `json:"name"`       // name
	Created   int64    `json:"created"`    // created unix timestamp
	ChangeLog string   `json:"change_log"` // change log
	Images    []*Image `json:"images"`     // images
}

// os/arch specified image may have multiple layers
type Image struct {
	Digest string   `json:"digest"` // digest
	Arch   string   `json:"arch"`   // cpu arch
	Os     string   `json:"os"`     // os
	Size   int64    `json:"size"`   // size B
	Layers []*Layer `json:"layers"` // layers
}

// layer
type Layer struct {
	Script string `json:"script"` // script
	Size   int64  `json:"size"`   // size B
}
