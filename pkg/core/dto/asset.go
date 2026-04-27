package dto

import (
	"github.com/aarondl/null/v8"
)

type AssetView struct {
	*models.AssetView
	Zone     string `json:"zone"`
	Division string `json:"division"`
	Station  string `json:"station"`
}

func (a *AssetView) GetOrganisationID() int {
	if a.OrganisationID.Valid {
		return a.OrganisationID.Int
	}
	return 0
}

func (a *AssetView) SetZone(z string) {
	a.Zone = z
}

func (a *AssetView) SetDivision(d string) {
	a.Division = d
}

func (a *AssetView) SetStation(s string) {
	a.Station = s
}

func NewAssetView(obj *models.AssetView) *AssetView {
	if obj == nil {
		return nil
	}

	res := &AssetView{AssetView: obj}
	// if res.AssetView.Info.Valid {
	// 	res.Info = parser.JsonToMapWithoutError[string, any]([]byte(res.AssetView.Info.String))
	// }
	return res
}

func NewAssetViewSlice(objSlice models.AssetViewSlice) []*AssetView {
	if objSlice == nil {
		return nil
	}

	res := make([]*AssetView, 0, len(objSlice))
	for _, obj := range objSlice {
		res = append(res, NewAssetView(obj))
	}
	return res
}

type AssetRequest struct {
	Code           string    `json:"code"`
	Alias          string    `json:"alias"`
	AssetTypeID    int       `json:"asset_type_id"`
	SmmsAssetCode  string    `json:"smms_asset_code"`
	Info           null.JSON `json:"info"`
	OrganisationID int       `json:"organisation_id"`
}

type AssetMakeDistribution struct {
	OrganisationID int     `json:"organisation_id"`
	AssetTypeID    int     `json:"asset_type_id"`
	AssetTypeCode  string  `json:"asset_type_code"`
	OrgCode        string  `json:"org_code"`
	Make           *string `json:"make"`
	MakeCount      int     `json:"make_count"`
	Zone           string  `json:"zone"`
	Division       string  `json:"division"`
	Station        string  `json:"station"`
}

func (a *AssetMakeDistribution) GetOrganisationID() int {
	return a.OrganisationID
}

func (a *AssetMakeDistribution) SetZone(z string) {
	a.Zone = z
}

func (a *AssetMakeDistribution) SetDivision(d string) {
	a.Division = d
}

func (a *AssetMakeDistribution) SetStation(s string) {
	a.Station = s
}

type AssetMakeFilters struct {
	OrganisationIDs []int    `json:"organisation_ids"`
	Makes           []string `json:"makes"`
	AssetTypeIDs    []int    `json:"asset_type_ids"`
	CodeContains    string   `json:"code_contains"`
	AliasContains   string   `json:"alias_contains"`
}
