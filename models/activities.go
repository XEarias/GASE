package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/gosimple/slug"
)

//Activities Model
type Activities struct {
	ID          int       `orm:"column(id);auto" json:"id"`
	Name        string    `orm:"column(name)" json:"name"`
	Description string    `orm:"column(description)" json:"description"`
	Sector      *Sectors  `orm:"column(sectors_id);rel(fk)" json:"sector"`
	Slug        string    `orm:"column(slug);size(255)"  json:"slug"`
	CreatedAt   time.Time `orm:"column(created_at);type(datetime);null;auto_now_add" json:"-"`
	UpdatedAt   time.Time `orm:"column(updated_at);type(datetime);null" json:"-"`
	DeletedAt   time.Time `orm:"column(deleted_at);type(datetime);null"  json:"-"`
}

//TableName =
func (t *Activities) TableName() string {
	return "activities"
}

// AddActivities insert a new Activities into database and returns
// last inserted Id on success.
func AddActivities(m *Activities) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetActivitiesByID retrieves Activities by Id. Returns error if
// Id doesn't exist
func GetActivitiesByID(id int) (v *Activities, err error) {
	o := orm.NewOrm()
	v = &Activities{ID: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllActivities retrieves all Activities matches certain condition. Returns empty list if
// no records exist
func GetAllActivities(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Activities))
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

	var l []Activities
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

// UpdateActivitiesByID updates Activities by Id and returns error if
// the record to be updated doesn't exist
func UpdateActivitiesByID(m *Activities) (err error) {
	o := orm.NewOrm()
	v := Activities{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteActivities deletes Activities by Id and returns error if
// the record to be deleted doesn't exist
func DeleteActivities(id int) (err error) {
	o := orm.NewOrm()
	v := Activities{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Activities{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func addDefaultDataActivities() (count int64, errors []error) {

	o := orm.NewOrm()

	dummyData := map[string][]Activities{
		/* 		"01": /* "figuras-geometricas-y-abstractas  {
			{
				Name:        "Figuras",
				Description: "",
			},
			{
				Name:        "Abstractos",
				Description: "",
			},
			{
				Name:        "3d",
				Description: "",
			},
			{
				Name:        "Esfera",
				Description: "",
			},
			{
				Name:        "Rectangulo",
				Description: "",
			},
			{
				Name:        "Piramide",
				Description: "",
			},
			{
				Name:        "Cuadrado",
				Description: "",
			},
			{
				Name:        "Rombo",
				Description: "",
			},
			{
				Name:        "Pentagono",
				Description: "",
			},
		}, */
		"02" /* agricultura-y-ganaderia */ : {
			{
				Name:        "Agricultura",
				Description: "",
			},
		},
	}

	for key, dummySector := range dummyData {

		sector := Sectors{Code: key}

		err := o.Read(&sector, "code")

		if err != nil {
			continue
		}

		for key, dummyActivity := range dummySector {

			dummySector[key].Slug = slug.Make(dummyActivity.Name)
			dummySector[key].Sector = &sector

		}

		result, err := o.InsertMulti(100, dummySector)

		if err != nil {

			errors = append(errors, err)
			continue
		}

		count += result

	}

	return
}
