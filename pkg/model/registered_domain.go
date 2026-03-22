package model

import (
	"fmt"
	"time"
)

type RegisteredDomain struct {
	ID                           string    `json:"_id" graphql:"_id"`
	Domain                       string    `json:"domain" graphql:"domain"`
	TLD                          string    `json:"tld" graphql:"tld"`
	Status                       string    `json:"status" graphql:"status"`
	AutoRenew                    bool      `json:"autoRenew" graphql:"autoRenew"`
	ExpiresAt                    time.Time `json:"expiresAt" graphql:"expiresAt"`
	RegisteredAt                 time.Time `json:"registeredAt" graphql:"registeredAt"`
	PurchasePrice                int       `json:"purchasePrice" graphql:"purchasePrice"`
	RenewalPrice                 int       `json:"renewalPrice" graphql:"renewalPrice"`
	RegistrantVerificationStatus *string   `json:"registrantVerificationStatus" graphql:"registrantVerificationStatus"`
}

func (d RegisteredDomain) Header() []string {
	return []string{"ID", "Domain", "Status", "Auto-Renew", "Expires", "Price/yr"}
}

func (d RegisteredDomain) Rows() [][]string {
	return [][]string{{
		d.ID,
		d.Domain,
		d.Status,
		fmt.Sprintf("%v", d.AutoRenew),
		d.ExpiresAt.Format("2006-01-02"),
		fmt.Sprintf("$%.2f", float64(d.RenewalPrice)/100),
	}}
}

type RegisteredDomains []RegisteredDomain

func (ds RegisteredDomains) Header() []string {
	return RegisteredDomain{}.Header()
}

func (ds RegisteredDomains) Rows() [][]string {
	rows := make([][]string, 0, len(ds))
	for _, d := range ds {
		rows = append(rows, d.Rows()[0])
	}
	return rows
}

type DomainSearchResult struct {
	Domain    string `json:"domain" graphql:"domain"`
	Available bool   `json:"available" graphql:"available"`
	Price     *int   `json:"price" graphql:"price"`
	TLD       string `json:"tld" graphql:"tld"`
}

type DomainSearchResults []DomainSearchResult

func (ds DomainSearchResults) Header() []string {
	return []string{"Domain", "Available", "Price/yr"}
}

func (ds DomainSearchResults) Rows() [][]string {
	rows := make([][]string, 0, len(ds))
	for _, d := range ds {
		avail := "✗"
		price := "-"
		if d.Available {
			avail = "✓"
			if d.Price != nil {
				price = fmt.Sprintf("$%.2f", float64(*d.Price)/100)
			}
		}
		rows = append(rows, []string{d.Domain, avail, price})
	}
	return rows
}

func (d DomainSearchResult) Header() []string {
	return []string{"Domain", "Available", "Price/yr"}
}

func (d DomainSearchResult) Rows() [][]string {
	avail := "✗"
	price := "-"
	if d.Available {
		avail = "✓"
		if d.Price != nil {
			price = fmt.Sprintf("$%.2f", float64(*d.Price)/100)
		}
	}
	return [][]string{{d.Domain, avail, price}}
}

type DNSRecord struct {
	ID       string `json:"id" graphql:"id"`
	Type     string `json:"type" graphql:"type"`
	Name     string `json:"name" graphql:"name"`
	Content  string `json:"content" graphql:"content"`
	TTL      int    `json:"ttl" graphql:"ttl"`
	Priority int    `json:"priority" graphql:"priority"`
	Proxied  bool   `json:"proxied" graphql:"proxied"`
}

type DNSRecords []DNSRecord

func (ds DNSRecords) Header() []string {
	return []string{"ID", "Type", "Name", "Content", "TTL", "Priority", "Proxied"}
}

func (ds DNSRecords) Rows() [][]string {
	rows := make([][]string, len(ds))
	for i, d := range ds {
		ttl := "Auto"
		if d.TTL > 1 {
			ttl = fmt.Sprintf("%d", d.TTL)
		}
		proxied := "No"
		if d.Proxied {
			proxied = "Yes"
		}
		priority := "-"
		if d.Priority > 0 {
			priority = fmt.Sprintf("%d", d.Priority)
		}
		rows[i] = []string{d.ID, d.Type, d.Name, d.Content, ttl, priority, proxied}
	}
	return rows
}

type RegistrantProfile struct {
	ID           string `json:"_id" graphql:"_id"`
	FirstName    string `json:"firstName" graphql:"firstName"`
	LastName     string `json:"lastName" graphql:"lastName"`
	Email        string `json:"email" graphql:"email"`
	Phone        string `json:"phone" graphql:"phone"`
	Address1     string `json:"address1" graphql:"address1"`
	City         string `json:"city" graphql:"city"`
	State        string `json:"state" graphql:"state"`
	Country      string `json:"country" graphql:"country"`
	PostalCode   string `json:"postalCode" graphql:"postalCode"`
	Organization string `json:"organization" graphql:"organization"`
	IsDefault    bool   `json:"isDefault" graphql:"isDefault"`
}

type RegistrantProfiles []RegistrantProfile

func (ps RegistrantProfiles) Header() []string {
	return []string{"ID", "Name", "Email", "Phone", "Country", "Default"}
}

