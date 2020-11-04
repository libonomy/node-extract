package controllers

import (
	"crypto/rsa"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/libonomy/node-extract/constants"
	"github.com/libonomy/node-extract/utils/helper"

	"github.com/libonomy/libonomy-gota/dataframe"
	"github.com/libonomy/node-extract/dto"
	"github.com/libonomy/node-extract/models"
	"github.com/libonomy/node-extract/utils"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

type bodyVariables struct {
	MachineID     string  `json:"machineId"`
	ComputerPower float64 `json:"computerPower"`
	DownloadSpeed float64 `json:"downSpeed"`
	Ylabels       string  `json:"yLabels"`
}

type check1 struct {
	Name string
}

type check2 struct {
	Name string
}

type check3 struct {
	Name string
}

type fileStats struct {
	LabelIndex int `json:"labelIndex"`
}

//Testing Function To Test its working
func Testing(w http.ResponseWriter, r *http.Request) {
	s := "This is for testing function only"

	dto.SendResponse(w, r, http.StatusOK, "Success", map[string]interface{}{"testing": s})
}

//GenerateCSV function to generate csv file from json data.
func GenerateCSV(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	data := []byte(body)
	fmt.Println(string(body))
	decData := helper.DecryptDatafromClient(data)

	//variables := []bodyVariables{}
	fileVariables := []bodyVariables{}
	var variables []bodyVariables
	json.Unmarshal(decData, &variables)
	filename := "./datasets/testing/dummyDataset.json"
	csvFilename := "./datasets/testing/dummyDataset.csv"

	/* err := json.NewDecoder(r.Body).Decode(&variables)
	if err != nil {
		fmt.Println("There is some Error in Decoding Body Request", err.Error())
		dto.SendResponse(w, r, http.StatusInternalServerError, "Bad", map[string]interface{}{"Error": err.Error()})
		return
	} */

	//fmt.Println(variables)
	_, err := os.Stat(filename)

	if err != nil {
		//checking if file does not exists
		fmt.Println("Error is os.stat", err.Error())
		//fmt.Println("File Info ", fileInfo)

		jsonFile, _ := os.Create(filename)
		csvFile, _ := os.Create(csvFilename)

		defer jsonFile.Close()
		file, _ := json.MarshalIndent(variables, "", " ")
		_ = ioutil.WriteFile(filename, file, 0644)

		fileOpened, _ := os.Open(filename)
		dataFrame := dataframe.ReadJSON(fileOpened)
		csvWriter := csv.NewWriter(csvFile)
		csvWriter.WriteAll(dataFrame.Records())
		// csvData := dataframe.ReadJSON(jsonFile)
		// csvData.WriteCSV(csvFile)
		dto.SendResponse(w, r, http.StatusOK, "Successfully Created CSV File From JSON", nil)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Cannot create file", nil)
	}
	fileBytes, _ := ioutil.ReadAll(file)
	json.Unmarshal(fileBytes, &fileVariables)
	var alreadyExistsCheck bool = false
	// fmt.Println("Already Exists Check ", alreadyExistsCheck)

	for index, element := range fileVariables {
		if element.MachineID == variables[0].MachineID {
			alreadyExistsCheck = true
			// fmt.Println("Machine Id Already Exists")
			fileVariables[index].ComputerPower = variables[0].ComputerPower
			fileVariables[index].DownloadSpeed = variables[0].DownloadSpeed
			fileVariables[index].Ylabels = variables[0].Ylabels
		}
		// fmt.Println("Index", index, "Data", element)
	}

	if alreadyExistsCheck == true {
		alreadyExistsCheck = false
		writeF, _ := json.MarshalIndent(fileVariables, "", " ")
		_ = ioutil.WriteFile(filename, writeF, 0644)

		csvFile, _ := os.Create(csvFilename)
		fileOpened, _ := os.Open(filename)
		dataFrame := dataframe.ReadJSON(fileOpened)
		csvWriter := csv.NewWriter(csvFile)
		csvWriter.WriteAll(dataFrame.Records())

		dto.SendResponse(w, r, http.StatusOK, "Successfully Updated CSV File", nil)
		return

	}
	fileVariables = append(fileVariables, variables[0])
	writeF, _ := json.MarshalIndent(fileVariables, "", " ")
	_ = ioutil.WriteFile(filename, writeF, 0644)

	csvFile, _ := os.Create(csvFilename)
	fileOpened, _ := os.Open(filename)
	dataFrame := dataframe.ReadJSON(fileOpened)
	csvWriter := csv.NewWriter(csvFile)
	csvWriter.WriteAll(dataFrame.Records())

	dto.SendResponse(w, r, http.StatusOK, "Successfully Created CSV File From JSON", nil)
}

//CleanData function to clean data and convert data to specific format
func CleanData(w http.ResponseWriter, r *http.Request) {
	filename := "./datasets/testing/dummyDataset.csv"
	fileSetStat := "./datasets/testing/dataStats.json"

	f, err := os.Open(filename)
	newFile := f
	defer newFile.Close()
	if err != nil {
		fmt.Println("Error in opening file", f.Name(), err.Error())
		dto.SendResponse(w, r, http.StatusInternalServerError, "Error in opening file "+filename+" Error is \n"+err.Error(), nil)
		return
	}
	reader := csv.NewReader(f)
	rawData, err := reader.ReadAll()
	if err != nil {
		dto.SendResponse(w, r, http.StatusInternalServerError, "Error in opening file 2"+err.Error(), nil)
		return
	}
	var dataWithoutID [][]string
	var machineID int
	// fmt.Println(dataWithoutID, machineID)
	for index, row := range rawData {
		var rowRawData []string
		if index == 0 {
			for rowIndex, rowValues := range row {
				if rowValues != "machineId" {
					rowRawData = append(rowRawData, rowValues)
				} else {
					machineID = rowIndex
				}
			}
			dataWithoutID = append(dataWithoutID, rowRawData)
			continue
		}
		for rowIndex, rowValues := range row {
			if rowIndex != machineID {
				rowRawData = append(rowRawData, rowValues)
			}
		}
		dataWithoutID = append(dataWithoutID, rowRawData)
	}
	rawData = dataWithoutID
	// fmt.Println("Data Without Machine Id", dataWithoutID)

	data := dataframe.LoadRecords(dataWithoutID)
	// fmt.Println("Printing Data Before Cleaning", data.String())
	// fmt.Println(data)

	// check := []int{1, 1, 1, 1, 1, 1}

	labelExtraction := []string{}

	labels := data.Col("yLabels").Records()

	var index int
	dataRecords := data.Records()
	for i, record := range dataRecords[0] {
		// fmt.Println("Index ", i, "Value", record)
		if record == "yLabels" {
			index = i
		}
	}
	// fmt.Println(fileSetStat)
	stats := fileStats{
		LabelIndex: index,
	}
	statsBytes, err := json.MarshalIndent(stats, "", "")
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Contact Support", map[string]interface{}{"Error": err.Error()})
		return
	}
	err = ioutil.WriteFile(fileSetStat, statsBytes, 0644)
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Contact Support", map[string]interface{}{"Error": err.Error()})
		return
	}
	// fmt.Println("Index Is", index)

	// fmt.Println("Labels From Dataset", labels)
	//Loop for extracting Labels
	for _, label := range labels {
		result := helper.StringContains(labelExtraction, label)
		if result == false {
			labelExtraction = append(labelExtraction, label)
		}
		// fmt.Println("String Contains", result)
	}
	// fmt.Println("Extracted Labels", labelExtraction)
	// fmt.Println(rawData)
	sort.Strings(labelExtraction)

	var newRecords [][]string
	// fmt.Println(labelExtraction)
	length := len(rawData[0])
	for i, record := range rawData {

		// fmt.Println("Length", length)
		if i == 0 {
			for _, label := range labelExtraction {
				record = append(record, label)
			}
			newRecords = append(newRecords, record)

			continue
		}
		for _, label := range labelExtraction {
			_ = label
			// fmt.Println(label)
			record = append(record, "0.0")

		}

		for index, label := range labelExtraction {
			if record[length-1] == label {
				record[index+length] = "1.0"
			}

			// fmt.Println("Record", record[len(record)-(length)])
		}
		newRecords = append(newRecords, record)

	}

	// fmt.Println(newRecords)
	var finalRecords [][]string
	for _, record := range newRecords {
		// fmt.Println(record[:length-1])
		// fmt.Println(record[length:])

		modified := append(record[:length-1], record[length:]...)
		finalRecords = append(finalRecords, modified)
		// finalRecords = append(finalRecords, record[length:])
	}

	// fmt.Println("After Cleaning Dataset", finalRecords)
	writer, _ := os.Create("./datasets/dataset.csv")
	wr := csv.NewWriter(writer)
	wr.WriteAll(finalRecords)
	wr.Flush()
	dto.SendResponse(w, r, http.StatusOK, "Success", map[string]interface{}{"raw csv data": finalRecords, "Index ": index, "Length": len(newRecords)})
}

