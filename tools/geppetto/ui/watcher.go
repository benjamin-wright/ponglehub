package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/builder"
	"ponglehub.co.uk/geppetto/scanner"
	"ponglehub.co.uk/geppetto/types"
)

type Watcher struct {
	builder *builder.Builder
	devices *devices
	scanner *scanner.Scanner
}

func NewWatcher(chartRepo string) (*Watcher, error) {
	devices, err := makeDevices()
	if err != nil {
		return nil, err
	}

	scanner := scanner.New()
	builder := builder.New(chartRepo)

	return &Watcher{
		builder: builder,
		devices: devices,
		scanner: scanner,
	}, nil
}

func (w *Watcher) Destroy() {
	w.devices.destroy()
}

func (w *Watcher) Start(target string) error {
	repos, err := w.scanner.ScanDir(target)
	if err != nil {
		return err
	}

	var state []types.RepoState
	building := true
	scroll := 0
	maxScroll := 0
	selected := -1
	highlighted := -1

	watchEvents, errorEvents := w.scanner.WatchDir(repos)
	commandEvents := w.devices.listen()
	progressEvents := w.builder.Build(repos, watchEvents)

	for {
		select {
		case cmd := <-commandEvents:
			switch cmd {
			case quitCommand:
				if selected != -1 {
					selected = -1
				} else {
					return nil
				}
			case upCommand:
				if highlighted == selected && maxScroll > 0 && scroll > 0 {
					scroll--
				} else if highlighted > 0 {
					highlighted--
				}
			case downCommand:
				if highlighted == selected && maxScroll > 0 && scroll < maxScroll {
					scroll++
				} else if highlighted < len(state)-1 {
					highlighted++
				}
			case selectCommand:
				if selected == highlighted {
					selected = -1
					maxScroll = 0
					scroll = 0
				} else {
					selected = highlighted
					maxScroll = 0
					scroll = 0
				}
			}
		case event := <-progressEvents:
			if event == nil {
				building = false
			} else {
				building = true
				state = event
			}
		case err := <-errorEvents:
			logrus.Fatalf("Error during run: %+v", err)
		}

		width, height := w.devices.getSize()

		w.devices.clear()
		w.devices.drawBorder(width, height)
		w.devices.drawTitle("GEPPETTO", width, building)

		offset := 3
		spareLines := height - 6 - len(state)

		for line, repo := range state {
			style := tcell.StyleDefault

			if line == highlighted {
				style = style.Background(tcell.ColorDarkGreen)
				w.devices.highlightLine(offset+line, width, style)
			}

			if line == selected {
				style = style.Foreground(tcell.ColorLightSlateGray)
			}

			w.devices.drawIcon(repo.Repo().RepoType, 2, line+offset, style)
			w.devices.drawText(repo.Repo().Name, 5, line+offset, 50, style)
			if repo.Built() {
				w.devices.drawText("✅", 60, line+offset, 5, style)
			} else if repo.Skipped() {
				w.devices.drawText("🔄", 60, line+offset, 5, style)
			} else if repo.Blocked() {
				w.devices.drawText("❌", 60, line+offset, 5, style)
			} else if repo.Errored() != nil {
				w.devices.drawText("🔥", 60, line+offset, 5, style)
				if line == selected {
					errorMsg := repo.Errored().Error()
					lines := w.devices.getNumLines(errorMsg, width-6)

					logrus.Infof("Message length: %d", lines)
					logrus.Infof("Screen height: %d", height)

					if lines > spareLines {
						maxScroll = lines - spareLines
					} else {
						maxScroll = 0
					}

					logrus.Infof("Max Scroll: %d", maxScroll)
					logrus.Infof("Lines: %d", lines)

					if lines > spareLines {
						lines = spareLines
					}

					logrus.Infof("Spare lines: %d", spareLines)
					logrus.Infof("Lines: %d", lines)

					w.devices.drawMultiline(errorMsg, 3, line+4, width-6, spareLines, scroll, tcell.StyleDefault)
					offset = offset + 1 + lines
				}
			} else if repo.Building() {
				if repo.Phase() == "check" {
					w.devices.drawText("💡", 60, line+offset, 5, style)
				} else {
					w.devices.drawText("🏗", 60, line+offset, 5, style)
					w.devices.drawText(repo.Phase(), 64, line+offset, 20, style)
				}
			} else {
				w.devices.drawText("⏳", 60, line+offset, 5, style)
			}
		}

		w.devices.flush()
	}
}
