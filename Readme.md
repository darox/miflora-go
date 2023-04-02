# miflora-go

A Golang application for reading data from Xiaomi Mi Flora plant sensors.

```
📡  Scanning for Monstera
✅  Connected to Monstera
👋  Disconnected from Monstera

🪴  Name: Monstera 
🔋  Battery Level: 33% 
⚙️   Firmware: 3.2.2 
🌡️  Temperature: 24.1°C 
⚡   Light: 15533 Lux 
💧  Moisture: 21% 
🌱  Conductivity: 190 µS/cm 
```

## Features

- Configurable by YAML file
- Retrieves data from the sensor, such as battery, humidity, conductivity, soil moisture and temperature.

## Installation

1. Install Golang on your computer if you don't have it already installed
2. Clone the repository: `git clone https://github.com/darox/miflora-go`
3. Build the application: `cd cmd/miflora-go && go build`
4. Add capabilities to run as none-root user: `sudo setcap 'cap_net_raw,cap_net_admin+eip' miflora-go`


## Usage

Under Linux, the application uses the mac address to connect to devices; under MacOs the UUID.


1. Copy `config/config.yaml` to the same folder where miflora-go will run or use the param `--config-path` to specify the path of the config file
2. The application will now scan the device and printout the result

## Contributing

1. Fork the repository
2. Create your feature branch: git checkout -b my-new-feature
3. Commit your changes: git commit -am 'Add some feature'
4. Push to the branch: git push origin my-new-feature
5. Submit a pull request

### License

BSD-3-Clause license

### Acknowledgments

- Xiaomi for creating the sensor
- [Creators of go-ble](https://github.com/go-ble/ble)
- [Creators of miflora wiki](https://github.com/ChrisScheffler/miflora/wiki/The-Basics)