func (ps RegistrantProfiles) Rows() [][]string {
	rows := make([][]string, len(ps))
	for i, p := range ps {
		isDefault := "No"
		if p.IsDefault {
			isDefault = "Yes"
		}
		rows[i] = []string{
			p.ID,
			p.FirstName + " " + p.LastName,
			p.Email,
			p.Phone,
			p.Country,
			isDefault,
		}
	}
	return rows
}

type CreateRegistrantProfileInput struct {
	FirstName    string  `json:"firstName" graphql:"firstName"`
	LastName     string  `json:"lastName" graphql:"lastName"`
	Email        string  `json:"email" graphql:"email"`
	Phone        string  `json:"phone" graphql:"phone"`
	Address1     string  `json:"address1" graphql:"address1"`
	Address2     *string `json:"address2,omitempty" graphql:"address2"`
	City         string  `json:"city" graphql:"city"`
	State        string  `json:"state" graphql:"state"`
	Country      string  `json:"country" graphql:"country"`
	PostalCode   string  `json:"postalCode" graphql:"postalCode"`
	Organization *string `json:"organization,omitempty" graphql:"organization"`
}

type UpdateRegistrantProfileInput struct {
	FirstName    *string `json:"firstName,omitempty" graphql:"firstName"`
	LastName     *string `json:"lastName,omitempty" graphql:"lastName"`
	Email        *string `json:"email,omitempty" graphql:"email"`
	Phone        *string `json:"phone,omitempty" graphql:"phone"`
	Address1     *string `json:"address1,omitempty" graphql:"address1"`
	Address2     *string `json:"address2,omitempty" graphql:"address2"`
	City         *string `json:"city,omitempty" graphql:"city"`
	State        *string `json:"state,omitempty" graphql:"state"`
	Country      *string `json:"country,omitempty" graphql:"country"`
	PostalCode   *string `json:"postalCode,omitempty" graphql:"postalCode"`
	Organization *string `json:"organization,omitempty" graphql:"organization"`
}

type UpdateRegistrantContactInput struct {
	FirstName    string  `json:"firstName" graphql:"firstName"`
	LastName     string  `json:"lastName" graphql:"lastName"`
	Email        string  `json:"email" graphql:"email"`
	Phone        string  `json:"phone" graphql:"phone"`
	Address1     string  `json:"address1" graphql:"address1"`
	Address2     *string `json:"address2,omitempty" graphql:"address2"`
	City         string  `json:"city" graphql:"city"`
	State        string  `json:"state" graphql:"state"`
	Country      string  `json:"country" graphql:"country"`
	PostalCode   string  `json:"postalCode" graphql:"postalCode"`
	Organization *string `json:"organization,omitempty" graphql:"organization"`
}

type CreateDNSRecordInput struct {
	Type     RegisteredDomainDNSRecordType `json:"type" graphql:"type"`
	Name     string        `json:"name" graphql:"name"`
	Content  string        `json:"content" graphql:"content"`
	TTL      *int          `json:"ttl,omitempty" graphql:"ttl"`
	Priority *int          `json:"priority,omitempty" graphql:"priority"`
	Proxied  *bool         `json:"proxied,omitempty" graphql:"proxied"`
}

func (CreateDNSRecordInput) GetGraphQLType() string {
	return "CreateRegisteredDomainDNSRecordInput"
}

type UpdateDNSRecordInput struct {
	Content  *string `json:"content,omitempty" graphql:"content"`
	TTL      *int    `json:"ttl,omitempty" graphql:"ttl"`
	Priority *int    `json:"priority,omitempty" graphql:"priority"`
	Proxied  *bool   `json:"proxied,omitempty" graphql:"proxied"`
}

func (UpdateDNSRecordInput) GetGraphQLType() string {
	return "UpdateRegisteredDomainDNSRecordInput"
}

type RegisteredDomainDNSRecordType string

const (
	RegisteredDomainDNSRecordTypeA     RegisteredDomainDNSRecordType = "A"
	RegisteredDomainDNSRecordTypeAAAA  RegisteredDomainDNSRecordType = "AAAA"
	RegisteredDomainDNSRecordTypeCNAME RegisteredDomainDNSRecordType = "CNAME"
	RegisteredDomainDNSRecordTypeMX    RegisteredDomainDNSRecordType = "MX"
	RegisteredDomainDNSRecordTypeTXT   RegisteredDomainDNSRecordType = "TXT"
	RegisteredDomainDNSRecordTypeSRV   RegisteredDomainDNSRecordType = "SRV"
	RegisteredDomainDNSRecordTypeCAA   RegisteredDomainDNSRecordType = "CAA"
	RegisteredDomainDNSRecordTypeNS    RegisteredDomainDNSRecordType = "NS"
)

func (t RegisteredDomainDNSRecordType) GetGraphQLType() string {
	return "RegisteredDomainDNSRecordType"
}

type PurchaseDomainResult struct {
	RegisteredDomain               RegisteredDomain `json:"registeredDomain" graphql:"registeredDomain"`
	PaymentAmountFromBalance       *int             `json:"paymentAmountFromBalance" graphql:"paymentAmountFromBalance"`
	PaymentAmountFromPaymentMethod *int             `json:"paymentAmountFromPaymentMethod" graphql:"paymentAmountFromPaymentMethod"`
	InvoiceID                      *string          `json:"invoiceID" graphql:"invoiceID"`
}