//NormalizeData function for data normalization from 0-1
func NormalizeData(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("./datasets/dataset.csv")

	if err != nil {

		dto.SendResponse(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}

	csvReader := csv.NewReader(f)
	rawCSVdata, err := csvReader.ReadAll()
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}
	dataFrame := dataframe.LoadRecords(rawCSVdata)
	normalize := [][]string{}
	// fmt.Println("Data Name ", dataFrame.Col(string(rawCSVdata[0][0])))
	for i, record := range rawCSVdata {
		if i == 0 {
			normalize = append(normalize, record)
			continue
		}
		noVal := []string{}
		for x, values := range record {
			// fmt.Println("Type ", reflect.TypeOf(values))

			val, err := strconv.ParseFloat(values, 64)
			if err != nil {
				dto.SendResponse(w, r, http.StatusBadRequest, err.Error(), nil)
				return
			}
			// fmt.Println("x", x, "Values", values)
			// fmt.Println("val", val)
			// fmt.Println("Min", dataFrame.Col(rawCSVdata[0][x]).Min())
			// fmt.Println("Max", dataFrame.Col(rawCSVdata[0][x]).Max())
			nVal := (val - dataFrame.Col(rawCSVdata[0][x]).Min()) / (dataFrame.Col(rawCSVdata[0][x]).Max() - dataFrame.Col(rawCSVdata[0][x]).Min())
			// fmt.Println("Number Value", nVal)
			sVal := strconv.FormatFloat(nVal, 'f', 6, 64)
			noVal = append(noVal, sVal)
			// fmt.Println(nVal)
		}
		normalize = append(normalize, noVal)
	}

	headers := rawCSVdata[0]
	wr, err := os.Create("./datasets/normalized.csv")
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, err.Error(), nil)
		return
	}
	csvWriter := csv.NewWriter(wr)
	csvWriter.WriteAll(normalize)
	dto.SendResponse(w, r, http.StatusOK, "Everthings Fine", map[string]interface{}{"data": rawCSVdata, "summary": dataFrame.Describe().Records(),
		"Headers": headers, "Normalixe": normalize, "Original": rawCSVdata})
}

