package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/builder"
	"ponglehub.co.uk/geppetto/scanner"
)

type Rollback struct {
	builder *builder.Builder
	devices *devices
	scanner *scanner.Scanner
}

func NewRollback() (*Rollback, error) {
	devices, err := makeDevices()
	if err != nil {
		return nil, err
	}

	scanner := scanner.New()
	builder := builder.New()

	return &Rollback{
		builder: builder,
		devices: devices,
		scanner: scanner,
	}, nil
}

func (r *Rollback) Destroy() {
	w.devices.destroy()
}

func (r *Rollback) Start(target string) error {
	repos, err := r.scanner.ScanDir(target)
	if err != nil {
		return err
	}

	for {
		width, height := r.devices.getSize()

		r.devices.clear()
		r.devices.drawBorder(width, height)
		r.devices.drawTitle("GEPPETTO", width, true)

		offset := 3

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
					w.devices.drawText("🏗️", 60, line+offset, 7, style)
					w.devices.drawText(repo.Phase(), 64, line+offset, 20, style)
				}
			} else {
				w.devices.drawText("⏳", 60, line+offset, 5, style)
			}
		}

		w.devices.flush()
	}
}
