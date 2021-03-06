package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

//Clients Model
type Clients struct {
	ID        int         `orm:"column(id);auto" json:"id"`
	Name      string      `orm:"column(name);size(255)" json:"name,omitempty" valid:"Required"`
	Email     string      `orm:"column(email);size(255)" json:"email,omitempty" valid:"Required; Email"`
	Password  string      `orm:"column(password);size(255)" json:"password,omitempty" valid:"Required"`
	Phone     string      `orm:"column(phone);size(255)" json:"phone,omitempty" valid:"Required"`
	Company   string      `orm:"column(company);size(255)" json:"company,omitempty" valid:"Required"`
	Country   *Countries  `orm:"column(countries_id);rel(fk)" json:"countries,omitempty" valid:"Required"`
	Orders    []*Orders   `orm:"reverse(many)" json:"orders,omitempty"`
	Projects  []*Projects `orm:"reverse(many)" json:"projects,omitempty"`
	Token     string      `orm:"-" json:"token,omitempty"`
	CreatedAt time.Time   `orm:"column(created_at);type(datetime);null;auto_now_add" json:"-"`
	UpdatedAt time.Time   `orm:"column(updated_at);type(datetime);null" json:"-"`
	DeletedAt time.Time   `orm:"column(deleted_at);type(datetime);null" json:"-"`
}

//TableName =
func (t *Clients) TableName() string {
	return "clients"
}

func (t *Clients) loadRelations() {

	o := orm.NewOrm()

	relations := []string{"Orders", "Projects"}

	for _, relation := range relations {
		o.LoadRelated(t, relation)
	}

	return

}

// AddClients insert a new Clients into database and returns
// last inserted Id on success.
func AddClients(m *Clients) (id int64, err error) {
	o := orm.NewOrm()
	m.Password = GetMD5Hash(m.Password)
	id, err = o.Insert(m)
	m.Password = ""
	return
}

// GetClientsByEmail retrieves Clients by Email. Returns error if Id doesn't exist
func GetClientsByEmail(Email string) (v *Clients, err error) {
	o := orm.NewOrm()

	v = &Clients{Email: Email}

	err = o.Read(v, "email")

	if err != nil {
		return nil, err
	}

	v.loadRelations()

	return
}

// LoginClients login a Clients, returns
// if Exists.
func LoginClients(m *Clients) (id int, err error) {
	o := orm.NewOrm()

	m.Password = GetMD5Hash(m.Password)

	err = o.QueryTable(m.TableName()).Filter("deleted_at__isnull", true).Filter("email", m.Email).Filter("password", m.Password).One(m)

	if err != nil {
		return 0, err
	}

	m.Password = ""

	return m.ID, err
}

// CreateOrUpdateUser login a Clients, returns
// if Exists.
func CreateOrUpdateUser(m *Clients) (id int, err error) {
	o := orm.NewOrm()

	params := orm.Params{"name": m.Name, "company": m.Company, "phone": m.Phone}

	query := o.QueryTable(m.TableName()).Filter("deleted_at__isnull", true).Filter("email", m.Email)
	err = query.One(m)

	if err == nil {
		query.Update(params)

		err = query.One(m)
		return m.ID, err
	}

	//password, err := password.Generate(5, 3, 0, false, false)
	m.Password = GetMD5Hash(m.Phone)

	newID, err := o.Insert(m)
	m.ID = int(newID)

	return m.ID, err

}

// GetClientsByID retrieves Clients by Id. Returns error if
// Id doesn't exist
func GetClientsByID(id int) (v *Clients, err error) {
	v = &Clients{ID: id}
	err = searchFK(v.TableName(), v.ID).One(v)

	if err != nil {
		return nil, err
	}

	v.loadRelations()

	return
}

// GetAllClients retrieves all Clients matches certain condition. Returns empty list if
// no records exist
func GetAllClients(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Clients))
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

	var l []Clients
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).Filter("deleted_at__isnull", true).RelatedSel().All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				v.loadRelations()
				v.Password = ""
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
				v.loadRelations()
				v.Password = ""
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateClientsByID updates Clients by Id and returns error if
// the record to be updated doesn't exist
func UpdateClientsByID(m *Clients) (err error) {
	o := orm.NewOrm()
	v := Clients{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64

		if m.Password != "" {
			m.Password = GetMD5Hash(m.Password)
		}

		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteClients deletes Clients by Id and returns error if
// the record to be deleted doesn't exist
func DeleteClients(id int, trash bool) (err error) {
	o := orm.NewOrm()
	v := Clients{ID: id}
	// ascertain id exists in the database
	err = o.Read(&v)

	if err != nil {
		return
	}

	if trash {
		_, err = o.Delete(&v)
	} else {
		v.DeletedAt = time.Now()
		_, err = o.Update(&v)
	}

	if err != nil {
		return
	}

	return
}

//GetClientsFromTrash return Clients soft Deleted
func GetClientsFromTrash() (clients []*Clients, err error) {

	o := orm.NewOrm()

	var v []*Clients

	_, err = o.QueryTable("clients").Filter("deleted_at__isnull", false).RelatedSel().All(&v)

	if err != nil {
		return
	}

	for _, currency := range v {
		currency.loadRelations()
	}

	clients = v

	return

}
