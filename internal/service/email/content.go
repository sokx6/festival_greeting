package email

import (
	"festival_greeting/internal/service/client"
	"festival_greeting/internal/service/config"
	"festival_greeting/internal/utils"
	"strings"

	"fmt"
)

func cleanHTMLContent(content string) string {
	content = strings.ReplaceAll(content, "```html\n", "")
	content = strings.ReplaceAll(content, "```html", "")
	content = strings.ReplaceAll(content, "```\n", "")
	content = strings.ReplaceAll(content, "```", "")
	return strings.TrimSpace(content)
}

func GetEmailContent(festivalName, friendName, senderName, avatarURL string, aiModel config.Model) (string, error) {
	prompt := fmt.Sprintf("现在是%s,今天是%s,请你为我生成一个给%s的节日问候邮件正文，使用html格式，只需要html内容，不需要除了html之外的所有内容，要求尽量不使用emoji,有丰富的动态效果，贴合节日氛围，邮件文字部分要清晰可见，不能因为美观而难以辨认，各种建议等额外内容都不需要，禁止使用markdown语法，发送者是%s,发送者头像图床链接是%s", "2025年10月6号", festivalName, friendName, senderName, avatarURL)

	apiClient := client.NewClient(aiModel)
	content, err := apiClient.GetResponse(prompt, aiModel)
	if err != nil {
		fmt.Printf("获取AI生成的邮件内容失败: %v\n", err)
		content, err = utils.GetFestivalEmail(festivalName, friendName, senderName, avatarURL)
		if err != nil {
			return "", fmt.Errorf("获取默认邮件内容失败: %w", err)
		}
	}
	fmt.Printf("AI生成的邮件内容: %s\n", content)
	content = cleanHTMLContent(content)
	fmt.Printf("提取后的邮件内容: %s\n", content)
	return content, nil
}
