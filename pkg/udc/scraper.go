package udc

import (
	"context"
	"fmt"
	"math/rand"
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

	// 1. Navigate to the main page
	if err := chromedp.Run(ctx,
		chromedp.Navigate(BaseURL),
		chromedp.Sleep(time.Duration(1000+rand.Intn(1000))*time.Millisecond),
	); err != nil {
		return fmt.Errorf("navigation failed: %v", err)
	}

	// 2. Extract all top-level menu links (except .vacant)
	var menuLinks []string
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`Array.from(document.querySelectorAll('ul.menu.boldmenu > li > a:not(.vacant)')).map(a => a.href)`, &menuLinks),
	); err != nil {
		return fmt.Errorf("failed to extract menu links: %v", err)
	}

	fmt.Printf("[INFO] Found %d menu links to process\n", len(menuLinks))

	// Create a global root node (TOP)
	globalRoot := &RawNode{
		ID:     "TOP",
		Parent: "",
		Code:   "TOP",
		Title:  "UDC Summary Root",
	}

	// Track which root nodes we've already added to avoid duplicates
	addedRootCodes := make(map[string]bool)

	totalNodes := 0
	for i, link := range menuLinks {
		fmt.Printf("[INFO] Processing page %d/%d: %s\n", i+1, len(menuLinks), link)

		if err := chromedp.Run(ctx,
			chromedp.Navigate(link),
			chromedp.Sleep(time.Duration(1000+rand.Intn(1000))*time.Millisecond),
			chromedp.Evaluate(`d.openAll();`, nil),
			chromedp.Sleep(time.Duration(2000+rand.Intn(2000))*time.Millisecond),
		); err != nil {
			fmt.Printf("[WARN] Failed to process link %s: %v\n", link, err)
			continue
		}

		// Scrape and parse tree for this root node
		rawRoots, err := ScrapeRawTree(ctx)
		if err != nil {
			fmt.Printf("[WARN] Failed to scrape tree for link %s: %v\n", link, err)
			continue
		}

		// Process each root node from this page
		for _, root := range rawRoots {
			if addedRootCodes[root.Code] {
				// This root node already exists, merge its children into the existing one
				if DebugMode {
					fmt.Printf("[DEBUG] Merging children from duplicate root: %s\n", root.Code)
				}
				existingRoot := findRootByCode(globalRoot.Children, root.Code)
				if existingRoot != nil {
					existingRoot.Children = append(existingRoot.Children, root.Children...)
				}
			} else {
				// New root node, add it
				if DebugMode {
					fmt.Printf("[DEBUG] Adding new root node: %s\n", root.Code)
				}
				globalRoot.Children = append(globalRoot.Children, root)
				addedRootCodes[root.Code] = true
			}
		}

		// Count nodes from this page
		pageNodeCount := countNodes(rawRoots)
		totalNodes += pageNodeCount
		fmt.Printf("[INFO] Page %d: collected %d nodes\n", i+1, pageNodeCount)
	}

	fmt.Printf("[INFO] Total nodes collected: %d\n", totalNodes)
	fmt.Printf("[INFO] Converting to YAML format...\n")

	modelRoots := ConvertRawToModel([]*RawNode{globalRoot})
	return WriteFullYAML(modelRoots, output)
}

// countNodes recursively counts all nodes in a tree
func countNodes(nodes []*RawNode) int {
	count := len(nodes)
	for _, node := range nodes {
		count += countNodes(node.Children)
	}
	return count
}

// findRootByCode finds a root node by its code
func findRootByCode(roots []*RawNode, code string) *RawNode {
	for _, root := range roots {
		if root.Code == code {
			return root
		}
	}
	return nil
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
