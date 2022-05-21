package email

import (
	"strconv"
	"time"
)

type Header struct {
	Company        string
	Subject        string
	CompanyLogoSrc string
}

type Footer struct {
	Company        string
	CompanyUrl     string
	ContactUrl     string
	CompanyLogoSrc string
	Copyright      string
}

type Layout struct {
	Header
	Footer
}

func (e *service) newLayoutData(subject string) Layout {

	layout := Layout{}

	companyLogoSrc := "https://static.wixstatic.com/media/847cc7_23ab3dd28e6147369154ca5691f5f2e2~mv2.png/v1/fill/w_72,h_72,al_c,usm_0.66_1.00_0.01/logo-640.png"
	contactUrl := e.config.Hosts.Vec + "contact"
	company := e.config.Companies.Brand
	companyUrl := e.config.Hosts.Vec
	copyright := "Copyright (c) " + strconv.Itoa(time.Now().Year())

	// define header
	layout.Header.Company = company
	layout.Header.CompanyLogoSrc = companyLogoSrc
	layout.Header.Subject = subject
	// define footer
	layout.Footer.Company = company
	layout.Footer.CompanyUrl = companyUrl
	layout.Footer.ContactUrl = contactUrl
	layout.Footer.Copyright = copyright

	return layout
}
