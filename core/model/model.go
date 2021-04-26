package model

// ModelInit model字段设置
type ModelInit struct {
	State   []string `ini:"state"`
	Created []string `ini:"created"`
	Updated []string `ini:"updated"`
	Deleted []string `ini:"deleted"`
}

// is 用来判断是否符合model的ini设置参数
func (m ModelInit) is(fieldType string, s string) bool {
	var arrTmp []string
	if fieldType == "State" {
		arrTmp = m.State
	} else if fieldType == "Created" {
		arrTmp = m.Created
	} else if fieldType == "Updated" {
		arrTmp = m.Updated
	} else if fieldType == "Deleted" {
		arrTmp = m.Deleted
	} else {
		return false
	}
	for _, v := range arrTmp {
		if s == v {
			return true
		}
	}
	return false
}

// IsDeleted 是否符合状态字段别名
func (m ModelInit) IsState(state string) bool {
	return m.is("State", state)
}

// IsCreated 是否符合创建时间字段别名
func (m ModelInit) IsCreated(state string) bool {
	return m.is("Created", state)
}

// IsUpdated 是否符合更新时间字段别名
func (m ModelInit) IsUpdated(state string) bool {
	return m.is("Updated", state)
}

// IsDeleted 是否符合删除时间字段别名
func (m ModelInit) IsDeleted(state string) bool {
	return m.is("Deleted", state)
}
