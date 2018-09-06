package mongo

import (
	"errors"
	"strconv"
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
	DBSales      = "gales-sales"
	colStations  = "station-nodes"
	colFuelSales = "fuel-sales-export"
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
		s.StationName = nms[s.StationID.Hex()]
		s.Date = fmtDateToTime(s.RecordDate)
	}

	if len(fs.Sales) == 0 {
		err = errors.New("no data for specified date range")
	}

	return fs, err
}

// GetFuelSalesCount method
func (db *DB) GetFuelSalesCount() (count int, err error) {
	s := db.getFreshSession()
	defer s.Close()

	count, err = s.DB(DBSales).C(colFuelSales).Find(nil).Count()

	return count, err
}

func (db *DB) fetchFuelSales(qp *model.QueryParams) (fs *model.FuelSales, err error) {

	empt := time.Time{}
	if qp.DateEnd == empt || qp.DateStart == empt {
		return fs, errors.New("Missing start and end dates in fetchFuelSales")
	}
	s := db.getFreshSession()
	defer s.Close()

	dateStart := fmtDateToInt(qp.DateStart)
	dateEnd := fmtDateToInt(qp.DateEnd)

	fs = &model.FuelSales{}
	col := s.DB(DBSales).C(colFuelSales)

	q := bson.M{"recordDate": bson.M{"$gte": dateStart, "$lt": dateEnd}}
	col.Find(q).Sort("stationID", "recordDate").All(&fs.Sales)

	return fs, err
}

func (db *DB) getStationNames() (m model.StationNameMap, err error) {

	s := db.getFreshSession()
	defer s.Close()

	sns := &model.StationNames{}

	col := s.DB(DBSales).C(colStations)
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

func fmtDateToInt(dte time.Time) (date int) {
	dteStr := dte.Format("20060102")
	date, _ = strconv.Atoi(dteStr)
	return date
}

func fmtDateToTime(dte int) (date time.Time) {
	date, _ = time.Parse("20060102", strconv.Itoa(dte))
	return date
}
