package bootstrap

import "fmt"

type startupProgress struct {
	total     int
	completed int
}

func newStartupProgress(total int) *startupProgress {
	return &startupProgress{total: total}
}

func (p *startupProgress) begin(stage string) {
	fmt.Printf("\r启动进度 [%s] %d/%d 正在%s", progressBar(p.completed, p.total, 24), p.completed, p.total, stage)
}

func (p *startupProgress) complete() {
	p.completed++
	if p.completed > p.total {
		p.completed = p.total
	}
	fmt.Printf("\r启动进度 [%s] %d/%d\n", progressBar(p.completed, p.total, 24), p.completed, p.total)
}

func progressBar(completed, total, width int) string {
	filled := 0
	if total > 0 {
		filled = completed * width / total
	}
	bar := make([]byte, width)
	for i := range bar {
		bar[i] = '-'
		if i < filled {
			bar[i] = '#'
		}
	}
	return string(bar)
}
