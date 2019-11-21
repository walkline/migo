package sqlscanner

import (
	"strings"
	"testing"
)

func TestTable(t *testing.T) {
	for _, testCase := range []struct {
		str      string
		expected []string
	}{
		{"SELECT name FROM users WHERE name ILIKE '%;%';", []string{"SELECT name FROM users WHERE name ILIKE '%;%';"}},
		{"", []string{}},
		{"nothing", []string{}},
		{"SELECT name FROM users", []string{}},
		{" SELECT name FROM users; ", []string{"SELECT name FROM users;"}},
		{`
			SELECT name FROM users;
			`, []string{"SELECT name FROM users;"}},
		{" SELECT name FROM users; ", []string{"SELECT name FROM users;"}},
		{`
			-- SELECT name FROM users;
			SELECT name FROM users;
			`, []string{"SELECT name FROM users;"}},
	} {
		actualResult := []string{}

		sqlScanner := NewSQLScanner(strings.NewReader(testCase.str))
		var query string
		for sqlScanner.Next(&query) {
			actualResult = append(actualResult, query)
		}

		if sqlScanner.Error != nil {
			t.Error(sqlScanner.Error)
		}

		if len(actualResult) != len(testCase.expected) {
			t.Error("actualResult != expectedResult")
		}

		for i := 0; i < len(actualResult); i++ {
			if actualResult[i] != testCase.expected[i] {
				t.Errorf("expected: %v; actual: %v", testCase.expected[i], actualResult[i])
			}
		}
	}
}

func TestReadSingleSQLLine(t *testing.T) {
	expectedResult := "SELECT name FROM users WHERE name ILIKE '%;%';"
	reader := strings.NewReader(expectedResult)
	sqlScanner := NewSQLScanner(reader)

	var query string
	for sqlScanner.Next(&query) {
		if query != expectedResult {
			t.Errorf("actual: %s; expected: %s", query, expectedResult)
		}
	}

	if sqlScanner.Error != nil {
		t.Error(sqlScanner.Error)
	}
}

func TestReadSQLDump(t *testing.T) {
	reader := strings.NewReader(`--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.5
-- Dumped by pg_dump version 9.6.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural; language';


--
-- Name: pg_trgm; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;
`)

	expectedResult := []string{
		"SET statement_timeout = 0;",
		"SET lock_timeout = 0;",
		"SET idle_in_transaction_session_timeout = 0;",
		"SET client_encoding = 'UTF8';",
		"SET standard_conforming_strings = on;",
		"SET check_function_bodies = false;",
		"SET client_min_messages = warning;",
		"SET row_security = off;",
		"CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;",
		"COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural; language';",
		"CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;",
	}

	actualResult := []string{}

	sqlScanner := NewSQLScanner(reader)

	var query string
	for sqlScanner.Next(&query) {
		actualResult = append(actualResult, query)
	}

	if sqlScanner.Error != nil {
		t.Error(sqlScanner.Error)
	}

	if len(actualResult) != len(expectedResult) {
		t.Error("actualResult != expectedResult")
	}

	for i := 0; i < len(actualResult); i++ {
		if actualResult[i] != expectedResult[i] {
			t.Errorf("expected: %v; actual: %v", expectedResult[i], actualResult[i])
		}
	}
}
