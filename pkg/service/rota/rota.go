package rota

import "math"

const raioTerraKm = 6371.0

func grausParaRadianos(graus float64) float64 {
	return graus * math.Pi / 180.0
}

// CalcularDistancia calcula a distância entre dois pontos especificados pelas suas latitudes e longitudes.
// As coordenadas devem estar em graus decimais.
func CalcularDistancia(lat1, lon1, lat2, lon2 float64) float64 {
	// Converter coordenadas para radianos
	lat1Rad := grausParaRadianos(lat1)
	lon1Rad := grausParaRadianos(lon1)
	lat2Rad := grausParaRadianos(lat2)
	lon2Rad := grausParaRadianos(lon2)

	// Calcular diferenças nas coordenadas
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	// Fórmula de Haversine
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distancia := raioTerraKm * c

	return distancia
}
