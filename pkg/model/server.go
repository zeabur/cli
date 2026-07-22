package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// OS labels reported for a server. Renting and reinstalling provision the base
// OS only, so a machine without Zeabur services is a first-class outcome — it
// just cannot host Zeabur projects until they are installed.
const (
	OSZeaburOS = "ZeaburOS"
	OSUbuntu   = "Ubuntu"
)

// serverOS derives the machine's OS from hasK3s.
//
// hasK3s is three-state, and the third state is the trap:
//
//   - true  — Zeabur services are installed.
//   - false — explicitly none.
//   - nil   — a legacy server, which *does* have them. The backend infers this
//     from the server's certificate data, but the GraphQL field exposes the raw
//     column, so the inference never reaches the wire.
//
// So the test is `== false`, never `!hasK3s`: the latter would sweep every
// legacy server in and mislabel it as a plain VPS. Mirrors isCleanMachine() in
// the dashboard.
//
// Ubuntu is what Zeabur provisions when renting or reinstalling, so it is
// accurate for every machine that reached this state through Zeabur. A
// self-registered server running some other distribution would be labelled
// Ubuntu too, but such a server has Zeabur services installed at registration
// and therefore reports ZeaburOS instead.
func serverOS(hasK3s *bool) string {
	if hasK3s != nil && !*hasK3s {
		return OSUbuntu
	}
	return OSZeaburOS
}

// formatUsage renders a used/total pair, or an em dash when the total is zero.
//
// Zero total is the contract for "not measured" — a real machine cannot have
// zero cores, zero memory or zero disk. Printing "0/0" would read as a machine
// in trouble rather than one whose metrics have not been collected yet.
// The unit is always spelled out: the backend reports CPU in millicores, so a
// 4-core machine reads as 4000 and would otherwise look like 4000 cores.
func formatUsage(used, total int, unit string) string {
	if total == 0 {
		return "—"
	}
	return fmt.Sprintf("%d/%d %s", used, total, unit)
}

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
	SSHUsername        *string             `graphql:"sshUsername"`
	Country            *string             `graphql:"country"`
	City               *string             `graphql:"city"`
	Continent          *string             `graphql:"continent"`
	IsManaged          bool                `graphql:"isManaged"`
	ProvisioningStatus *string             `graphql:"provisioningStatus"`
	CreatedAt          time.Time           `graphql:"createdAt"`
	Status             ServerStatus        `graphql:"status"`
	ProviderInfo       *ServerProviderInfo `graphql:"providerInfo"`
	Events             []ServerEvent       `graphql:"events"`
	HasK3s             *bool               `graphql:"hasK3s"`
}

// MarshalJSON adds the derived `os` field so machine-readable output states the
// machine's kind outright, instead of leaving every caller to rediscover that
// hasK3s is three-state.
func (s ServerDetail) MarshalJSON() ([]byte, error) {
	type alias ServerDetail
	return json.Marshal(struct {
		alias
		OS string `json:"os"`
	}{alias(s), serverOS(s.HasK3s)})
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
	return []string{"ID", "Name", "IP", "Provider", "Location", "Status", "VM Status", "OS", "CPU", "Memory", "Disk", "Managed", "Created At"}
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

	cpu := formatUsage(s.Status.UsedCPU, s.Status.TotalCPU, "m")
	memory := formatUsage(s.Status.UsedMemory, s.Status.TotalMemory, "MB")
	disk := formatUsage(s.Status.UsedDisk, s.Status.TotalDisk, "MB")

	return [][]string{
		{
			s.ID,
			s.Name,
			s.IP,
			provider,
			location,
			status,
			s.Status.VMStatus,
			serverOS(s.HasK3s),
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
	HasK3s             *bool               `graphql:"hasK3s"`
}

// MarshalJSON adds the derived `os` field. See ServerDetail.MarshalJSON.
func (s ServerListItem) MarshalJSON() ([]byte, error) {
	type alias ServerListItem
	return json.Marshal(struct {
		alias
		OS string `json:"os"`
	}{alias(s), serverOS(s.HasK3s)})
}

type ServerListItems []ServerListItem

func (s ServerListItems) Header() []string {
	return []string{"ID", "Name", "IP", "Provider", "Location", "Status", "VM Status", "OS"}
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
			serverOS(item.HasK3s),
		}
	}
	return rows
}
