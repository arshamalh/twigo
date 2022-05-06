package entities

import "time"

type Space struct {
	ID               string    `json:"id"`
	State            string    `json:"state"` // It's a enum actually, not a string, so maybe we should parse it
	CreatorID        string    `json:"creator_id,omitempty"`
	Title            string    `json:"title,omitempty"`
	Lang             string    `json:"lang,omitempty"`
	HostIDs          []string  `json:"host_ids,omitempty"`
	InvitedUserIDs   []string  `json:"invited_user_ids,omitempty"`
	SpeakerIDs       []string  `json:"speaker_ids,omitempty"`
	IsTicketed       bool      `json:"is_ticketed,omitempty"`
	SubscriberCount  int       `json:"subscriber_count,omitempty"`
	ParticipantCount int       `json:"participant_count,omitempty"`
	ScheduledStart   time.Time `json:"scheduled_start,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	StartedAt        time.Time `json:"started_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
	EndedAt          time.Time `json:"ended_at,omitempty"`
}
