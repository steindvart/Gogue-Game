package cli

import "github.com/gdamore/tcell/v2"

// @todo - перенести остальные цвета

const (
	// Базовые цвета
	colorBlack      = "#000000"
	colorWhite      = "#FFFFFF"
	colorGray       = "#808080"
	colorStrongGray = "#515152"
	colorLightGray  = "#E0E0E0"
	colorSilver     = "#C0C0C0"
	colorGold       = "#FFD700"
	colorOrange     = "#FFA500"
	colorRed        = "#FF0000"
	colorTomato     = "#FF6347"
	colorGreen      = "#00FF00"
	colorLimeGreen  = "#32CD32"
	colorLightGreen = "#90EE90"
	colorPurple     = "#FF00FF"
	colorMedPurple  = "#9370DB"
	colorRoyalBlue  = "#4169E1"
	colorSkyBlue    = "#87CEEB"
	colorAqua       = "#0194A7"
	colorTurquoise  = "#00CED1"
	colorBrown      = "#8B4513"
)

const (
	statColorTreasures  = colorGold
	statColorLevel      = colorSkyBlue
	statColorEnemies    = colorRed
	statColorFood       = colorOrange
	statColorElixirs    = colorMedPurple
	statColorScrolls    = colorTurquoise
	statColorHitsDealt  = colorLimeGreen
	statColorHitsMissed = colorGray
	statColorCellsMoved = colorSilver
)

func StatTcellColor(hex string) tcell.Color {
	return tcell.GetColor(hex)
}
