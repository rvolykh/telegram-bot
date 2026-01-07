package scrap

type PowerOnMenuItem struct {
	Name          string `json:"name"`
	RawMobileHTML string `json:"rawMobileHtml"`
}

type PowerOnMember struct {
	MenuItems []PowerOnMenuItem `json:"menuItems"`
}

type PowerOnResponse struct {
	Member []PowerOnMember `json:"hydra:member"`
}

type Schedule struct {
	Today    string `json:"today"`
	Tomorrow string `json:"tomorrow"`
}
