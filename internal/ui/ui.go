package ui

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/romanitalian/rss-reader/internal/feed"
	"github.com/romanitalian/rss-reader/internal/models"
	"github.com/romanitalian/rss-reader/internal/storage"
)

// MaxWidthLayout restricts content width
type MaxWidthLayout struct {
	MaxWidth float32
}

// Layout organizes objects with limited width
func (m *MaxWidthLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) == 0 {
		return
	}

	pos := fyne.NewPos(0, 0)
	if containerSize.Width > m.MaxWidth {
		pos.X = (containerSize.Width - m.MaxWidth) / 2
	}

	// Set content width
	width := m.MaxWidth
	if containerSize.Width < m.MaxWidth {
		width = containerSize.Width
	}

	// Position each object
	for _, o := range objects {
		size := o.MinSize()
		size.Width = width
		o.Resize(fyne.NewSize(width, size.Height))
		o.Move(pos)
	}
}

// MinSize calculates the minimum container size
func (m *MaxWidthLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}

	minSize := fyne.NewSize(0, 0)
	for _, o := range objects {
		objMin := o.MinSize()
		minSize.Width = fyne.Max(minSize.Width, objMin.Width)
		minSize.Height = fyne.Max(minSize.Height, objMin.Height)
	}

	minSize.Width = fyne.Min(minSize.Width, m.MaxWidth)
	return minSize
}

// NewMaxWidthLayout creates a new layout with width restriction
func NewMaxWidthLayout(maxWidth float32) *MaxWidthLayout {
	return &MaxWidthLayout{MaxWidth: maxWidth}
}

// UI represents the application user interface
type UI struct {
	app         fyne.App
	mainWindow  fyne.Window
	storage     storage.Storage
	feedService *feed.FeedService
	feeds       []*models.Feed
	currentFeed *models.Feed

	// Widgets
	feedList      *widget.List
	itemList      *widget.List
	contentView   *widget.RichText
	itemContainer *fyne.Container
	splitView     *container.Split
}

// NewUI creates a new user interface instance
func NewUI() (*UI, error) {
	// Create data directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dataDir := filepath.Join(homeDir, ".rssreader")

	// Initialize storage
	store, err := storage.NewFileStorage(dataDir)
	if err != nil {
		return nil, err
	}

	// Create and configure UI
	ui := &UI{
		app:         app.New(),
		storage:     store,
		feedService: feed.NewFeedService(),
	}

	return ui, nil
}

// Run starts the application
func (u *UI) Run() {
	u.mainWindow = u.app.NewWindow("RSS Reader")
	u.mainWindow.Resize(fyne.NewSize(1000, 600))

	// Load saved feeds
	feeds, err := u.storage.GetAllFeeds()
	if err == nil {
		u.feeds = feeds
	}

	// Create UI elements
	u.createWidgets()
	u.createLayout()

	// Show window and run application
	u.mainWindow.ShowAndRun()
}

