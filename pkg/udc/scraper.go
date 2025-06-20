package udc

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"gopkg.in/yaml.v3"
)

const BaseURL = "https://udcsummary.info/php/index.php"

func ScrapeFullHierarchy(output string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(BaseURL),
		chromedp.Sleep(3*time.Second),
	); err != nil {
		return fmt.Errorf("navigation failed: %v", err)
	}

	rawRoots, err := ScrapeRawTree(ctx)
	if err != nil {
		return err
	}

	modelRoots := ConvertRawToModel(rawRoots)
	return WriteFullYAML(modelRoots, output)
}

func WriteFullYAML(nodes []*UDCNode, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	return encoder.Encode(nodes)
}
