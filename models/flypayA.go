package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//FlyPay Base model
type FlyPay struct {
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	StatusCode     int    `json:"statusCode"`
	OrderReference string `json:"orderReference"`
	TransactionID  string `json:"transactionId"`
}

//Get flyPays for FlyPayA model
func (flypay *FlyPay) GetFlyPays(filter Filter, ch chan *Response) {
	// we initialize our filtered FlyPayAList
	var flyPayFilteredList FlyPayList
	// we initialize our full FlyPayAList
	var flyPayFullList FlyPayList

	//get file data
	byteValue := flypay.GetFileData(ch)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'flyPayFullList' which we defined above
	err := json.Unmarshal(byteValue, &flyPayFullList)
	if err != nil {
		ch <- &Response{List: &FlyPayList{}, Err: err}
	}

	//calculate list length
	length := len(flyPayFullList.FlyPay)

	//calculate the number of goroutines depending on number of items per routine
	//we can set the number as we see is suitable for the data size
	routinesCount := length / ItemsPerRoutine

	//we make channel for internal filter processing
	processChannel := make(chan []FlyPay)

	//if the data size fits in one goroutine
	//we send all list to be proccessed in one go routine
	if routinesCount == 0 {
		go flypay.FilterDataChunk(flyPayFullList.FlyPay, filter, processChannel)
		flyPayFilteredList.FlyPay = append(flyPayFilteredList.FlyPay, <-processChannel...)
	} else {
		//else if there is more than one goroutine we split up the items into chunks to be processed
		for count, i := 1, 0; count <= routinesCount; i++ {
			//last iteration we take the rist of data
			//we don't constrain to a fixed size
			if i == routinesCount-1 {
				go flypay.FilterDataChunk(flyPayFullList.FlyPay[i*ItemsPerRoutine:length-1], filter, processChannel)
				flyPayFilteredList.FlyPay = append(flyPayFilteredList.FlyPay, <-processChannel...)
			} else {
				go flypay.FilterDataChunk(flyPayFullList.FlyPay[i*ItemsPerRoutine:(count*ItemsPerRoutine)-1], filter, processChannel)
				flyPayFilteredList.FlyPay = append(flyPayFilteredList.FlyPay, <-processChannel...)
			}
			count++
		}
	}

	//we return the filtered list to the controller
	ch <- &Response{List: &flyPayFilteredList, Err: err}
}

//each chunk of data gets filtered in a goroutine
func (flypay *FlyPay) FilterDataChunk(flypays []FlyPay, filter Filter, pCh chan []FlyPay) {
	for i := len(flypays) - 1; i >= 0; i-- {
		if !flypays[i].IsAMatch(&filter) {
			flypays = append(flypays[:i], flypays[i+1:]...)
			continue
		}
	}
	pCh <- flypays
}

//read json data from data source file
func (flypay *FlyPay) GetFileData(ch chan *Response) []byte {
	//open our JSON file
	flypayAFile, err := os.Open(ProviderAFilePath)

	// if we os.Open returns an error then handle it
	if err != nil {
		//check if error is file does nor exist
		//check if it exist in another path
		//this step is done because of the unit testing
		if os.IsNotExist(err) {
			flypayAFile, err = os.Open(ProviderAAlternativeFilePath)
			if err != nil {
				ch <- &Response{List: &FlyPayList{}, Err: err}
			}
		} else {
			ch <- &Response{List: &FlyPayList{}, Err: err}
		}
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer flypayAFile.Close()

	// read our opened JSON File as a byte array.
	byteValue, err := ioutil.ReadAll(flypayAFile)
	if err != nil {
		ch <- &Response{List: &FlyPayList{}, Err: err}
	}

	return byteValue
}

//to determine if object matches the filter
func (flypay *FlyPay) IsAMatch(filter *Filter) bool {
	if filter.AmountMax != 0 && flypay.Amount > filter.AmountMax {
		return false
	}
	if filter.AmountMin != 0 && flypay.Amount < filter.AmountMin {
		return false
	}
	if filter.Currency != "" && flypay.Currency != filter.Currency {
		return false
	}
	if filter.AStatusCode != 0 && flypay.StatusCode != filter.AStatusCode {
		return false
	}
	return true
}

//we use this function to Cast back to base struct type
//this function can be customized as to fit for new added providers
func (fromObject *FlyPay) ConvertFlypayToBase() FlyPay {
	return *fromObject
}

// Another way of implementing GetFlyPays
// It will save some memory as it reads objects from file
// One by one but it would have to go synchronously through the file
//it can become handy if the size of the file is bigger than the available memory
/*
//Get flyPays for FlyPayA model
func GetFlyPaysA(filter Filter, ch chan *Response) {
	//open our JSON file
	flypayAFile, err := os.Open("./models/data/flypayA.json")

	// if we os.Open returns an error then handle it
	if err != nil  {
		//check if error is file does nor exist
		//check if it exist in another path
		//this step is done because of the unit testing
		if os.IsNotExist(err){
			flypayAFile, err = os.Open("./data/flypayA.json")
			if err != nil  {
				ch <- &Response{List: &FlyPayList{}, Err: err}
			}
		}else{
			ch <- &Response{List: &FlyPayList{}, Err: err}
		}
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer flypayAFile.Close()

	// we initialize our empty FlyPayAList
	var flyPayFilteredList FlyPayList

	//we initialize our json decoder
	dec := json.NewDecoder(flypayAFile)

	//we discard the first 3 tokens
	// to get to the start of the flypays list
	for i:= 0 ; i< 3 ; i++ {
		_ , err = dec.Token()
		if err != nil {
			ch <- &Response{List: &FlyPayList{}, Err: err}
		}
	}

	// we iterate over json objects until we reach the end of data
	for dec.More() {
		//initialize flypay struct to decode into
		var flypay FlyPay
		//decode the next json object into flypay
		err = dec.Decode(&flypay)
		//check for decoding errors
		if err != nil {
			ch <- &Response{List: &FlyPayList{}, Err: err}
		}
		//check if flypay match the filter
		if flypay.IsAMatch(&filter){
			//if matched we append the object to our list
			flyPayFilteredList.FlyPay = append(flyPayFilteredList.FlyPay, flypay)
		}
	}

	//we return the filtered list to the controller
	ch <- &Response{List: &flyPayFilteredList, Err: err}
}
*/
