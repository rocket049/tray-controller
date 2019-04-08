package main

import (
	"encoding/json"
	"io/ioutil"
)

type myCfg struct {
	Exec string
	Args string
	Envs string
	Wd   string
}

func GetCfgFromJSON(filename string) (*myCfg, error) {
	var r1 = new(myCfg)
	buf1, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf1, &r1)
	if err != nil {
		return nil, err
	}
	return r1, nil
}
