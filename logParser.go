package main

import (
	"bufio"   //To read input by line
	"errors"  //Custom error
	"flag"    //Imported for easy flag management.
	"fmt"     //To output to console
	"os"      //Arg finder and file open
	"sort"    //Sort array
	"strings" //String manipulation
)

//See initFlags for description or run with -h
var httpLog, verbose, sorting, merge, outputLog, showTime, multipleFiles, keyvaluepair bool
var fileName, input, outputName, output, outputLogName string

func initFlags() {
	output = ""
	const (
		defaultHttpLog = false
		//	defaultFileName = "artifactory.log"
		defaultInput         = "ERROR"
		defaultVerbose       = false
		defaultSort          = true
		defaultMerged        = true
		defaultOutputLog     = false
		defaultOutputLogName = "output.txt"
		defaultShowTime      = false
		defaultMultipleFiles = false
		defaultKVP           = false
	)
	flag.BoolVar(&httpLog, "http", defaultHttpLog, "Used to set targeted formatting style (i.e. request.log). (default false)")
	flag.StringVar(&fileName, "file", "", "(Optional flag) Set filename to read from. If flag is not included, will default to first non-flag argument, otherwise it will hit flag default. (default artifactory.log)")
	flag.StringVar(&fileName, "f", "", "Shorthand filename. (default (artifactory.log)")
	flag.StringVar(&input, "input", defaultInput, "Set input ($INPUT) text.")
	flag.StringVar(&input, "i", defaultInput, "Shorthand input.")
	//&output = input + ".txt"
	flag.StringVar(&outputName, "output", "", "Set output filename. No filename outputs to terminal. (default \"\")")
	flag.StringVar(&outputName, "o", "", "Shorthand output. (default \"\")")
	flag.BoolVar(&verbose, "verbose", defaultVerbose, "Set verbose output. Leaves lines mostly unformatted. (default false)")
	flag.BoolVar(&verbose, "v", defaultVerbose, "Shorthand verbose. (default false)")
	flag.BoolVar(&sorting, "sort", defaultSort, "Set sorted output.")
	flag.BoolVar(&sorting, "s", defaultSort, "Set sorted output.")
	flag.BoolVar(&merge, "merge", defaultMerged, "Set merged duplicate output.")
	flag.BoolVar(&merge, "m", defaultMerged, "Shorthand merge.")
	flag.BoolVar(&merge, "merge-duplicates", defaultMerged, "Longhand merge.")
	flag.BoolVar(&outputLog, "log", defaultOutputLog, "Send output log to file. Can be mixed with verbose for print and save to file. (default false)")
	flag.BoolVar(&outputLog, "l", defaultOutputLog, "Shorthand output log (default false)")
	flag.StringVar(&outputLogName, "logname", defaultOutputLogName, "Set filename for output log.")
	flag.StringVar(&outputLogName, "ln", defaultOutputLogName, "Shorthand log filename.")
	flag.BoolVar(&showTime, "time", defaultShowTime, "Show timestamps in output. (default false)")
	flag.BoolVar(&showTime, "t", defaultShowTime, "Shorthand timestamp. (default false)")
	flag.BoolVar(&multipleFiles, "many", defaultMultipleFiles, "Allow multiple filenames to be input. (default false)")
	flag.BoolVar(&keyvaluepair, "kvp", defaultKVP, "Set extra args as inputfile input pairs. Currently does not work with the -output flag. (default false)")

	flag.Parse()
	tempArr := flag.Args()
	ifVerbose(fmt.Sprintln("Read flags and initialized with the following parameters:\n================\nFilename:\t", fileName, "\nHttp Style:\t", httpLog, "\nInput:\t\t", input, "\nOutput:\t\t", outputName, "\nVerbose:\t", verbose, "\nSort:\t\t", sorting, "\nMerge:\t\t", merge, "\nRemaining Args:\t", tempArr))
}

func ifVerbose(in string) {
	if verbose {
		fmt.Print(in)
		/*		if outputName != "" {
				output += in
			}*/
	}
	if outputLog {
		output += in
	}
}

