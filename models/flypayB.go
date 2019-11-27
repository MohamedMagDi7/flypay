package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type FlyPayBList struct {
	FlyPayB []FlyPayB `json:"transactions"`
}

//FlyPayB model
type FlyPayB struct {
	Amount         int    `json:"value"`
	Currency       string `json:"transactionCurrency"`
	StatusCode     int    `json:"statusCode"`
	OrderReference string `json:"orderInfo"`
	TransactionID  string `json:"paymentId"`
}

//Get flyPays for FlyPayB model
func (flypay *FlyPayB) GetFlyPays(filter Filter, ch chan *Response) {
	// we initialize our filtered FlyPayAList
	var flyPayFilteredList FlyPayList
	// we initialize our full FlyPayAList
	var flyPayFullList FlyPayBList

	//get file data
	byteValue := flypay.GetFileData(ch)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'flypayBData' which we defined above
	err := json.Unmarshal(byteValue, &flyPayFullList)
	if err != nil {
		ch <- &Response{List: &FlyPayList{}, Err: err}
	}

	//calculate list length
	length := len(flyPayFullList.FlyPayB)

	//calculate the number of goroutines depending on number of items per routine
	//we can set the number as we see is suitable for the data size
	routinesCount := length / ItemsPerRoutine

	//we make channel for internal filter processing
	processChannel := make(chan []FlyPay)

	//if the data size fits in one goroutine
	//we send all list to be proccessed in one go routine
	if routinesCount == 0 {
		go flypay.FilterDataChunk(flyPayFullList.FlyPayB, filter, processChannel)
		flyPayFilteredList.FlyPay = append(flyPayFilteredList.FlyPay, <-processChannel...)
	} else {
		//else if there is more than one goroutine we split up the items into chunks to be processed
		for count, i := 1, 0; count <= routinesCount; i++ {
			//last iteration we take the rist of data
			//we don't constrain to a fixed size
			if i == routinesCount-1 {
				go flypay.FilterDataChunk(flyPayFullList.FlyPayB[i*ItemsPerRoutine:length-1], filter, processChannel)
				flyPayFilteredList.FlyPay = append(flyPayFilteredList.FlyPay, <-processChannel...)
			} else {
				go flypay.FilterDataChunk(flyPayFullList.FlyPayB[i*ItemsPerRoutine:(count*ItemsPerRoutine)-1], filter, processChannel)
				flyPayFilteredList.FlyPay = append(flyPayFilteredList.FlyPay, <-processChannel...)
			}
			count++
		}
	}

	//we return filtered list to the controller
	ch <- &Response{List: &flyPayFilteredList, Err: err}
}

//each chunk of data gets filtered in a goroutine
func (flypay *FlyPayB) FilterDataChunk(flypays []FlyPayB, filter Filter, pCh chan []FlyPay) {
	var flypayList []FlyPay
	for i := len(flypays) - 1; i >= 0; i-- {
		if flypays[i].IsAMatch(&filter) {
			//if item matches filter we Convert from FlyPayB to FlyPay and append to list
			flypayList = append(flypayList, flypays[i].ConvertFlypayToBase())
		}
	}
	pCh <- flypayList
}

//read json data from data source file
func (flypay *FlyPayB) GetFileData(ch chan *Response) []byte {
	// Open our JSON File
	flypayBFile, err := os.Open(ProviderBFilePath)

	// if we os.Open returns an error then handle it
	if err != nil {
		//check if error is file does not exist
		//check if it exist in another path
		//this step is done because of the unit testing paths differ
		if os.IsNotExist(err) {
			flypayBFile, err = os.Open(ProviderBAlternativeFilePath)
			if err != nil {
				ch <- &Response{List: &FlyPayList{}, Err: err}
			}
		} else {
			ch <- &Response{List: &FlyPayList{}, Err: err}
		}
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer flypayBFile.Close()

	// read our opened JSON File as a byte array.
	byteValue, err := ioutil.ReadAll(flypayBFile)
	if err != nil {
		ch <- &Response{List: &FlyPayList{}, Err: err}
	}

	return byteValue

}

//to Filter Our Data
func (flypay *FlyPayB) IsAMatch(filter *Filter) bool {
	if filter.AmountMax != 0 && flypay.Amount > filter.AmountMax {
		return false
	}
	if filter.AmountMin != 0 && flypay.Amount < filter.AmountMin {
		return false
	}
	if filter.Currency != "" && flypay.Currency != filter.Currency {
		return false
	}
	if filter.BStatusCode != 0 && flypay.StatusCode != filter.BStatusCode {
		return false
	}
	return true
}

//we use this function to Cast back to base struct type
//this function can be customized as to fit for new added providers
func (fromObject *FlyPayB) ConvertFlypayToBase() FlyPay {
	return FlyPay(*fromObject)
}
