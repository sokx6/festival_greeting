package utils

import (
	"html/template"
	"strings"
	"time"

	"github.com/6tail/lunar-go/calendar"
)

type FestivalEmailData struct {
	FestivalName string
	FriendName   string
	FriendEmail  string
}

func IsTodayFestival() (bool, string) {
	now := time.Now()
	year, month, day := now.Date()

	solar := calendar.NewSolar(year, int(month), day, 0, 0, 0)
	solarFestival := solar.GetFestivals()
	if solarFestival.Len() > 0 {
		return true, solarFestival.Front().Value.(string)
	} else {
		lunar := solar.GetLunar()
		lunarFestival := lunar.GetFestivals()
		if lunarFestival.Len() > 0 {
			return true, lunarFestival.Front().Value.(string)
		}
	}
	return false, ""
}

func GetFestivalEmail(festivalName string, friendName string, friendEmail string) (string, error) {

	tmpl, err := template.ParseFiles("template/festival.html")
	if err != nil {
		return "", err
	}
	data := FestivalEmailData{
		FestivalName: festivalName,
		FriendName:   friendName,
		FriendEmail:  friendEmail,
	}
	var emailBuilder strings.Builder
	err = tmpl.Execute(&emailBuilder, data)
	if err != nil {
		return "", err
	}
	return emailBuilder.String(), err
}
