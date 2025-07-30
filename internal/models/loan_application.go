package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoanApplication struct {
	gorm.Model
	UUID                     *uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();not null"`
	IncomingOrganizationUuid *uuid.UUID   `gorm:"type:uuid;not null;index"`
	IssueOrganizationUuid    *uuid.UUID   `gorm:"type:uuid;not null;index"`
	IncomingOrganization     Organization `gorm:"foreignKey:IncomingOrganizationUuid;references:UUID"`
	IssueOrganization        Organization `gorm:"foreignKey:IssueOrganizationUuid;references:UUID"`
	Value                    int64        `gorm:"not null;check:value >= 1000"`
	Phone                    string       `gorm:"not null;size:20"`
	Comment                  string       `gorm:"type:text"`
	//TODO: Добавить комментарий(Новый клиент, повторный клиент с активным займом, повторный клиент с закрытытым займом, просрочник ОКР(взыскание), просрочник СДВ(юрист)...)
}
