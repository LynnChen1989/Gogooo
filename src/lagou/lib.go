package main

import (
	"github.com/jinzhu/now"
	"github.com/satori/go.uuid"
	"math"
	"strconv"
	"strings"
	"time"
)

func GetUUID() string {
	uid := uuid.Must(uuid.NewV4())
	return uid.String()
}

func CalculateTotalPage(totalCount, pageSize float64) int {
	totalPage := float64(totalCount) / float64(pageSize)

	return int(math.Ceil(totalPage))
}

func ToPipelineJobs(dJobs []Result) []LgJob {
	var pJobs []LgJob
	for _, v := range dJobs {
		longitude, _ := strconv.ParseFloat(v.Longitude, 64)
		latitude, _ := strconv.ParseFloat(v.Latitude, 64)
		pJobs = append(pJobs, LgJob{
			City:     v.City,
			District: v.District,

			CompanyShortName: v.CompanyShortName,
			CompanyFullName:  v.CompanyFullName,
			CompanyLabelList: strings.Join(v.CompanyLabelList, ","),
			CompanySize:      v.CompanySize,
			FinanceStage:     v.FinanceStage,

			PositionName:      v.PositionName,
			PositionLables:    strings.Join(v.PositionLables, ","),
			PositionAdvantage: v.PositionAdvantage,
			WorkYear:          v.WorkYear,
			Education:         v.Education,
			Salary:            v.Salary,

			IndustryField:  v.IndustryField,
			IndustryLables: strings.Join(v.IndustryLables, ","),

			Longitude:  longitude,
			Latitude:   latitude,
			Linestaion: v.Linestaion,

			CreateTime: MustDateToUnix(v.CreateTime),
			AddTime:    time.Now().Unix(),
		})
	}

	return pJobs
}

func MustDateToUnix(date string) int64 {
	if date == `` || date == `0000-00-00` {
		return 0
	}

	loc, _ := time.LoadLocation("Asia/Chongqing")
	t, err := now.ParseInLocation(loc, date)
	if err != nil {
		return 0
	}

	return t.Unix()
}
