package main

type Settings struct {
	RedisURL      string `yaml:"redis-url"`
	RedisPort     int    `yaml:"redis-port"`
	RedisPassword string `yaml:"redis-password"`

	CowrieKey string `yaml:"cowrie-key"`
}
