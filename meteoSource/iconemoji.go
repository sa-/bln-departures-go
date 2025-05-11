package meteoSource

func GetEmoji(iconId int) string {
	switch iconId {
	case 1:
		return "❓" // Not available
	case 2:
		return "☀️" // Sunny
	case 3:
		return "🌤️" // Mostly sunny
	case 4:
		return "⛅" // Partly sunny
	case 5:
		return "🌥️" // Mostly cloudy
	case 6:
		return "☁️" // Cloudy
	case 7:
		return "☁️" // Overcast
	case 8:
		return "☁️" // Overcast with low clouds
	case 9:
		return "🌫️" // Fog
	case 10:
		return "🌦️" // Light rain
	case 11:
		return "🌧️" // Rain
	case 12:
		return "🌧️" // Possible rain
	case 13:
		return "🌧️" // Rain shower
	case 14:
		return "⛈️" // Thunderstorm
	case 15:
		return "🌩️" // Local thunderstorms
	case 16:
		return "🌨️" // Light snow
	case 17:
		return "❄️" // Snow
	case 18:
		return "🌨️" // Possible snow
	case 19:
		return "🌨️" // Snow shower
	case 20:
		return "🌨️" // Rain and snow
	case 21:
		return "🌨️" // Possible rain and snow
	case 22:
		return "🌨️" // Rain and snow (duplicate)
	case 23:
		return "🧊" // Freezing rain
	case 24:
		return "🧊" // Possible freezing rain
	case 25:
		return "🌨️" // Hail
	case 26:
		return "🌙" // Clear (night)
	case 27:
		return "🌙" // Mostly clear (night)
	case 28:
		return "☁️" // Partly clear (night)
	case 29:
		return "☁️" // Mostly cloudy (night)
	case 30:
		return "☁️" // Cloudy (night)
	case 31:
		return "☁️" // Overcast with low clouds (night)
	case 32:
		return "🌧️" // Rain shower (night)
	case 33:
		return "🌩️" // Local thunderstorms (night)
	case 34:
		return "🌨️" // Snow shower (night)
	case 35:
		return "🌨️" // Rain and snow (night)
	case 36:
		return "🧊" // Possible freezing rain (night)
	default:
		return "❓" // Unknown weather condition
	}
}