// createWidgets creates all application widgets
func (u *UI) createWidgets() {
	// Feed list
	u.feedList = widget.NewList(
		func() int {
			return len(u.feeds)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template Feed Item")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			feed := u.feeds[id]
			obj.(*widget.Label).SetText(feed.Title)
		},
	)

	// When a feed is selected
	u.feedList.OnSelected = func(id widget.ListItemID) {
		u.currentFeed = u.feeds[id]
		u.updateItemList()
		u.itemList.Select(0)
	}

	// Article list
	u.itemList = widget.NewList(
		func() int {
			if u.currentFeed == nil {
				return 0
			}
			return len(u.currentFeed.Items)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Title"),
				widget.NewLabel("Date"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if u.currentFeed == nil {
				return
			}
			item := u.currentFeed.Items[id]

			vbox := obj.(*fyne.Container)
			title := vbox.Objects[0].(*widget.Label)
			date := vbox.Objects[1].(*widget.Label)

			title.SetText(item.Title)
			date.SetText(item.Published.Format("02 Jan 2006 15:04"))

			if item.Read {
				title.TextStyle = fyne.TextStyle{Italic: true}
			} else {
				title.TextStyle = fyne.TextStyle{}
			}
		},
	)

	// When an article is selected
	u.itemList.OnSelected = func(id widget.ListItemID) {
		if u.currentFeed == nil || id >= len(u.currentFeed.Items) {
			return
		}

		item := u.currentFeed.Items[id]
		u.showContent(item)

		// Mark the article as read
		if !item.Read {
			item.Read = true
			u.storage.MarkItemAsRead(u.currentFeed.ID, item.ID)
			u.itemList.Refresh()
		}
	}

	// Article content view area
	u.contentView = widget.NewRichTextFromMarkdown("")
}

// createLayout creates the UI layout
func (u *UI) createLayout() {
	// Action buttons for feeds
	addFeedBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), u.showAddFeedDialog)
	refreshFeedsBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), u.refreshFeeds)
	deleteFeedBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), u.deleteSelectedFeed)

	feedActions := container.NewHBox(
		addFeedBtn,
		refreshFeedsBtn,
		deleteFeedBtn,
	)

	feedPanel := container.NewBorder(
		widget.NewLabel("RSS Feeds"),
		feedActions,
		nil, nil,
		u.feedList,
	)

	// Container for article details
	u.itemContainer = container.NewBorder(
		nil, nil, nil, nil,
		container.NewScroll(u.contentView),
	)

	// Article list with header
	itemPanel := container.NewBorder(
		widget.NewLabel("Articles"),
		nil, nil, nil,
		u.itemList,
	)

	// Split container for article list and content
	u.splitView = container.NewHSplit(
		itemPanel,
		u.itemContainer,
	)
	u.splitView.Offset = 0.3

	// Main split container
	mainSplit := container.NewHSplit(
		feedPanel,
		u.splitView,
	)
	mainSplit.Offset = 0.2

	u.mainWindow.SetContent(container.NewPadded(mainSplit))
}

// showContent displays the article content
func (u *UI) showContent(item models.Item) {
	// Title
	title := widget.NewRichTextFromMarkdown("# " + item.Title)
	title.Wrapping = fyne.TextWrapWord

	// Publication date
	date := widget.NewLabel(fmt.Sprintf("Published: %s", item.Published.Format("02 Jan 2006 15:04")))

	// Open in Browser button
	openBtn := widget.NewButton("Open in Browser", func() {
		u.openURL(item.Link)
	})

	// Limit content width and enable word wrapping
	content := widget.NewRichTextFromMarkdown(wrapText(item.Content, 80))
	content.Wrapping = fyne.TextWrapWord

	// Create container for article display
	header := container.NewVBox(
		title,
		date,
		openBtn,
		widget.NewSeparator(),
	)

	// Create VBox for all content
	articleContent := container.NewVBox(
		header,
		content,
	)

	// Wrap content in layout with limited width
	maxWidthContainer := container.New(
		NewMaxWidthLayout(650),
		articleContent,
	)

	// Place everything in a scroll container
	scrollContainer := container.NewScroll(maxWidthContainer)

	// Set new content
	u.itemContainer.Objects[0] = scrollContainer
	u.itemContainer.Refresh()
}

// wrapText breaks long strings into shorter ones
func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if len(line) > maxWidth {
			// Don't break text inside HTML tags
			if strings.Contains(line, "<") && strings.Contains(line, ">") {
				continue
			}

			var wrapped strings.Builder
			remaining := line

			for len(remaining) > maxWidth {
				idx := maxWidth
				// Look for a space to break
				for idx > 0 && remaining[idx] != ' ' {
					idx--
				}
				if idx == 0 {
					// If no space, break at maxWidth
					idx = maxWidth
				}

				wrapped.WriteString(remaining[:idx])
				wrapped.WriteString("\n")
				remaining = remaining[idx:]
				if len(remaining) > 0 && remaining[0] == ' ' {
					remaining = remaining[1:]
				}
			}

			wrapped.WriteString(remaining)
			lines[i] = wrapped.String()
		}
	}

	return strings.Join(lines, "\n")
}

