package email

import (
	"festival_greeting/internal/service/client"
	"festival_greeting/internal/service/config"
	"festival_greeting/internal/utils"
	"strings"
	"time"

	"fmt"
)

const (
	festivalPrompt = "现在是%s,今天是%s,请你为我生成一个给%s的节日问候邮件正文，使用html格式，只需要html内容，不需要除了html之外的所有内容，要求尽量不使用emoji,有丰富的动态效果，贴合节日氛围，风格简约柔和，优雅、高级、沉稳，邮件文字部分要清晰可见，不能因为美观而难以辨认，各种建议等额外内容都不需要，禁止使用markdown语法，发送者是%s,发送者头像图床链接是%s"

	birthdayPrompt = "今天是公历%s，请你为我生成一个给%s的公历生日问候邮件正文，使用html格式，只需要html内容，不需要除了html之外的所有内容，要求尽量不使用emoji,有丰富的动态效果，风格简约柔和，优雅、高级、沉稳，贴合生日氛围，邮件文字部分要清晰可见，不能因为美观而难以辨认，各种建议等额外内容都不需要，禁止使用markdown语法，发送者是%s,发送者头像图床链接是%s"

	lunarBirthdayPrompt = "今天是农历%02d月%02d日，阳历%s，请你为我生成一个给%s的农历生日问候邮件正文，使用html格式，只需要html内容，不需要除了html之外的所有内容，要求尽量不使用emoji,有丰富的动态效果，风格简约柔和，优雅、高级、沉稳，贴合生日氛围，邮件文字部分要清晰可见，不能因为美观而难以辨认，各种建议等额外内容都不需要，禁止使用markdown语法，发送者是%s,发送者头像图床链接是%s"
)

func cleanHTMLContent(content string) string {
	content = strings.ReplaceAll(content, "```html\n", "")
	content = strings.ReplaceAll(content, "```html", "")
	content = strings.ReplaceAll(content, "```\n", "")
	content = strings.ReplaceAll(content, "```", "")
	return strings.TrimSpace(content)
}

func GetFestivalEmail(festivalName, friendName, senderName, avatarURL string, aiModel config.Model) (string, error) {
	prompt := fmt.Sprintf(festivalPrompt, time.Now().Format("2006-01-02 14:45:14"), festivalName, friendName, senderName, avatarURL)

	apiClient := client.NewClient(aiModel)
	content, err := apiClient.GetResponse(prompt, aiModel)
	if err != nil {
		fmt.Printf("获取AI生成的邮件内容失败: %v\n", err)
		content, err = utils.GetFesTmpl(festivalName, friendName, senderName, avatarURL)
		if err != nil {
			return "", fmt.Errorf("获取默认邮件内容失败: %w", err)
		}
	}
	fmt.Printf("AI生成的邮件内容: %s\n", content)
	content = cleanHTMLContent(content)
	fmt.Printf("提取后的邮件内容: %s\n", content)
	return content, nil
}

func GetBirthdayEmail(friendName, senderName, avatarURL string, month int, day int, aiModel config.Model, isLunar bool) (string, error) {
	if isLunar {
		prompt := fmt.Sprintf(lunarBirthdayPrompt, month, day, time.Now().Format("2006-01-02 14:45:11"), friendName, senderName, avatarURL)

		apiClient := client.NewClient(aiModel)

		content, err := apiClient.GetResponse(prompt, aiModel)
		if err != nil {
			content, err = utils.GetBirthTmpl(friendName, senderName, avatarURL, month, day, isLunar)
			if err != nil {
				return "", fmt.Errorf("获取默认邮件内容失败: %w", err)
			}
		}
		return content, nil
	}

	prompt := fmt.Sprintf(birthdayPrompt, time.Now().Format("2006-01-02 14:45:11"), friendName, senderName, avatarURL)

	apiClient := client.NewClient(aiModel)
	content, err := apiClient.GetResponse(prompt, aiModel)
	if err != nil {
		content, err = utils.GetBirthTmpl(friendName, senderName, avatarURL, month, day, isLunar)
		if err != nil {
			return "", fmt.Errorf("获取默认邮件内容失败: %w", err)
		}
	}

	content = cleanHTMLContent(content)
	return content, nil
}
