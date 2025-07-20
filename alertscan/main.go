package main

import (
	"encoding/json"
	"flag"
	"log"
	"math"
	"os"
	"time"
)

type Alert struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Service     string    `json:"service"`
	Component   string    `json:"component"`
	Severity    string    `json:"severity"`
	Metric      string    `json:"metric"`
	Value       int       `json:"value"`
	Threshold   int       `json:"threshold"`
	Description string    `json:"description"`
}

func parseTimeArg(t string) time.Time {
	r, err := time.Parse(time.RFC3339, t)
	if err != nil {
		log.Fatalf("unable to parse time: %v", err)
	}
	return r
}

func parseSeverity(s string) int {
	var r int
	switch s {
	case "critical":
		r = 10
	case "warning":
		r = 5
	case "info":
		r = 1
	default:
		r = 0
	}
	return r
}

func calculateDeviation(value, threshold int) float64 {
	if threshold == 0 {
		return 0
	}
	return ((float64(value - threshold)) / float64(threshold)) * 100
}

func calculateWeightedPriority(severity, value, threshold, affectedComponents int) float64 {
	deviation := calculateDeviation(value, threshold)
	normalizedSeverity := float64(severity) / float64(10)
	affectedComponentsScore := math.Min(float64(affectedComponents)/10, 1.0)

	severityComponent := normalizedSeverity * 0.4            // 40% on severity
	deviationComponent := math.Min(deviation/100, 1.0) * 0.4 // 40% on deviation
	componentsComponent := affectedComponentsScore * 0.2     // 20% on components affected

	// Calculate final priority score (0-1 range)
	priority := severityComponent + deviationComponent + componentsComponent
	return priority
}

func groupRelatedAlerts(alerts []Alert) map[string][]Alert {
	r := make(map[string][]Alert)
	for _, alert := range alerts {
		r[alert.Component] = append(r[alert.Component], alert)
	}
	return r
}

func loadAlertsFromFile(filePath, filterService string, minSev int, filterStart, filterEnd time.Time) []Alert {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	var data struct {
		Alerts []Alert `json:"alerts"`
	}

	err = json.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatalf("unable to parse alerts file: %v", err)
	}

	var alerts []Alert

	for _, alert := range data.Alerts {
		if filterService != "" && alert.Service != filterService {
			// doesn't match filter, skip
			continue
		}
		if parseSeverity(alert.Severity) < minSev {
			// alert severity lower than minSev
			continue
		}
		if !filterStart.IsZero() && filterStart.Compare(alert.Timestamp) < 0 {
			// alert happens before start
			continue
		}
		if !filterEnd.IsZero() && filterEnd.Compare(alert.Timestamp) > 0 {
			//alert happens later than finish
			continue
		}
		alerts = append(alerts, alert)
	}

	return alerts
}

func main() {
	log.Println("Starting alertscan")
	filePathArg := flag.String("file", "alerts.json", "Path to alerts file")
	filterStartArg := flag.String("start", "", "Filter time start")
	filterEndArg := flag.String("end", "", "Filter time end")
	filterServiceArg := flag.String("service", "", "Filter service name (empty means no filter)")
	filterSeverityArg := flag.String("severity", "debug", "Minimal alert severity (critical|warning|info|debug)")
	debugArg := flag.Bool("debug", false, "enable debug")
	flag.Parse()

	var filterStart, filterEnd time.Time

	if *filterStartArg != "" {
		filterStart = parseTimeArg(*filterStartArg)
	}

	if *filterEndArg != "" {
		filterEnd = parseTimeArg(*filterEndArg)
	}

	if filterStart.Compare(filterEnd) > 0 {
		log.Fatalf("Scan filter ends earlier (%s) than starts (%s), please check the input", filterEnd, filterStart)
	}

	if *debugArg {
		log.Print("File path:", *filePathArg)
		log.Printf("Parsing log from %s to %s", filterStart, filterEnd)
		if *filterSeverityArg != "" {
			log.Println("Severity filter:", *filterSeverityArg)
		}
		if *filterServiceArg != "" {
			log.Println("Service filter:", *filterServiceArg)
		}

	}

	filterMinSeverity := parseSeverity(*filterSeverityArg)

	alerts := loadAlertsFromFile(*filePathArg, *filterServiceArg, filterMinSeverity, filterStart, filterEnd)
	groupedAlerts := groupRelatedAlerts(alerts)

	log.Printf("Found %d alerts", len(alerts))

	log.Println("Alerts by component:")
	for k, v := range groupedAlerts {
		log.Printf("[%s] %d alerts", k, len(v))
	}

}