// showAddFeedDialog shows the dialog for adding a new feed
func (u *UI) showAddFeedDialog() {
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("https://example.com/rss")

	dialog.ShowForm("Add RSS Feed", "Add", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Feed URL", urlEntry),
		},
		func(confirm bool) {
			if confirm {
				feedURL := urlEntry.Text
				if feedURL != "" {
					u.addFeed(feedURL)
				}
			}
		},
		u.mainWindow,
	)
}

// addFeed adds a new feed by URL
func (u *UI) addFeed(feedURL string) {
	// Check URL validity
	_, err := url.Parse(feedURL)
	if err != nil {
		dialog.ShowError(err, u.mainWindow)
		return
	}

	// Show loading indicator
	progress := dialog.NewProgress("Loading Feed", "Fetching RSS feed...", u.mainWindow)
	progress.Show()

	// Run loading in a separate goroutine
	go func() {
		feed, err := u.feedService.FetchFeed(feedURL)
		if err != nil {
			progress.Hide()
			dialog.ShowError(err, u.mainWindow)
			return
		}

		// Save the feed
		err = u.storage.SaveFeed(feed)
		if err != nil {
			progress.Hide()
			dialog.ShowError(err, u.mainWindow)
			return
		}

		// Update UI
		u.feeds = append(u.feeds, feed)
		u.feedList.Refresh()

		// Select the new feed
		newID := len(u.feeds) - 1
		u.feedList.Select(newID)

		progress.Hide()
	}()
}

// deleteSelectedFeed deletes the selected feed
func (u *UI) deleteSelectedFeed() {
	if u.currentFeed == nil {
		return
	}

	dialog.ShowConfirm("Delete Feed",
		fmt.Sprintf("Are you sure you want to delete '%s'?", u.currentFeed.Title),
		func(confirm bool) {
			if confirm {
				feedID := u.currentFeed.ID
				u.storage.DeleteFeed(feedID)

				// Update feed list
				for i, feed := range u.feeds {
					if feed.ID == feedID {
						u.feeds = append(u.feeds[:i], u.feeds[i+1:]...)
						break
					}
				}

				u.currentFeed = nil
				u.feedList.Refresh()
				u.itemList.Refresh()
				u.itemContainer.Objects[0] = container.NewScroll(widget.NewRichTextFromMarkdown(""))
				u.itemContainer.Refresh()
			}
		},
		u.mainWindow,
	)
}

// refreshFeeds updates all feeds
func (u *UI) refreshFeeds() {
	if len(u.feeds) == 0 {
		return
	}

	// Show loading indicator
	progress := dialog.NewProgress("Refreshing Feeds", "Updating RSS feeds...", u.mainWindow)
	progress.Show()

	go func() {
		for i, oldFeed := range u.feeds {
			newFeed, err := u.feedService.UpdateFeed(oldFeed)
			if err != nil {
				continue
			}

			u.feeds[i] = newFeed
			u.storage.UpdateFeed(newFeed)

			// If this is the current feed, update its display
			if u.currentFeed != nil && u.currentFeed.ID == newFeed.ID {
				u.currentFeed = newFeed
				u.updateItemList()
			}

			progress.SetValue(float64(i+1) / float64(len(u.feeds)))
		}

		u.feedList.Refresh()
		progress.Hide()
	}()
}

// updateItemList updates the article list for the current feed
func (u *UI) updateItemList() {
	u.itemList.Refresh()
}

// openURL opens a link in the default browser
func (u *UI) openURL(urlStr string) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		dialog.ShowError(err, u.mainWindow)
		return
	}

	u.app.OpenURL(parsedURL)
}
