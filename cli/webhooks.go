package cli

import (
	"encoding/json"

	mfxsdk "github.com/MainfluxLabs/mainflux/pkg/sdk/go"
	"github.com/spf13/cobra"
)

var cmdWebhooks = []cobra.Command{
	{
		Use:   "create <JSON_webhooks> <thing_id> <user_token>",
		Short: "Create webhooks",
		Long:  `Create webhooks for certain thing.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 3 {
				logUsage(cmd.Use)
				return
			}
			var webhooks []mfxsdk.Webhook
			if err := json.Unmarshal([]byte(args[0]), &webhooks); err != nil {
				logError(err)
				return
			}
			whs, err := sdk.CreateWebhooks(webhooks, args[1], args[2])
			if err != nil {
				logError(err)
				return
			}
			logJSON(whs)
		},
	},
	{
		Use:   "get <by-thing | by-group | by-id> <id> <user_token>",
		Short: "Get webhooks",
		Long: `Get all webhooks by group or get webhook by id:
		by-group - lists all webhooks by group by provided <id>
		by-id - shows webhook by provided <id>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 3 {
				logUsage(cmd.Use)
				return
			}

			switch args[0] {
			case "by-group":
				l, err := sdk.ListWebhooksByGroup(args[1], args[2])
				if err != nil {
					logError(err)
					return
				}
				logJSON(l)
			case "by-thing":
				l, err := sdk.ListWebhooksByThing(args[1], args[2])
				if err != nil {
					logError(err)
					return
				}
				logJSON(l)
			case "by-id":
				w, err := sdk.GetWebhook(args[1], args[2])
				if err != nil {
					logError(err)
					return
				}

				logJSON(w)
			default:
				return
			}
		},
	},
	{
		Use:   "update <JSON_webhook> <webhook_id> <user_token>",
		Short: "Update webhook",
		Long:  `Update webhook record`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 3 {
				logUsage(cmd.Use)
				return
			}

			var wh mfxsdk.Webhook
			if err := json.Unmarshal([]byte(args[0]), &wh); err != nil {
				logError(err)
				return
			}

			if err := sdk.UpdateWebhook(wh, args[1], args[2]); err != nil {
				logError(err)
				return
			}

			logOK()
		},
	},
	{
		Use:   "delete <JSON_ids> <user_token>",
		Short: "Delete webhooks",
		Long:  `Delete webhooks by provided IDs`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 3 {
				logUsage(cmd.Use)
				return
			}
			var ids []string
			if err := json.Unmarshal([]byte(args[0]), &ids); err != nil {
				logError(err)
				return
			}
			if err := sdk.DeleteWebhooks(ids, args[1]); err != nil {
				logError(err)
				return
			}
			logOK()
		},
	},
}

// NewWebhooksCmd returns users command.
func NewWebhooksCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "webhooks [create | get | delete]",
		Short: "Webhooks management",
		Long:  `Webhooks management: create, update, delete and get webhooks`,
	}

	for i := range cmdWebhooks {
		cmd.AddCommand(&cmdWebhooks[i])
	}

	return &cmd
}
