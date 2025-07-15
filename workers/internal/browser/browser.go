package browser

import (
	"context"
	"math/rand"
	"time"

	"github.com/chromedp/chromedp"
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

// WatchRutubeHuman — с эмуляцией движений и подменой отпечатков
func WatchRutubeHuman(link string, duration time.Duration, userAgent, proxyAddr string) error {
	allocOpts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.UserAgent(userAgent),
	}
	if proxyAddr != "" {
		allocOpts = append(allocOpts, chromedp.ProxyServer(proxyAddr))
	}
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	defer cancel()
	ctx, cancel2 := chromedp.NewContext(allocCtx)
	defer cancel2()

	var videoSelector = `button[data-testid="player-play"]`
	var videoTag = `video`

	// Подмена отпечатков (JS spoof)
	spoofTasks := chromedp.Tasks{
		chromedp.Navigate(link),
		chromedp.Sleep(time.Second * 2),
		chromedp.Evaluate(`(() => {try{Intl.DateTimeFormat = function(){return {resolvedOptions:()=>({timeZone:"Europe/Moscow"})}}}catch(e){}})()`, nil),
		chromedp.Evaluate(`(() => {const toDataURL=HTMLCanvasElement.prototype.toDataURL;HTMLCanvasElement.prototype.toDataURL=function(){return toDataURL.apply(this,arguments)+"canvas_spoof"}})()`, nil),
		chromedp.Evaluate(`(() => {const getp=WebGLRenderingContext.prototype.getParameter;WebGLRenderingContext.prototype.getParameter=function(p){if(p===37445)return "Intel Inc.";if(p===37446)return "Intel Iris OpenGL Engine";return getp.apply(this,arguments)}})()`, nil),
	}
	_ = chromedp.Run(ctx, spoofTasks)

	// Движение, клик, скролл
	humanTasks := chromedp.Tasks{
		chromedp.WaitVisible(videoTag, chromedp.ByQuery),
		chromedp.Sleep(time.Second * 2),
		chromedp.MouseClickXY(200, 200),
		chromedp.Sleep(time.Millisecond * 500),
		chromedp.MouseClickXY(800, 400),
		chromedp.Sleep(time.Millisecond * 500),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < 4; i++ {
				y := 300 + i*100
				chromedp.MouseClickXY(300, float64(y)).Do(ctx)
				time.Sleep(time.Millisecond * 500)
			}
			return nil
		}),
		chromedp.Sleep(time.Second * 1),
		chromedp.KeyEvent("\uE00F"), // PageDown
		chromedp.Sleep(time.Second * 2),
		chromedp.KeyEvent("\uE010"), // PageUp
		chromedp.Sleep(duration - 7*time.Second),
	}
	return chromedp.Run(ctx, humanTasks)
}
