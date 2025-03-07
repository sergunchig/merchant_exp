package models

type Dao struct {
	connectStr string
}

func NewDao(connStr string) Dao {
	return Dao{
		connectStr: connStr,
	}
}
func (dao Dao) ConnectStr() string {
	return dao.connectStr
}
