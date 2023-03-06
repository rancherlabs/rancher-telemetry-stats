package main

import (
	"testing"
)

func TestNewDates_List(t *testing.T) {
	datesList, err := NewDates("2022-12-24,2023-12-24,2024-12-24")
	if err != nil {
		t.Error(err)
	}
	stringifiedDates := datesList.String()
	expected := "2022-12-24,2023-12-24,2024-12-24"
	if stringifiedDates != expected {
		t.Errorf("%s != %s", stringifiedDates, expected)
	}
}

func TestNewDates_Range(t *testing.T) {
	datesList, err := NewDates("2022-12-24:2022-12-27")
	if err != nil {
		t.Error(err)
	}

	stringifiedDatesList := datesList.String()
	expected := "2022-12-24,2022-12-25,2022-12-26,2022-12-27"
	if stringifiedDatesList != expected {
		t.Errorf("%s != %s", stringifiedDatesList, expected)
	}
}

func TestNewDates_ListAndRange(t *testing.T) {
	datesList, err := NewDates("2022-12-24,2033-12-01:2033-12-03,2022-12-25")
	if err != nil {
		t.Error(err)
	}

	stringifiedDatesList := datesList.String()
	expected := "2022-12-24,2033-12-01,2033-12-02,2033-12-03,2022-12-25"
	if stringifiedDatesList != expected {
		t.Fail()
		t.Errorf("%s != %s", stringifiedDatesList, expected)
	}
}

func TestParseDate(t *testing.T) {
	date, err := parseDate("2022-12-1")
	if err != nil {
		t.Error(err)
	}

	if date.Format("2006-01-02") != "2022-12-01" {
		t.Fail()
	}

	date, err = parseDate("2022-12-01")
	if err != nil {
		t.Error(err)
	}

	if date.Format("2006-01-02") != "2022-12-01" {
		t.Fail()
	}
}
