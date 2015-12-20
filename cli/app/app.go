package main

import (
	"fmt"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/gxfont"
	"github.com/google/gxui/samples/flags"
	"github.com/martinlindhe/rogue"
)

func main() {
	gl.StartDriver(appMain)
}

func appMain(driver gxui.Driver) {

	//seed := time.Now().Unix()
	seed := int64(1450549167)

	island := rogue.GenerateIsland(seed, 220, 140)
	island.FillWithCritters()
	/*
		islandColImgFile, _ := os.Create("island_col.png")
		png.Encode(islandColImgFile, islandColImage)

		islandImage := island.HeightMapAsImage()
		islandImgFile, _ := os.Create("island.png")
		png.Encode(islandImgFile, islandImage)
	*/
	/*
		for i := 0; i < 10; i++ {
			island.Tick()
		}
	*/
	theme := flags.CreateTheme(driver)

	window := theme.CreateWindow(800, 600, "rogue")
	window.SetBackgroundBrush(gxui.CreateBrush(gxui.Gray50))
	window.SetScale(flags.DefaultScaleFactor)

	font, _ := driver.CreateFont(gxfont.Monospace, 25)

	splitterAB := theme.CreateSplitterLayout()
	splitterAB.SetOrientation(gxui.Horizontal)
	splitterAB.AddChild(topLeftPanelHolder(&theme, &driver, &island))
	splitterAB.AddChild(panelHolder("B", &theme, &font))

	splitterCD := theme.CreateSplitterLayout()
	splitterCD.SetOrientation(gxui.Horizontal)
	splitterCD.AddChild(panelHolder("C", &theme, &font))
	splitterCD.AddChild(panelHolder("D", &theme, &font))

	vSplitter := theme.CreateSplitterLayout()
	vSplitter.SetOrientation(gxui.Vertical)
	vSplitter.AddChild(splitterAB)
	vSplitter.AddChild(splitterCD)
	window.AddChild(vSplitter)

	window.OnClose(driver.Terminate)
}

// Create a PanelHolder with a 3 panels
func topLeftPanelHolder(theme *gxui.Theme, driver *gxui.Driver, island *rogue.Island) gxui.PanelHolder {

	label := func(text string) gxui.Label {
		label := (*theme).CreateLabel()
		font, _ := (*driver).CreateFont(gxfont.Monospace, 25)
		label.SetFont(font)
		label.SetText(text)
		return label
	}

	// tab 1: map
	islandColImage := island.ColoredHeightMapAsImage()
	img := (*theme).CreateImage()
	texture := (*driver).CreateTexture(islandColImage, 1)
	img.SetTexture(texture)

	// tab 2: spawn list XXXXX
	mapSettings := fmt.Sprintf("seed: %d", island.Seed)

	holder := (*theme).CreatePanelHolder()
	holder.AddPanel(img, "map")
	holder.AddPanel(label(mapSettings), "map settings")
	return holder
}

// Create a PanelHolder with a 3 panels
func panelHolder(name string, theme *gxui.Theme, font *gxui.Font) gxui.PanelHolder {
	label := func(text string) gxui.Label {
		label := (*theme).CreateLabel()
		label.SetFont(*font)
		label.SetText(text)
		return label
	}

	holder := (*theme).CreatePanelHolder()
	holder.AddPanel(label(name+" 0 content"), name+" 0 panel")
	holder.AddPanel(label(name+" 1 content"), name+" 1 panel")
	holder.AddPanel(label(name+" 2 content"), name+" 2 panel")
	return holder
}
