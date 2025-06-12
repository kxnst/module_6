package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	sws "github.com/gofiber/websocket/v2"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"guitar_processor/cmd/server/handlers"
	"guitar_processor/cmd/server/middlewares"
	"guitar_processor/cmd/server/request"
	"guitar_processor/cmd/server/utils"
	pctx "guitar_processor/internal/context"
	"guitar_processor/internal/effect/dsp"
	"guitar_processor/internal/effect/service"
	"guitar_processor/internal/entity"
	"guitar_processor/internal/repository"
	"net/http"
	"testing"
	"time"
)

type DummyStream struct{}

type MockedAuthService struct {
	utils.AuthService
}

func (m *MockedAuthService) ValidateTokenAndGetUser(_ string) (*entity.User, error) {
	return &entity.User{Name: "test", Login: "test"}, nil
}

func (m *MockedAuthService) GenerateToken(_ string) (string, error) {
	return "dummy-token", nil
}

func (d *DummyStream) Start() error { return nil }
func (d *DummyStream) Stop() error  { return nil }
func (d *DummyStream) Close() error { return nil }

type DummyEffect struct {
	dsp.ParameterStorage
	val dsp.Parameter
}

func (e *DummyEffect) GetSlug() string { return "dummy" }
func (e *DummyEffect) Process(samples []float32, _ *pctx.ProcessingContext) {
	for i := range samples {
		samples[i] = e.val.GetValue()
	}
}

func (e *DummyEffect) GetParameterDefinitions() []dsp.ParameterDefinition {
	return []dsp.ParameterDefinition{
		&dsp.SliderParameterDefinition{Name: "Value", Min: -1, Max: 1, Default: 0.75},
	}
}

func (e *DummyEffect) Init(_ *pctx.ProcessingContext) {
	for _, param := range e.Parameters {
		if def, ok := param.GetDefinition().(*dsp.SliderParameterDefinition); ok && def.Name == "Value" {
			e.val = param
		}
	}
}
func (e *DummyEffect) GetInfo(_ *entity.Effect) dsp.EffectInfo {
	return dsp.EffectInfo{
		Slug:                 e.GetSlug(),
		ParameterDefinitions: e.GetParameterDefinitions(),
	}
}

type MockAudioService struct {
	reg *service.EffectsRegistry
}

func (m *MockAudioService) GetDevices() ([]*dsp.Device, error) { return nil, nil }
func (m *MockAudioService) Start(*request.DeviceInfoRequest) (service.Stream, error) {
	return &DummyStream{}, nil
}
func (m *MockAudioService) GetEffectsInfo() []*dsp.EffectInfo {
	var result []*dsp.EffectInfo
	for _, e := range m.reg.GetEffects() {
		result = append(
			result,
			&dsp.EffectInfo{
				Slug:                 e.GetSlug(),
				ParameterDefinitions: e.GetParameterDefinitions(),
				Name:                 e.GetSlug(),
				DspType:              e.GetSlug(),
			},
		)
	}

	return result
}
func (m *MockAudioService) AddEffectToChain(slug string) (dsp.Effect, error) {
	ctx := &pctx.ProcessingContext{}
	return m.reg.CreateBySlug(slug, ctx)
}
func (m *MockAudioService) SetEffectParameter(_ string, _ float32) error {
	return nil
}
func (m *MockAudioService) Terminate() {}

func setupApp(as service.AudioService) *fiber.App {
	app := fiber.New()

	auth := &MockedAuthService{}

	authMiddleware := middlewares.NewAuthMiddleware(auth)
	authClosure := func(c *fiber.Ctx) error { return authMiddleware.RequireAuth(c) }
	handler := handlers.NewWebSocketHandler(as, auth)

	app.Post("/ws", authClosure, func(c *fiber.Ctx) error {
		return handler.HandleInit(c)
	})

	app.Get("/ws", sws.New(func(c *sws.Conn) {
		_ = handler.HandleStream(c)
	}))

	return app

}

func TestE2E_AudioServiceWithDummyEffect(t *testing.T) {
	fxApp := fx.New(
		fx.Provide(func() service.AudioService {
			reg := service.NewEffectsRegistry(&repository.EffectRepository{})

			return &MockAudioService{reg: reg}
		}),
		fx.Invoke(func(as service.AudioService) {
			app := setupApp(as)

			go func() {
				if err := app.Listen(":3001"); err != nil {
					panic(err)
				}
			}()
			time.Sleep(500 * time.Millisecond)
			authReq := request.AuthRequest{Login: "test", Password: "<PASSWORD>"}

			body, _ := json.Marshal(authReq)

			req, err := http.NewRequest("POST", "http://localhost:3001/auth", bytes.NewReader(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			payload := request.DeviceInfoRequest{
				Device1Index: 1,
				Device2Index: 2,
				SampleRate:   44100,
				BufferSize:   128,
			}
			body, _ = json.Marshal(payload)

			req, err = http.NewRequest("POST", "http://localhost:3001/ws", bytes.NewReader(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "dummy-token")

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			ws, _, err := connectWebSocket("ws://localhost:3001/ws?token=dummy-token")
			require.NoError(t, err)
			defer ws.Close()

			err = ws.WriteJSON(map[string]string{"type": "effects_list"})
			require.NoError(t, err)

			_, msg, err := ws.ReadMessage()
			require.NoError(t, err)

			var response map[string]interface{}
			err = json.Unmarshal(msg, &response)
			require.NoError(t, err)
			require.Equal(t, "effect_list_response", response["type"])

			err = ws.WriteJSON(map[string]interface{}{
				"type": "add_effect",
				"data": map[string]interface{}{
					"slug": "distortion",
				},
			})
			require.NoError(t, err)

			_, msg, err = ws.ReadMessage()
			require.NoError(t, err)

			var addEffectResp map[string]interface{}
			err = json.Unmarshal(msg, &addEffectResp)
			require.NoError(t, err)
			require.Equal(t, "effect_added", addEffectResp["type"])
			var gainUUID string
			params := addEffectResp["parameters"].([]interface{})
			for _, p := range params {
				param := p.(map[string]interface{})
				def := param["definition"].(map[string]interface{})
				if def["name"] == "Gain" {
					gainUUID = param["uuid"].(string)
					break
				}
			}
			require.NotEmpty(t, gainUUID)

			err = ws.WriteJSON(map[string]interface{}{
				"type": "set_param",
				"data": map[string]interface{}{
					"uuid":  gainUUID,
					"value": 5,
				},
			})
			require.NoError(t, err)

			_, msg, err = ws.ReadMessage()
			require.NoError(t, err)

			var paramResp map[string]interface{}
			err = json.Unmarshal(msg, &paramResp)
			require.NoError(t, err)
			require.Equal(t, "param_updated", paramResp["type"])
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	require.NoError(t, fxApp.Start(ctx))
	defer fxApp.Stop(ctx)
}

func connectWebSocket(url string) (*websocket.Conn, *http.Response, error) {
	dialer := websocket.Dialer{}
	conn, resp, err := dialer.Dial(url, nil)
	return conn, resp, err
}
