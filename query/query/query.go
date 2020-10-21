package query

import (
	"encoding/json"

	"github.com/grafadruid/go-druid/query"
	"github.com/grafadruid/go-druid/query/datasource"
	"github.com/grafadruid/go-druid/query/types"
)

type Base struct {
	ID         string                 `json:"ID"`
	QueryType  string                 `json:"queryType"`
	DataSource query.DataSource       `json:"dataSource"`
	Intervals  []types.Interval       `json:"intervals"`
	Context    map[string]interface{} `json:"context"`
}

func NewBase() *Base {
	b := &Base{}
	b.SetQueryType("base")
	return b
}

func (b *Base) SetID(ID string) *Base {
	b.ID = ID
	return b
}

func (b *Base) SetQueryType(queryType string) *Base {
	b.QueryType = queryType
	return b
}

func (b *Base) SetDataSource(dataSource query.DataSource) *Base {
	b.DataSource = dataSource
	return b
}

func (b *Base) SetIntervals(intervals []types.Interval) *Base {
	b.Intervals = intervals
	return b
}

func (b *Base) SetContext(context map[string]interface{}) *Base {
	b.Context = context
	return b
}

func (b *Base) Type() query.ComponentType {
	return b.QueryType
}

func (b *Base) Language() query.QueryLanguage {
	return query.NativeLanguage
}

func (b *Base) UnmarshalJSON(data []byte) error {
	var tmp struct {
		ID         string                 `json:"ID"`
		QueryType  string                 `json:"queryType"`
		DataSource json.RawMessage        `json:"dataSource"`
		Intervals  []types.Interval       `json:"intervals"`
		Context    map[string]interface{} `json:"context"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	d, err := datasource.Load(tmp.DataSource)
	if err != nil {
		return err
	}
	if d.Type() == "query" {
		d.(*datasource.Query).UnmarshalJSONWithQueryLoader(tmp.DataSource, Load)
	}
	b.ID = tmp.ID
	b.QueryType = tmp.QueryType
	b.DataSource = d
	b.Intervals = tmp.Intervals
	b.Context = tmp.Context
	return nil
}

func Load(data []byte) (query.Query, error) {
	var t struct {
		Typ string `json:"queryType"`
	}
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	var q query.Query
	switch t.Typ {
	case "datasourceMetadata":
		q = NewDataSourceMetadata()
	case "groupBy":
		q = NewGroupBy()
	case "scan":
		q = NewScan()
	case "search":
		q = NewSearchSearch()
	case "segmentMetadata":
		q = NewSegmentMetadata()
	case "sql":
		q = NewSQLQuery()
	case "timeBoundary":
		q = NewTimeBoundary()
	case "timeseries":
		q = NewTimeseries()
	case "topN":
		q = NewTopN()
	}
	return q, json.Unmarshal(data, &q)
}