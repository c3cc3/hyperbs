package hyperbs

import (
    "fmt"
    "os"
    "bytes"
    "bufio"
    "strings"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    GitGoIpfsApi "github.com/ipfs/go-ipfs-api" // GitGoIpfsApi is alias
    // "github.com/hyperledger/fabric/protos/peer"
)

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func Set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

// Get returns the value of the specified asset key
func Get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}

// ipfs calling 
func Set_addipfs(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	var sender, receiver, filename string

	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting filename")
	}

	key := args[0]
	// logger.Info( "key: " + key )
	value := args[1]
	// logger.Info( "value: " + value )

	str_Slice := strings.Split( value, "|")
	// logger.Info(str_Slice)

	for i, str_Slice := range str_Slice {
		if i==0 {
			sender = str_Slice
		} else if i == 1  {
			receiver = str_Slice
		} else if i == 2  {
			filename = str_Slice
		}
	}

	// logger.Info( key, sender, receiver, filename)
	// logger.Info("Add to ipfs: " + filename)
	fmt.Println( key, sender, receiver, filename)
	fmt.Println("Add to ipfs: " + filename)

// search with container name (ipfs0)
	mhash, err := AddIpfs( "ipfs0", "5001", filename)

	if err != nil {
		// logger.Info("AddIpfs() error")
		jsonResp := "{\"Error\":\"Failed to add to IPFS" + "\"}"
		return "", fmt.Errorf(jsonResp)
	}
    // logger.Info( "Success to add on ipfs: " +  mhash)
	fmt.Println( "Success to add on ipfs: " +  mhash )
	value = value + "|" +  mhash

	stub_err := stub.PutState(key, []byte(value))
	if stub_err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", key)
	}
    // logger.Info( "Success to set on ledger. key:" +  key + ", value: " + value)
	fmt.Println( "Success to set on ledger. key:" +  key + ", value: " + value)

	jsonResp := "{\"Name\":\"" + key  + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	return key, nil
}

func Get_catipfs(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var sender, receiver, filename, mhash string

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting document number.")
	}
	key :=  args[0]

	// logger.Info( "document number(key):" + key )

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", key, err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", key)
	}

	str_Slice := strings.Split( string(value), "|")
	// logger.Info(str_Slice)

	for i, str_Slice := range str_Slice {
		if i==0 {
			sender = str_Slice
		} else if i == 1  {
			receiver = str_Slice
		} else if i == 2  {
			filename = str_Slice
		} else if i == 3 {
			mhash = str_Slice
		}
	}
	// logger.Info( "From ledger: key=",  key)
	// logger.Info( "From ledger: sender=",  sender)
	// logger.Info( "From ledger: receiver=",  receiver)
	// logger.Info( "From ledger: filename=",  filename)
	// logger.Info( "From ledger: mhash=",  mhash)

	fmt.Println( "From ledger: key=",  key)
	fmt.Println( "From ledger: sender=",  sender)
	fmt.Println( "From ledger: receiver=",  receiver)
	fmt.Println( "From ledger: filename=",  filename)
	fmt.Println( "From ledger: mhash=",  mhash)

	contents, err := CatIpfs( "ipfs0", "5001", mhash)
	if err != nil {
		// logger.Info("CatIpfs() error")
        jsonResp := "{\"Error\":\"Failed to add to IPFS" + "\"}"
        return "", fmt.Errorf(jsonResp)
	}

	jsonResp := "{\"contents\":\"" + contents  + "\"}"
    fmt.Printf("Query Response:%s\n", jsonResp)

	return mhash, nil
}

func CatIpfs(Ip string, Port string, mhash string) (string, error) {

	
	UrlPort := Ip + ":" + Port
	shell := GitGoIpfsApi.NewShell(UrlPort)

	reader, err := shell.Cat(mhash)
	if err != nil {
		// logger.Error("shell.Cat() error.")
		return mhash, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	str_buf := buf.String()
	// logger.Info("buf: " + str_buf)

	return str_buf, err
}

func AddIpfs(Ip string, Port string, filename string) (string, error) {

	UrlPort := Ip + ":" + Port
	shell := GitGoIpfsApi.NewShell(UrlPort)
	bytedata, err := RetrieveROM( filename )
	if err != nil  {
    	// logger.Info("file open error:" + filename);
		return filename, err
	}

	s := string(bytedata[:])
	bufferExample := bytes.NewBufferString(s)

	mhash, err := shell.Add(bufferExample)
	if err != nil {
		// logger.Error("shell.Add() error.")
		return filename, err
	}

/*/
	file_mhash = "/ipfs" +  mhash
	buf, err = shell.Cat( file_mhash)
	if err != nil {
		// logger.Error("shell.Cat() error.")
		return filename, err
	}
*/
	return mhash, err
}

func AddNoPinIpfs(Ip string, Port string, filename string) (string, error) {

	UrlPort := Ip + ":" + Port
	shell := GitGoIpfsApi.NewShell(UrlPort)
	bytedata, err := RetrieveROM( filename )
	if err != nil {
		return "RetriveROM() error", err
	}

	s := string(bytedata[:])
	bufferExample := bytes.NewBufferString(s)

	mhash, err := shell.AddNoPin(bufferExample)
	if err != nil {
		return "shell.AddNoPin() error", err
	}

	return mhash, err
}

// Loading data of a file to byte memory
func RetrieveROM(filename string) ([]byte, error) {
    file, err := os.Open(filename)

    if err != nil {
        return nil, err
    }
    defer file.Close()

    stats, statsErr := file.Stat()
    if statsErr != nil {
        return nil, statsErr
    }

    var size int64 = stats.Size()
    bytes := make([]byte, size)

	fmt.Println("file size : ", size);

    bufr := bufio.NewReader(file)
    _,err = bufr.Read(bytes)

    return bytes, err
}

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
    if e != nil {
        panic(e)
    }
}
