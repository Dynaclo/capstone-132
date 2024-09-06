package main

//import (
//    "image/color"
//    "log"
//
//
//    "gonum.org/v1/plot"
//    "gonum.org/v1/plot/palette"
//    "gonum.org/v1/plot/plotter"
//
//)
//func main() {
//    data := map[string]float64{
//        "Category A": 30,
//        "Category B": 20,
//        "Category C": 25,
//        "Category D": 25,
//    }
//
//    // Convert map to slices for plotting
//    var labels []string
//    var values []float64
//    for k, v := range data {
//        labels = append(labels, k)
//        values = append(values, v)
//    }
//
//    // Create the pie chart
//    p := plot.New()
//
//
//    p.Title.Text = "Pie Chart Example"
//    p.X.Label.Text = "Category"
//    p.Y.Label.Text = "Value"
//
//    pie, err := plotter.NewPieChart(values)
//    if err != nil {
//        log.Fatalf("could not create pie chart: %v", err)
//    }
//
//    pie.Labels = labels
//    pie.Palette = palette.Palette{
//        color.RGBA{R: 255, G: 0, B: 0, A: 255},  // Red
//        color.RGBA{R: 0, G: 255, B: 0, A: 255},  // Green
//        color.RGBA{R: 0, G: 0, B: 255, A: 255},  // Blue
//        color.RGBA{R: 255, G: 255, B: 0, A: 255},// Yellow
//    }
//
//    p.Add(pie)
//
//    // Save the plot to a PNG file
//    if err := p.Save(6*vg.Inch, 6*vg.Inch, "pie_chart.png"); err != nil {
//        log.Fatalf("could not save plot: %v", err)
//    }
//}
