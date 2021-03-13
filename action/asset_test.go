package action

import (
	"testing"

	"github.com/mitchellh/mapstructure"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/usecase/mock_usecase"

	"github.com/golang/mock/gomock"
	"github.com/mpppk/imagine/testutil"

	fsa "github.com/mpppk/lorca-fsa"
)

var assetAction = &assetActionCreator{}

func newAssetChannel(assets []*model.Asset) chan *model.Asset {
	c := make(chan *model.Asset, len(assets))
	for _, asset := range assets {
		c <- asset
	}
	return c
}

func decodePayload(t *testing.T, inputPayload, outputPayload interface{}) {
	if err := mapstructure.Decode(inputPayload, outputPayload); err != nil {
		t.Fatalf("failed to parse action payload: %#v", inputPayload)
	}
}

func decodeAssetScanRequestPayload(t *testing.T, inputPayload interface{}) *assetScanRequestPayload {
	var payload assetScanRequestPayload
	decodePayload(t, inputPayload, &payload)
	return &payload
}

func Test_assetScanHandler_Do(t *testing.T) {
	type args struct {
		action *fsa.Action
	}
	tests := []struct {
		name        string
		args        args
		existAssets []*model.Asset
		wantActions []*fsa.Action
		wantErr     bool
	}{
		{
			name: "",
			args: args{
				&fsa.Action{
					Type: "",
					Payload: &assetScanRequestPayload{
						WsPayload:  WsPayload{"default-workspace"},
						RequestNum: 2,
						Queries:    nil,
						Reset:      false,
					},
					Error: false,
					Meta:  nil,
				},
			},
			existAssets: []*model.Asset{
				{ID: 1, Path: "path1", Name: "path1"},
				{ID: 2, Path: "path2", Name: "path2"},
			},
			wantActions: []*fsa.Action{
				assetAction.scanRunning("default-workspace", []*model.Asset{
					{ID: 1, Path: "path1", Name: "path1"},
					{ID: 2, Path: "path2", Name: "path2"},
				}, 2),
			},
			wantErr: false,
		},
		// TODO: scanのcancel、queries、2回目のscanでcntが引き継がれているかなどをテスト
		{
			name: "",
			args: args{
				&fsa.Action{
					Type: "",
					Payload: &assetScanRequestPayload{
						WsPayload:  WsPayload{"default-workspace"},
						RequestNum: 1,
						Queries:    nil,
						Reset:      false,
					},
					Error: false,
					Meta:  nil,
				},
			},
			existAssets: []*model.Asset{
				{ID: 1, Path: "path1", Name: "path1"},
				{ID: 2, Path: "path2", Name: "path2"},
			},
			wantActions: []*fsa.Action{
				assetAction.scanRunning("default-workspace", []*model.Asset{
					{ID: 1, Path: "path1", Name: "path1"},
				}, 1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dispatcher := testutil.NewMockDispatcher(t)
			ctrl := gomock.NewController(t)
			assetUseCase := mock_usecase.NewMockAsset(ctrl)
			payload := decodeAssetScanRequestPayload(t, tt.args.action.Payload)
			assetUseCase.EXPECT().ListAsyncByQueries(
				gomock.Any(),
				gomock.Eq(payload.WorkSpaceName),
				gomock.Eq(payload.Queries),
			).Return(newAssetChannel(tt.existAssets), nil)
			handler := newAssetHandlerCreator(assetUseCase).Scan()
			if err := handler.Do(tt.args.action, dispatcher.Dispatch); (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
			dispatcher.Test(tt.wantActions)
		})
	}
}

func Test_assetScanHandler_Do_Twice(t *testing.T) {
	type args struct {
		action *fsa.Action
	}
	tests := []struct {
		name         string
		args1        args
		args2        args
		existAssets  []*model.Asset
		wantActions1 []*fsa.Action
		wantActions2 []*fsa.Action
		wantErr1     bool
		wantErr2     bool
	}{
		{
			name: "",
			args1: args{
				&fsa.Action{
					Type: "",
					Payload: &assetScanRequestPayload{
						WsPayload:  WsPayload{"default-workspace"},
						RequestNum: 1,
						Queries:    nil,
						Reset:      false,
					},
					Error: false,
					Meta:  nil,
				},
			},
			args2: args{
				&fsa.Action{
					Type: "",
					Payload: &assetScanRequestPayload{
						WsPayload:  WsPayload{"default-workspace"},
						RequestNum: 1,
						Queries:    nil,
						Reset:      false,
					},
					Error: false,
					Meta:  nil,
				},
			},
			existAssets: []*model.Asset{
				{ID: 1, Path: "path1", Name: "path1"},
				{ID: 2, Path: "path2", Name: "path2"},
			},
			wantActions1: []*fsa.Action{
				assetAction.scanRunning("default-workspace", []*model.Asset{
					{ID: 1, Path: "path1", Name: "path1"},
				}, 1),
			},
			wantActions2: []*fsa.Action{
				assetAction.scanRunning("default-workspace", []*model.Asset{
					{ID: 2, Path: "path2", Name: "path2"},
				}, 2),
			},
			wantErr1: false,
			wantErr2: false,
		},
		// TODO: scanのcancel、queries、2回目のscanでcntが引き継がれているかなどをテスト
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dispatcher := testutil.NewMockDispatcher(t)
			ctrl := gomock.NewController(t)
			assetUseCase := mock_usecase.NewMockAsset(ctrl)

			payload1 := decodeAssetScanRequestPayload(t, tt.args1.action.Payload)
			assetUseCase.EXPECT().ListAsyncByQueries(
				gomock.Any(),
				gomock.Eq(payload1.WorkSpaceName),
				gomock.Eq(payload1.Queries),
			).Return(newAssetChannel(tt.existAssets), nil)

			payload2 := decodeAssetScanRequestPayload(t, tt.args1.action.Payload)
			assetUseCase.EXPECT().ListAsyncByQueries(
				gomock.Any(),
				gomock.Eq(payload2.WorkSpaceName),
				gomock.Eq(payload2.Queries),
			).Return(newAssetChannel(tt.existAssets), nil)

			handler := newAssetHandlerCreator(assetUseCase).Scan()

			if err := handler.Do(tt.args1.action, dispatcher.Dispatch); (err != nil) != tt.wantErr1 {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr1)
			}
			dispatcher.TestAndClean(tt.wantActions1)

			if err := handler.Do(tt.args2.action, dispatcher.Dispatch); (err != nil) != tt.wantErr2 {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr2)
			}
			dispatcher.TestAndClean(tt.wantActions2)
		})
	}
}
