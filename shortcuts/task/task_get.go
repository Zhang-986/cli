// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package task

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"

	"github.com/larksuite/cli/shortcuts/common"
)

var GetTask = common.Shortcut{
	Service:     "task",
	Command:     "+get",
	Description: "get a single task by id",
	Risk:        "read",
	Scopes:      []string{"task:task:read"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,

	Flags: []common.Flag{
		{Name: "task-id", Desc: "task id (guid or applink URL)", Required: true},
	},

	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		taskId := url.PathEscape(extractTaskGuid(runtime.Str("task-id")))
		return common.NewDryRunAPI().
			GET("/open-apis/task/v2/tasks/" + taskId).
			Params(map[string]interface{}{"user_id_type": "open_id"})
	},

	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		taskId := extractTaskGuid(runtime.Str("task-id"))

		task, err := getTaskDetail(runtime, taskId)
		if err != nil {
			return err
		}

		guid, _ := task["guid"].(string)
		summary, _ := task["summary"].(string)
		urlVal, _ := task["url"].(string)
		urlVal = truncateTaskURL(urlVal)

		outData := map[string]interface{}{
			"guid":    guid,
			"summary": summary,
			"url":     urlVal,
		}
		if description, ok := task["description"].(string); ok {
			outData["description"] = description
		}
		if dueObj, ok := task["due"].(map[string]interface{}); ok {
			if tsStr, ok := dueObj["timestamp"].(string); ok {
				if ts, pErr := strconv.ParseInt(tsStr, 10, 64); pErr == nil {
					outData["due_at"] = time.UnixMilli(ts).Local().Format(time.RFC3339)
				}
			}
		}
		if createdAtStr, ok := task["created_at"].(string); ok {
			if ts, pErr := strconv.ParseInt(createdAtStr, 10, 64); pErr == nil {
				outData["created_at"] = time.UnixMilli(ts).Local().Format(time.RFC3339)
			}
		}
		if status, ok := task["status"].(string); ok {
			outData["status"] = status
		}

		runtime.OutFormat(outData, nil, func(w io.Writer) {
			fmt.Fprintf(w, "📋 Task Detail\n")
			if guid != "" {
				fmt.Fprintf(w, "  ID: %s\n", guid)
			}
			if summary != "" {
				fmt.Fprintf(w, "  Summary: %s\n", summary)
			}
			if desc, ok := task["description"].(string); ok && desc != "" {
				fmt.Fprintf(w, "  Description: %s\n", desc)
			}
			if status, ok := task["status"].(string); ok && status != "" {
				fmt.Fprintf(w, "  Status: %s\n", status)
			}
			if dueAt, ok := outData["due_at"].(string); ok {
				fmt.Fprintf(w, "  Due: %s\n", dueAt)
			}
			if createdAt, ok := outData["created_at"].(string); ok {
				fmt.Fprintf(w, "  Created: %s\n", createdAt)
			}
			if urlVal != "" {
				fmt.Fprintf(w, "  URL: %s\n", urlVal)
			}
		})
		return nil
	},
}
