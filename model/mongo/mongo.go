package mongo

import (
	"errors"
	"time"

	"github.com/pulpfree/gales-fuelsale-report/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DB struct
type DB struct {
	session *mgo.Session
}

// DB Constants
const (
	DBDips            = "gales-dips"
	DBSales           = "gales-sales"
	colGdipsStations  = "store"
	colGdipsFuelSales = "fuelsale"
)

// Time form constant
const (
	timeShortForm  = "20060102"
	timeRecordForm = "2006-01-02"
)

// NewDB sets up new MongoDB struct
func NewDB(connection string) (model.DBHandler, error) {

	s, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}

	return &DB{
		session: s,
	}, err
}

// Close method
func (db *DB) Close() {
	db.session.Close()
}

// GetFuelSales method
func (db *DB) GetFuelSales(sDte, eDte string) (fs *model.FuelSales, err error) {

	qp, err := db.setQueryParams(sDte, eDte)
	if err != nil {
		return fs, err
	}

	fs, err = db.fetchFuelSales(qp)
	if err != nil {
		return fs, err
	}
	// These dates, different than the model.QueryParams, are required in the report
	fs.DateEnd, _ = time.Parse(timeRecordForm, eDte)
	fs.DateStart, _ = time.Parse(timeRecordForm, sDte)

	nms, err := db.getStationNames()
	if err != nil {
		return fs, err
	}

	for _, s := range fs.Sales {
		s.StationName = nms[s.StationID]
	}

	return fs, err
}

// GetFuelSalesCount method
func (db *DB) GetFuelSalesCount() (count int, err error) {
	s := db.getFreshSession()
	defer s.Close()

	count, err = s.DB(DBDips).C(colGdipsFuelSales).Find(nil).Count()

	return count, err
}

func (db *DB) fetchFuelSales2(qp *model.QueryParams) (fs *model.FuelSales, err error) {

	empt := time.Time{}
	if qp.DateEnd == empt || qp.DateStart == empt {
		return fs, errors.New("Missing start and end dates in fetchFuelSales")
	}
	s := db.getFreshSession()
	defer s.Close()

	fs = &model.FuelSales{}

	col := s.DB(DBDips).C(colGdipsFuelSales)
	match := bson.M{
		"$match": bson.M{"record_ts": bson.M{"$gte": qp.DateStart, "$lt": qp.DateEnd}},
	}
	group := bson.M{
		"$group": bson.M{
			"_id":  bson.M{"stationID": "$store_id"},
			"NL":   bson.M{"$sum": "$fuel_sales.NL"},
			"SNL":  bson.M{"$sum": "$fuel_sales.SNL"},
			"DSL":  bson.M{"$sum": "$fuel_sales.DSL"},
			"CDSL": bson.M{"$sum": "$fuel_sales.CDSL"},
		},
	}
	project := bson.M{
		"$project": bson.M{
			"stationID": "$_id.stationID",
			"NL":        1,
			"SNL":       1,
			"DSL":       1,
			"CDSL":      1,
		},
	}
	sort := bson.M{
		"$sort": bson.M{"_id.stationID": 1},
	}

	pipe := col.Pipe([]bson.M{match, group, project, sort})
	pipe.All(&fs.Sales)

	return fs, err
}

func (db *DB) fetchFuelSales(qp *model.QueryParams) (fs *model.FuelSales, err error) {

	empt := time.Time{}
	if qp.DateEnd == empt || qp.DateStart == empt {
		return fs, errors.New("Missing start and end dates in fetchFuelSales")
	}
	s := db.getFreshSession()
	defer s.Close()

	fs = &model.FuelSales{}
	col := s.DB(DBDips).C(colGdipsFuelSales)

	q := bson.M{"record_ts": bson.M{"$gte": qp.DateStart, "$lt": qp.DateEnd}}
	col.Find(q).Sort("store_id", "record_date").All(&fs.Sales)

	return fs, err
}

func (db *DB) getStationNames() (m model.StationNameMap, err error) {

	s := db.getFreshSession()
	defer s.Close()

	sns := &model.StationNames{}

	col := s.DB(DBDips).C(colGdipsStations)
	iter := col.Find(nil).Iter()
	err = iter.All(&sns.Names)

	// put everything in a map for easy access
	m = make(model.StationNameMap)
	for _, nm := range sns.Names {
		m[nm.ID.Hex()] = nm.Name
	}

	return m, err
}

// Helper methods

func (db *DB) getFreshSession() *mgo.Session {
	return db.session.Copy()
}

func (db *DB) setQueryParams(sDte, eDte string) (qp *model.QueryParams, err error) {

	st, _ := time.Parse(timeRecordForm, sDte)
	et, _ := time.Parse(timeRecordForm, eDte)
	met := et.Add(time.Hour * 24)

	qp = &model.QueryParams{DateStart: st, DateEnd: met}

	return qp, err
}