//Struct to store string and num of occurrences
type LogString struct {
	count int
	text  string
	times []string //timestamp array
	//timestamp array
	//first occurrence already sorted haha
	//last occurrence
	//duration array - httpLog only
	//shortest duration
	//longest duration
	//line array
}

func getOutput(strArr []LogString) string {
	returnString := ""
	//Iterate through the array
	for _, i := range strArr {
		if showTime {
			//Print each line.
			returnString += fmt.Sprintln(i)
		} else {
			//Do some manipulation to not have the times shown.
			returnString += fmt.Sprintf("{%d %s}\n", i.count, i.text)
		}
	}
	return returnString
}

func outputLogs(strArr []LogString) {
	ifVerbose(fmt.Sprintln("Outputting data with desired medium."))
	//Grab this now for posterity.
	writeOut := getOutput(strArr)
	//If outputName has been set
	if outputName != "" {
		//Open file
		f, err := os.Create(outputName)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		//Write to file.
		_, err = f.WriteString(writeOut)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		//Otherwise just print out the output.
		fmt.Print(writeOut)
	}
}

//Funny name, but this just outputs the logs. Metalogs?
func logOut() {
	//If we are outputting to a file.
	if outputLog {
		//Open file
		f, err := os.Create(outputLogName)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		//Write to file
		_, err = f.WriteString(output)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	//For verbose mode, it prints as the strings come in, so we don't need to print out here.
}

// func (l LogString) createLogString(){

// }

//Check to see if it exists. Also output index to reduce iterations.
func stringExists(str string, strArr []LogString) (int, bool) {
	//For every element (starting from 0)
	for n, i := range strArr {
		//If it's the same string
		if i.text == str {
			//Return location, and that it exists.
			return n, true
		}
	}
	//Otherwise, return it doesn't exist. Numeric value here doesn't matter.
	return 0, false
}

//Type used for "custom" sort.
//Also includes required functions.
type ByCount []LogString

func (a ByCount) Len() int {
	return len(a)
}
func (a ByCount) Less(i, j int) bool {
	return a[i].count < a[j].count
}
func (a ByCount) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func parseLogs(fileName string) {
	//Potentially verbose the above part? Seems irrelevant outside of script debugging.
	//Print that we're opening files.
	ifVerbose(fmt.Sprintln("Opening file:", fileName))
	file, err := os.Open(fileName)
	if err != nil {
		panic(errors.New("error: invalid filename. please ensure file exists. filename: \"" + fileName + "\"\nNote: file can be specified with -f or -file flags. See -h for help."))
	}
	//	fmt.Println("We opened: ", fileName)
	//Defer is great for files. Automatically closes when "done".
	//Works on error captures if there's a parent stack. Although we won't have a parent stack.
	defer file.Close()

	//"Scan" file
	scanner := bufio.NewScanner(file)
	//Split it by lines
	scanner.Split(bufio.ScanLines)

	//Where we'll store the strings and occurrences.
	textArray := []LogString{}

	//Bring these here to remember what flags to include in ops.
	// var httpLog, verbose, sorting, merge bool
	// var fileName, input, outputName, output string

	for scanner.Scan() {
		//Get the next line.
		text := scanner.Text()
		ifVerbose(fmt.Sprintln("Scanning text:", text))
		//Find the index where input begins.
		found := strings.Contains(text, input)
		ifVerbose(fmt.Sprintln("Index of Input found:", found))
		//If input exists in the text...
		if found == true {
			var textTime string
			var remainingText string
			if httpLog {

				//Create a dynamic stripper. httpLog would then only specify timestamps
				textTime = text[:14]
				remainingText = text[15:]
				tempStr := strings.SplitN(remainingText, "|", 2)
				remainingText = tempStr[1]
				textTime += " " + tempStr[0]
				ifVerbose(fmt.Sprintln("Text:", remainingText, "\tTimestamp:", textTime))
				// fmt.Println(textTime, remainingText)
			} else {
				textTime = text[:23]
				remainingText = text[24:]
				ifVerbose(fmt.Sprintln("Text:", remainingText, "\tTimestamp:", textTime))
			}
			//If we don't want to merge, don't merge.
			if !merge {
				ifVerbose(fmt.Sprintln("Adding values as new element to slice."))
				textArray = append(textArray, LogString{1, remainingText, []string{textTime}})
			} else {
				//Find if it exists in our array and get location.
				oldIndex, exists := stringExists(remainingText, textArray)
				//If it exists...
				if exists {
					ifVerbose(fmt.Sprintln("Input found in slice. Updating slice."))
					//Add to the existing count.
					textArray[oldIndex].count++
					textArray[oldIndex].times = append(textArray[oldIndex].times, textTime)
					ifVerbose(fmt.Sprintln("Data in slice updated. Current count of this error:", textArray[oldIndex].count))
				} else { //Otherwise...
					ifVerbose(fmt.Sprintln("Input not found in slice. Adding to slice."))
					//Add new value and start count at 1. It does occur once right now.
					textArray = append(textArray, LogString{1, remainingText, []string{textTime}})
				}
			}
		} //If input doesn't exist, we don't care. Move on to next line.
	}
	if sorting {
		ifVerbose(fmt.Sprintln("Sorting slice."))
		//Sort the array "ByCount" aka, by duplicateString.count
		sort.Sort(ByCount(textArray))
	}

	outputLogs(textArray)
}

func main() {
	//Flag setup.
	initFlags()

	remainingFlags := flag.Args()
	ifVerbose(fmt.Sprintln("Remaining flags:", remainingFlags))
	if keyvaluepair {
		//Check if someone included a filename flag.
		if fileName != "" {
			//If so, prepend the flags with that file.
			remainingFlags = append([]string{fileName, input}, remainingFlags...)
			ifVerbose(fmt.Sprintln("Manual filename found, adding to remainingflags."))
			//Otherwise, check if there aren't enough flags.
		} else if len(remainingFlags) < 1 {
			//If so, hit the default.
			remainingFlags = append(remainingFlags, "artifactory.log")
			remainingFlags = append(remainingFlags, input)
			ifVerbose(fmt.Sprintln("Not enough arguments found, reverting to default."))
		}

		//If we're looking for multiple files..
	} else if multipleFiles {
		//Check if someone included a filename flag.
		if fileName != "" {
			//If so, prepend the flags with that file.
			remainingFlags = append([]string{fileName}, remainingFlags...)
			ifVerbose(fmt.Sprintln("Manual filename found, setting as remainingflags."))
			//Otherwise, check if there aren't enough flags.
		} else if len(remainingFlags) < 1 {
			//If so, hit the default.
			remainingFlags = append(remainingFlags, "artifactory.log")
			ifVerbose(fmt.Sprintln("Not enough arguments found, reverting to default."))
		}
		//If we're not looking for multiple files...
	} else {
		//Check if there are enough flags
		if len(remainingFlags) < 1 {
			// If not, hit the default.
			remainingFlags = append(remainingFlags, "artifactory.log")
			ifVerbose(fmt.Sprintln("Not enough arguments found, reverting to default."))
			//If there are enough inputs...
		} else {
			//Strip them all away since we're only checking one file.
			remainingFlags = remainingFlags[:1]
			ifVerbose(fmt.Sprintln("Stripping extra arguments."))
		}
	}
	if keyvaluepair {
		ifVerbose(fmt.Sprintln("Parsing through KVP."))
		length := len(remainingFlags)
		for i := 0; i < length-1; i += 2 {
			ifVerbose(fmt.Sprintln("Parsing KVP:", i/2))
			input = remainingFlags[i+1]
			parseLogs(remainingFlags[i])
		}
	} else {
		ifVerbose(fmt.Sprintln("Parsing through arguments."))
		for c, i := range remainingFlags {
			ifVerbose(fmt.Sprintln("Parsing through file number:", c+1))
			parseLogs(i)
		}
	}
	ifVerbose(fmt.Sprintln("Completed log parsing."))
	ifVerbose(fmt.Sprintln("Finishing up logs."))
	logOut()
}
