Option Explicit

' 固定長データ生成用VBAスクリプト
'
' 下記のような表を作ることで、指定したBytes数の固定長レコードを作る。
' 例えば下記では、F列とG列がそれぞれ一つのレコードとして作成され、
' 改行CRLF区切りでテキストファイルが生成される。
'
' 列の先頭のデータが入っている場合だけ、ファイル出力対象にしている。
' このデータを Function HasKeyNo で6桁の数値かどうかチェックしている。
' 実データに合わせて修正する必要がある。
'
' パディングはAttr列として指定した場所に
' 数値 と入っていた場合はゼロ埋め
' 半角 と入っていた場合は半角空白埋め
' 全角 と入っていた場合は全角空白埋め
' として処理する
' 
' 
' （作成する表のイメージ）
' 
'       A    B      C      D       E     F        G 
'     +----+------+------+-------+----+---------+----------
'   1 | No | Name | Attr | Bytes | .. | 例1     | 例2
'   2 |  1 | 番号 | 数値 |    10 | .. | 202001  | 202002
'   3 |  2 | 名前 | 全角 |    10 | .. | 山田    | 鈴木
'   4 |  3 | カナ | 全角 |    10 | .. | ヤマダ  | スズキ
'   5 |  4 | 英字 | 半角 |    10 | .. | Yamada  | Suzuki
'   6 |  5 | 年齢 | 数値 |    10 | .. | 21      | 42
'   . | .. | .... | .... | ..... | .. | ...     | ...
'  21 | 20 | 備考 | 全角 |    10 | .. | なし    | あり
'     +----+------+------+-------+----+---------+----------
'  シート名：データシート
' 
' 
' （出力イメージ）
'
' 0000202001山田　　　ヤマダ　　Yamada    0000000021 ....なし　　　[CRLF]
' 0000202002鈴木　　　スズキ　　Suzuki    0000000042 ....あり　　　[CRLF]
' 


Const DataSheetName As String = "データシート"
Const ColumnNoOfName  As Integer = 2
Const ColumnNoOfAttr  As Integer = 3
Const ColumnNoOfBytes As Integer = 4
Const RowNoOfDataBegin As Integer = 2
Const RowSizeOfData As Integer = 20
Const ColumnNoOfDataBegin = 6
Const ColumnSizeOfData = 16


''' SJIS計算でのByte数を返す
Function WidthBySJIS(s As String) As Integer

  WidthBySJIS = LenB(StrConv(s, vbFromUnicode))

End Function


''' 数字0でパディングした文字列を返す、文字列が指定数より長い場合は、パディングせずにそのまま返す
Function PaddingWithZero(str As String, length As Integer) As String

  If WidthBySJIS(str) >= length Then
    PaddingWithZero = str
  Else
    Dim rest As Integer: rest = length - WidthBySJIS(str)
    PaddingWithZero = String(rest, "0") & str
  End If

End Function


''' 半角空白でパディングした文字列を返す、文字列が指定数より長い場合は、パディングせずにそのまま返す
Function PaddingWithHalfSpace(str As String, length As Integer) As String

  If WidthBySJIS(str) >= length Then
    PaddingWithHalfSpace = str
  Else
    Dim rest As Integer: rest = length - WidthBySJIS(str)
    PaddingWithHalfSpace = str & String(rest, " ")
  End If

End Function


''' 全角空白でパディングした文字列を返す、文字列が指定数より長い場合は、パディングせずにそのまま返す
Function PaddingWithFullSpace(str As String, length As Integer) As String

  If WidthBySJIS(str) >= length Then
    PaddingWithFullSpace = str
  Else
    Dim rest As Integer: rest = (length - WidthBySJIS(str)) / 2
    PaddingWithFullSpace = str & String(rest, "　")
  End If

End Function


''' 指定した属性でパディングした文字列を返す、文字列が指定数より長い場合は、パディングせずにそのまま返す
Function Padding(data As String, attr As String, length As Integer) As String

  If attr = "数値" Then
    Padding = PaddingWithZero(data, length)
  ElseIf attr = "半角" Then
    Padding = PaddingWithHalfSpace(data, length)
  ElseIf attr = "全角" Then
    Padding = PaddingWithFullSpace(data, length)
  Else
    Padding = "Error: 不明なデータ属性：" & attr
  End If


End Function


''' 指定した列において、各行ごとのByte数以下に収まっていない場合、メッセージを生成、収まっている場合は、空文字列を返す
Function CheckDataLength(cellcolumnno As Integer) As String

  Dim result As String
  result = ""

  Dim i As Integer
  For i = RowNoOfDataBegin To RowNoOfDataBegin + RowSizeOfData
  
    With Sheets(DataSheetName)
    
      Dim attr  As String
      Dim bytes As Integer
      
      attr = CStr(Cells(i, ColumnNoOfAttr).Value)
      bytes = CInt(Cells(i, ColumnNoOfBytes).Value)
      
      Dim data  As String
      Dim b As Integer
     
      data = CStr(Cells(i, cellcolumnno).Value)
      b = WidthBySJIS(data)
      
      If b > bytes Then
      
        Dim name As String
        name = CStr(Cells(i, ColumnNoOfName).Value)
        
        result = result & "文字数超過：" & name & " / 入力byte数:" & CStr(b) & " / 入力可能byte数" & CStr(bytes) & vbCrLf
        
      End If
  
    End With
  
  Next

  CheckDataLength = result

