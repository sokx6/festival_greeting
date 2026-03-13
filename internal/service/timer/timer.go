package timer

import (
	"festival_greeting/internal/service/config"
	"festival_greeting/internal/service/email"
	"festival_greeting/internal/utils"
	"fmt"
	"time"
)

func StartDailyTask() {
	for {
		now := time.Now()
		config, err := config.LoadConfig("config.toml")
		if err != nil {
			fmt.Printf("加载配置失败: %v\n", err)
			continue
		}

		sendTime := config.SendTime
		next := time.Date(now.Year(), now.Month(), now.Day(), sendTime.Hour, sendTime.Minute, sendTime.Second, 0, now.Location())

		// 如果当前时间已经过了今天的 9 点，那么下一次执行就是明天的 9 点
		if next.Before(now) {
			next = next.AddDate(0, 0, 1)
		}

		// 计算需要等待的时长
		duration := next.Sub(now)
		fmt.Printf("下次执行任务将在: %v (等待: %v)\n", next.Format("2006-01-02 11:45:14"), duration)

		// 创建定时器并等待
		timer := time.NewTimer(duration)
		<-timer.C // 阻塞直到时间到达
		if isFestival, festivalName := utils.IsTodayFestival(); isFestival {
			fmt.Println("今天是节日，开始发送邮件...")
			for _, friend := range config.Friends {
				emailContent, err := email.GetFestivalEmail(festivalName, friend.Name, config.SenderName, config.AvatarURL, config.Model)
				if err != nil {
					fmt.Printf("获取邮件内容失败: %v\n", err)
					continue
				}

				sender := email.NewEmailSender(config.Email.From, friend.Email, fmt.Sprintf("%s的祝福", festivalName), config.Email.Host, config.Email.Port, config.Email.Username, config.Email.Password, emailContent)
				err = sender.Send()
				if err != nil {
					fmt.Printf("发送邮件失败: %v\n", err)
				}
			}
		}

		for _, friend := range config.Friends {
			if utils.IsTodayBirthday(friend.Birthday.Month, friend.Birthday.Day, friend.Birthday.IsLunar) {
				fmt.Printf("今天是%s的生日，开始发送邮件...\n", friend.Name)
				emailContent, err := email.GetBirthdayEmail(friend.Name, config.SenderName, config.AvatarURL, friend.Birthday.Month, friend.Birthday.Day, config.Model, friend.Birthday.IsLunar)
				if err != nil {
					fmt.Printf("获取生日邮件内容失败: %v\n", err)
					continue
				}

				sender := email.NewEmailSender(config.Email.From, friend.Email, fmt.Sprintf("%s的生日祝福", friend.Name), config.Email.Host, config.Email.Port, config.Email.Username, config.Email.Password, emailContent)
				err = sender.Send()
				if err != nil {
					fmt.Printf("发送生日邮件失败: %v\n", err)
				}
			}
		}

		fmt.Printf("任务完成，等待下一次执行...\n")
	}
}
