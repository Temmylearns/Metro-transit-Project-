package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"strings"
	"time"
)

/*
*Api variables for the routes
 */
type Route struct {
	RouteID    string `json:"route_id"`
	AgencyID   int    `json:"agency_id"`
	RouteLabel string `json:"route_label"`
}

/*
*Api variables for the directions
 */

type Direction struct {
	DirectionID   int    `json:"direction_id"`
	DirectionName string `json:"direction_name"`
}
type Stop struct {
	PlaceCode   string `json:"place_code"`
	Description string `json:"description"`
}
type Departure struct {
	DepartureText string `json:"departure_text"`
	DirectionText string `json:"direction_text"`
}
type DepartureResponse struct {
	Departures []Departure `json:"departures"`
}

//-----------------------
//substring getters
//----------------------

func getRoutes(apiURL string) ([]Route, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data from API: status code %d", resp.StatusCode)
	}

	var routes []Route
	if err := json.NewDecoder(resp.Body).Decode(&routes); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return routes, nil
}
func getDirections(apiURL2 string) ([]Direction, error) {

	resp, err := http.Get(apiURL2)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data from API: status code %d", resp.StatusCode)
	}

	var directions []Direction

	if err := json.NewDecoder(resp.Body).Decode(&directions); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return directions, nil
}
func getStops(apiURL3 string) ([]Stop, error) {

	resp, err := http.Get(apiURL3)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data from API: status code %d", resp.StatusCode)
	}

	var stops []Stop

	if err := json.NewDecoder(resp.Body).Decode(&stops); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return stops, nil
}

func getTime(apiURL4 string) ([]Departure, error) {
	resp, err := http.Get(apiURL4)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data from API: status code %d", resp.StatusCode)
	}

	var response DepartureResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return response.Departures, nil
}

/*
* substring finders
 */
func findRoutes(routes []Route, input string) string {
	for _, route := range routes {
		if route.RouteID == input || route.RouteLabel == input {
			//fmt.Printf("The Route ID for '%s' is: %s, Route Label is: %s\n", input, route.RouteID, route.RouteLabel)
			return route.RouteID
		}
	}
	_ = routes
	fmt.Printf("No matching route found for the input '%s'\n", input)

	return ""
}

/*
* substring finders
@param routeID is the input from the findRoutes method
*/
func findDirection(routeID string) string {
	apiURL2 := fmt.Sprintf("http://svc.metrotransit.org/NexTrip/Directions/%s", routeID)

	directions, err := getDirections(apiURL2)
	_ = directions

	if err != nil {
		log.Fatalf("Error fetching directions: %v", err)
	}
	return ""
}

/*
* substring finders
@param routeID is the input from the findRoutes for the API
@param directionID is the input for the API
this method returns stop and and error if it cant find it.
*/
func findStop(routeID string, directionID int) ([]Stop, error) {
	apiURL3 := fmt.Sprintf("http://svc.metrotransit.org/NexTrip/Stops/%s/%d", routeID, directionID)

	stops, err := getStops(apiURL3)
	if err != nil {
		log.Fatalf("Error fetching stops: %v", err)
	}
	return stops, err

}

/*
@param []Stop  the returns for the find stops APi
@param descriptions the input from the stops to get the Placecode
this method takes stops and the description from the stops to return the place code.
*/
func findPC(stops []Stop, description string) string {
	for _, stop := range stops {
		if strings.Contains(stop.Description, description) {
			fmt.Printf("PlaceCode: %s, Description: %s\n", stop.PlaceCode, stop.Description)
			return stop.PlaceCode
		}
	}

	fmt.Println("wrong Busstop", description)
	return ""
}

/*
* substring finders
@param routeID is the input from the findRoutes for the API
@param directionID is the input for the API
@param PlaceCode is the input from the findPC
this method returns stop and and error if it cant find it.
*/
func findDepartTime(routeID string, directionID int, PlaceCode string) string {

	apiURL4 := fmt.Sprintf("http://svc.metrotransit.org/NexTrip/%s/%d/%s", routeID, directionID, PlaceCode)

	departures, err := getTime(apiURL4)
	if err != nil {
		log.Fatalf("Error fetching departure times: %v", err)

	}
	//_ = departures

	for _, departure := range departures {
		//fmt.Printf("The departure time is '%s'", departure.DepartureText)

		return departure.DepartureText

	}

	return ""
}

