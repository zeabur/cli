package send

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	apiKey      string
	from        string
	to          []string
	cc          []string
	bcc         []string
	replyTo     []string
	subject     string
	html        string
	text        string
	scheduledAt string // RFC3339; if set → schedule
	batchFile   string // JSON file path; if set → batch
}

func NewCmdSend(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send an email (single, scheduled, or batch)",
		Long: `Send an email via Zeabur Email.

Modes:
  - Default:        sends immediately (POST /emails)
  - --scheduled-at: schedules the email (POST /emails/schedule)
  - --batch-file:   sends a batch from a JSON file (POST /emails/batch, max 100)

Authentication:
  Pass your Z-Send API key via --api-key or the ZSEND_API_KEY environment variable.

Batch file format (JSON array):
  [
    {"from":"a@d.com","to":["b@e.com"],"subject":"Hi","html":"<p>Hello</p>"},
    ...
  ]`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSend(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.apiKey, "api-key", "", "Z-Send API key (or set ZSEND_API_KEY)")
	cmd.Flags().StringVar(&opts.from, "from", "", "Sender address (required unless --batch-file)")
	cmd.Flags().StringArrayVar(&opts.to, "to", nil, "Recipient(s), can be repeated")
	cmd.Flags().StringArrayVar(&opts.cc, "cc", nil, "CC recipient(s)")
	cmd.Flags().StringArrayVar(&opts.bcc, "bcc", nil, "BCC recipient(s)")
	cmd.Flags().StringArrayVar(&opts.replyTo, "reply-to", nil, "Reply-To address(es)")
	cmd.Flags().StringVar(&opts.subject, "subject", "", "Email subject (required unless --batch-file)")
	cmd.Flags().StringVar(&opts.html, "html", "", "HTML body")
	cmd.Flags().StringVar(&opts.text, "text", "", "Plain-text body")
	cmd.Flags().StringVar(&opts.scheduledAt, "scheduled-at", "", "Schedule time in RFC3339 format (e.g. 2026-04-01T10:00:00Z)")
	cmd.Flags().StringVar(&opts.batchFile, "batch-file", "", "Path to JSON file containing an array of email objects")

	return cmd
}

func runSend(f *cmdutil.Factory, opts Options) error {
	if opts.apiKey == "" {
		opts.apiKey = os.Getenv("ZSEND_API_KEY")
	}
	if opts.apiKey == "" {
		return fmt.Errorf("Z-Send API key is required (--api-key or ZSEND_API_KEY)")
	}
	if !strings.HasPrefix(opts.apiKey, "zs_") {
		return fmt.Errorf("invalid API key format: must start with zs_")
	}

	if opts.batchFile != "" {
		return runBatch(f, opts)
	}
	if opts.scheduledAt != "" {
		return runSchedule(f, opts)
	}
	return runSingle(f, opts)
}

func buildSingleReq(opts Options) (model.ZSendSendEmailRequest, error) {
	if opts.from == "" {
		return model.ZSendSendEmailRequest{}, fmt.Errorf("--from is required")
	}
	if len(opts.to) == 0 {
		return model.ZSendSendEmailRequest{}, fmt.Errorf("--to is required")
	}
	if opts.subject == "" {
		return model.ZSendSendEmailRequest{}, fmt.Errorf("--subject is required")
	}
	if opts.html == "" && opts.text == "" {
		return model.ZSendSendEmailRequest{}, fmt.Errorf("either --html or --text is required")
	}
	return model.ZSendSendEmailRequest{
		From:    opts.from,
		To:      opts.to,
		Cc:      opts.cc,
		Bcc:     opts.bcc,
		ReplyTo: opts.replyTo,
		Subject: opts.subject,
		HTML:    opts.html,
		Text:    opts.text,
	}, nil
}

func runSingle(f *cmdutil.Factory, opts Options) error {
	req, err := buildSingleReq(opts)
	if err != nil {
		return err
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Sending email..."),
	)
	s.Start()
	reply, err := f.ApiClient.SendZSendEmail(context.Background(), opts.apiKey, req)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(reply)
	}
	fmt.Printf("Email sent successfully.\nID:     %s\nStatus: %s\n", reply.ID, reply.Status)
	return nil
}

func runSchedule(f *cmdutil.Factory, opts Options) error {
	base, err := buildSingleReq(opts)
	if err != nil {
		return err
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Scheduling email..."),
	)
	s.Start()
	reply, err := f.ApiClient.ScheduleZSendEmail(context.Background(), opts.apiKey, model.ZSendScheduleEmailRequest{
		ZSendSendEmailRequest: base,
		ScheduledAt:           opts.scheduledAt,
	})
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(reply)
	}
	fmt.Printf("Email scheduled successfully.\nID:     %s\nStatus: %s\n", reply.ID, reply.Status)
	return nil
}

func runBatch(f *cmdutil.Factory, opts Options) error {
	data, err := os.ReadFile(opts.batchFile)
	if err != nil {
		return fmt.Errorf("read batch file: %w", err)
	}

	var emails []model.ZSendSendEmailRequest
	if err := json.Unmarshal(data, &emails); err != nil {
		return fmt.Errorf("parse batch file: %w", err)
	}
	if len(emails) == 0 {
		return fmt.Errorf("batch file contains no emails")
	}
	if len(emails) > 100 {
		return fmt.Errorf("batch file contains %d emails; maximum is 100", len(emails))
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Sending batch of %d emails...", len(emails))),
	)
	s.Start()
	reply, err := f.ApiClient.SendZSendBatchEmail(context.Background(), opts.apiKey, model.ZSendBatchEmailRequest{
		Emails: emails,
	})
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(reply)
	}
	fmt.Printf("Batch submitted successfully.\nJob ID: %s\nStatus: %s\nTotal:  %d\n", reply.JobID, reply.Status, reply.TotalCount)
	return nil
}
