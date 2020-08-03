Sub 按钮1_Click()

Dim w1 As Workbook
Dim s1 As Worksheet
Dim w2 As Workbook
Dim s2 As Worksheet


Set w1 = Workbooks.Open("D:\copy\test.xlsx")
Set w2 = Workbooks.Open("D:\copy\main.xlsx")

Set s1 = w1.Worksheets("test")



Set s2 = w2.Worksheets("jieguo")

's1.Range("A2:D20").Copy Destination:=s2.Range("B2:E20")

s1.Range(s1.Cells(2, 1), s1.Cells(20, 4)).Copy Destination:=s2.Range(s2.Cells(2, 2), s2.Cells(20, 5))
s2.Range(s2.Cells(2, 1), s2.Cells(20, 1)).Value = "890"
s2.Range(s2.Cells(2, 1), s2.Cells(20, 1)).Borders.LineStyle = 1

End Sub

