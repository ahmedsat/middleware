package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	data        = map[string]any{}
	out         = map[string]any{}
	attachments = map[string]string{}
	errorsList  = []string{}
)

func CreateFarmApplicationFromKoboData(koboDataFile string) []string {

	kf, err := os.Open(koboDataFile)
	if err != nil {
		errorsList = append(errorsList, err.Error())
		return errorsList
	}
	defer kf.Close()

	of, err := os.Create("out.json")
	if err != nil {
		errorsList = append(errorsList, err.Error())
		return errorsList
	}
	defer of.Close()

	err = json.NewDecoder(kf).Decode(&data)
	if err != nil {
		errorsList = append(errorsList, err.Error())
		return errorsList
	}

	err = parsAttachments()
	if err != nil {
		errorsList = append(errorsList, err.Error())
		return errorsList
	}

	IsUnique(copyM2MS("farm_name", "farm_name"))
	if len(copyM2MS("farm_owner_name", "farm_owner")) < 4 {
		errorsList = append(errorsList, "farm owner name is too short\n")
	}
	if len(copyM2MS("woman_farm_owner_name", "women_owner_name")) < 4 {
		errorsList = append(errorsList, "farm owner's wife name is too short\n")
	}
	if copyM2MAtt("farm_owner_national_id", "owner_id") == "" {
		errorsList = append(errorsList, "farm owner national id is empty\n")
	}
	if copyM2MAtt("farm_owner_photo", "farm_owner_photo") == "" {
		errorsList = append(errorsList, "farm owner photo is empty\n")
	}
	if copyM2MAtt("woman_farm_owner_photo", "women_owner_photo") == "" {
		errorsList = append(errorsList, "farm owner's wife photo is empty\n")
	}
	copyM2MI("average_number_children", "childern_averg")
	copyM2MS("registration_date", "registration_date")
	copyM2MS("farm_operator", "farm_operator")
	copyM2MI("year_of_reclamation", "year_reclamation")
	copyM2MF("total_farm_area_in_feddan", "farm_area")
	copyM2MF("cultivated_farm_area_in_feddan", "cultivated_area")
	// TODO: fix phone number format
	copyM2MS("phone", "farm_owner_phone")
	copyM2MS("region", "region")
	copyM2MS("citytown", "city")
	copyM2MS("village", "village")
	copyM2MS("farm_address", "farm_address")
	copyM2MCheckBox("leading_engineers", "leading_engineers")
	copyM2MS("engineer_name", "engineer_name")
	copyM2MS("remarks", "other_details")
	copyM2MLatLong("location_table", "Farm_coordinates_")
	copyM2MFarmers("farmers", "farmer")
	copyM2MAnimals("table_8", "animal")
	copyM2MWorkers("workers", "workers")

	err = json.NewEncoder(of).Encode(out)
	if err != nil {
		errorsList = append(errorsList, err.Error())
		return errorsList
	}
	return errorsList
}

func IsUnique(farmName string) {
	url := strings.ReplaceAll(fmt.Sprintf(
		"/api/resource/Farm Application?filters=[[\"Farm Application\",\"farm_name\",\"=\",\"%s\"]]",
		farmName), " ", "%20")

	res, err := ERPRequest("GET", url, nil)
	if err != nil {
		errorsList = append(errorsList, err.Error())
		return
	}
	fmt.Println(res.StatusCode)
	if res.StatusCode == 200 {
		errorsList = append(errorsList, "farm name already exists\n")
	}
}

