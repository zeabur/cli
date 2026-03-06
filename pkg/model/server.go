package model

import (
	"fmt"
	"time"
)

type Server struct {
	ID      string  `graphql:"_id"`
	Country *string `graphql:"country"`
	City    *string `graphql:"city"`
	IP      string  `graphql:"ip"`
	Name    string  `graphql:"name"`
}

func (s Server) GetID() string {
	return "server-" + s.ID
}

func (s Server) String() string {
	var identifier string

	if s.Country != nil {
		identifier = *s.Country
		if s.City != nil {
			identifier = fmt.Sprintf("%s, %s", *s.City, *s.Country)
		}
	} else {
		identifier = s.IP
	}

	return fmt.Sprintf("%s (%s)", s.Name, identifier)
}

func (s Server) IsAvailable() bool {
	return true
}

type Servers []Server

func (s Servers) Header() []string {
	return []string{"ID", "Name", "Location", "IP"}
}

func (s Servers) Rows() [][]string {
	rows := make([][]string, len(s))
	for i, server := range s {
		location := server.IP
		if server.Country != nil {
			location = *server.Country
			if server.City != nil {
				location = fmt.Sprintf("%s, %s", *server.City, *server.Country)
			}
		}
		rows[i] = []string{server.GetID(), server.Name, location, server.IP}
	}
	return rows
}

type ServerStatus struct {
	IsOnline    bool    `graphql:"isOnline"`
	TotalCPU    int     `graphql:"totalCPU"`
	UsedCPU     int     `graphql:"usedCPU"`
	TotalMemory int     `graphql:"totalMemory"`
	UsedMemory  int     `graphql:"usedMemory"`
	TotalDisk   int     `graphql:"totalDisk"`
	UsedDisk    int     `graphql:"usedDisk"`
	VMStatus    string  `graphql:"vmStatus"`
	Latency     float64 `graphql:"latency"`
}

type ServerDetail struct {
	ID                 string              `graphql:"_id"`
	Name               string              `graphql:"name"`
	IP                 string              `graphql:"ip"`
	SSHPort            int                 `graphql:"sshPort"`
	SSHUsername         *string             `graphql:"sshUsername"`
	Country            *string             `graphql:"country"`
	City               *string             `graphql:"city"`
	Continent          *string             `graphql:"continent"`
	IsManaged          bool                `graphql:"isManaged"`
	ProvisioningStatus *string             `graphql:"provisioningStatus"`
	CreatedAt          time.Time           `graphql:"createdAt"`
	Status             ServerStatus        `graphql:"status"`
	ProviderInfo       *ServerProviderInfo `graphql:"providerInfo"`
	Events             []ServerEvent       `graphql:"events"`
}

type ServerProviderInfo struct {
	Code string `graphql:"code"`
	Name string `graphql:"name"`
}

type ServerEvent struct {
	Message  string    `graphql:"message"`
	Time     time.Time `graphql:"time"`
	Severity string    `graphql:"severity"`
}

func (s *ServerDetail) Header() []string {
	return []string{"ID", "Name", "IP", "Provider", "Location", "Status", "VM Status", "CPU", "Memory", "Disk", "Managed", "Created At"}
}

