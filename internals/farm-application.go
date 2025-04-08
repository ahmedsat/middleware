package internals

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/ahmedsat/middleware/helpers"
	"github.com/ahmedsat/utils"
)

// FarmApplication implements Mapper interface to map data from kobo format to erp format
type FarmApplication struct {
	WorkflowState              string  `json:"workflow_state"`
	FarmName                   string  `json:"farm_name"`
	FarmOwnerName              string  `json:"farm_owner_name"`
	WomanFarmOwnerName         string  `json:"woman_farm_owner_name"`
	FarmOwnerNationalID        string  `json:"farm_owner_national_id"`
	WomanFarmOwnerNationalID   string  `json:"woman_farm_owner_national_id"`
	FarmOwnerPhoto             string  `json:"farm_owner_photo"`
	WomanFarmOwnerPhoto        string  `json:"woman_farm_owner_photo"`
	AverageNumberOfChildren    int     `json:"average_number_of_children"`
	RegistrationDate           string  `json:"registration_date"`
	FarmOperator               string  `json:"farm_operator"`
	YearOfReclamation          int     `json:"year_of_reclamation"`
	TotalFarmAreaInFeddan      float64 `json:"total_farm_area_in_feddan"`
	CultivatedFarmAreaInFeddan float64 `json:"cultivated_farm_area_in_feddan"`
	Phone                      string  `json:"phone"`
	Region                     string  `json:"region"`
	CityTown                   string  `json:"citytown"`
	Village                    string  `json:"village"`
	FarmAddress                string  `json:"farm_address"`
	LeadingEngineers           int     `json:"leading_engineers"`
	EngineerName               string  `json:"engineer_name"`
	TotalFarmers               int     `json:"total_farmers"`
	Remarks                    string  `json:"remarks"`
	NamingSeries               string  `json:"naming_series"`

	FarmOwnershipDocument string `json:"farm_ownership_document"`
	FarmSample            string `json:"farm_sample"`

	Attachs       []Attachment    `json:"attachs"`
	LocationTable []LocationEntry `json:"location_table"`
	Workers       []WorkerEntry   `json:"workers"`
	Animals       []AnimalStock   `json:"table_8"`
	Farmers       []FarmerEntry   `json:"farmers"`
}

type Attachment struct {
	Attach string `json:"attach"`
}

type LocationEntry struct {
	LatitudeNums  string `json:"latitude_nums"`
	LongitudeNums string `json:"longitude_nums"`
}

type WorkerEntry struct {
	Worker int    `json:"worker"`
	Age    string `json:"age"`
	Gender string `json:"gender"`
	Parent string `json:"parent"`
}

type AnimalStock struct {
	Animal string `json:"animal"`
	Number int    `json:"number"`
}

type FarmerEntry struct {
	Farmer            string  `json:"farmer"`
	OwnedAreaInFeddan float64 `json:"owned_area_in_feddan"`
	Gender            string  `json:"gender"`
	FarmerPhoto       string  `json:"farmer_photo"`
	FarmerNationalID  string  `json:"farmer_national_id"`
	FarmerPhoneNumber string  `json:"farmer_phone_number"`
}