End Function


''' 指定した列において、固定長のテキストデータを返す
Function GetFixedTextFrom(cellcolumnno As Integer) As String

  Dim text As String

  Dim i As Integer
  For i = RowNoOfDataBegin To RowNoOfDataBegin + RowSizeOfData
  
    With Sheets(DataSheetName)
    
      Dim attr  As String
      Dim bytes As Integer
      
      attr = CStr(Cells(i, ColumnNoOfAttr).Value)
      bytes = CInt(Cells(i, ColumnNoOfBytes).Value)
      
      Dim data  As String
      Dim fixed As String
     
      data = CStr(Cells(i, cellcolumnno).Value)
    
      fixed = Padding(data, attr, bytes)
      
      text = text & fixed
    
    End With
  
  Next

  GetFixedTextFrom = text

End Function


''' すべての列において、各行ごとのByte数以下に収まっていない場合、メッセージを生成、収まっている場合は、空文字列を返す
Function CheckAllData() As String

  Dim i As Integer
  
  For i = ColumnNoOfDataBegin To ColumnNoOfDataBegin + ColumnSizeOfData

    Dim check As String
    check = CheckDataLength(i)

    If check <> "" Then
      
      Dim no As Integer
      no = i - ColumnNoOfDataBegin + 1
      
      CheckAllData = "列" & CStr(no) & vbCrLf & check
      Exit Function
      
    End If

  Next

  CheckAllData = ""

End Function


''' 指定した列において、先頭のキーデータが6桁の数字列かどうかをチェックする
Function HasKeyNo(cellcolumnno As Integer) As Boolean

  HasKeyNo = False

  Dim keyno As String
      
  keyno = CStr(Cells(RowNoOfDataBegin, cellcolumnno).Value)
  
  If keyno Like "######" Then
    
    HasKeyNo = True
    
  End If

End Function


''' このままGetAllData()を実行した場合、どの列が対象になるかを返す、一件も対象が無い場合、空文字列を返す
Function ListOutputableColumn() As String

  Dim result As String
  result = ""
  
  Dim i As Integer
  For i = ColumnNoOfDataBegin To ColumnNoOfDataBegin + ColumnSizeOfData

    If HasKeyNo(i) Then
      
      Dim no As Integer
      no = i - ColumnNoOfDataBegin + 1
      
      result = result & "列" & CStr(no) & vbCrLf
    
    End If

  Next

  ListOutputableColumn = result

End Function


''' 全ての列を対象にした、最終的な固定長データを返す
Function GetAllData() As String

  Dim text As String
  Dim i As Integer
  
  For i = ColumnNoOfDataBegin To ColumnNoOfDataBegin + ColumnSizeOfData

    If HasKeyNo(i) Then

      Dim record As String
      record = GetFixedTextFrom(i) & vbCrLf
      
      text = text & record
    
    End If

  Next

  GetAllData = text
  
End Function



Sub ファイル作成()

  ''' 入力データチェック
  Dim check As String

  check = CheckAllData()
  
  If check <> "" Then
  
    MsgBox check, vbOKOnly, "入力チェックでエラーが見つかりました"
    Exit Sub
    
  End If
  
  
  ''' ダイアログ表示
  Dim filepath As Variant
  
  filepath = Application.GetSaveAsFilename(Title:="保存先をファイル名を指定してください", _
                                            FileFilter:="テキスト形式,*.txt")
  
  If VarType(filepath) = vbBoolean Then
    ' Do nothing
    Exit Sub
  End If

  ''' 出力対象番号のチェック
  Dim target As String
  target = ListOutputableColumn()
  
  If target = "" Then
    MsgBox "出力対象が見つかりません" & vbCrLf & "入力済の場合は、キーデータが正しいかチェックして下さい", vbOKOnly, "対象なし"
    Exit Sub
  Else
    MsgBox target & vbCrLf & "を出力します", vbOKOnly, "ファイル書き込み処理"
  End If

  ''' ファイル保存
  Dim text As String
  text = GetAllData()

  Open filepath For Output As #1
    Print #1, text;
  Close #1

  MsgBox "完了しました", vbOKOnly, "ファイル書き込み処理"
  
End Sub
