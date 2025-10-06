package building

// BuildingType 建筑类型
type BuildingType string

const (
	BuildingTypeResidential  BuildingType = "residential"
	BuildingTypeCommercial   BuildingType = "commercial"
	BuildingTypeIndustrial   BuildingType = "industrial"
	BuildingTypePublic       BuildingType = "public"
	BuildingTypeRecreational BuildingType = "recreational"
)

// BuildingConfig 建筑配置
type BuildingConfig struct {
	Type         BuildingType           `json:"type"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	MaxLevel     int                    `json:"max_level"`
	BaseCost     int64                  `json:"base_cost"`
	UpgradeCost  int64                  `json:"upgrade_cost"`
	Capacity     int                    `json:"capacity"`
	Efficiency   float64                `json:"efficiency"`
	Requirements map[string]interface{} `json:"requirements"`
}
