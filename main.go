package main

import (
	"image/color"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var max int = 0
var pulled []int
var duplicateCheck bool = true
var sortOut bool = true
var pulledstr []string

// check if int i is in slice s
func contains(s []int, i int) bool {
	for _, x := range s {
		if x == i {
			return true
		}
	}
	return false
}

func generate(n int) []int {
	var currentSelect []int
	var current int
	for iterator := 0; iterator < n; iterator++ {
		if len(pulled) == max {
			// all generated
			break
		} else {
			// generate rand number
			// if in list already -> again
			current = rand.Intn(max) + 1
			if contains(pulled, current) && duplicateCheck {
				iterator--
			} else {
				currentSelect = append(currentSelect, current)
				pulled = append(pulled, current)
			}
		}
	}
	if sortOut {
		sort.Ints(currentSelect)
	}
	return currentSelect
}

func updateOut(out *widget.Label, in int) string {
	if max <= 0 {
		out.SetText("First set a max amount of numbers in the settings!")
		return ""
	} else {
		generated := generate(in)
		generatedAsString := []string{}

		for i := range generated {
			j := generated[i]
			text := strconv.Itoa(j)
			generatedAsString = append(generatedAsString, text)
		}

		resultGenerated := strings.Join(generatedAsString, "; ")
		resultGeneratedOut := strings.Join(generatedAsString, ";\n")
		if len(resultGenerated) == 0 {
			out.SetText("Max amount of numbers has been generated.")
			return ""
		} else {
			out.SetText(resultGeneratedOut)
			return resultGenerated
		}
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Random Number Generator")

	// TAB 1 SETTINGS

	text1 := canvas.NewText("Set max amount of numbers to generate in total:", color.White)

	input1 := widget.NewEntry()
	input1.SetPlaceHolder("0")

	content1 := container.NewVBox(input1, widget.NewButton("Save", func() {
		p := strings.Split(input1.Text, " ")[0]
		max, _ = strconv.Atoi(p)
	}))

	text2 := canvas.NewText("Don't generate duplicates?", color.White)

	select1 := widget.NewSelect([]string{"Yes", "No"}, func(value string) {
		if value == "Yes" {
			duplicateCheck = true
		} else if value == "No" {
			duplicateCheck = false
		}
	})

	select1.PlaceHolder = "Yes"

	text5 := canvas.NewText("Sort output?", color.White)

	select2 := widget.NewSelect([]string{"Yes", "No"}, func(value string) {
		if value == "Yes" {
			sortOut = true
		} else if value == "No" {
			sortOut = false
		}
	})

	select2.PlaceHolder = "Yes"

	button1 := widget.NewButton("Reset", func() {
		pulled = nil
		pulledstr = nil
	})

	grid1Sub1 := container.New(layout.NewAdaptiveGridLayout(2),
		text1, content1)

	grid1Sub2 := container.New(layout.NewAdaptiveGridLayout(2),
		text2, select1)

	grid1Sub3 := container.New(layout.NewAdaptiveGridLayout(2), text5, select2)

	gridMain1 := container.New(layout.NewAdaptiveGridLayout(1), grid1Sub1, grid1Sub2, grid1Sub3, button1)

	// TAB 2 GENERATOR

	generatedOut := widget.NewLabel("")

	text3 := canvas.NewText("How many numbers shall be generated?", color.White)

	input2 := widget.NewEntry()
	input2.SetPlaceHolder("0")

	content2 := container.NewVBox(input2, widget.NewButton("Generate", func() {
		p, _ := strconv.Atoi(strings.Split(input2.Text, " ")[0])
		pulledstr = append(pulledstr, updateOut(generatedOut, p))
	}))

	grid2Sub1 := container.New(layout.NewAdaptiveGridLayout(2),
		text3, content2)

	gridMain2 := container.New(layout.NewVBoxLayout(),
		grid2Sub1, generatedOut)

	// TAB 3 HISTORY

	pulls := binding.NewString()
	content3 := widget.NewLabelWithData(pulls)
	pulls.Set("All generated numbers:")

	scroll1 := container.NewScroll(content3)

	// Put together

	tabs := container.NewAppTabs(
		container.NewTabItem("Settings", gridMain1),
		container.NewTabItem("Generator", gridMain2),
		container.NewTabItem("History", scroll1),
	)

	w.SetContent(tabs)

	go func() {
		for range time.Tick(time.Second) {
			allPulledAsString := []string{}

			for i := range pulledstr {
				j := "Round " + strconv.Itoa(i+1) + ":     " + pulledstr[i]
				allPulledAsString = append(allPulledAsString, j)
			}

			resultPulled := strings.Join(allPulledAsString, "\n")

			pulls.Set(resultPulled)
		}
	}()

	w.ShowAndRun()
}
