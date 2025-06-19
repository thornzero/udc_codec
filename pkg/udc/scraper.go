package udc

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"gopkg.in/yaml.v3"
)

type UDCNode struct {
	ID       string     `yaml:"id"`
	Code     string     `yaml:"code"`
	Title    string     `yaml:"title"`
	Children []*UDCNode `yaml:"children,omitempty"`
}

// Configurable parameters
const (
	BaseURL    = "https://udcsummary.info/php/index.php?lang=en"
	MaxDepth   = 10
	MaxRetries = 3
)

func ScrapeFullHierarchyChromedpProduction(output string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	// Setup local random generator
	localRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Visit root
	if err := chromedp.Run(ctx,
		chromedp.Navigate(BaseURL),
		chromedp.Sleep(3*time.Second),
	); err != nil {
		return fmt.Errorf("navigation failed: %v", err)
	}

	state := LoadPersistentState()
	root := &UDCNode{Title: "UDC Summary Root"}
	defer state.FinalFlush()

	if err := walkTree(ctx, root, state, 0, localRand); err != nil {
		return err
	}
	return WriteFullYAML(root.Children, output)
}

func walkTree(ctx context.Context, parent *UDCNode, state *PersistentState, depth int, localRand *rand.Rand) error {
	if depth > MaxDepth {
		log.Printf("Max recursion depth reached at: %s", parent.Title)
		return nil
	}

	// Extract HTML
	var html string
	err := withRetries(func() error {
		return chromedp.Run(ctx, chromedp.OuterHTML("#classtree", &html))
	})
	if err != nil {
		return err
	}

	nodes, nodeMap := parseTreeHTML(html)

	// Build children
	for _, node := range nodes {
		if node.Parent == parent.Code {
			udcNode := &UDCNode{
				Code:  node.Code,
				Title: node.Title,
				ID:    node.ID,
			}
			if nodeMap[node.Code] == nil {
				udcNode.Children = []*UDCNode{}
			}
			parent.Children = append(parent.Children, udcNode)
		}
	}
	// Recurse children
	for _, child := range parent.Children {
		_, expandable := expandableNodes(html, child.Code)
		if !expandable {
			continue
		}

		// Avoid duplicate visits
		state.Lock()
		if state.Visited[child.Code] {
			state.Unlock()
			continue
		}
		state.Visited[child.Code] = true
		state.Unlock()

		log.Printf("Expanding: %s - %s", child.Code, child.Title)

		script := fmt.Sprintf(`d.openTo(%s,true);`, child.ID)
		err := withRetries(func() error {
			return chromedp.Run(ctx,
				chromedp.Evaluate(script, nil),
				randomSleep(localRand),
			)
		})
		if err != nil {
			return fmt.Errorf("expand failed on %s: %v", child.Code, err)
		}

		// Recursive call
		if err := walkTree(ctx, child, state, depth+1, localRand); err != nil {
			return err
		}
	}
	return nil
}

// Retry wrapper
func withRetries(task func() error) error {
	var err error
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		err = task()
		if err == nil {
			return nil
		}
		log.Printf("Retry %d/%d after error: %v", attempt, MaxRetries, err)
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	return err
}

// Random polite delay
func randomSleep(localRand *rand.Rand) chromedp.Action {
	delay := 500 + localRand.Intn(800)
	return chromedp.Sleep(time.Duration(delay) * time.Millisecond)
}

// Parsing helpers — same as before
type parsedNode struct {
	ID     string
	Parent string
	Code   string
	Title  string
}

func parseTreeHTML(html string) ([]*parsedNode, map[string]*parsedNode) {
	re := regexp.MustCompile(`d\.add\((\d+),(\d+),'(.*?)','.*?&nbsp;&nbsp;(.*?)','`)
	matches := re.FindAllStringSubmatch(html, -1)

	nodes := []*parsedNode{}
	nodeMap := make(map[string]*parsedNode)
	for _, m := range matches {
		code := strings.TrimSpace(m[3])
		if code == "-" || code == "--" || code == "---" || code == "----" {
			continue
		}
		pn := &parsedNode{
			ID:     m[1],
			Parent: m[2],
			Code:   code,
			Title:  strings.TrimSpace(m[4]),
		}
		nodes = append(nodes, pn)
		nodeMap[pn.ID] = pn
	}
	return nodes, nodeMap
}

func expandableNodes(html, code string) (string, bool) {
	re := regexp.MustCompile(`d\.add\((\d+),(\d+),'` + regexp.QuoteMeta(code) + `','`)
	m := re.FindStringSubmatch(html)
	if len(m) > 0 {
		return m[1], true
	}
	return "", false
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
