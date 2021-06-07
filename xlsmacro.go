// +build windows

/*
Usage
  xlsmacro.exe ExcelFilePath.xls MacroName [Args ...]

  on the Excel side, write Macro like this

  Sub Something(arg1 As String, arg2 As String)
  ........
  End Sub
*/

package main

import (
  "fmt"
  "log"
  "os"
  ole "github.com/go-ole/go-ole"
  "github.com/go-ole/go-ole/oleutil"
)

func main() {

  if len(os.Args) < 3 {
    fmt.Println("Usage: xlsmacro XLSFile MacroName [MacroArgs...]")
    os.Exit(1)
  }

  xlsfile := os.Args[1]
  macroNameAndArgs := make([]interface{}, len(os.Args[2:]))
  for i, v := range os.Args[2:] {
    macroNameAndArgs[i] = v
  }

  ole.CoInitialize(0)
  defer ole.CoUninitialize()

  unknown, err := oleutil.CreateObject("Excel.Application")
  if err != nil {
    log.Fatal(err)
  }

  excel, err := unknown.QueryInterface(ole.IID_IDispatch)
  if err != nil {
    log.Fatal(err)
  }
  defer func() {
    oleutil.CallMethod(excel, "Quit")
    excel.Release()
  }()

  oleutil.PutProperty(excel, "Visible", false)

  workbooks := oleutil.MustGetProperty(excel, "Workbooks").ToIDispatch()
  defer func() {
    oleutil.CallMethod(workbooks, "Close")
    workbooks.Release()
  }()

  workbook, err := oleutil.CallMethod(workbooks, "Open", xlsfile)
  if err != nil {
    log.Fatal(err)
  }
  defer workbook.ToIDispatch().Release()

  _, errrun := oleutil.CallMethod(excel, "Run", macroNameAndArgs...)
  if errrun != nil {
    log.Fatal(errrun)
  }
}


