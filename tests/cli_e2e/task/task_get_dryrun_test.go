// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package task

import (
	"context"
	"testing"
	"time"

	clie2e "github.com/larksuite/cli/tests/cli_e2e"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestTask_GetDryRun(t *testing.T) {
	t.Setenv("LARKSUITE_CLI_CONFIG_DIR", t.TempDir())
	t.Setenv("LARKSUITE_CLI_APP_ID", "task_dryrun_test")
	t.Setenv("LARKSUITE_CLI_APP_SECRET", "task_dryrun_secret")
	t.Setenv("LARKSUITE_CLI_BRAND", "feishu")

	tests := []struct {
		name       string
		args       []string
		wantMethod string
		wantURL    string
		wantTaskID string
	}{
		{
			name: "get task by guid",
			args: []string{
				"task", "+get",
				"--task-id", "task-guid-123",
				"--dry-run",
			},
			wantMethod: "GET",
			wantURL:    "/open-apis/task/v2/tasks/task-guid-123",
			wantTaskID: "task-guid-123",
		},
		{
			name: "get task by applink URL resolves to guid",
			args: []string{
				"task", "+get",
				"--task-id", "https://applink.feishu.cn/client/todo/task?guid=task-from-url",
				"--dry-run",
			},
			wantMethod: "GET",
			wantURL:    "/open-apis/task/v2/tasks/task-from-url",
			wantTaskID: "task-from-url",
		},
	}

	for _, temp := range tests {
		tt := temp
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			t.Cleanup(cancel)

			result, err := clie2e.RunCmd(ctx, clie2e.Request{
				Args:      tt.args,
				DefaultAs: "bot",
			})
			require.NoError(t, err)
			result.AssertExitCode(t, 0)

			out := result.Stdout
			if count := gjson.Get(out, "api.#").Int(); count != 1 {
				t.Fatalf("expected 1 API call, got %d\nstdout:\n%s", count, out)
			}
			if method := gjson.Get(out, "api.0.method").String(); method != tt.wantMethod {
				t.Fatalf("api[0].method = %q, want %q\nstdout:\n%s", method, tt.wantMethod, out)
			}
			if url := gjson.Get(out, "api.0.url").String(); url != tt.wantURL {
				t.Fatalf("api[0].url = %q, want %q\nstdout:\n%s", url, tt.wantURL, out)
			}
			if got := gjson.Get(out, "api.0.params.user_id_type").String(); got != "open_id" {
				t.Fatalf("api[0].params.user_id_type = %q, want open_id\nstdout:\n%s", got, out)
			}
		})
	}
}