type Submission struct {
	ID                    int    `json:"_id"`
	Start                 string `json:"start"`
	End                   string `json:"end"`
	Deviceid              string `json:"deviceid"`
	EngineerName          string `json:"engineer_name"`
	EngineerPhoto         string `json:"Engineer_photo"`
	LeadingEngineers      string `json:"leading_engineers"`
	FarmName              string `json:"farm_name"`
	FarmOwner             string `json:"farm_owner"`
	WomenOwnerName        string `json:"women_owner_name"`
	FarmOwnerPhoto        string `json:"farm_owner_photo"`
	WomenOwnerPhoto       string `json:"women_owner_photo"`
	OwnerID               string `json:"owner_id"`
	WomenOwnerID          string `json:"women_owner_id"`
	ChildrenAverg         string `json:"childern_averg"`
	FarmOwnershipDocument string `json:"farm_ownership_document"`
	FarmOwnerPhone        string `json:"farm_owner_phone"`
	Region                string `json:"region"`
	City                  string `json:"city"`
	Village               string `json:"village"`
	RegistrationDate      string `json:"registration_date"`
	FarmAddress           string `json:"farm_address"`
	FarmOperator          string `json:"farm_operator"`
	YearReclamation       string `json:"year_reclamation"`
	FarmCoordinates       string `json:"Farm_coordinates_"`
	FarmArea              string `json:"farm_area"`
	CultivatedArea        string `json:"cultivated_area"`
	FarmOutline           string `json:"farm_outline"`
	Farmer                []struct {
		Name              string `json:"farmer/farmer_name"`
		Area              string `json:"farmer/farmer_area"`
		Gender            string `json:"farmer/farmer_gender"`
		Image             string `json:"farmer/image_jw6yt19"`
		NationalId        string `json:"farmer/file_ei2sh05"`
		FarmerPhoneNumber string `json:"farmer/farmer_phone_number"`
	} `json:"farmer"`
	Animal []struct {
		Name     string `json:"animal/animal_name"`
		Quantity string `json:"animal/quantity"`
	} `json:"animal"`
	OtherDetails string `json:"other_details"`
	Attachments  []struct {
		File string `json:"group_rw4sq75/attachments"`
	} `json:"group_rw4sq75"`
	Workers []struct {
		ID     string `json:"workers/worker"`
		Age    string `json:"workers/worker_age"`
		Gender string `json:"workers/worker_gender"`
	} `json:"workers"`
	AnalysisSample string `json:"analysis_sample"`
	Image          string `json:"image_tj9wu34"`
	SignaturePhoto string `json:"signature_photo"`

	Meta struct {
		InstanceID   string `json:"meta/instanceID"`
		RootUUID     string `json:"meta/rootUuid"`
		DeprecatedID string `json:"meta/deprecatedID"`
	} `json:"meta"`
	Formhub struct {
		UUID string `json:"formhub/uuid"`
	} `json:"formhub"`

	Version string `json:"__version__"`
	XFormID string `json:"_xform_id_string"`
	UUID    string `json:"_uuid"`

	AttachmentsInfo []AttachmentInfo `json:"_attachments"`
}

type AttachmentInfo struct {
	DownloadURL       string `json:"download_url"`
	DownloadLargeURL  string `json:"download_large_url"`
	DownloadMediumURL string `json:"download_medium_url"`
	DownloadSmallURL  string `json:"download_small_url"`
	MimeType          string `json:"mimetype"`
	Filename          string `json:"filename"`
	Instance          int64  `json:"instance"`
	XForm             int64  `json:"xform"`
	ID                int64  `json:"id"`
}

func (f *FarmApplication) Scan(r io.Reader) (err error) {

	Submission := Submission{}

	err = json.NewDecoder(r).Decode(&Submission)
	if err != nil {
		return
	}

	attachmentsMap := map[string]string{}

	for _, attachment := range Submission.AttachmentsInfo {
		fileNameParts := strings.Split(attachment.Filename, "/")
		attachmentsMap[fileNameParts[len(fileNameParts)-1]] = attachment.DownloadURL
		f.Attachs = append(f.Attachs, Attachment{
			Attach: attachment.DownloadURL,
		})
	}

	f.EngineerName = Submission.EngineerName
	if Submission.LeadingEngineers == "OK" {
		f.LeadingEngineers = 1
	}
	f.FarmName = Submission.FarmName
	f.FarmOwnerName = Submission.FarmOwner
	f.WomanFarmOwnerName = Submission.WomenOwnerName
	f.FarmOwnerPhoto = attachmentsMap[Submission.FarmOwnerPhoto]
	f.WomanFarmOwnerPhoto = attachmentsMap[Submission.WomenOwnerPhoto]
	f.FarmOwnerNationalID = attachmentsMap[Submission.OwnerID]
	f.WomanFarmOwnerNationalID = attachmentsMap[Submission.WomenOwnerID]
	f.FarmOwnershipDocument = attachmentsMap[Submission.FarmOwnershipDocument]
	f.Phone = fmt.Sprintf("01%09s", Submission.FarmOwnerPhone)
	f.Region = Submission.Region
	f.CityTown = Submission.City
	f.Village = Submission.Village
	f.RegistrationDate = Submission.RegistrationDate
	f.FarmAddress = Submission.FarmAddress
	f.FarmOperator = Submission.FarmOperator
	f.YearOfReclamation, err = strconv.Atoi(Submission.YearReclamation)
	if err != nil {
		return
	}
	locationParts := strings.Split(Submission.FarmCoordinates, " ")
	f.LocationTable = append(f.LocationTable, LocationEntry{
		LatitudeNums:  locationParts[0],
		LongitudeNums: locationParts[1],
	})

	f.TotalFarmAreaInFeddan, err = strconv.ParseFloat(Submission.FarmArea, 64)
	if err != nil {
		return
	}

	f.CultivatedFarmAreaInFeddan, err = strconv.ParseFloat(Submission.CultivatedArea, 64)
	if err != nil {
		return
	}

	for _, farmer := range Submission.Farmer {
		area, err := strconv.ParseFloat(farmer.Area, 64)
		if err != nil {
			return err
		}
		f.Farmers = append(f.Farmers, FarmerEntry{
			Farmer:            farmer.Name,
			OwnedAreaInFeddan: area,
			Gender:            farmer.Gender,
			FarmerNationalID:  attachmentsMap[farmer.NationalId],
			FarmerPhoto:       attachmentsMap[farmer.Image],
			FarmerPhoneNumber: farmer.FarmerPhoneNumber,
		})
	}

	for _, animal := range Submission.Animal {
		quantity, err := strconv.Atoi(animal.Quantity)
		if err != nil {
			return err
		}
		f.Animals = append(f.Animals, AnimalStock{
			Animal: animal.Name,
			Number: quantity,
		})
	}

	f.Remarks = Submission.OtherDetails

	for _, attach := range Submission.Attachments {
		f.Attachs = append(f.Attachs, Attachment{
			Attach: attachmentsMap[attach.File],
		})
	}

	f.FarmSample = Submission.AnalysisSample
	f.NamingSeries = "From Kobo.DD./.MM./.YYYY.-.####"
	utils.TODO("update workflow status")
	return
}

