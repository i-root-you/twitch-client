package obsws

// This file is automatically generated.
// https://github.com/christopher-dG/go-obs-websocket/blob/master/codegen/protocol.py

// SwitchScenesEvent : Indicates a scene change.
//
// Since obs-websocket version: 0.3.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#switchscenes
type SwitchScenesEvent struct {
	// The new scene.
	// Required: Yes.
	SceneName string `json:"scene-name"`
	// List of sources in the new scene.
	// Required: Yes.
	Sources []interface{} `json:"sources"`
	_event  `json:",squash"`
}

// ScenesChangedEvent : The scene list has been modified.
// Scenes have been added, removed, or renamed.
//
// Since obs-websocket version: 0.3.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#sceneschanged
type ScenesChangedEvent struct {
	_event `json:",squash"`
}

// SceneCollectionChangedEvent : Triggered when switching to another scene collection or when renaming the current scene collection.
//
// Since obs-websocket version: 4.0.0.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#scenecollectionchanged
type SceneCollectionChangedEvent struct {
	_event `json:",squash"`
}

// SceneCollectionListChangedEvent : Triggered when a scene collection is created, added, renamed, or removed.
//
// Since obs-websocket version: 4.0.0.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#scenecollectionlistchanged
type SceneCollectionListChangedEvent struct {
	_event `json:",squash"`
}
