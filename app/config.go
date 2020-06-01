package main

type AppConfig struct {
	TrollBin string `config:"troll-bin,required"`
	Mosmllib string `config:"mosmllib"`
	Listen   string `config:"listen"`
}
