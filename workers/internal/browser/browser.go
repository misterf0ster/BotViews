package browser

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

func RandomProxy(proxies []string) string {
	if len(proxies) == 0 || proxies[0] == "" {
		return ""
	}
	return proxies[rand.Intn(len(proxies))]
}
func RandomUserAgent() string {
	uaList := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15",
		"Mozilla/5.0 (Linux; Android 11; SM-A505F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Mobile Safari/537.36",
	}
	return uaList[rand.Intn(len(uaList))]
}

// Cмотрит видео до рекламы, ждет и досматривает рекламу полностью, потом еще 3-5 сек ролик и выходит
func WatchRutubeWithAdLogic(videoURL string, userAgent string, proxyAddr string) error {

	if videoURL == "" {
		return fmt.Errorf("empty video URL")
	}
	parsed, err := url.Parse(strings.TrimSpace(videoURL))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("invalid video URL after trim: %q, err: %w", videoURL, err)
	}

	log.Printf("WatchRutubeWithAdLogic URL: %q", videoURL)
	if videoURL == "" {
		return fmt.Errorf("empty video URL")
	}

	pw, err := playwright.Run()
	if err != nil {
		return err
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Proxy:    &playwright.Proxy{Server: proxyAddr},
		// если прокси пустой, можно Proxy не передавать
	})
	if err != nil {
		return err
	}
	defer browser.Close()

	contextOptions := playwright.BrowserNewContextOptions{}
	if userAgent != "" {
		contextOptions.UserAgent = playwright.String(userAgent)
	}

	context, err := browser.NewContext(contextOptions)
	if err != nil {
		return err
	}

	page, err := context.NewPage()
	if err != nil {
		return err
	}

	videoSelector := `video`

	videoURL = strings.TrimSpace(videoURL)
	fmt.Printf("Playwright final URL: [%q]\n", videoURL)

	_, err = page.Goto(videoURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return err
	}

	_, err = page.WaitForSelector(videoSelector)
	if err != nil {
		return err
	}

	adAppeared := false
	adWasVisible := false

	log.Println("Watching main video until ad appears...")

mainLoop:
	for {
		time.Sleep(750 * time.Millisecond)

		result, err := page.Evaluate(`() => {
			const adElem = document.querySelector('div[data-testid="player-advertising"]');
			return adElem && window.getComputedStyle(adElem).display !== "none";
		}`)
		if err != nil {
			log.Printf("Evaluate error: %v", err)
			continue
		}

		adVisible, ok := result.(bool)
		if !ok {
			log.Printf("Unexpected eval result type")
			continue
		}

		if adVisible && !adWasVisible {
			log.Println("Ad has appeared, watching ad now...")
			adAppeared = true
			adWasVisible = true
		}
		if adAppeared && !adVisible && adWasVisible {
			log.Println("Ad finished, resuming main video for 3-5 seconds...")
			break mainLoop
		}
	}

	afterAdWatch := time.Duration(rand.Intn(3)+3) * time.Second
	time.Sleep(afterAdWatch)

	_, err = page.Goto("about:blank")
	if err != nil {
		return err
	}

	return nil
}

// WatchRutubeHuman — с эмуляцией движений и подменой отпечатков
//func WatchRutubeHuman(link string, duration time.Duration, userAgent, proxyAddr string) error {
//	allocOpts := []chromedp.ExecAllocatorOption{
//		chromedp.Flag("headless", true),
//		chromedp.Flag("disable-gpu", true),
//		chromedp.UserAgent(userAgent),
//	}
//	if proxyAddr != "" {
//		allocOpts = append(allocOpts, chromedp.ProxyServer(proxyAddr))
//	}
//	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)
//	defer cancel()
//	ctx, cancel2 := chromedp.NewContext(allocCtx)
//	defer cancel2()

//var videoSelector = `button[data-testid="player-play"]`
//	var videoTag = `video`

// Подмена отпечатков (JS spoof)
//	spoofTasks := chromedp.Tasks{
//		chromedp.Navigate(link),
//		chromedp.Sleep(time.Second * 2),
//		chromedp.Evaluate(`(() => {try{Intl.DateTimeFormat = function(){return {resolvedOptions:()=>({timeZone:"Europe/Moscow"})}}}catch(e){}})()`, nil),
//		chromedp.Evaluate(`(() => {const toDataURL=HTMLCanvasElement.prototype.toDataURL;HTMLCanvasElement.prototype.toDataURL=function(){return toDataURL.apply(this,arguments)+"canvas_spoof"}})()`, nil),
//		chromedp.Evaluate(`(() => {const getp=WebGLRenderingContext.prototype.getParameter;WebGLRenderingContext.prototype.getParameter=function(p){if(p===37445)return "Intel Inc.";if(p===37446)return "Intel Iris OpenGL Engine";return getp.apply(this,arguments)}})()`, nil),
//	}
//	_ = chromedp.Run(ctx, spoofTasks)

// Движение, клик, скролл
//	humanTasks := chromedp.Tasks{
//		chromedp.WaitVisible(videoTag, chromedp.ByQuery),
//		chromedp.Sleep(time.Second * 2),
//		chromedp.MouseClickXY(200, 200),
//		chromedp.Sleep(time.Millisecond * 500),
//		chromedp.MouseClickXY(800, 400),
//		chromedp.Sleep(time.Millisecond * 500),
//		chromedp.ActionFunc(func(ctx context.Context) error {
//			for i := 0; i < 4; i++ {
//				y := 300 + i*100
//				chromedp.MouseClickXY(300, float64(y)).Do(ctx)
//				time.Sleep(time.Millisecond * 500)
//			}
//			return nil
//		}),
//		chromedp.Sleep(time.Second * 1),
//		chromedp.KeyEvent("\uE00F"), // PageDown
//		chromedp.Sleep(time.Second * 2),
//		chromedp.KeyEvent("\uE010"), // PageUp
//		chromedp.Sleep(duration - 7*time.Second),
//	}
//	return chromedp.Run(ctx, humanTasks)
//}
