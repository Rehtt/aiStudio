package repository

import (
	"aiStudio/internal/mysql"
	"aiStudio/internal/repository/model"
)

func CreateRecord(token, prompt, genId string) error {
	return mysql.DB.Create(&model.RecordTable{
		Token:  token,
		Prompt: prompt,
		GenID:  genId,
	}).Error
}
