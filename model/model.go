package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// DBHandler interface
type DBHandler interface {
	Close()
	GetFuelSales(sDte, eDte string) (fs *FuelSales, err error)
	GetFuelSalesCount() (count int, err error)
}

// FuelSales struct
type FuelSales struct {
	DateStart time.Time
	DateEnd   time.Time
	Sales     []*StationSales
}

// StationSales struct
type StationSales struct {
	Date        time.Time `bson:"record_ts" json:"date"`
	RecordDate  int       `bson:"record_date" json:"recordDate"`
	StationID   string    `bson:"store_id" json:"stationID"`
	StationName string    `json:"stationName"`
	Fuel        *Fuel     `bson:"fuel_sales" json:"fuelSales"`
}

// Fuel struct
type Fuel struct {
	NL   float64 `bson:"NL" json:"NL"`
	SNL  float64 `bson:"SNL" json:"SNL"`
	DSL  float64 `bson:"DSL" json:"DSL"`
	CDSL float64 `bson:"CDSL" json:"CDSL"`
}

// StationNames struct
type StationNames struct {
	Names []*StationName
}

// StationNameMap map
type StationNameMap map[string]string

// StationName struct
type StationName struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
}

// QueryParams struct
type QueryParams struct {
	DateStart time.Time
	DateEnd   time.Time
}