func parsAttachments() (err error) {
	att, ok := data["_attachments"]
	if !ok {
		return errors.New("input map doesn't has attachments")
	}

	attList, ok := att.([]any)
	if !ok {
		return fmt.Errorf(`expected: %T, got: %T`, attList, att)
	}

	for i, att := range attList {

		attMap, ok := att.(map[string]any)
		if !ok {
			return fmt.Errorf(`expected: %T, got: %T`, attMap, att)
		}

		download_url, ok := attMap["download_url"]
		if !ok {
			return fmt.Errorf("attachment(%d) doesn't has download_url", i)
		}

		download_str, ok := download_url.(string)
		if !ok {
			return fmt.Errorf(`expected: %T, got: %T`, download_str, download_url)
		}

		strs := strings.Split(download_str, "%2F")
		filename := strs[len(strs)-1]
		attachments[filename] = download_str

	}

	return
}

func copyM2MS(dstKey, srcKey string) (result string) {
	val, ok := data[srcKey]
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("source map doesn't has key: %s\n", srcKey))
		return
	}

	result, ok = val.(string)
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", result, val))
		return
	}

	out[dstKey] = result
	return
}

func copyM2MI(dstKey, srcKey string) (val int) {
	data, ok := data[srcKey]
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("source map doesn't has key: %s\n", srcKey))
		return
	}

	// i, ok := val.(int)
	// if !ok {
	// 	i, err := strconv.Atoi(val.(string))
	// }
	switch v := data.(type) {
	case int:
		val = v
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", i, val))
			return
		}
		val = i
	default:
		errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", val, data))
		return
	}

	out[dstKey] = val
	return
}

func copyM2MF(dstKey, srcKey string) (val any) {
	val, ok := data[srcKey]
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("source map doesn't has key: %s\n", srcKey))
		return
	}

	var f float64
	switch v := val.(type) {
	case float64:
		f = v
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", f, val))
			return
		}
	default:
		errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", f, val))
		return
	}

	out[dstKey] = f
	return
}

func copyM2MAtt(dstKey, srcKey string) (val string) {
	fName, ok := data[srcKey]
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("input map doesn't has %s\n", srcKey))
		return
	}
	fNameStr, ok := fName.(string)
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", fNameStr, fName))
		return
	}

	val, ok = attachments[fNameStr]
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("attachment list doesn't has %s => %s\n", srcKey, fNameStr))
		return
	}

	out[dstKey] = val

	return
}

func copyM2MCheckBox(dstKey, srcKey string) (val any) {
	val, ok := data[srcKey]
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("source map doesn't has key: %s\n", srcKey))
		return
	}
	out[dstKey] = Ternary(reflect.DeepEqual(val, "OK"), 1, 0)
	return
}

func copyM2MLatLong(dstKey, srcKey string) (val any) {

	val, ok := data[srcKey]
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("source map doesn't has key: %s\n", srcKey))
		return
	}

	valStr, ok := val.(string)
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", valStr, val))
		return
	}

	strs := strings.Split(valStr, " ")
	if len(strs) != 4 {
		errorsList = append(errorsList, fmt.Sprintf("expected: 4, got: %d\n", len(strs)))
		return
	}

	// template := `{"latitude": "%s","longitude": "%s"}`

	out[dstKey] = []map[string]any{
		{
			"latitude_nums":  strs[0],
			"longitude_nums": strs[1],
		},
	}
	return
}

