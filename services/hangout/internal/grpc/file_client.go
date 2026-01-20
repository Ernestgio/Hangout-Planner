package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type FileService interface {
	GenerateUploadURLs(ctx context.Context, baseStoragePath string, files []*filepb.FileUploadIntent) (*filepb.GenerateUploadURLsResponse, error)
	ConfirmUpload(ctx context.Context, fileIDs []string) error
	GetFileByMemoryID(ctx context.Context, memoryID string) (*filepb.FileWithURL, error)
	GetFilesByMemoryIDs(ctx context.Context, memoryIDs []string) (map[string]*filepb.FileWithURL, error)
	DeleteFile(ctx context.Context, memoryID string) error
	Close() error
}

type fileServiceClient struct {
	client filepb.FileServiceClient
	conn   *grpc.ClientConn
}

func NewFileServiceClient(cfg *config.GRPCClientConfig) (FileService, error) {
	var opts []grpc.DialOption

	if cfg.MTLSEnabled {
		tlsConfig, err := loadClientTLSConfig(cfg)
		if err != nil {
			return nil, err
		}
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(cfg.FileServiceURL, opts...)
	if err != nil {
		return nil, apperrors.ErrConnectFileService
	}

	return &fileServiceClient{
		client: filepb.NewFileServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *fileServiceClient) GenerateUploadURLs(ctx context.Context, baseStoragePath string, files []*filepb.FileUploadIntent) (*filepb.GenerateUploadURLsResponse, error) {
	req := &filepb.GenerateUploadURLsRequest{
		BaseStoragePath: baseStoragePath,
		Files:           files,
	}
	return c.client.GenerateUploadURLs(ctx, req)
}

func (c *fileServiceClient) ConfirmUpload(ctx context.Context, fileIDs []string) error {
	req := &filepb.ConfirmUploadRequest{
		FileIds: fileIDs,
	}
	_, err := c.client.ConfirmUpload(ctx, req)
	return err
}

func (c *fileServiceClient) GetFileByMemoryID(ctx context.Context, memoryID string) (*filepb.FileWithURL, error) {
	req := &filepb.GetFileByMemoryIDRequest{
		MemoryId: memoryID,
	}
	resp, err := c.client.GetFileByMemoryID(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.File, nil
}

func (c *fileServiceClient) GetFilesByMemoryIDs(ctx context.Context, memoryIDs []string) (map[string]*filepb.FileWithURL, error) {
	req := &filepb.GetFilesByMemoryIDsRequest{
		MemoryIds: memoryIDs,
	}
	resp, err := c.client.GetFilesByMemoryIDs(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Files, nil
}

func (c *fileServiceClient) DeleteFile(ctx context.Context, memoryID string) error {
	req := &filepb.DeleteFileRequest{
		MemoryId: memoryID,
	}
	_, err := c.client.DeleteFile(ctx, req)
	return err
}

func (c *fileServiceClient) Close() error {
	return c.conn.Close()
}

func loadClientTLSConfig(cfg *config.GRPCClientConfig) (*tls.Config, error) {
	clientCert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate (cert: %s, key: %s): %w", cfg.CertFile, cfg.KeyFile, err)
	}

	caCert, err := os.ReadFile(cfg.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load CA certificate (%s): %w", cfg.CAFile, err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate to pool (%s)", cfg.CAFile)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
		ServerName:   "file",
		MinVersion:   tls.VersionTLS12,
	}

	return tlsConfig, nil
}
