package query

import (
	"encoding/json"

	"github.com/grafadruid/go-druid/query"
	"github.com/grafadruid/go-druid/query/filter"
	"github.com/grafadruid/go-druid/query/types"
	"github.com/grafadruid/go-druid/query/virtualcolumn"
)

type Order string

const (
	Ascending  Order = "ASCENDING"
	Descending       = "DESCENDING"
	None             = "NONE"
)

type Scan struct {
	Base
	VirtualColumns []query.VirtualColumn `json:"virtualColumns"`
	ResultFormat   string                `json:"resultFormat"`
	BatchSize      int64                 `json:"batchSize"`
	Limit          int64                 `json:"limit"`
	Order          Order                 `json:"order"`
	Filter         query.Filter          `json:"filter"`
	Columns        []string              `json:"columns"`
	Legacy         bool                  `json:"legacy"`
}

func NewScan() *Scan {
	s := &Scan{}
	s.Base.SetQueryType("scan")
	return s
}

func (s *Scan) SetDataSource(dataSource query.DataSource) *Scan {
	s.Base.SetDataSource(dataSource)
	return s
}

func (s *Scan) SetIntervals(intervals []types.Interval) *Scan {
	s.Base.SetIntervals(intervals)
	return s
}

func (s *Scan) SetContext(context map[string]interface{}) *Scan {
	s.Base.SetContext(context)
	return s
}

func (s *Scan) SetVirtualColumns(virtualColumns []query.VirtualColumn) *Scan {
	s.VirtualColumns = virtualColumns
	return s
}

func (s *Scan) SetResultFormat(resultFormat string) *Scan {
	s.ResultFormat = resultFormat
	return s
}

func (s *Scan) SetBatchSize(batchSize int64) *Scan {
	s.BatchSize = batchSize
	return s
}

func (s *Scan) SetLimit(limit int64) *Scan {
	s.Limit = limit
	return s
}

func (s *Scan) SetOrder(order Order) *Scan {
	s.Order = order
	return s
}

func (s *Scan) SetFilter(filter query.Filter) *Scan {
	s.Filter = filter
	return s
}

func (s *Scan) SetColumns(columns []string) *Scan {
	s.Columns = columns
	return s
}

func (s *Scan) SetLegacy(legacy bool) *Scan {
	s.Legacy = legacy
	return s
}

func (s *Scan) UnmarshalJSON(data []byte) error {
	var b Base
	var tmp struct {
		VirtualColumns []json.RawMessage `json:"virtualColumns"`
		ResultFormat   string            `json:"resultFormat"`
		BatchSize      int64             `json:"batchSize"`
		Limit          int64             `json:"limit"`
		Order          Order             `json:"order"`
		Filter         json.RawMessage   `json:"filter"`
		Columns        []string          `json:"columns"`
		Legacy         bool              `json:"legacy"`
	}
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	var err error
	var v query.VirtualColumn
	vv := make([]query.VirtualColumn, len(tmp.VirtualColumns))
	for i := range tmp.VirtualColumns {
		if v, err = virtualcolumn.Load(tmp.VirtualColumns[i]); err != nil {
			return err
		}
		vv[i] = v
	}
	var f query.Filter
	if tmp.Filter != nil {
		f, err = filter.Load(tmp.Filter)
		if err != nil {
			return err
		}
	}
	s.Base = b
	s.VirtualColumns = vv
	s.ResultFormat = tmp.ResultFormat
	s.BatchSize = tmp.BatchSize
	s.Limit = tmp.Limit
	s.Order = tmp.Order
	s.Filter = f
	s.Columns = tmp.Columns
	s.Legacy = tmp.Legacy
	return nil
}