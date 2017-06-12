# goRex [![Build Status](https://travis-ci.org/amller/goRex.svg?branch=master)](https://travis-ci.org/amller/goRex) [![GoDoc](https://godoc.org/github.com/amller/goRex?status.svg)](https://godoc.org/github.com/amller/goRex)
GoRex is a tool which allows to convert textual data into JSON format.  
Note that goRex was build with the idea in mind of converting logs.

The data is read from os.Stdin and written to os.Stdout.  
The file containing the regular expression is set with the flag `regexp`. The default value is ".regexp".  

The conversion is based on a regular expression as following:  
Every line will be evaluated for itself, creating a new JSON object.  
GoRex evaluates then the regular expression and checks for the capture groups.  
The JSON has one field for every caputure group and the matches on the textual data
will be assigned to their corresponding JSON fields.  

For more information about capture groups read https://golang.org/pkg/regexp/syntax  

### Example
.regexp:  
```
\[(?P<Tool>[a-zA-Z0-9]+)\]\[(?P<Date>[0-3][0-9]\.[01][0-9]\.[0-9][0-9]:[0-2][0-9]:[0-6][0-9])\] (?P<MessageType>[a-zA-Z]*): (?P<Message>.*)
```

os.Stdin:
```
[Tool1][07.04.12:12:25] Warning: A warning message!
[Tool2][07.04.12:12:30] Error: An error message!
[Tool3][07.04.12:13:01] Info: An info message!
[Tool2][07.04.12:16:30] Error: An error message!
[Tool2][07.04.12:16:30] Warning: A warning message!
[Tool2][07.04.12:16:30] Warning: A warning message!
[Tool2][07.04.12:16:30] Warning: A warning message!
[Tool3][07.04.12:16:44] Info: An info message!
[Tool3][07.04.12:18:30] Update: An update message!
[Tool3][07.04.12:18:34] Warning: A warning message!
[Tool1][07.04.12:18:44] Info: An info message!
```

os.Stdout:
```json
[
	{
		"Date": "07.04.12:12:25",
		"Message": "A warning message!",
		"MessageType": "Warning",
		"Tool": "Tool1"
	},
	{
		"Date": "07.04.12:12:30",
		"Message": "An error message!",
		"MessageType": "Error",
		"Tool": "Tool2"
	},
	{
		"Date": "07.04.12:13:01",
		"Message": "An info message!",
		"MessageType": "Info",
		"Tool": "Tool3"
	},
	{
		"Date": "07.04.12:16:30",
		"Message": "An error message!",
		"MessageType": "Error",
		"Tool": "Tool2"
	},
	{
		"Date": "07.04.12:16:30",
		"Message": "A warning message!",
		"MessageType": "Warning",
		"Tool": "Tool2"
	},
	{
		"Date": "07.04.12:16:30",
		"Message": "A warning message!",
		"MessageType": "Warning",
		"Tool": "Tool2"
	},
	{
		"Date": "07.04.12:16:30",
		"Message": "A warning message!",
		"MessageType": "Warning",
		"Tool": "Tool2"
	},
	{
		"Date": "07.04.12:16:44",
		"Message": "An info message!",
		"MessageType": "Info",
		"Tool": "Tool3"
	},
	{
		"Date": "07.04.12:18:30",
		"Message": "An update message!",
		"MessageType": "Update",
		"Tool": "Tool3"
	},
	{
		"Date": "07.04.12:18:34",
		"Message": "A warning message!",
		"MessageType": "Warning",
		"Tool": "Tool3"
	}
]
```
