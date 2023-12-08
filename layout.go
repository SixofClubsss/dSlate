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
		contractEdit())

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

	walletSelect := container.NewVBox(
		rpcWalletEdit(),
		rpcLoginEdit())

	walletConnect := container.NewVBox(
		rpcConnectButton(),
	)

	wallet := container.NewVBox(
		walletSelect,
		walletConnect,
	)

	box1 := container.NewVBox(
		daemonSelectOption(),
		heightDisplay())

	box2 := container.NewVBox(
		daemonConnectBox(),
		walletConnectBox(),
		debugEnabled)

	box3 := container.NewVBox(
		layout.NewSpacer(),
		wallet,
		box1,
		box2,
		layout.NewSpacer(),
	)
	box4 := container.NewHBox(
		layout.NewSpacer(),
		box3,
		layout.NewSpacer(),
	)

	box5 := container.NewVBox(
		layout.NewSpacer(),
		builtOnImage(),
		container.NewGridWithColumns(3,
			layout.NewSpacer(),
			balanceDisplay(),
			layout.NewSpacer(),
		),
		cardImage(),

		layout.NewSpacer())

	// box6 := container.NewHBox(
	// 	layout.NewSpacer(),
	// 	box5,
	// 	layout.NewSpacer(),
	// )
	// search_content := container.NewVBox(
	// 	contractEdit(),
	// 	searchButton(),
	// 	layout.NewSpacer(),
	// 	contractCode(),
	//
	// 	enableGnomon(),
	// 	gnomonOpts())

	// scroll := container.NewScroll(search_content)

	tabs := container.NewAppTabs(
		container.NewTabItem("DERO Bags", box5),
		container.NewTabItem("Settings", box4),
		// container.NewTabItem("DERO Address", imageContent),
		// container.NewTabItem("Blank Slate", blankWidget()),

		//	container.NewTabItem("Search", scroll),

	//	container.NewTabItem("NFA", nfaOpts()),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	return tabs

}