//SplitAndShuffle to split and shuffle data
func SplitAndShuffle(w http.ResponseWriter, r *http.Request) {
	percentage := r.FormValue("trainPercentage")
	percent, err := strconv.ParseFloat(percentage, 64)
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Contact Support", map[string]interface{}{"Error Message": err.Error()})
		return
	}
	filePath := "./datasets/normalized.csv"
	trainFilePath := "./datasets/train_normalized.csv"
	testFilePath := "./datasets/test_normalized.csv"

	percent = percent / 100
	fmt.Println("Percent", percent)
	normalizedFile, err := os.Open(filePath)
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Contact Support", map[string]interface{}{"Error Message": err.Error()})
		return
	}

	normalizedFileReader := csv.NewReader(normalizedFile)
	normalizedFileData, err := normalizedFileReader.ReadAll()
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Contact Support", map[string]interface{}{"Error Message": err.Error()})
		return
	}

	headers := normalizedFileData[0]
	fmt.Println("Headers Of Normalized Data", headers)

	splittingLength := int(math.Floor(float64(len(normalizedFileData)) * percent))
	fmt.Println("Splitting Length", splittingLength)

	var trainRawData [][]string
	var testRawData [][]string
	trainRawData = append(trainRawData, headers)
	testRawData = append(testRawData, headers)
	for i := 1; i < splittingLength; i++ {
		trainRawData = append(trainRawData, normalizedFileData[i])
	}
	for i := splittingLength; i < len(normalizedFileData); i++ {
		testRawData = append(testRawData, normalizedFileData[i])
	}
	// fmt.Println("Train Raw Data", trainRawData)
	// fmt.Println("Test Raw Data", testRawData)

	err = helper.WriteCSVFile(trainFilePath, trainRawData)
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Contact Support", map[string]interface{}{"Error Message": err.Error()})
		return
	}

	err = helper.WriteCSVFile(testFilePath, testRawData)
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Contact Support", map[string]interface{}{"Error Message": err.Error()})
		return
	}
	// trainPercent := percent / 100
	// testPercent := ((100 - percent) / 100) / 2
	// var validatePercent int
	// var even bool
	// fmt.Println("Train Percent", trainPercent)
	// fmt.Println("Test Percent", testPercent)

	// f, _ := os.Open("./datasets/normalized.csv")
	// reader := csv.NewReader(f)
	// rawCSV, _ := reader.ReadAll()
	// headers := rawCSV[0]
	// // length := len(rawCSV) / 2
	// fmt.Println("Train Data Length With Mat", math.Ceil(float64(len(rawCSV))*testPercent))
	// fmt.Println("Train Data Length", float64(len(rawCSV))*testPercent)
	// // fmt.Println("Length", length)
	// fmt.Println("Headers", headers)
	// fmt.Println(len(rawCSV))
	// if int(float64(len(rawCSV))*testPercent)%2 == 0 {
	// 	fmt.Println("Yes its even")
	// 	even = true
	// } else {
	// 	fmt.Println("No ")
	// 	even = false
	// }

	// csvPath := "./datasets/train_normalized.csv"
	// fmt.Println(csvPath)
	// if even == true {
	// 	// for i := 0; i < length; i++ {
	// 	// 	fmt.Println(csvPath)
	// 	// }
	// }

	// fmt.Println(testPercent, even, validatePercent, trainPercent)
	dto.SendResponse(w, r, http.StatusOK, "Successfully Splitted Data", map[string]interface{}{})
}

