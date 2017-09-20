[![Build Status](https://travis-ci.org/wgoudsbloem/logdb.svg?branch=master)](https://travis-ci.org/wgoudsbloem/logdb)

This is a command line application.

This application maintains no state, everytime a view is request by the get command it 
 rebuilds the view with the state from the file.

to run the application from a linux machine simple execute the binary, or install go on your machine and execute go run logdb.go

Example:

```
goudsbloem@goudsbloem-MS-7597 ~/go/src/sample $ ./logdb  
welcome to logDB

> add an entry:                   put={"key":"<entry name>","value":<integer value>}`  
> get an sum of the values:       get{"key":"<entry name>", "value":"sum"}  
> get the last entry value:       get{"key":"<entry name>", "value":"last"}  

> put={"key":"mydb","value":123}  
ok {"key":"position","value":0}  
> put={"key":"mydb","value":27}  
ok {"key":"position","value":4}  
> get={"key":"mydb","value":"sum"}  
ok {"key":"sum","value":150}  
> put={"key":"mydb","value":-100}  
ok {"key":"position","value":7}  
> get={"key":"mydb","value":"sum"}  
ok {"key":"sum","value":50}  
> get={"key":"mydb","value":"last"}  
ok {"key":"last","value":-100}


the mydb file contains:

goudsbloem@goudsbloem-MS-7597 ~/go/src/sample $ cat mydb  
123  
27  
-100
```

code is hosted on https://github.com/wgoudsbloem/logdb