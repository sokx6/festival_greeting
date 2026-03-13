package utils

import (
	"html/template"
	"strings"
	"time"

	"github.com/Lofanmi/chinese-calendar-golang/calendar"
)

type BirthdayEmailData struct {
	FriendName string
	Date       string
	Time       string
	SenderName string
	AvatarURL  string
}

func IsTodayBirthday(month int, day int, isLunar bool) bool {
	now := time.Now()
	calender := calendar.ByTimestamp(now.Unix())
	if isLunar {
		return calender.Lunar.GetMonth() == int64(month) && calender.Lunar.GetDay() == int64(day)
	}
	return calender.Solar.GetMonth() == int64(month) && calender.Solar.GetDay() == int64(day)
}

func GetBirthTmpl(friendName, senderName, avatarURL string, month int, day int, isLunar bool) (string, error) {
	tmpl, err := template.ParseFiles("./templates/birthday.html")
	if err != nil {
		return "", err
	}
	date := time.Now().Format("2006-01-02")
	time := time.Now().Format("2006-01-02 14:45:11")
	data := BirthdayEmailData{
		FriendName: friendName,
		Date:       date,
		Time:       time,
		SenderName: senderName,
		AvatarURL:  avatarURL,
	}

	var emailBuilder strings.Builder
	err = tmpl.Execute(&emailBuilder, data)
	if err != nil {
		return "", err
	}
	return emailBuilder.String(), err
}
