package grpc

import (
	"testing"
	"time"

	"github.com/derticom/home_work/hw12_13_14_15_calendar/internal/model"
	"github.com/derticom/home_work/hw12_13_14_15_calendar/pb"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_modelToPb(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	type args struct {
		events []model.Event
	}
	tests := []struct {
		name string
		args args
		want []*pb.Event
	}{
		{
			name: "case first",
			args: args{
				events: []model.Event{
					{
						ID:           model.EventUUID(id),
						Header:       "head",
						Date:         now,
						Duration:     100 * time.Second,
						Description:  "description",
						NotifyBefore: 300 * time.Second,
					},
				},
			},
			want: []*pb.Event{
				{
					Uuid:         id.String(),
					Header:       "head",
					Date:         timestamppb.New(now),
					Duration:     100,
					Description:  "description",
					NotifyBefore: 300,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := modelToPb(tt.args.events)
			assert.Equal(t, tt.want, got)
		})
	}
}
