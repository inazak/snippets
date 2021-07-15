Option Explicit

Private regexp As Variant

'Initializing a regular expression object
Function InitRegexp()
  Set regexp = CreateObject("VBScript.RegExp")
End Function

'Returns whether the regular expression matches or not.
Function IsMatch(re As String, text As String) As Boolean

  regexp.IgnoreCase = False
  regexp.Global = False
  regexp.Pattern = re

  IsMatch = regexp.Test(text)

End Function

'Return reference objects matched by regular expression.
Function GetSubmatch(re As String, text As String) As Object

  regexp.IgnoreCase = False
  regexp.Global = True
  regexp.Pattern = re

  Dim matches As Object
  Set matches = regexp.Execute(text)
  Set GetSubmatch = matches

End Function

'Replace with a regular expression
Function ReplaceString(re As String, text As String, rep As String)

  regexp.IgnoreCase = False
  regexp.Global = True
  regexp.Pattern = re

  ReplaceString = regexp.Replace(text, rep)

End Function

'Returns the numeric value of year by the string YYYYMM, or -1 if it cannot be interpreted.
Function ParseYearFromYYYYMM(s As String) As Long
  Dim y As Integer
  Dim matches As Object
  Set matches = GetSubmatch("^([0-9]{4})[-./]?[0-9]{2}$", s)
  If matches.Count = 0 Then
    y = -1
  Else
    y = CInt(matches.Item(0).SubMatches(0))
  End If
  ParseYearFromYYYYMM = y
End Function

'Returns the numeric value of month by the string YYYYMM, or -1 if it cannot be interpreted.
Function ParseMonthFromYYYYMM(s As String) As Long
  Dim m As Integer
  Dim matches As Object
  Set matches = GetSubmatch("^[0-9]{4}[-./]?([0-9]{2})$", s)
  If matches.Count = 0 Then
    m = -1
  Else
    m = CInt(matches.Item(0).SubMatches(0))
  End If
  ParseMonthFromYYYYMM = m
End Function

'Returns the date of the first day of the current month of the year specified by the string YYYYMM.
Function GetFirstDayOfMonthFromYYYYMM(s As String) As Date
  Dim y As Integer
  Dim m As Integer
  y = ParseYearFromYYYYMM(s)
  m = ParseMonthFromYYYYMM(s)

  If y = -1 Or m = -1 Then
    GetFirstDayOfMonthFromYYYYMM = Empty
  End If
  
  GetFirstDayOfMonthFromYYYYMM = DateSerial(y, m, 1)
End Function

'Returns the date of the last day of the current month of the year specified by the string YYYYMM.
Function GetEndOfDayOfThisMonthFromYYYYMM(s As String) As Date
  Dim d As Date
  Dim e As Date
  d = GetFirstDayOfMonthFromYYYYMM(s)
  e = DateSerial(Year(d), Month(d) + 1, 1)
  GetEndOfDayOfThisMonthFromYYYYMM = e - 1
End Function

'Returns the date of the last day of the previous month of the year specified by the string YYYYMM.
Function GetEndOfDayOfLastMonthFromYYYYMM(s As String) As Date
  GetEndOfDayOfLastMonthFromYYYYMM = GetFirstDayOfMonthFromYYYYMM(s) - 1
End Function


Sub Test()
 
  InitRegexp
  
  MsgBox "Expect:True,  got:" & IsMatch("^[0-9]", "1234")
  MsgBox "Expect:False, got:" & IsMatch("^[0-9]", "abcd")

  MsgBox "Expect:1234,  got:" & ReplaceString("[a-z]", "a1b2c3d4", "")

  MsgBox "Expect:2019/12/31,  got:" & Format(GetEndOfDayOfLastMonthFromYYYYMM("202001"), "YYYY/MM/DD")
  MsgBox "Expect:2020/02/29,  got:" & Format(GetEndOfDayOfThisMonthFromYYYYMM("202002"), "YYYY/MM/DD")

End Sub

