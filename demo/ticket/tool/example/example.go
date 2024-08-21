package example

import (
	"context"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

var (
	Port = 9999
)

func Do() {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // 不开启图像界面
		//chromedp.Flag("disable-gpu", true),
		//chromedp.Flag("no-sandbox", true),
		//chromedp.Flag("disable-dev-shm-usage", true),
		//chromedp.Flag("mute-audio", false), // 关闭声音
		chromedp.ExecPath("D:\\software\\chrome-win\\chrome.exe"),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// 点击
	Client(ctx)

	// 缓存
	//Cookie(ctx)

	// 可见
	//Visible()

	//time.Sleep(3 * time.Second)

	ClosePages(ctx)

	time.Sleep(30 * time.Second)
}

func ClosePages(ctx context.Context) {
	if err := chromedp.Run(ctx, page.Close()); err != nil {
		log.Fatal(err)
	}
}
