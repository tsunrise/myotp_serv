package test

import (
	"myotp_serv/util/urlUtil"
	"testing"
)

func TestURLMatchExact(t *testing.T) {
	if !urlUtil.MatchExact("", "") {
		t.Fail()
	}
	if !urlUtil.MatchExact("/", "") {
		t.Fail()
	}
	if !urlUtil.MatchExact("", "/") {
		t.Fail()
	}
	if !urlUtil.MatchExact("ab", "ab/") {
		t.Fail()
	}
	if !urlUtil.MatchExact("ab/c/", "ab/c") {
		t.Fail()
	}
	if !urlUtil.MatchExact("a/ab/cd/", "/a/ab/cd") {
		t.Fail()
	}
	if !!urlUtil.MatchExact("ab/c/d", "ab/c") {
		t.Fail()
	}
	if !!urlUtil.MatchExact("ab/d/c", "ab/d/c/e") {
		t.Fail()
	}
	if !!urlUtil.MatchExact("ab/c/d", "ab/c/dd") {
		t.Fail()
	}
}

func TestURLMatch(t *testing.T) {
	if !urlUtil.Match("/ab/cd/e/f", "/ab/cd") {
		t.Fail()
	}
	if !urlUtil.Match("/ab/cd/e/f", "/ab/cd/") {
		t.Fail()
	}
	if !urlUtil.Match("/", "") {
		t.Fail()
	}
	if !urlUtil.Match("/ab/cd/e/f", "/ab/cd/e/f/") {
		t.Fail()
	}
	if !urlUtil.Match("/ab/12/fd/56", "/ab/12/fd") {
		t.Fail()
	}
	if !!urlUtil.Match("/ab/", "/ab/cd") {
		t.Fail()
	}
	if !!urlUtil.Match("/ab/cd/e/f", "/ab/cd/e/f/g") {
		t.Fail()
	}
	if !!urlUtil.Match("/ab/cd/e/f", "/ab/cd/d/f") {
		t.Fail()
	}
}
