package mapper

import (
	"jwt-project/database/model"
	"jwt-project/dto"
)

/*
func MapperSignUp(p model.Person) dto.DtoSignUp {
	return dto.DtoSignUp{
		Id:             p.ID,
		PassW:          *p.Password,
		Tokennn:        *p.Token,
		RefreshTokennn: *p.RefreshToken,
		CreatedAttt:    p.CreatedAt,
		UpdatedAttt:    p.UpdatedAt,
		Userid:         p.UserId,
	}
}
*/

func MapperSignUp(d dto.DtoSignUp) *model.Person {
	return &model.Person{
		ID:           d.ID,
		Password:     &d.Password,
		Token:        &d.Token,
		RefreshToken: &d.RefreshToken,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		UserId:       d.UserId,
	}
}
