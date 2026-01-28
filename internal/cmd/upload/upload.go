package deploy

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/constant"
)

type Options struct{}

func NewCmdUpload(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "upload",
		Short:   "Upload local project to Zeabur",
		PreRunE: util.NeedProjectContextWhenNonInteractive(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpload(f, opts)
		},
	}

	return cmd
}

func runUpload(f *cmdutil.Factory, opts *Options) error {
	var err error

	bytes, _, err := util.PackZip()
	if err != nil {
		return fmt.Errorf("packing zip: %w", err)
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Uploading codes to Zeabur ..."),
	)
	s.Start()

	uploadID, err := UploadZipToService(context.Background(), bytes)
	if err != nil {
		return err
	}
	s.Stop()

	fmt.Println(constant.ZeaburDashURL + "/uploads/" + uploadID)
	return nil
}

func UploadZipToService(ctx context.Context, zipBytes []byte) (string, error) {
	// Step 1: Calculate SHA256 hash of content
	h := sha256.New()
	if _, err := h.Write(zipBytes); err != nil {
		return "", fmt.Errorf("failed to calculate content hash: %w", err)
	}
	contentHash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Step 2: Create upload session
	createUploadReq := struct {
		ContentHash          string `json:"content_hash"`
		ContentHashAlgorithm string `json:"content_hash_algorithm"`
		ContentLength        int64  `json:"content_length"`
	}{
		ContentHash:          contentHash,
		ContentHashAlgorithm: "sha256",
		ContentLength:        int64(len(zipBytes)),
	}

	createUploadBody, err := json.Marshal(createUploadReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal create upload request: %w", err)
	}

	createUploadResp, err := http.NewRequestWithContext(ctx, "POST", constant.ZeaburServerURL+"/v2/upload", bytes.NewReader(createUploadBody))
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}

	token := viper.GetString("token")
	createUploadResp.Header.Set("Content-Type", "application/json")
	createUploadResp.Header.Set("Cookie", "token="+token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(createUploadResp)
	if err != nil {
		return "", fmt.Errorf("failed to create upload session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create upload session: status code %d", resp.StatusCode)
	}

	var uploadSession struct {
		PresignHeader struct {
			ContentType string `json:"Content-Type"`
		} `json:"presign_header"`
		PresignMethod string `json:"presign_method"`
		PresignURL    string `json:"presign_url"`
		UploadID      string `json:"upload_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&uploadSession); err != nil {
		return "", fmt.Errorf("failed to decode upload session response: %w", err)
	}

	// Step 3: Upload file to S3
	uploadReq, err := http.NewRequestWithContext(ctx, uploadSession.PresignMethod, uploadSession.PresignURL, bytes.NewReader(zipBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create S3 upload request: %w", err)
	}

	uploadReq.Header.Set("Content-Type", uploadSession.PresignHeader.ContentType)
	uploadReq.Header.Set("Content-Length", strconv.FormatInt(int64(len(zipBytes)), 10))

	uploadResp, err := client.Do(uploadReq)
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to upload to S3: status code %d", uploadResp.StatusCode)
	}

	// Step 4: Prepare upload for deployment
	prepareReq := struct {
		UploadType string `json:"upload_type"`
	}{
		UploadType: "new_project",
	}

	prepareBody, err := json.Marshal(prepareReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal prepare request: %w", err)
	}

	prepareResp, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/v2/upload/%s/prepare", constant.ZeaburServerURL, uploadSession.UploadID),
		bytes.NewReader(prepareBody))
	if err != nil {
		return "", fmt.Errorf("failed to create prepare request: %w", err)
	}

	prepareResp.Header.Set("Content-Type", "application/json")
	prepareResp.Header.Set("Cookie", "token="+token)

	resp, err = client.Do(prepareResp)
	if err != nil {
		return "", fmt.Errorf("failed to prepare upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to prepare upload: status code %d", resp.StatusCode)
	}

	return uploadSession.UploadID, nil
}
