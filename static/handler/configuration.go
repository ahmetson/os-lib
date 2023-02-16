package handler

import (
	"github.com/blocklords/gosds/db"
	"github.com/blocklords/gosds/static/configuration"
	"github.com/blocklords/gosds/static/smartcontract"

	"github.com/blocklords/gosds/message"
)

// Register a new smartcontract in the configuration.
// It requires smartcontract address. First call smartcontract_register command.
func ConfigurationRegister(db *db.Database, request message.Request) message.Reply {
	if _, err := message.GetString(request.Parameters, "id"); err == nil {
		return message.Fail("parameter should not have the 'id' parameter. Its generated by database.")
	}
	conf, err := configuration.New(request.Parameters)
	if err != nil {
		return message.Fail(err.Error())
	}

	if configuration.ExistInDatabase(db, conf) {
		return message.Fail("Smartcontract found in the config")
	}

	if err = configuration.SetInDatabase(db, conf); err != nil {
		return message.Fail("Configuration saving in the database failed: " + err.Error())
	}

	return message.Reply{
		Status:  "OK",
		Message: "",
		Params:  conf.ToJSON(),
	}
}

// Returns configuration and smartcontract information related to the configuration
func ConfigurationGet(db *db.Database, request message.Request) message.Reply {
	if _, err := message.GetString(request.Parameters, "id"); err == nil {
		return message.Fail("parameter should not have the 'id' parameter. Its generated by database.")
	}
	if _, err := message.GetString(request.Parameters, "address"); err == nil {
		return message.Fail("parameter should not have the 'address' parameter. Its not needed for reading configuration.")
	}
	conf, err := configuration.New(request.Parameters)
	if err != nil {
		return message.Fail(err.Error())
	}

	if !configuration.ExistInDatabase(db, conf) {
		return message.Fail("Configuration not registered in the database")
	}

	err = configuration.LoadDatabaseParts(db, conf)
	if err != nil {
		return message.Fail("Configuration loading in the database failed: " + err.Error())
	}

	s, getErr := smartcontract.GetFromDatabase(db, conf.NetworkId, conf.Address)
	if getErr != nil {
		return message.Fail("Failed to get smartcontract from database: " + getErr.Error())
	}

	reply := message.Reply{
		Status:  "OK",
		Message: "",
		Params: map[string]interface{}{
			"configuration": conf.ToJSON(),
			"smartcontract": s.ToJSON(),
		},
	}
	return reply
}
