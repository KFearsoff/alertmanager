// Copyright 2022 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telegram

import (
	"fmt"
	"net/url"
	//"os"
	"testing"

	"github.com/go-kit/log"
	commoncfg "github.com/prometheus/common/config"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/notify/test"
)

func TestTelegramUnmarshal(t *testing.T) {
	in := `
route:
  receiver: test
receivers:
- name: test
  telegram_configs:
  - chat_id: 1234
    bot_token: secret
`
	var c config.Config
	err := yaml.Unmarshal([]byte(in), &c)
	require.NoError(t, err)

	require.Len(t, c.Receivers, 1)
	require.Len(t, c.Receivers[0].TelegramConfigs, 1)

	require.Equal(t, "https://api.telegram.org", c.Receivers[0].TelegramConfigs[0].APIUrl.String())
	require.Equal(t, config.Secret("secret"), c.Receivers[0].TelegramConfigs[0].BotToken)
	require.Equal(t, int64(1234), c.Receivers[0].TelegramConfigs[0].ChatID)
	require.Equal(t, "HTML", c.Receivers[0].TelegramConfigs[0].ParseMode)
}

func TestTelegramRetry(t *testing.T) {
	// Fake url for testing purposes
	fakeURL := config.URL{
		URL: &url.URL{
			Scheme: "https",
			Host:   "FAKE_API",
		},
	}
	notifier, err := New(
		&config.TelegramConfig{
			HTTPConfig: &commoncfg.HTTPClientConfig{},
			APIUrl:     &fakeURL,
		},
		test.CreateTmpl(t),
		log.NewNopLogger(),
	)
	require.NoError(t, err)

	for statusCode, expected := range test.RetryTests(test.DefaultRetryCodes()) {
		actual, _ := notifier.retrier.Check(statusCode, nil)
		require.Equal(t, expected, actual, fmt.Sprintf("error on status %d", statusCode))
	}
}

/*
func TestTelegramRedactedBotToken(t *testing.T) {
	ctx, u, fn := test.GetContextWithCancelingURL()
	defer fn()

	botToken := "test"

	notifier, err := New(
		&config.TelegramConfig{
			APIUrl: &config.URL{URL: u},
			HTTPConfig: &commoncfg.HTTPClientConfig{},
			ChatID:   1234,
			BotToken: config.Secret(botToken),
		},
		test.CreateTmpl(t),
		log.NewNopLogger(),
	)
	require.NoError(t, err)

	test.AssertNotifyLeaksNoSecret(ctx, t, notifier, botToken)
}

func TestTelegramBotTokenFromFile(t *testing.T) {
	ctx, u, fn := test.GetContextWithCancelingURL()
	defer fn()

	botToken := "test"

	f, err := os.CreateTemp("", "telegram_test")
	require.NoError(t, err, "creating temp file failed")
	_, err = f.WriteString(botToken)
	require.NoError(t, err, "writing to temp file failed")

	notifier, err := New(
		&config.TelegramConfig{
			APIUrl: &config.URL{URL: u},
			HTTPConfig: &commoncfg.HTTPClientConfig{},
			ChatID:   1234,
			BotTokenFile: f.Name(),
		},
		test.CreateTmpl(t),
		log.NewNopLogger(),
	)
	require.NoError(t, err)

	test.AssertNotifyLeaksNoSecret(ctx, t, notifier, botToken)
}
*/