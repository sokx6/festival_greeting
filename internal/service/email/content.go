package email

import (
	"festival_greeting/internal/service/client"
	"festival_greeting/internal/service/config"
	"festival_greeting/internal/utils"
	"strings"
	"time"

	"fmt"
)

type Template struct {
	Date         string
	FestivalName string
	FriendName   string
	SenderName   string
	AvatarURL    string
	Relation     string
	Original     string
}

func (t *Template) Fill() string {
	content := t.Original
	content = strings.ReplaceAll(content, "{{Date}}", t.Date)
	content = strings.ReplaceAll(content, "{{FestivalName}}", t.FestivalName)
	content = strings.ReplaceAll(content, "{{FriendName}}", t.FriendName)
	content = strings.ReplaceAll(content, "{{SenderName}}", t.SenderName)
	content = strings.ReplaceAll(content, "{{AvatarURL}}", t.AvatarURL)
	content = strings.ReplaceAll(content, "{{Relation}}", t.Relation)
	return content
}

const (
	festivalPrompt = "现在是%s,今天是%s,请你为我生成一个给%s的节日问候邮件正文，使用html格式，只需要html内容，不需要除了html之外的所有内容，要求尽量不使用emoji,正文部分不要有动态效果，横幅和末尾可以有动态效果，贴合节日氛围，风格简约柔和，优雅、高级、沉稳，邮件文字部分要清晰可见，不能因为美观而难以辨认，各种建议等额外内容都不需要，禁止使用markdown语法，发送者是%s,发送者头像图床链接是%s\n**绝对禁止输出除了html以外的其他所有内容**"

	birthdayPrompt = "今天是公历%s，请你为我生成一个给%s的公历生日问候邮件正文，使用html格式，只需要html内容，不需要除了html之外的所有内容，要求尽量不使用emoji,正文部分不要有动态效果，横幅和末尾可以有动态效果，风格简约柔和，优雅、高级、沉稳，贴合生日氛围，邮件文字部分要清晰可见，不能因为美观而难以辨认，各种建议等额外内容都不需要，禁止使用markdown语法，发送者是%s,发送者头像图床链接是%s\n**绝对禁止输出除了html以外的其他所有内容**"

	lunarBirthdayPrompt = "今天是农历%02d月%02d日，阳历%s，请你为我生成一个给%s的农历生日问候邮件正文，使用html格式，只需要html内容，不需要除了html之外的所有内容，要求尽量不使用emoji,正文部分不要有动态效果，横幅和末尾可以有动态效果，风格简约柔和，优雅、高级、沉稳，贴合生日氛围，邮件文字部分要清晰可见，不能因为美观而难以辨认，各种建议等额外内容都不需要，禁止使用markdown语法，发送者是%s,发送者头像图床链接是%s\n**绝对禁止输出除了html以外的其他所有内容**"
)

func normalizeHTMLContent(content string) string {
	lines := strings.Split(content, "\n")
	isBlock := false
	var normalizedLines []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if !isBlock && strings.HasPrefix(trimmedLine, "```") {
			isBlock = true
			continue
		}
		if isBlock && trimmedLine == "```" {
			isBlock = false
			continue
		}
		if isBlock {
			normalizedLines = append(normalizedLines, trimmedLine)
		}
	}
	return strings.Join(normalizedLines, "\n")
}

func GetFestivalEmail(festivalName, friendName, senderName, avatarURL string, aiModel config.Model) (string, error) {
	prompt := fmt.Sprintf(festivalPrompt, time.Now().Format("2006-01-02 15:04:05"), festivalName, friendName, senderName, avatarURL)

	apiClient := client.NewClient(aiModel)
	content, err := apiClient.GetResponse(prompt)
	if err != nil {
		fmt.Printf("获取AI生成的邮件内容失败: %v\n", err)
		content, err = utils.GetFesTmpl(festivalName, friendName, senderName, avatarURL)
		if err != nil {
			return "", fmt.Errorf("获取默认邮件内容失败: %w", err)
		}
	}
	fmt.Printf("AI生成的邮件内容: %s\n", content)
	content = normalizeHTMLContent(content)
	fmt.Printf("提取后的邮件内容: %s\n", content)
	return content, nil
}

func GetBirthdayEmail(friendName, senderName, avatarURL string, month int, day int, aiModel config.Model, isLunar bool) (string, error) {
	if isLunar {
		prompt := fmt.Sprintf(lunarBirthdayPrompt, month, day, time.Now().Format("2006-01-02 15:04:05"), friendName, senderName, avatarURL)

		apiClient := client.NewClient(aiModel)

		content, err := apiClient.GetResponse(prompt)
		if err != nil {
			content, err = utils.GetBirthTmpl(friendName, senderName, avatarURL, month, day, isLunar)
			if err != nil {
				return "", fmt.Errorf("获取默认邮件内容失败: %w", err)
			}
		}
		return content, nil
	}

	prompt := fmt.Sprintf(birthdayPrompt, time.Now().Format("2006-01-02 15:04:05"), friendName, senderName, avatarURL)

	apiClient := client.NewClient(aiModel)
	content, err := apiClient.GetResponse(prompt)
	if err != nil {
		content, err = utils.GetBirthTmpl(friendName, senderName, avatarURL, month, day, isLunar)
		if err != nil {
			return "", fmt.Errorf("获取默认邮件内容失败: %w", err)
		}
	}

	content = normalizeHTMLContent(content)
	return content, nil
}
