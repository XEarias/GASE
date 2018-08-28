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

//Sectors model
type Sectors struct {
	ID        int       `orm:"column(id);pk" json:"id"`
	Name      string    `orm:"column(name);size(255)" json:"name"`
	Slug      string    `orm:"column(slug);size(255)" json:"slug"`
	Code      string    `orm:"column(code);size(255)" json:"-"`
	CreatedAt time.Time `orm:"column(created_at);type(datetime);null;auto_now_add"`
	UpdatedAt time.Time `orm:"column(updated_at);type(datetime);null"`
	DeletedAt time.Time `orm:"column(deleted_at);type(datetime);null"`
}

//TableName define Name
func (t *Sectors) TableName() string {
	return "sectors"
}

// AddSectors insert a new Sectors into database and returns last inserted Id on success.
func AddSectors(m *Sectors) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

//GetSectorsById retrieves Sectors by Id. Returns error if Id doesn't exist
func GetSectorsById(id int) (v *Sectors, err error) {
	o := orm.NewOrm()
	v = &Sectors{ID: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

//GetAllSectors retrieves all Sectors matches certain condition. Returns empty list if no records exist
func GetAllSectors(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Sectors))
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

	var l []Sectors
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

//UpdateSectorsById updates Sectors by Id and returns error if the record to be updated doesn't exist
func UpdateSectorsById(m *Sectors) (err error) {
	o := orm.NewOrm()
	v := Sectors{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

//DeleteSectors deletes Sectors by Id and returns error if the record to be deleted doesn't exist
func DeleteSectors(id int) (err error) {
	o := orm.NewOrm()
	v := Sectors{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Sectors{ID: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//AddDefaultDataSectors on init app
func AddDefaultDataSectors() (count int64, err error) {

	o := orm.NewOrm()

	dummyData := []*Sectors{
		{
			Name: "Figuras Geometricas y Abstractas",
			Code: "01",
		},
		{
			Name: "Agricultura y ganaderia",
			Code: "02",
		},
		{
			Name: "Animales y Mascotas",
			Code: "03",
		},
		{
			Name: "Arte y fotografia",
			Code: "04",
		},
		{
			Name: "Industria automotriz",
			Code: "05",
		},
		{
			Name: "Accesorios y glamour",
			Code: "06",
		},
		{
			Name: "Transporte y logistica",
			Code: "07",
		},
		{
			Name: "Asesoria y Consultoria",
			Code: "08",
		},
		{
			Name: "Construccion y arquitectura",
			Code: "09",
		},
		{
			Name: "Ropa y Moda",
			Code: "10",
		},
		{
			Name: "Educacion y Formacion",
			Code: "01",
		},
		{
			Name: "Hogar y Jardin",
			Code: "11",
		},
		{
			Name: "Alimentos y Bebidas",
			Code: "12",
		},
		{
			Name: "Belleza y Cuidado personal",
			Code: "13",
		},
		{
			Name: "Salud y Medicina",
			Code: "14",
		},
		{
			Name: "Deportes y ejercicios",
			Code: "15",
		},
		{
			Name: "Hobies y Entretenimiento",
			Code: "16",
		},
		{
			Name: "Organización sin fin de lucro",
			Code: "17",
		},
		{
			Name: "Tecnologia y telecomunicacion",
			Code: "18",
		},
		{
			Name: "Turismo y viajes",
			Code: "19",
		},
		{
			Name: "Reparaciones y Mantenimiento",
			Code: "20",
		},
		{
			Name: "Niños",
			Code: "21",
		},
		{
			Name: "Medio ambiente",
			Code: "22",
		},
		{
			Name: "Social Media",
			Code: "23",
		},
		{
			Name: "Seguridad",
			Code: "24",
		},
		{
			Name: "Personas",
			Code: "25",
		},
		{
			Name: "Espiritualidad",
			Code: "26",
		},
	}

	for _, v := range dummyData {
		v.Slug = slug.Make(v.Name)
	}

	count, err = o.InsertMulti(100, dummyData)

	return
}