func main() {
	p := fmt.Println
	//---------------
	//Time
	//------------------

	// the first API for the routes
	apiURL := "http://svc.metrotransit.org/NexTrip/routes"
	routes, err := getRoutes(apiURL)
	if err != nil {
		log.Fatalf("Error fetching routes: %v", err)
	}
	//---------------
	// Input the route ID or label from the user
	//------------------
	for {
		reader := bufio.NewReader(os.Stdin) // scanner is the bufio
		p("Enter the route ID or label:")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		routeID := findRoutes(routes, input)
		if routeID == "" {
			return
		}

		//---------------
		// Input the directionID from the User
		//------------------

		p("Enter the directionID:")
		var directionID int
		directionInput, _ := reader.ReadString('\n')
		directionInput = strings.TrimSpace(directionInput)
		yes := true
		if yes {
			switch directionInput {
			case "north":
				directionID = 0

			case "south":
				directionID = 1

			default:
				fmt.Println("Unknown direction, defaulting to 0")
				directionID = 0
			}
		}
		//---------------
		// Input the busstop  from the User
		//------------------
		p("Enter the bustop:")
		busstopInput, _ := reader.ReadString('\n')
		busstopInput = strings.TrimSpace(busstopInput)

		stops, err := findStop(routeID, directionID)
		if err != nil {
			log.Fatalf("Error fetching stops: %v", err)
		}
		placeCode := findPC(stops, busstopInput)
		if placeCode == "" {
			return
		}
		//---------------
		// getting the departure time from the findDepartTime method and calculate the diffence from the minutes
		//------------------
		s2 := findDepartTime(routeID, directionID, placeCode)

		/*if s2 == "" {
			fmt.Println("No departures found")
			return
		}*/

		if strings.Contains(s2, "Min") || strings.Contains(s2, "Due") {
			if strings.Contains(s2, "Min") {
				p(fmt.Sprintf("The next bus will arrive in %sutes", s2))
			} else if strings.Contains(s2, "Due") {
				p(fmt.Sprintf("The next bus will arrive in %s", s2))
			}

		} else {

			ss := strings.Split(s2, ":")

			s := strings.Join(ss, "")

			// converting the depauture text to an int
			currentTime := time.Now()
			currentHour := currentTime.Hour()
			currentMin := currentTime.Minute()

			departureTime, err := strconv.Atoi(s)

			if err != nil {
				//  handle error(eyes roll)
				log.Fatalf("Error fetching departures times: %v", err)
			}

			departureHour := departureTime / 100
			departureMin := departureTime % 100

			if currentHour == 0 && departureHour == 12 {
				departureHourAM := 0
				diffHour := (departureHourAM - currentHour)
				diffMinutes := (departureMin - currentMin)

				if diffHour == 0 {
					p("The next bus will arrive in:", diffMinutes, "minutes")
				}

			}
			// if its less than 12 so if its in the morning
			if currentHour < 12 {

				departureHourAM := (departureTime / 100)
				diffHour := (departureHourAM - currentHour)
				diffMinutes := (departureMin - currentMin)

				if diffHour > 0 {
					subReminder := 60 + diffMinutes
					p("The next bus will arrive in:", subReminder, "minutes")

				} else {
					p("The next bus will arrive in:", diffMinutes, "minutes")
				}

			}

			if currentHour == 12 && departureHour == 13 {

				departureHour2 := 1
				//diffHour := (departureHourAM - currentHour)
				diffMinutes := (currentMin - departureMin)

				if departureHour2 > 0 {
					subReminder := 60 - diffMinutes
					p(" The next bus will arrive in:", subReminder, "minutes")

				}
			}

			// if its the greater than 12  if its in the afternoon
			if currentHour > 12 {

				departureHourPM := departureHour + 12
				diffHour := (departureHourPM - currentHour)
				diffMinutes := (departureMin - currentMin)
				if diffHour > 0 {
					subReminder := 60 + diffMinutes
					p("The next bus will arrive in:", subReminder, "minutes")

				} else {
					p("The next bus will arrive in:", diffMinutes, "minutes")
				}
			}
		}
	}

}