//Train funciton for training Neural Network
func Train(w http.ResponseWriter, r *http.Request) {
	rateString := r.URL.Query().Get("rate")
	epochsString := r.URL.Query().Get("epochs")
	hiddenString := r.URL.Query().Get("hidden")

	rate, _ := strconv.ParseFloat(rateString, 64)
	epochs, _ := strconv.ParseFloat(epochsString, 64)
	hidden, _ := strconv.ParseFloat(hiddenString, 64)

	// if rate == 0 || epochs == 0 || hidden == 0 {
	// 	dto.SendResponse(w, r, http.StatusBadGateway, "Rate,Epochs & hidden Cannot be Empty", nil)
	// 	return
	// }
	if rate == 0 {
		rate = 0.2
	}
	if epochs == 0 {
		epochs = 10000
	}
	if hidden == 0 {
		hidden = 5
	}
	// fmt.Println("Rate is :\t", rate)
	// fmt.Println("Type of Rate is :\t", reflect.TypeOf(rate))

	// fmt.Println("Epochs is :\t", epochs)
	// fmt.Println("hidden", hidden)

	// f, err := os.Open("./datasets/normalized.csv")
	labelIndexStruct := fileStats{}
	labelIndexFile, err := os.Open("./datasets/testing/dataStats.json")
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Contact Support", map[string]interface{}{"Error Message": err.Error()})
		return
	}
	labelIndexBytes, err := ioutil.ReadAll(labelIndexFile)
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Contact Support", map[string]interface{}{"Error Message": err.Error()})
		return
	}
	err = json.Unmarshal(labelIndexBytes, &labelIndexStruct)
	fmt.Println("Label Structure", labelIndexStruct)
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Contact Support", map[string]interface{}{"Error Message": err.Error()})
		return
	}

	// f, err := os.Open("./datasets/iris_train.csv")
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Contact Support", map[string]interface{}{"Error Message": err.Error()})
		return
	}

	//normalized train file for training
	normalizedFile, err := os.Open("./datasets/train_normalized.csv")
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Contact Support", map[string]interface{}{"Error Message": err.Error()})
		return
	}
	fileReader := csv.NewReader(normalizedFile)
	rawNormalizedData, err := fileReader.ReadAll()
	if err != nil {
		dto.SendResponse(w, r, http.StatusBadRequest, "Contact Support", map[string]interface{}{"Error Message": err.Error()})
		return
	}

	_ = rawNormalizedData
	// reader := csv.NewReader(f)
	// rawCSVdata, err := reader.ReadAll()
	featuresLength := len(rawNormalizedData[0]) - labelIndexStruct.LabelIndex
	// fmt.Println("Features Length", featuresLength, "Labels Length", labelIndexStruct.LabelIndex)
	inputsData := make([]float64, featuresLength*len(rawNormalizedData))
	labelsData := make([]float64, labelIndexStruct.LabelIndex*len(rawNormalizedData))
	fmt.Println("Features Data Length", len(inputsData), "Labels Data Length", len(labelsData))

	// inputsData := make([]float64, 4*len(rawCSVdata))
	// labelsData := make([]float64, 3*len(rawCSVdata))

	var inputIdx int
	var labelIdx int

	for idx, record := range rawNormalizedData {
		if idx == 0 {
			continue
		}

		for i, val := range record {
			parsedVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				// fmt.Println("Error in Parsing Float Value", err.Error())
				return
			}

			//if i >= labelIndex. This Condition Is Should be Implemented
			if i >= labelIndexStruct.LabelIndex {
				labelsData[labelIdx] = parsedVal
				labelIdx++
			} else {
				inputsData[inputIdx] = parsedVal
				inputIdx++
			}
		}
	}

	// inputs := mat.NewDense(len(rawCSVdata), 4, inputsData)
	// labels := mat.NewDense(len(rawCSVdata), 3, labelsData)

	inputs := mat.NewDense(len(rawNormalizedData), featuresLength, inputsData)
	labels := mat.NewDense(len(rawNormalizedData), labelIndexStruct.LabelIndex, labelsData)

	// config := models.NeuralNetConfig{
	// 	InputNeurons:  4,
	// 	OutputNeurons: 3,
	// 	HiddenNeurons: 5,
	// 	NumEpochs:     int(epochs),
	// 	LearningRate:  rate,
	// }
	config := models.NeuralNetConfig{
		InputNeurons:  featuresLength,
		OutputNeurons: labelIndexStruct.LabelIndex,
		HiddenNeurons: 5,
		NumEpochs:     int(epochs),
		LearningRate:  rate,
	}

	network := utils.NewNetwork(config)
	// fmt.Println("Printing Network Before Training", network)

	trainOutput, err := utils.Train(inputs, labels, network)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Printing Network After Training", network)

	// fmt.Println(trainOutput)
	// fmt.Println(network.BHidden.RawMatrix())
	// fmt.Print("Printing B hidden \t")
	// fmt.Println(network.BHidden.Dims())
	rowsBHidden, colsBHidden := network.BHidden.Dims()

	// fmt.Print("Printing W hidden \t")
	// fmt.Println(network.WHidden.Dims())
	rowsWHidden, colsWHidden := network.WHidden.Dims()

	// fmt.Print("Printing W out  \t")
	// fmt.Println(network.WOut.Dims())
	rowsWOut, colsWOut := network.WOut.Dims()

	// fmt.Print("Printing B out \t")
	// fmt.Println(network.BOut.Dims())
	rowsBOut, colsBOut := network.BOut.Dims()

	configuration := models.ModelConfig{}

	for i := 0; i < rowsBHidden; i++ {
		configuration.BHidden = append(configuration.BHidden, network.BHidden.RawRowView(i))
	}
	for i := 0; i < rowsWHidden; i++ {
		configuration.WHidden = append(configuration.WHidden, network.WHidden.RawRowView(i))
	}
	for i := 0; i < rowsWOut; i++ {
		configuration.WOut = append(configuration.WOut, network.WOut.RawRowView(i))
	}
	for i := 0; i < rowsBOut; i++ {
		configuration.BOut = append(configuration.BOut, network.BOut.RawRowView(i))
	}

	configuration.BHiddenDims = append(configuration.BHiddenDims, rowsBHidden)
	configuration.BHiddenDims = append(configuration.BHiddenDims, colsBHidden)
	configuration.WHiddenDims = append(configuration.WHiddenDims, rowsWHidden)
	configuration.WHiddenDims = append(configuration.WHiddenDims, colsWHidden)
	configuration.WOutDims = append(configuration.WOutDims, rowsWOut)
	configuration.WOutDims = append(configuration.WOutDims, colsWOut)
	configuration.BOutDims = append(configuration.BOutDims, rowsBOut)
	configuration.BOutDims = append(configuration.BOutDims, colsBOut)
	configuration.InputNeurons = network.Config.InputNeurons
	configuration.HiddenNeurons = network.Config.HiddenNeurons
	configuration.LearningRate = network.Config.LearningRate
	configuration.NumEpochs = network.Config.NumEpochs
	configuration.OutputNeurons = network.Config.OutputNeurons

	// fmt.Println("Configuration Values ", configuration)

	var truePosNeg int
	numPreds, _ := trainOutput.Dims()
	for i := 0; i < numPreds; i++ {
		// fmt.Println("Prediction Index ", i, "\t", trainOutput.RowView(i))
		// Get the label.
		labelRow := mat.Row(nil, i, labels)
		var species int
		for idx, label := range labelRow {
			if label == 1.0 {
				species = idx
				break
			}
		}

		// Accumulate the true positive/negative count.
		if trainOutput.At(i, species) == floats.Max(mat.Row(nil, i, trainOutput)) {
			truePosNeg++
		}
	}

	// Calculate the accuracy (subset accuracy).
	accuracy := float64(truePosNeg) / float64(numPreds)

	// Output the Accuracy value to standard out.
	fmt.Printf("\nAccuracy of Testing = %0.2f %%\n", accuracy*100)
	stats := models.ModelStats{}
	stats.TrainAccuracy = accuracy * 100

	file, _ := json.MarshalIndent(configuration, "", " ")

	_ = ioutil.WriteFile("./models/test.json", file, 0644)

	file, _ = json.MarshalIndent(stats, "", " ")

	_ = ioutil.WriteFile("./models/stats.json", file, 0644)

	dto.SendResponse(w, r, http.StatusOK, "Training Result ", map[string]interface{}{"output": configuration, "Accuracy": accuracy * 100})

}

