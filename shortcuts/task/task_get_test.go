// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package task

import (
	"strings"
	"testing"

	"github.com/larksuite/cli/internal/httpmock"
)

func TestGetTask(t *testing.T) {
	tests := []struct {
		name           string
		taskId         string
		formatFlag     string
		expectedOutput []string
	}{
		{
			name:       "pretty format",
			taskId:     "task-123",
			formatFlag: "pretty",
			expectedOutput: []string{
				"📋 Task Detail",
				"task-123",
				"Buy groceries",
			},
		},
		{
			name:       "json format",
			taskId:     "task-456",
			formatFlag: "json",
			expectedOutput: []string{
				`"guid": "task-456"`,
				`"summary": "Review PR"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, stdout, _, reg := taskShortcutTestFactory(t)
			warmTenantToken(t, f, reg)

			reg.Register(&httpmock.Stub{
				Method: "GET",
				URL:    "/open-apis/task/v2/tasks/" + tt.taskId,
				Body: map[string]interface{}{
					"code": 0, "msg": "success",
					"data": map[string]interface{}{
						"task": map[string]interface{}{
							"guid":        tt.taskId,
							"summary":     map[string]string{"task-123": "Buy groceries", "task-456": "Review PR"}[tt.taskId],
							"description": "task description here",
							"status":      "in_progress",
							"created_at":  "1775174400000",
							"due": map[string]interface{}{
								"timestamp": "1775347200000",
							},
							"url": "https://example.com/" + tt.taskId,
						},
					},
				},
			})

			err := runMountedTaskShortcut(t, GetTask, []string{"+get", "--task-id", tt.taskId, "--format", tt.formatFlag, "--as", "bot"}, f, stdout)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			out := stdout.String()
			outNorm := strings.ReplaceAll(out, `":"`, `": "`)

			for _, expected := range tt.expectedOutput {
				if !strings.Contains(outNorm, expected) && !strings.Contains(out, expected) {
					t.Errorf("output missing expected string (%s), got: %s", expected, out)
				}
			}
		})
	}
}
