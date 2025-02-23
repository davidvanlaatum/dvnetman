package dto

import (
	"context"
	mongoadapt "dvnetman/pkg/mongo/adapt/mock"
	"dvnetman/pkg/mongo/modal"
	"dvnetman/pkg/openapi"
	"dvnetman/pkg/testutils"
	"dvnetman/pkg/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var testUUID1 = uuid.MustParse("8793D858-5D4C-4D64-866E-7F80DC89F2B1")
var testUUID2 = uuid.MustParse("B42B0114-D69D-4EE8-9681-260C2ECFD7F6")
var testUUID3 = uuid.MustParse("3F469621-8779-4C2D-BCEB-801EA8BD74FB")

func TestConverter_DeviceToOpenAPI(t *testing.T) {
	deviceTypeSetup := func(db *mongoadapt.MockMongoDatabase) {
		manufacturerCursor := mongoadapt.NewMockMongoCursor(t)
		manufacturerCollection := mongoadapt.NewMockMongoCollection(t)
		deviceTypeCursor := mongoadapt.NewMockMongoCursor(t)
		deviceTypeCollection := mongoadapt.NewMockMongoCollection(t)
		db.EXPECT().Collection("device_types").Return(deviceTypeCollection)
		db.EXPECT().Collection("manufacturers").Return(manufacturerCollection)
		deviceTypeCollection.EXPECT().Find(mock.Anything, mock.Anything).Return(deviceTypeCursor, nil)
		deviceTypeCursor.EXPECT().All(mock.Anything, mock.Anything).RunAndReturn(
			func(ctx context.Context, i interface{}) error {
				*(i.(*[]*modal.DeviceType)) = []*modal.DeviceType{
					{
						Base: modal.Base{
							ID: (*modal.UUID)(&testUUID2),
						},
						Manufacturer: (*modal.UUID)(&testUUID3),
						Model:        "test device type",
					},
				}
				return nil
			},
		)
		manufacturerCollection.EXPECT().Find(mock.Anything, mock.Anything, mock.Anything).Return(
			manufacturerCursor, nil,
		)
		manufacturerCursor.EXPECT().All(mock.Anything, mock.Anything).RunAndReturn(
			func(ctx context.Context, i interface{}) error {
				*(i.(*[]*modal.Manufacturer)) = []*modal.Manufacturer{
					{
						Base: modal.Base{
							ID: (*modal.UUID)(&testUUID3),
						},
						Name: "test manufacturer",
					},
				}
				return nil
			},
		)
	}

	tests := []struct {
		name     string
		in       *modal.Device
		expected *openapi.Device
		err      error
		setup    func(database *mongoadapt.MockMongoDatabase)
	}{
		{
			name: "basic fields should be copied",
			in: &modal.Device{
				Base: modal.Base{
					ID:      (*modal.UUID)(&testUUID1),
					Version: 1,
				},
				Name: utils.ToPtr("test name"),
			},
			expected: &openapi.Device{
				Id:      testUUID1,
				Name:    utils.ToPtr("test name"),
				Version: 1,
			},
		},
		{
			name: "device type should be resolved",
			in: &modal.Device{
				Base: modal.Base{
					ID: (*modal.UUID)(&testUUID1),
				},
				DeviceType: (*modal.UUID)(&testUUID2),
			},
			expected: &openapi.Device{
				Id: testUUID1,
				DeviceType: &openapi.ObjectReference{
					Id:          testUUID2,
					DisplayName: utils.ToPtr("test manufacturer test device type"),
				},
				Version: 0,
			},
			setup: deviceTypeSetup,
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				r := require.New(t)
				database := mongoadapt.NewMockMongoDatabase(t)
				if test.setup != nil {
					test.setup(database)
				}
				client := modal.NewDBClient(database)
				converter := NewConverter(client)
				ctx := testutils.GetTestContext(t)
				d, err := converter.DeviceToOpenAPI(ctx, test.in)
				r.Equal(test.err, err)
				r.EqualExportedValues(test.expected, d)
			},
		)
	}
}

func TestConverter_DeviceToOpenAPISearchResults(t *testing.T) {

}

func TestConverter_UpdateDeviceFromOpenAPI(t *testing.T) {
	tests := []struct {
		name     string
		in       *openapi.Device
		existing *modal.Device
		expected *modal.Device
		err      error
		setup    func(database *mongoadapt.MockMongoDatabase)
	}{
		{
			name: "id should not be updated",
			in: &openapi.Device{
				Id: testUUID1,
			},
			existing: &modal.Device{
				Base: modal.Base{
					ID: (*modal.UUID)(&testUUID2),
				},
			},
			expected: &modal.Device{
				Base: modal.Base{
					ID: (*modal.UUID)(&testUUID2),
				},
			},
		},
		{
			name: "device type should be updated",
			in: &openapi.Device{
				DeviceType: &openapi.ObjectReference{
					Id: testUUID1,
				},
			},
			existing: &modal.Device{
				DeviceType: (*modal.UUID)(&testUUID2),
			},
			expected: &modal.Device{
				DeviceType: (*modal.UUID)(&testUUID1),
			},
			setup: func(db *mongoadapt.MockMongoDatabase) {
				deviceTypeCursor := mongoadapt.NewMockMongoCursor(t)
				deviceTypeCollection := mongoadapt.NewMockMongoCollection(t)
				db.EXPECT().Collection("device_types").Return(deviceTypeCollection)
				deviceTypeCollection.EXPECT().Find(mock.Anything, mock.Anything).Return(deviceTypeCursor, nil)
				deviceTypeCursor.EXPECT().All(mock.Anything, mock.Anything).RunAndReturn(
					func(ctx context.Context, i interface{}) error {
						*(i.(*[]*modal.DeviceType)) = []*modal.DeviceType{
							{
								Base: modal.Base{
									ID: (*modal.UUID)(&testUUID1),
								},
							},
						}
						return nil
					},
				)
			},
		},
		{
			name: "basic fields should be updated if set",
			in: &openapi.Device{
				Name:         utils.ToPtr("test name"),
				Description:  utils.ToPtr("test description"),
				RackPosition: utils.ToPtr(1.0),
			},
			expected: &modal.Device{
				Name:         utils.ToPtr("test name"),
				Description:  utils.ToPtr("test description"),
				RackPosition: utils.ToPtr(2),
			},
		},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				r := require.New(t)
				database := mongoadapt.NewMockMongoDatabase(t)
				if test.setup != nil {
					test.setup(database)
				}
				client := modal.NewDBClient(database)
				converter := NewConverter(client)
				ctx := testutils.GetTestContext(t)
				mod := test.existing
				if mod == nil {
					mod = &modal.Device{}
				}
				err := converter.UpdateDeviceFromOpenAPI(ctx, test.in, mod)
				r.Equal(test.err, err)
				r.EqualExportedValues(test.expected, mod)
			},
		)
	}
}
