package service

import (
	"fmt"

	"github.com/IDEA/SERVER/pkg/dto"
	"github.com/IDEA/SERVER/pkg/gateway"
)

const (
	MemberSheetID = "1RHeYiAgOs3sQsREj5duzhUQUfNe3W6wiiCdZWi5rcTk"
)

type ManageMemberService struct {
	gg gateway.GoogleOAuthGateway
}

func NewManageMemberService(gg gateway.GoogleOAuthGateway) *ManageMemberService {
	return &ManageMemberService{gg: gg}
}

func (s *ManageMemberService) StoreApplicationMemberData(application *dto.Application) error {
	values := make([][]interface{}, 1)
	values[0] = make([]interface{}, len(application.Surveys)+2)
	values[0][0] = application.Name
	values[0][1] = application.Email
	values[0][2] = application.Surveys[0].Answer //学校名
	values[0][3] = application.Surveys[6].Answer // 学年
	values[0][4] = application.Surveys[1].Answer //専攻
	values[0][5] = application.Surveys[5].Answer //希望班
	values[0][6] = application.Surveys[2].Answer //技術スキル
	values[0][7] = application.Surveys[3].Answer //自己PR
	values[0][8] = application.Surveys[4].Answer //ハンドルネーム
	rangeAA := "A:A"
	valueA, err := s.gg.GetSpreadSheetValues(MemberSheetID, rangeAA)
	if err != nil {
		return err
	}
	existNum := len(valueA) - 1
	rangeAI := fmt.Sprintf("A%d:I%d", existNum+2, existNum+3)
	if err := s.gg.UpdateSpreadSheetValues(MemberSheetID, rangeAI, values); err != nil {
		return err
	}
	return nil
}
