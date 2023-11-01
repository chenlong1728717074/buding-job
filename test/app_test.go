package test

import (
	"fmt"
	"testing"
	"time"
)

func TestCobra(t *testing.T) {
	//_ = do.JobInfoDo{}
	//_ = do.JobManagementDo{}
}

func TestTime(t *testing.T) {
	now := time.Now().Add(-time.Second * 90)
	fmt.Println(now)
	fmt.Println(now.Before(time.Now()))
}
func TestArr(t *testing.T) {
	fmt.Println(len(getSince()))

}
func getSince() []int {
	return nil
}