func (s *ServerDetail) Rows() [][]string {
	location := s.IP
	if s.City != nil && s.Country != nil {
		location = fmt.Sprintf("%s, %s", *s.City, *s.Country)
	} else if s.Country != nil {
		location = *s.Country
	} else if s.City != nil {
		location = *s.City
	}

	provider := ""
	if s.ProviderInfo != nil {
		provider = s.ProviderInfo.Name
	}

	status := "Offline"
	if s.Status.IsOnline {
		status = "Online"
	}

	managed := "No"
	if s.IsManaged {
		managed = "Yes"
	}

	cpu := fmt.Sprintf("%d/%d", s.Status.UsedCPU, s.Status.TotalCPU)
	memory := fmt.Sprintf("%d/%d MB", s.Status.UsedMemory, s.Status.TotalMemory)
	disk := fmt.Sprintf("%d/%d MB", s.Status.UsedDisk, s.Status.TotalDisk)

	return [][]string{
		{
			s.ID,
			s.Name,
			s.IP,
			provider,
			location,
			status,
			s.Status.VMStatus,
			cpu,
			memory,
			disk,
			managed,
			s.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}
}

type CloudProvider struct {
	Code     string `graphql:"code"`
	Name     string `graphql:"name"`
	Icon     string `graphql:"icon"`
	Homepage string `graphql:"homepage"`
	Console  string `graphql:"console"`
}

type DedicatedServerRegion struct {
	ID        string `graphql:"id"`
	Name      string `graphql:"name"`
	Continent string `graphql:"continent"`
	Country   string `graphql:"country"`
	City      string `graphql:"city"`
}

type DedicatedServerPlan struct {
	Name                 string   `graphql:"name" json:"name"`
	CPU                  int      `graphql:"cpu" json:"cpu"`
	Memory               int      `graphql:"memory" json:"memory"`
	Disk                 int      `graphql:"disk" json:"disk"`
	Egress               int      `graphql:"egress" json:"egress"`
	Price                int      `graphql:"price" json:"price"`
	OriginalPrice        *float64 `graphql:"originalPrice" json:"originalPrice,omitempty"`
	GPU                  *string  `graphql:"gpu" json:"gpu,omitempty"`
	Available            bool     `graphql:"available" json:"available"`
	MaxOutboundBandwidth *int     `graphql:"maxOutboundBandwidth" json:"maxOutboundBandwidth,omitempty"`
}

type DedicatedServerPlans []DedicatedServerPlan

func (p DedicatedServerPlans) Header() []string {
	return []string{"Name", "CPU", "Memory", "Disk", "Egress", "GPU", "Price (USD/mo)", "Available"}
}

func (p DedicatedServerPlans) Rows() [][]string {
	rows := make([][]string, len(p))
	for i, plan := range p {
		gpu := "-"
		if plan.GPU != nil {
			gpu = *plan.GPU
		}
		available := "Yes"
		if !plan.Available {
			available = "No"
		}
		rows[i] = []string{
			plan.Name,
			fmt.Sprintf("%d cores", plan.CPU),
			fmt.Sprintf("%d GB", plan.Memory),
			fmt.Sprintf("%d GB", plan.Disk),
			fmt.Sprintf("%d GB", plan.Egress),
			gpu,
			fmt.Sprintf("$%d", plan.Price),
			available,
		}
	}
	return rows
}

type ServerListItem struct {
	ID                 string              `graphql:"_id"`
	Name               string              `graphql:"name"`
	IP                 string              `graphql:"ip"`
	Country            *string             `graphql:"country"`
	City               *string             `graphql:"city"`
	ProvisioningStatus *string             `graphql:"provisioningStatus"`
	Status             ServerStatus        `graphql:"status"`
	ProviderInfo       *ServerProviderInfo `graphql:"providerInfo"`
}

type ServerListItems []ServerListItem

func (s ServerListItems) Header() []string {
	return []string{"ID", "Name", "IP", "Provider", "Location", "Status", "VM Status"}
}

func (s ServerListItems) Rows() [][]string {
	rows := make([][]string, len(s))
	for i, item := range s {
		location := item.IP
		if item.City != nil && item.Country != nil {
			location = fmt.Sprintf("%s, %s", *item.City, *item.Country)
		} else if item.Country != nil {
			location = *item.Country
		} else if item.City != nil {
			location = *item.City
		}

		provider := ""
		if item.ProviderInfo != nil {
			provider = item.ProviderInfo.Name
		}

		status := "Offline"
		if item.Status.IsOnline {
			status = "Online"
		}

		rows[i] = []string{
			item.ID,
			item.Name,
			item.IP,
			provider,
			location,
			status,
			item.Status.VMStatus,
		}
	}
	return rows
}
