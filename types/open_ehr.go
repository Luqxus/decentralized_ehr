package types

import "time"

type OpenEHR struct {
	ArchetypeNodeID string        `json:"archetype_node_id"`
	Name            Name          `json:"name"`
	Patient         Patient       `json:"patient"`
	Context         Context       `json:"context"`
	Content         []ContentItem `json:"content"`
}

type Context struct {
	StartTime          time.Time          `json:"start_time"`
	EndTime            time.Time          `json:"end_time"`
	Location           string             `json:"location"`
	HealthcareFacility HealthcareFacility `json:"health_care_facility"`
}

type HealthcareFacility struct {
	Name string `json:"name"`
}

type ContentItem struct {
	ArchetypeNodeID string   `json:"archetype_node_id"`
	Name            Name     `json:"name"`
	Data            Data     `json:"data,omitempty"`
	Protocol        Protocol `json:"protocol,omitempty"`
}

type Data struct {
	ArchetypeNodeID string  `json:"archetype_node_id"`
	Events          []Event `json:"events"`
}

type Event struct {
	ArchetypeNodeID string    `json:"archetype_node_id"`
	Data            EventData `json:"data"`
}

type EventData struct {
	ArchetypeNodeID string `json:"archetype_node_id"`
	Items           []Item `json:"items"`
}

type Item struct {
	ArchetypeNodeID string           `json:"archetype_node_id"`
	Name            Name             `json:"name"`
	Value           MeasurementValue `json:"value"`
}

type MeasurementValue struct {
	Magnitude float64 `json:"magnitude"`
	Units     string  `json:"units"`
}

type Protocol struct {
	ArchetypeNodeID string          `json:"archetype_node_id"`
	Items           []ProtocolItems `json:"items"`
}

type ProtocolItems struct {
	ArchetypeNodeID string `json:"archetype_node_id"`
	Name            Name   `json:"name"`
	Value           Value  `json:"value"`
}

type Value struct {
	Value string `json:"value"`
}

type Patient struct{}

type Name struct {
	Value string `json:"value"`
}