//Predict Function for prediction
func Predict(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	stringf1 := r.FormValue("f1")
	stringf2 := r.FormValue("f2")
	stringf3 := r.FormValue("f3")
	stringf4 := r.FormValue("f4")

	f1, _ := strconv.ParseFloat(stringf1, 64)
	f2, _ := strconv.ParseFloat(stringf2, 64)
	f3, _ := strconv.ParseFloat(stringf3, 64)
	f4, _ := strconv.ParseFloat(stringf4, 64)

	input := []float64{f1, f2, f3, f4}
	configuration := models.ModelConfig{}
	file, _ := ioutil.ReadFile("./models/test.json")
	_ = json.Unmarshal([]byte(file), &configuration)

	// fmt.Println("Printing Configurations", configuration)
	config := models.NeuralNetConfig{
		InputNeurons:  configuration.InputNeurons,
		OutputNeurons: configuration.OutputNeurons,
		HiddenNeurons: configuration.HiddenNeurons,
		NumEpochs:     configuration.NumEpochs,
		LearningRate:  configuration.LearningRate,
	}

	network := utils.NewNetwork(config)

	network.BHidden = mat.NewDense(configuration.BHiddenDims[0], configuration.BHiddenDims[1], nil)
	for i := 0; i < configuration.BHiddenDims[0]; i++ {
		fmt.Println(i)
		network.BHidden.SetRow(i, configuration.BHidden[i])
	}

	network.WHidden = mat.NewDense(configuration.WHiddenDims[0], configuration.WHiddenDims[1], nil)
	for i := 0; i < configuration.WHiddenDims[0]; i++ {
		network.WHidden.SetRow(i, configuration.WHidden[i])
	}

	network.BOut = mat.NewDense(configuration.BOutDims[0], configuration.BOutDims[1], nil)
	for i := 0; i < configuration.BOutDims[0]; i++ {
		network.BOut.SetRow(i, configuration.BOut[i])
	}

	network.WOut = mat.NewDense(configuration.WOutDims[0], configuration.WOutDims[1], nil)
	for i := 0; i < configuration.WOutDims[0]; i++ {
		network.WOut.SetRow(i, configuration.WOut[i])
	}

	// fmt.Println("Printing Configurations", configuration)
	features := mat.NewDense(1, 4, input)
	predictions, err := utils.Predict(features, network)

	if err != nil {
		fmt.Println("Error in something", err.Error())
		dto.SendResponse(w, r, http.StatusInternalServerError, "Error in Prediction", map[string]interface{}{"Error": err.Error()})
	}

	fmt.Println(predictions)
	fmt.Println(floats.MaxIdx(mat.Row(nil, 0, predictions)))

	dto.SendResponse(w, r, http.StatusOK, "Success ", map[string]interface{}{"Class Name": models.ClassNames[floats.MaxIdx(mat.Row(nil, 0, predictions))], "Prediction": predictions.RawMatrix(), "Max Index": floats.MaxIdx(mat.Row(nil, 0, predictions)) + 1})

}

type proResponse struct {
	Code int           `json:"code"`
	Msg  string        `json:"message"`
	Key  rsa.PublicKey `json:"key"`
}

//GettingPublicKey server
func GettingPublicKey(w http.ResponseWriter, r *http.Request) {

	publicKey := constants.PublicKey
	fmt.Println("showing key", publicKey)
	objResp := proResponse{
		Msg:  "Success",
		Code: 200,
		Key:  publicKey,
	}
	userJSON, _ := json.Marshal(objResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
	return
}
