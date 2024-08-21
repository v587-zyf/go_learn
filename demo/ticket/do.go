package ticket

import (
	"context"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func Do() {
	// 新建上下文
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// 添加浏览器选项
	exePath := "D:\\software\\chrome-win\\chrome.exe"
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // 不开启图像界面
		//chromedp.Flag("disable-gpu", true),
		//chromedp.Flag("no-sandbox", true),
		//chromedp.Flag("disable-dev-shm-usage", true),
		//chromedp.Flag("mute-audio", false), // 关闭声音
		chromedp.ExecPath(exePath),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// 新建浏览器实例
	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	// 上下文设置超时时间
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	Login(ctx)

	account := "18927586075"
	password := "zyf990819"
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

		// 等3秒 关页面
		chromedp.Sleep(time.Second * 10),
		page.Close(),
	}
	err := chromedp.Run(ctx, actions...)
	if err != nil {
		log.Fatal(err)
	}

	ClosePage(ctx)
}

func ClosePage(ctx context.Context) {
	err := chromedp.Run(ctx, page.Close())
	if err != nil {
		log.Fatal(err)
	}
}
