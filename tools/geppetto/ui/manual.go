package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/builder"
	"ponglehub.co.uk/geppetto/scanner"
	"ponglehub.co.uk/geppetto/types"
)

type Manual struct {
	builder *builder.Builder
	devices *devices
	scanner *scanner.Scanner
}

func NewManual(chartRepo string) (*Manual, error) {
	devices, err := makeDevices()
	if err != nil {
		return nil, err
	}

	scanner := scanner.New()
	builder := builder.New(chartRepo)

	return &Manual{
		builder: builder,
		devices: devices,
		scanner: scanner,
	}, nil
}

func (m *Manual) Destroy() {
	m.devices.destroy()
}

func (m *Manual) Start(target string) error {
	repos, err := m.scanner.ScanDir(target)
	if err != nil {
		return err
	}

	var state []types.RepoState
	building := true
	scroll := 0
	maxScroll := 0
	selected := -1
	highlighted := 0

	commandEvents := m.devices.listen()
	inputEvents := make(chan builder.InputSignal, 3)
	progressEvents := m.builder.Build(repos, inputEvents)

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
			case rebuildCommand:
				if highlighted != -1 {
					logrus.Infof("Rebuilding repo: %s", state[highlighted].Repo().Name)
					inputEvents <- builder.InputSignal{
						Repo:       state[highlighted].Repo().Name,
						Invalidate: true,
						Reinstall:  false,
					}
				}
			case reinstallCommand:
				if highlighted != -1 {
					logrus.Infof("Rebuilding repo from scratch: %s", state[highlighted].Repo().Name)
					inputEvents <- builder.InputSignal{
						Repo:       state[highlighted].Repo().Name,
						Invalidate: true,
						Reinstall:  true,
					}
				}
			case rebuildAllCommand:
				logrus.Infof("Rebuilding repo: %s", state[highlighted].Repo().Name)
				inputEvents <- builder.InputSignal{
					Nuke: true,
				}
			}

		case event := <-progressEvents:
			if event == nil {
				building = false
			} else {
				building = true
				state = event
			}
		}

		width, height := m.devices.getSize()

		m.devices.clear()
		m.devices.drawBorder(width, height)
		m.devices.drawTitle("GEPPETTO", width, building)

		offset := 3
		spareLines := height - 6 - len(state)

		for line, repo := range state {
			style := tcell.StyleDefault

			if line == highlighted {
				style = style.Background(tcell.ColorDarkGreen)
				m.devices.highlightLine(offset+line, width, style)
			}

			if line == selected {
				style = style.Foreground(tcell.ColorLightSlateGray)
			}

			m.devices.drawIcon(repo.Repo().RepoType, 2, line+offset, style)
			m.devices.drawText(repo.Repo().Name, 5, line+offset, 50, style)
			if repo.Built() {
				m.devices.drawText("✅", 60, line+offset, 5, style)
			} else if repo.Skipped() {
				m.devices.drawText("🔄", 60, line+offset, 5, style)
			} else if repo.Blocked() {
				m.devices.drawText("❌", 60, line+offset, 5, style)
			} else if repo.Errored() != nil {
				m.devices.drawText("🔥", 60, line+offset, 5, style)
				if line == selected {
					errorMsg := repo.Errored().Error()
					lines := m.devices.getNumLines(errorMsg, width-6)

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

					m.devices.drawMultiline(errorMsg, 3, line+4, width-6, spareLines, scroll, tcell.StyleDefault)
					offset = offset + 1 + lines
				}
			} else if repo.Building() {
				if repo.Phase() == "check" {
					m.devices.drawText("💡", 60, line+offset, 5, style)
				} else {
					m.devices.drawText("🏗", 60, line+offset, 5, style)
					m.devices.drawText(repo.Phase(), 64, line+offset, 20, style)
				}
			} else {
				m.devices.drawText("⏳", 60, line+offset, 5, style)
			}
		}

		m.devices.flush()
	}
}
