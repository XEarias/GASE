package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

// Currencies struct
type Currencies struct {
	ID        int         `orm:"column(id);auto"`
	Name      string      `orm:"column(name);size(255)"`
	Iso       string      `orm:"column(iso);size(3)"`
	Symbol    string      `orm:"column(symbol);size(3)"`
	Gateways  []*Gateways `orm:"reverse(many)"`
	CreatedAt time.Time   `orm:"column(created_at);type(datetime);null;auto_now_add"`
	UpdatedAt time.Time   `orm:"column(updated_at);type(datetime);null"`
	DeletedAt time.Time   `orm:"column(deleted_at);type(datetime);null"`
}

// TableName =
func (t *Currencies) TableName() string {
	return "currencies"
}

// AddCurrencies insert a new Currencies into database and returns
// last inserted Id on success.
func AddCurrencies(m *Currencies) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCurrenciesById retrieves Currencies by Id. Returns error if
// Id doesn't exist
func GetCurrenciesById(id int) (v *Currencies, err error) {
	o := orm.NewOrm()
	v = &Currencies{ID: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCurrencies retrieves all Currencies matches certain condition. Returns empty list if
// no records exist
func GetAllCurrencies(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Currencies))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Currencies
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateCurrencies updates Currencies by Id and returns error if
// the record to be updated doesn't exist
func UpdateCurrenciesById(m *Currencies) (err error) {
	o := orm.NewOrm()
	v := Currencies{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCurrencies deletes Currencies by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCurrencies(id int) (err error) {
	o := orm.NewOrm()
	v := Currencies{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Currencies{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//AddDefaultData on init app
func AddDefaultDataCurrencies() (err error) {

	o := orm.NewOrm()

	dummyData := []*Currencies{
		{
			Symbol: "€",
			Name:   "Euro",
			Iso:    "EUR",
		},
		{
			Symbol: "$",
			Name:   "Dólar estadounidense",
			Iso:    "USD",
		},
	}

	_, err = o.InsertMulti(100, dummyData)

	return err
}

func addRelationsGatewaysCurrencies() []error {

	o := orm.NewOrm()

	dummyData := map[string][]string{
		"01": /* Paypal */ {
			"USD",
		},
	}

	var errors []error

	for key, dummyGateway := range dummyData {

		gateway := Gateways{Code: key}

		err := o.Read(&gateway, "code")

		if err != nil {
			continue
		}

		m2m := o.QueryM2M(&gateway, "Currencies")

		var InsertManyCurrencies []*Currencies

		for _, iso := range dummyGateway {

			currency := Currencies{Iso: iso}

			err := o.Read(&currency, "iso")

			if err != nil {
				continue
			}

			InsertManyCurrencies = append(InsertManyCurrencies, &currency)

		}

		_, err = m2m.Add(InsertManyCurrencies)

		if err != nil {
			errors = append(errors, err)
		}

	}

	return errors
}
