package migrations

import (
	r "gopkg.in/dancannon/gorethink.v2"
)

type Migration struct{}

func (m *Migration) Is_the_table_settings_ready_3(up bool) (term r.Term) {
	term = r.TableList().Contains("settings")
	return term
}

func (m *Migration) Is_the_table_settings_ready_1(up bool) (term r.Term) {
	if up {
		term = r.TableList().Do(func(result r.Term) r.Term {
			return r.Branch(result.Contains("settings"),
				nil,
				r.TableCreate("settings"),
			)
		})
	} else {
		term = r.TableList().Do(func(result r.Term) r.Term {
			return r.Branch(result.Contains("settings"),
				r.TableDrop("settings"),
				nil,
			)
		})
	}
	return term
}

func (m *Migration) Create_settings_2(up bool) (term r.Term) {
	if up {
		term = r.Table("settings").Get("update_region").Eq(nil).Do(func(result r.Term) r.Term {
			return r.Table("settings").Insert(map[string]interface{}{"id": "update_region", "time": r.Now()})
		})
	} else {
		term = r.Table("settings").Get("update_region").Ne(nil).Do(func(result r.Term) r.Term {
			return r.Table("settings").Get("update_region").Delete()
		})
	}
	return term
}
