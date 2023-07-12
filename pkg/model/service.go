package model

import (
	"time"
)

type Service struct {
	ID   string `bson:"_id" json:"id" graphql:"_id"`
	Name string `bson:"name" json:"name" graphql:"name"`
	//Template  ServiceTemplate    `bson:"template" json:"template" graphql:"template"`
	//Project *Project `bson:"project" json:"project" graphql:"project"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt" graphql:"createdAt"`

	// MarketItemCode is the code of the item in the marketplace. Only used if Template is ServiceTemplateMarketplace.
	MarketItemCode *string `bson:"marketItemCode" json:"marketItemCode" graphql:"marketItemCode"`

	CustomBuildCommand *string  `bson:"customBuildCommand" json:"customBuildCommand" graphql:"customBuildCommand"`
	CustomStartCommand *string  `bson:"customStartCommand" json:"customStartCommand" graphql:"customStartCommand"`
	OutputDir          *string  `bson:"outputDir" json:"outputDir" graphql:"outputDir"`
	RootDirectory      string   `bson:"rootDirectory" json:"rootDirectory" graphql:"rootDirectory"`
	WatchPaths         []string `bson:"watchPaths" json:"watchPaths" graphql:"watchPaths"`

	// PrePrompt is a prompt appended to the beginning of each session between the user and the agent.
	// only used if Template is ServiceTemplateAgent.
	//PrePrompt *string `bson:"prePrompt" json:"prePrompt" graphql:"prePrompt"`

	// Plugins is a list of plugins that are installed in the service.
	// only used if Template is ServiceTemplateAgent.
	//Plugins []AgentPluginInstallation `bson:"plugins" json:"plugins" graphql:"plugins"`

	// Contexts is a list of contexts that are installed in the service.
	// only used if Template is ServiceTemplateAgent.
	//Contexts []AgentContextInstallation `bson:"contexts" json:"contexts" graphql:"contexts"`
}

type ServiceConnection struct {
	PageInfo *PageInfo      `json:"pageInfo" graphql:"pageInfo"`
	Edges    []*ServiceEdge `json:"edges" graphql:"edges"`
}

type ServiceEdge struct {
	Node   *Service `json:"node" graphql:"node"`
	Cursor string   `json:"cursor" graphql:"cursor"`
}
