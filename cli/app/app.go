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
	fmt.Printf("Using seed %d\n", seed)

	island := rogue.GenerateIsland(seed, 220, 100)

	/*
		islandColImgFile, _ := os.Create("island_col.png")
		png.Encode(islandColImgFile, islandColImage)

		islandImage := island.HeightMapAsImage()
		islandImgFile, _ := os.Create("island.png")
		png.Encode(islandImgFile, islandImage)
	*/

	island.FillWithCritters()
	for i := 0; i < 10; i++ {
		island.Tick()
	}

	theme := flags.CreateTheme(driver)

	//	font := createMonospaceFont(25, &driver)

	window := theme.CreateWindow(800, 600, "rogue")
	window.SetBackgroundBrush(gxui.CreateBrush(gxui.Gray50))
	window.SetScale(flags.DefaultScaleFactor)

	/*

		// XXX position!?
		label := theme.CreateLabel()
		label.SetFont(font)
		label.SetText("Hello world")
		window.AddChild(label)
	*/

	splitterAB := theme.CreateSplitterLayout()
	splitterAB.SetOrientation(gxui.Horizontal)
	splitterAB.AddChild(topLeftPanelHolder(&theme, &driver, &window, &island))
	splitterAB.AddChild(panelHolder("B", theme))

	splitterCD := theme.CreateSplitterLayout()
	splitterCD.SetOrientation(gxui.Horizontal)
	splitterCD.AddChild(panelHolder("C", theme))
	splitterCD.AddChild(panelHolder("D", theme))

	vSplitter := theme.CreateSplitterLayout()
	vSplitter.SetOrientation(gxui.Vertical)
	vSplitter.AddChild(splitterAB)
	vSplitter.AddChild(splitterCD)
	window.AddChild(vSplitter)

	window.OnClose(driver.Terminate)
}

func createMonospaceFont(size int, driver *gxui.Driver) gxui.Font {

	p := *driver

	font, err := p.CreateFont(gxfont.Monospace, size)
	if err != nil {
		panic(err)
	}
	return font
}

// Create a PanelHolder with a 3 panels
func topLeftPanelHolder(theme *gxui.Theme, driver *gxui.Driver, window *gxui.Window, island *rogue.Island) gxui.PanelHolder {

	name := "top left"
	label := func(text string) gxui.Label {
		label := (*theme).CreateLabel()
		label.SetText(text)
		return label
	}

	islandColImage := island.ColoredHeightMapAsImage()
	img := (*theme).CreateImage()
	texture := (*driver).CreateTexture(islandColImage, 1)
	img.SetTexture(texture)
	(*window).AddChild(img)

	holder := (*theme).CreatePanelHolder()
	holder.AddPanel(label(name+" 0 content"), name+" 0 panel")
	holder.AddPanel(label(name+" 1 content"), name+" 1 panel")
	return holder
}

// Create a PanelHolder with a 3 panels
func panelHolder(name string, theme gxui.Theme) gxui.PanelHolder {
	label := func(text string) gxui.Label {
		label := theme.CreateLabel()
		label.SetText(text)
		return label
	}

	holder := theme.CreatePanelHolder()
	holder.AddPanel(label(name+" 0 content"), name+" 0 panel")
	holder.AddPanel(label(name+" 1 content"), name+" 1 panel")
	holder.AddPanel(label(name+" 2 content"), name+" 2 panel")
	return holder
}
