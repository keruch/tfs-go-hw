package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
)

type (
	creationTime   string
	operationType  int8
	operationValue int
	operationID    string
	operation      struct {
		Type      interface{} `json:"type,omitempty"`
		Value     interface{} `json:"value,omitempty"`
		ID        interface{} `json:"ID,omitempty"`
		CreatedAt interface{} `json:"created_at,omitempty"`
	}
	operationData struct {
		Company   string `json:"company"`
		Operation struct {
			operation
		}
		operation
	}
	companyTransactions struct {
		Company              string        `json:"company"`
		ValidOperationsCount int           `json:"valid_operations_count"`
		Balance              int           `json:"balance"`
		InvalidOperations    []interface{} `json:"invalid_operations"`
	}
	companyTransactionsTable map[string]companyTransactions
	sortByTime               []operationData
)

func updateTransactions(ct companyTransactionsTable, o operationData) companyTransactionsTable {
	val, ok := ct[o.Company]
	if !ok {
		ct[o.Company] = companyTransactions{}
	}
	if o.Valid() {
		val = companyTransactions{
			Company:              o.Company,
			ValidOperationsCount: ct[o.Company].ValidOperationsCount + 1,
			Balance:              ct[o.Company].Balance + int(o.getOperationType())*int(o.getOperationValue()),
			InvalidOperations:    ct[o.Company].InvalidOperations,
		}
		ct[o.Company] = val
	} else {
		var tempID interface{}
		if intVal, err := strconv.Atoi(string(o.getOperationID())); err == nil {
			tempID = intVal
		} else {
			tempID = string(o.getOperationID())
		}
		val = companyTransactions{
			Company:              o.Company,
			ValidOperationsCount: ct[o.Company].ValidOperationsCount,
			Balance:              ct[o.Company].Balance,
			InvalidOperations:    append(ct[o.Company].InvalidOperations, tempID),
		}
		ct[val.Company] = val
	}
	return ct
}

func (a sortByTime) Len() int {
	return len(a)
}

func (a sortByTime) Less(i, j int) bool {
	var (
		second, first time.Time
		err           error
	)
	first, err = time.Parse(time.RFC3339, string(a[i].getCreationTime()))
	if err != nil {
		panic(err)
	}
	second, err = time.Parse(time.RFC3339, string(a[j].getCreationTime()))
	if err != nil {
		panic(err)
	}
	return first.Before(second)
}

func (a sortByTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (o operationData) Correct() bool {
	if o.Company == "" || o.getCreationTime() == "" || o.getOperationID() == "" {
		return false
	}
	return true
}

func (o operationData) Valid() bool {
	if o.getOperationType() == 0 || o.getOperationValue() == 0 {
		return false
	}
	return true
}

func (o operationData) getCreationTime() creationTime {
	var t creationTime
	tempCreationTime := o.Operation.CreatedAt
	if tempCreationTime == nil {
		tempCreationTime = o.CreatedAt
	}
	switch val := tempCreationTime.(type) {
	case string:
		_, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return ""
		}
		t = creationTime(val)
	default:
		t = ""
	}
	return t
}

func (o operationData) getOperationType() operationType {
	var t operationType
	tempOperationType := o.Operation.Type
	if tempOperationType == nil {
		tempOperationType = o.Type
	}
	switch val := tempOperationType.(type) {
	case string:
		switch {
		case val == `+` || val == `income`:
			t = 1
		case val == `-` || val == `outcome`:
			t = -1
		}
	default:
		t = 0
	}
	return t
}

func (o operationData) getOperationValue() operationValue {
	var t operationValue
	tempOperationValue := o.Operation.Value
	if tempOperationValue == nil {
		tempOperationValue = o.Value
	}
	switch val := tempOperationValue.(type) {
	case string:
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return 0
		}
		t = operationValue(intVal)
	case float64:
		if val != float64(int(val)) {
			return 0
		}
		t = operationValue(val)
	default:
		t = 0
	}
	return t
}

func (o operationData) getOperationID() operationID {
	var t operationID
	tempOperationID := o.Operation.ID
	if tempOperationID == nil {
		tempOperationID = o.ID
	}
	switch val := tempOperationID.(type) {
	case float64:
		t = operationID(strconv.Itoa(int(val)))
	case string:
		t = operationID(val)
	default:
		t = ""
	}
	return t
}

func (o operationData) String() string {
	return fmt.Sprintf("Company: %v\nType: %v\nValue: %v\nID: %v\nCreatedAt: %v",
		o.Company, o.Type, o.Value, o.ID, o.CreatedAt)
}

func openPath(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	return file, nil
}

func openPathCommandLine() (*os.File, error) {
	path := flag.String("filepath", "", "path of json with financial transactions")
	flag.Parse()
	if *path == "" {
		return nil, errors.New("command line filepath flag is empty")
	}
	return openPath(*path)
}

func openPathEnv() (*os.File, error) {
	path, ok := os.LookupEnv("FILE")
	if path == "" {
		return nil, errors.New("FILE variable specified but is empty")
	}
	if !ok {
		return nil, errors.New("FILE variable does not specified")
	}
	return openPath(path)
}

func openFile() (*os.File, error) {
	var (
		e, err error
		file   *os.File
	)
	if file, err = openPathCommandLine(); err == nil {
		return file, nil
	}
	e = fmt.Errorf("failed to read from command line: %s", err)
	if file, err = openPathEnv(); err == nil {
		return file, nil
	}
	e = fmt.Errorf("%w\nfailed to read from environment variables: %s", e, err)
	return nil, e
}

func processData(data []byte) ([]companyTransactions, error) {
	var (
		operations []operationData
		err        error
	)

	err = json.Unmarshal(data, &operations)
	if err != nil {
		return nil, err
	}

	operations = prepareData(operations)
	ct := make(companyTransactionsTable)
	for _, o := range operations {
		ct = updateTransactions(ct, o)
	}

	result := make([]companyTransactions, 0)
	for _, val := range ct {
		result = append(result, val)
	}
	return result, nil
}

func prepareData(o []operationData) []operationData {
	o = deleteIncorrectOperations(o)
	sort.Sort(sortByTime(o))
	return o
}

func deleteIncorrectOperations(o []operationData) []operationData {
	res := make([]operationData, 0)
	for _, val := range o {
		if val.Correct() {
			res = append(res, val)
		}
	}
	return res
}

func init() {
	if err := os.Setenv("FILE", "D:\\Go\\src\\tfs-go-hw\\hw2\\billing.json"); err != nil {
		return
	}
}

func main() {
	file, err := openFile()
	if err != nil {
		fmt.Printf("%sReading from stdin\n", err)
		file = os.Stdin
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error occupied: %s\n", err)
		return
	}

	result, err := processData(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	reportJSON, err := json.MarshalIndent(result, "", "   ")
	if err != nil {
		fmt.Println(err)
		return
	}

	outFile, err := os.Create("hw2/out.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outFile.Close()
	_, err = fmt.Fprint(outFile, string(reportJSON))
	if err != nil {
		fmt.Println(err)
		return
	}
}
