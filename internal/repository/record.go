package repository

import (
	"aiStudio/internal/mysql"
	"aiStudio/internal/repository/model"
)

func CreateRecord(data *model.RecordTable) error {
	return mysql.DB.Create(data).Error
}

func UpdateRecordByGenId(genId string, table model.RecordTable) error {
	return mysql.DB.Model(&model.RecordTable{}).Where("gen_id = ?", genId).Updates(table).Error
}

func GetRecord(genId string, token ...string) (*model.RecordTable, error) {
	var tmp = new(model.RecordTable)
	db := mysql.DB.Where("gen_id = ?", genId)
	if len(token) != 0 {
		db.Where("token = ?", token[0])
	}
	err := db.Find(tmp).Error
	if err != nil {
		return nil, err
	}
	if tmp.ID == 0 {
		return nil, nil
	}
	return tmp, nil
}
