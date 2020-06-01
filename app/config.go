package main

type AppConfig struct {
	Listen   string `config:"listen"`
	LogFile  string `config:"log-file"`
	Mosmllib string `config:"mosmllib"`
	TrollBin string `config:"troll-bin,required"`
}
