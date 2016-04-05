package util

import (
	"log"
	"testing"
)

func TestSave(t *testing.T) {
	c := Configuration{}
	c.ClientId = "Tom."
	c.save()
}

func TestLoad(t *testing.T) {
	c := Configuration{}
	c.load()
	if c.ClientId != "Tom." {
		log.Println("ClientId: " + c.ClientId)
		t.Fail()
	}
}