func (f *FarmApplication) Validate() (err error) {
	sb := strings.Builder{}

	count := 0

	required := func(failedName, failed string) {
		if failed == "" {
			sb.WriteString(fmt.Sprintf("%d. %s is required\n", count+1, failedName))
			count++
		}
	}

	url := strings.ReplaceAll(fmt.Sprintf(
		"/api/resource/Farm Application?filters=[[\"Farm Application\",\"farm_name\",\"=\",\"%s\"]]",
		f.FarmName), " ", "%20")
	res, err := helpers.ERPRequest("GET", url, nil)
	if err != nil {
		return
	}
	if res.StatusCode == 200 {
		if c, _ := io.ReadAll(res.Body); len(c) > 11 {
			sb.WriteString(fmt.Sprintf("%d. Farm name already exists\n", count+1))
			count++
		}
	}

	if len(strings.Split(f.FarmOwnerName, " ")) < 4 {
		sb.WriteString(fmt.Sprintf("%d. Farm owner should be at least 4 words\n", count+1))
		count++
	}

	required("farm owner photo", f.FarmOwnerPhoto)
	required("farm owner national id", f.FarmOwnerNationalID)
	required("farm owner image", f.FarmOwnerPhoto)
	required("registration date", f.RegistrationDate)

	if f.TotalFarmAreaInFeddan < 5 {
		sb.WriteString(fmt.Sprintf("%d. Total farm area should be at last 5 feddan's \n", count+1))
		count++
	}

	if f.TotalFarmAreaInFeddan < f.CultivatedFarmAreaInFeddan {
		sb.WriteString(fmt.Sprintf("%d. Cultivated farm area are greater than total farm area\n", count+1))
		count++
	}

	required("Owner's phone number", f.Phone)
	required("Region", f.Region)
	required("Engineer Name", f.EngineerName)

	if slices.ContainsFunc(f.Farmers, func(farmer FarmerEntry) bool {
		return farmer.Farmer != f.FarmOwnerName
	}) {
		sb.WriteString(fmt.Sprintf("%d. Farm owner should be one of the farmers\n", count+1))
		count++
	}

	if len(f.LocationTable) < 1 {
		sb.WriteString(fmt.Sprintf("%d. Farm location is required\n", count+1))
		count++
	}

	ownedArea := 0.0
	for _, farmer := range f.Farmers {
		ownedArea += farmer.OwnedAreaInFeddan
		if len(strings.Split(farmer.Farmer, " ")) < 4 {
			sb.WriteString(fmt.Sprintf("%d. Farmer name (%s) should be at least 4 words\n", count+1, farmer.Farmer))
			count++
		}
		required("gender", farmer.Gender)
		utils.TODO("require farmers photo")
		required(fmt.Sprintf("farmer %s national id\n", farmer.Farmer), farmer.FarmerNationalID)
		required("farmers phone number", farmer.FarmerPhoneNumber)
	}

	if ownedArea != f.TotalFarmAreaInFeddan {
		sb.WriteString(fmt.Sprintf("%d. area owned by farmers is %0.2f it should be %0.2f\n", count+1, ownedArea, f.TotalFarmAreaInFeddan))
		count++
	}

	if sb.String() != "" {
		return errors.New(sb.String())
	}
	return
}