func copyM2MFarmers(dstKey, srcKey string) (val any) {
	val, ok := data[srcKey]
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("source map doesn't has key: %s\n", srcKey))
		return
	}

	list, ok := val.([]any)
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", list, val))
		return
	}

	farmers := []any{}

	for i, farmer := range list {
		farmerMap, ok := farmer.(map[string]any)
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", farmerMap, farmer))

			return ok
		}

		farmerName, ok := farmerMap["farmer/farmer_name"]
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("farmer(%d) doesn't has farmer/farmer_name\n", i))
			return ok
		}

		farmerArea, ok := farmerMap["farmer/farmer_area"]
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("farmer(%d) doesn't has farmer/farmer_area\n", i))
			return ok
		}

		farmerGender, ok := farmerMap["farmer/farmer_gender"]
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("farmer(%d) doesn't has farmer/farmer_gender\n", i))
			return ok
		}

		farmerImageName, ok := farmerMap["farmer/image_jw6yt19"]
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("farmer(%d) doesn't has farmer/image_jw6yt19\n", i))
			return ok
		}

		farmerImage, ok := attachments[farmerImageName.(string)]
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("farmer(%d) doesn't has farmer/image_jw6yt19\n", i))
			return ok
		}

		farmerIdName, ok := farmerMap["farmer/file_ei2sh05"]
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("farmer(%d) doesn't has farmer/file_ei2sh05\n", i))
			return ok
		}

		farmerId, ok := attachments[farmerIdName.(string)]
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("farmer(%d) doesn't has farmer/file_ei2sh05\n", i))
			return ok
		}

		var farmerAreaF float64
		switch v := farmerArea.(type) {
		case float64:
			farmerAreaF = v
		case string:
			farmerAreaF, err := strconv.ParseFloat(v, 64)
			if err != nil {
				errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", farmerAreaF, farmerArea))
				return false
			}
		default:
			errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", farmerAreaF, farmerArea))
			return false
		}

		// template := "{'farmer': '%s','owned_area_in_feddan': %f,'gender': '%s','farmer_photo': '%s','farmer_id': '%s'}"
		farmers = append(farmers, map[string]any{
			"farmer":               farmerName,
			"owned_area_in_feddan": farmerAreaF,
			"gender":               farmerGender,
			"farmer_photo":         farmerImage,
			"farmer_id":            farmerId,
		})

	}
	out[dstKey] = farmers
	return
}

func copyM2MAnimals(dstKey, srcKey string) (val any) {
	val, ok := data[srcKey]
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("source map doesn't has key: %s\n", srcKey))
		return
	}

	list, ok := val.([]any)
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", list, val))
		return
	}

	animals := []any{}

	for i, animal := range list {
		animalMap, ok := animal.(map[string]any)
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", animalMap, animal))

			return ok
		}

		animalName, ok := animalMap["animal/animal_name"]
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("animal(%d) doesn't has animal/animal_name\n", i))
			return ok
		}

		quantity, ok := animalMap["animal/quantity"]
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("animal(%d) doesn't has animal/animal_area\n", i))
			return ok
		}

		var quantityI int
		switch v := quantity.(type) {
		case int:
			quantityI = v
		case string:
			quantityI, err := strconv.Atoi(v)
			if err != nil {
				errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", quantityI, quantity))
				return false
			}
		default:
			errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", quantityI, quantity))
			return false
		}

		// template := `{"animal": "%s","number": %d}`
		animals = append(animals, map[string]any{
			"animal": animalName,
			"number": quantityI,
		})

	}
	out[dstKey] = animals
	return
}

func copyM2MWorkers(dstKey, srcKey string) (val any) {

	val, ok := data[srcKey]
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("source map doesn't has key: %s\n", srcKey))
		return
	}

	list, ok := val.([]any)
	if !ok {
		errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", list, val))
		return
	}

	Workers := []any{}

	for i, worker := range list {
		workerMap, ok := worker.(map[string]any)
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("expected: %T, got: %T\n", workerMap, worker))

			return ok
		}

		workerId, ok := workerMap["workers/worker"]
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("worker(%d) doesn't has workers/worker\n", i))
		}

		workerAge, ok := workerMap["workers/worker_age"]
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("worker(%d) doesn't has workers/age\n", i))
			return ok
		}

		workerGender, ok := workerMap["workers/worker_gender"]
		if !ok {
			errorsList = append(errorsList, fmt.Sprintf("worker(%d) doesn't has workers/gender\n", i))
			return ok
		}

		//	template := `{"worker": %d,"age": "%s","gender": "%s"}`
		Workers = append(Workers, map[string]any{
			"worker": workerId,
			"age":    workerAge,
			"gender": workerGender,
		})
	}

	out[dstKey] = Workers
	return
}

func Ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
