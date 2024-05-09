package app

var Registry = map[string]ResourceMetadata{
	"broker": {
		Kind:   "broker",
		Render: nil,
	},

	"topic": {
		Kind:   "topic",
		Render: nil,
	},

	"consumer": {
		Kind:   "consumer",
		Render: nil,
	},

	"quota": {
		Kind:   "quota",
		Render: nil,
	},

	"schema": {
		Kind:   "schema",
		Render: nil,
	},

	"acl": {
		Kind:   "acl",
		Render: nil,
	},

	"group": {
		Kind:   "group",
		Render: nil,
	},
}
