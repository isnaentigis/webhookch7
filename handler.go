package webhook

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/aiteung/atapi"
	"github.com/aiteung/atmessage"
	"github.com/aiteung/module/model"
	"github.com/whatsauth/wa"
)

func PostBalasan(w http.ResponseWriter, r *http.Request) {
	var msg model.IteungMessage
	var resp atmessage.Response
	json.NewDecoder(r.Body).Decode(&msg)
	link := "https://medium.com/@isnaenti16/bikin-whatsauth-api-untuk-whatssapp-notif-dan-share-location-1ae186822e5a"
	if r.Header.Get("Secret") == os.Getenv("SECRET") {
		if msg.Message == "loc" || msg.Message == "Loc" || msg.Message == "lokasi" || msg.LiveLoc {
			location, err := ReverseGeocode(msg.Latitude, msg.Longitude)
			if err != nil {
				// Handle the error (e.g., log it) and set a default location name
				location = "Unknown Location"
			}

			reply := fmt.Sprintf("Hai hai haiii kamu pasti lagi di %s \n Koordinatenya : %s - %s\n Cara Penggunaan WhatsAuth Ada di link dibawah ini"+
				"yaa kak %s\n", location,
				strconv.Itoa(int(msg.Longitude)), strconv.Itoa(int(msg.Latitude)), link)
			dt := &wa.TextMessage{
				To:       msg.Phone_number,
				IsGroup:  false,
				Messages: reply,
			}
			resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")

		} else if msg.Message == "Bacot" || msg.Message == "Taii" || msg.Message == "somplak" {
			dt := &wa.TextMessage{
				To:       msg.Phone_number,
				IsGroup:  false,
				Messages: fmt.Sprintf("heleh %s apaan si ko kasar gitu, jangan kasar dong kakak", msg.Alias_name),
			}
			resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")

		} else if msg.Message == "cantik" || msg.Message == "ganteng" || msg.Message == "cakep" {
			dt := &wa.TextMessage{
				To:       msg.Phone_number,
				IsGroup:  false,
				Messages: fmt.Sprintf("makasiihh kakak %s kamu jugaa cakep kooo", msg.Alias_name),
			}
			resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")

		} else {
			randm := []string{
				"Hai Hai Haiii kamuuuui " + msg.Alias_name + "\nisnaenti lagi sibuk \n aku paijo salam kenall yaaaa \n Cara penggunaan WhatsAuth ada di link berikut ini ya kak...\n" + link,
				"Moon maap nih aku cape",
				"gimana ai kamu",
				"Ihhh kamu cantik banget",
				"Beliin aku makanan kek kali kali",
				"hem bisaa aja kamu",
				"Pagi menjelang sore selamat sore kak",
				"Cihuyyyy hari ini giseuu bahagiaaa banget",
				"Bercandyaaa berrcandyaaaa",
			}
			dt := &wa.TextMessage{
				To:       msg.Phone_number,
				IsGroup:  false,
				Messages: GetRandomString(randm),
			}
			resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")
		}
	} else {
		resp.Response = "Secret Salah"
	}
	fmt.Fprintf(w, resp.Response)
}

func ReverseGeocode(latitude, longitude float64) (string, error) {
	// OSM Nominatim API endpoint
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f", latitude, longitude)

	// Make a GET request to the API
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	// Extract the place name from the response
	displayName, ok := result["display_name"].(string)
	if !ok {
		return "", fmt.Errorf("unable to extract display_name from the API response")
	}

	return displayName, nil
}

func Liveloc(w http.ResponseWriter, r *http.Request) {
	var msg model.IteungMessage
	var resp atmessage.Response
	json.NewDecoder(r.Body).Decode(&msg)

	// Reverse geocode to get the place name
	location, err := ReverseGeocode(msg.Latitude, msg.Longitude)
	if err != nil {
		// Handle the error (e.g., log it) and set a default location name
		location = "Unknown Location"
	}

	reply := fmt.Sprintf("Hai hai haiii kamu pasti lagi di %s \n Koordinatenya : %s - %s\n", location,
		strconv.Itoa(int(msg.Longitude)), strconv.Itoa(int(msg.Latitude)))

	if r.Header.Get("Secret") == os.Getenv("SECRET") {
		dt := &wa.TextMessage{
			To:       msg.Phone_number,
			IsGroup:  false,
			Messages: reply,
		}
		resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://cloud.wa.my.id/api/send/message/text")
	} else {
		resp.Response = "Secret Salah"
	}
	fmt.Fprintf(w, resp.Response)
}

func GetRandomString(strings []string) string {
	randomIndex := rand.Intn(len(strings))
	return strings[randomIndex]
}
