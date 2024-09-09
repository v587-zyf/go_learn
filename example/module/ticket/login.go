package ticket

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func Login(ctx context.Context) {
	account := "18927586075"
	password := "123"
	actions := []chromedp.Action{
		// 网址
		chromedp.Navigate(`https://kyfw.12306.cn/otn/resources/login.html`),
		// 等1秒
		chromedp.Sleep(time.Second * 1),
		// 输入账号
		chromedp.WaitVisible(`#J-userName`, chromedp.ByID),
		chromedp.SendKeys(`#J-userName`, account, chromedp.ByID),
		// 等1秒
		chromedp.Sleep(time.Second * 1),
		// 输入密码
		chromedp.WaitVisible(`#J-password`, chromedp.ByID),
		chromedp.SendKeys(`#J-password`, password, chromedp.ByID),
		// 点击登录
		chromedp.Submit("#J-login", chromedp.ByID),
	}
	if err := chromedp.Run(ctx, actions...); err != nil {
		log.Fatal(err)
	}
}
