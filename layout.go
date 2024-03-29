package main

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// Organize content without layout
func placeWithout() *container.AppTabs {
	settings_content := container.NewWithoutLayout(
		rpcWalletEdit(),
		rpcLoginEdit(),
		rpcConnectButton(),
		daemonSelectOption(),
		daemonConnectBox(),
		walletConnectBox(),
		heightDisplay(),
		balanceDisplay(),
		builtOnImage())

	search_content := container.NewWithoutLayout(
		searchButton(),
		contractEntry())

	scroll := container.NewScroll(search_content)
	image_content := container.NewWithoutLayout(cardImage())

	tabs := container.NewAppTabs(
		container.NewTabItem("Settings", settings_content),
		container.NewTabItem("Search", scroll),
		container.NewTabItem("Blank Slate", blankWidget()),
		container.NewTabItem("?", image_content),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	return tabs
}

// Organize content with layout
func placeWith() *container.AppTabs {

	wallet := container.NewBorder(
		nil,
		nil,
		nil,
		rpcConnectButton(),
		rpcWalletEdit())

	checks := container.NewHBox(
		daemonConnectBox(),
		walletConnectBox())

	grid := container.NewAdaptiveGrid(2,
		heightDisplay(),
		balanceDisplay())

	hbox := container.NewAdaptiveGrid(2,
		daemonSelectOption(),
		debugEnabled)

	vbox := container.NewVBox(hbox, grid, checks, rpcLoginEdit(), wallet)

	settings_content := container.NewBorder(
		nil,
		vbox,
		nil,
		nil,
		layout.NewSpacer(),
		container.NewCenter(container.NewStack(builtOnImage())))

	search_content := container.NewVBox(
		contractEntry(),
		searchButton(),
		layout.NewSpacer(),
		codeButton(),
		layout.NewSpacer(),
		updateButton(),
		layout.NewSpacer(),
		enableGnomon(),
		gnomonOpts())

	scroll := container.NewScroll(search_content)
	imageContent := container.NewWithoutLayout(cardImage())

	tabs := container.NewAppTabs(
		container.NewTabItem("Settings", settings_content),
		container.NewTabItem("Search", scroll),
		container.NewTabItem("Blank Slate", blankWidget()),
		container.NewTabItem("?", imageContent),
		container.NewTabItem("NFA", nfaOpts()),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	return tabs

}
