package main

import (
  "fmt"
  "regexp"
  "strconv"
  "time"
)

func main() {

  ps := []string{
    "2001-1",
    "2002.11",
    "2003/12",
    "20042",
    "200503",
    "200612",
    "209901",
  }

  for _, p := range ps {

    this, err := GetEndOfThisMonthString(p)
    if err != nil {
      fmt.Printf("unexpected error: %s => %s\n", p, err)
    } else {
      fmt.Printf("end of this month: %s => %s\n", p, this)
    }

    last, err := GetEndOfLastMonthString(p)
    if err != nil {
      fmt.Printf("unexpected error: %s => %s\n", p, err)
    } else {
      fmt.Printf("end of last month: %s => %s\n", p, last)
    }

  }
}

func ParseYYYYMM(s string) (year, month int, err error) {
  match := regexp.MustCompile(`^(\d{4})[-./]?([1-9]|0[1-9]|1[0-2])$`).FindStringSubmatch(s)
  if match == nil {
    return -1, -1, fmt.Errorf("yyyymm matching failure")
  }
  year, e1 := strconv.Atoi(match[1])
  if e1 != nil {
    return -1, -1, e1
  }
  month, e2 := strconv.Atoi(match[2])
  if e2 != nil {
    return -1, -1, e2
  }
  return year, month, nil
}


func GetEndOfThisMonthString(yyyymm string) (string, error) {
	t, err := GetFirstDayOfThisMonth(yyyymm)
  if err != nil {
    return "", err
  }
  t = t.AddDate(0, 1, 0)
  t = t.AddDate(0, 0, -1)
  return t.Format("2006/01/02"), nil
}

func GetEndOfLastMonthString(yyyymm string) (string, error) {
	t, err := GetFirstDayOfThisMonth(yyyymm)
  if err != nil {
    return "", err
  }
  t = t.AddDate(0, 0, -1)
  return t.Format("2006/01/02"), nil
}

func GetFirstDayOfThisMonth(yyyymm string) (time.Time, error) {
  y, m, err := ParseYYYYMM(yyyymm)
  if err != nil {
    return time.Time{}, err
  }
	return time.Date(y, time.Month(m), 1, 0, 0, 0,0, time.UTC), nil
}

/*
end of this month: 2001-1 => 2001/01/31
end of last month: 2001-1 => 2000/12/31
end of this month: 2002.11 => 2002/11/30
end of last month: 2002.11 => 2002/10/31
end of this month: 2003/12 => 2003/12/31
end of last month: 2003/12 => 2003/11/30
end of this month: 20042 => 2004/02/29
end of last month: 20042 => 2004/01/31
end of this month: 200503 => 2005/03/31
end of last month: 200503 => 2005/02/28
end of this month: 200612 => 2006/12/31
end of last month: 200612 => 2006/11/30
end of this month: 209901 => 2099/01/31
end of last month: 209901 => 2098/12/31
*/

