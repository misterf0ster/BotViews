package browser

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/chromedp/chromedp"
)

// WatchRutubeWithAdLogic — смотрит видео до рекламы, ждет и досматривает рекламу полностью, потом еще 3-5 сек ролик и выходит
func WatchRutubeWithAdLogic(ctx context.Context, videoURL string) error {
	var adSelector = `div[data-testid="player-advertising"]`
	var videoSelector = `video`

	// 1. Заходим на страницу и ждём появления видео-плеера
	if err := chromedp.Run(ctx,
		chromedp.Navigate(videoURL),
		chromedp.WaitVisible(videoSelector, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second), // время на инициализацию скриптов
	); err != nil {
		return err
	}

	// 2. Смотрим ролик до появления рекламы (или до её окончания)
	adAppeared := false
	adCompleted := false
	adWasVisible := false

	log.Println("Watching main video until ad appears...")

mainLoop:
	for {
		// Проверяем каждые 0.75 секунды — появилась ли реклама
		time.Sleep(750 * time.Millisecond)
		var adVisible bool
		_ = chromedp.Run(ctx, chromedp.EvaluateAsDevTools(
			`!!document.querySelector('div[data-testid="player-advertising"]') && window.getComputedStyle(document.querySelector('div[data-testid="player-advertising"]')).display !== "none"`, &adVisible,
		))
		if adVisible && !adWasVisible {
			log.Println("Ad has appeared, watching ad now...")
			adAppeared = true
			adWasVisible = true
		}
		if adAppeared && !adVisible && adWasVisible {
			// Реклама закончилась
			log.Println("Ad finished, resuming main video for 3-5 seconds...")
			adCompleted = true
			break mainLoop
		}
	}

	// 3. После рекламы смотрим еще 3–5 секунд основной ролик
	afterAdWatch := time.Duration(rand.Intn(3)+3) * time.Second // 3-5 сек
	chromedp.Run(ctx, chromedp.Sleep(afterAdWatch))

	// 4. Покидаем страницу (например, about:blank)
	chromedp.Run(ctx, chromedp.Navigate("about:blank"))

	return nil
}
