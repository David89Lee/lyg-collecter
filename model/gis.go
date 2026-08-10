package model

func FindWater(minLng, minLat, maxLng, maxLat float64) ([]string, error) {
	query := `
		SELECT ST_AsText(ST_Transform(way, 4326))
		FROM planet_osm_polygon
		WHERE water IS NOT NULL
		AND way && ST_Transform(ST_MakeEnvelope(?,?,?,?, 4326), 3857)
	`
	wktList, err := findWktList(query, minLng, minLat, maxLng, maxLat)
	if err != nil {
		return nil, err
	}
	return wktList, err
}

func FindBuilding(minLng, minLat, maxLng, maxLat float64) ([]string, error) {
	query := `
		SELECT ST_AsText(ST_Transform(way, 4326))
		FROM planet_osm_polygon
		WHERE building IS NOT NULL
		AND way && ST_Transform(ST_MakeEnvelope(?,?,?,?, 4326), 3857)
	`
	wktList, err := findWktList(query, minLng, minLat, maxLng, maxLat)
	if err != nil {
		return nil, err
	}
	return wktList, err
}

func FindHighway(minLng, minLat, maxLng, maxLat float64) ([]string, error) {
	query := `
		SELECT ST_AsText(ST_Transform(way, 4326))
		FROM planet_osm_line
		WHERE highway IS NOT NULL
		AND way && ST_Transform(ST_MakeEnvelope(?,?,?,?, 4326), 3857)
	`
	wktList, err := findWktList(query, minLng, minLat, maxLng, maxLat)
	if err != nil {
		return nil, err
	}
	return wktList, err
}

func FindWetland(minLng, minLat, maxLng, maxLat float64) ([]string, error) {
	query := `
		SELECT ST_AsText(ST_Transform(way, 4326))
		FROM planet_osm_polygon
		WHERE wetland IS NOT NULL
		AND way && ST_Transform(ST_MakeEnvelope(?,?,?,?, 4326), 3857)
	`
	wktList, err := findWktList(query, minLng, minLat, maxLng, maxLat)
	if err != nil {
		return nil, err
	}
	return wktList, err
}

func findWktList(query string, minLng, minLat, maxLng, maxLat float64) ([]string, error) {
	var wktList []string
	err := PostgreSQLDB.Raw(query, minLng, minLat, maxLng, maxLat).Find(&wktList).Error
	if err != nil {
		return nil, err
	}
	return wktList, err
}
