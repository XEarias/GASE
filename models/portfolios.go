package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

//Portfolios Model
type Portfolios struct {
	ID          int         `orm:"column(id);auto" json:"id"`
	Name        string      `orm:"column(name);size(255)" json:"name"`
	Activity    *Activities `orm:"column(activities_id);rel(fk)" json:"activity"`
	Description string      `orm:"column(description)" json:"description"`
	Client      string      `orm:"column(client);size(255)" json:"client"`
	Location    *Locations  `orm:"column(locations_id);rel(fk)" json:"location"`
	Service     *Services   `orm:"column(services_id);rel(fk)" json:"service"`
	Slug        string      `orm:"column(slug);size(255)" json:"slug"`
	CreatedAt   time.Time   `orm:"column(created_at);type(datetime);null;auto_now_add"`
	UpdatedAt   time.Time   `orm:"column(updated_at);type(datetime);null"`
	DeletedAt   time.Time   `orm:"column(deleted_at);type(datetime);null"`
}

//TableName define Name
func (t *Portfolios) TableName() string {
	return "portfolios"
}

// AddPortfolios insert a new Portfolios into database and returns last inserted Id on success.
func AddPortfolios(m *Portfolios) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPortfoliosById retrieves Portfolios by Id. Returns error if Id doesn't exist
func GetPortfoliosById(id int) (v *Portfolios, err error) {
	o := orm.NewOrm()
	v = &Portfolios{ID: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPortfolios retrieves all Portfolios matches certain condition. Returns empty list if no records exist
func GetAllPortfolios(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Portfolios))
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

	var l []Portfolios
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

//UpdatePortfoliosById updates Portfolios by Id and returns error if the record to be updated doesn't exist
func UpdatePortfoliosById(m *Portfolios) (err error) {
	o := orm.NewOrm()
	v := Portfolios{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePortfolios deletes Portfolios by Id and returns error if
// the record to be deleted doesn't exist
func DeletePortfolios(id int) (err error) {
	o := orm.NewOrm()
	v := Portfolios{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Portfolios{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
