package syncer

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/xylonx/x-telebot/util"
)

var (
	ErrorNilOption = errors.New("the option is not set")
)

type S3Syncer struct {
	client *s3.Client
	bucket string
}

var _ Synchronizer = &S3Syncer{}

type Option struct {
	Endpoint     string
	AccessID     string
	AccessSecret string
	BucketName   string
	Region       string
}

// customResolver := aws.EndpointResolverW(func(service, region string) (aws.Endpoint, error) {
// 	return aws.Endpoint{
// 		URL:           opt.Endpoint,
// 		SigningRegion: opt.Region,
// 	}, nil
// })

// cfg, err := config.LoadDefaultConfig(context.TODO(),
// 	config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
// 		opt.AccessID,
// 		opt.SecretKey,
// 		""),
// 	),
// 	config.WithRegion(client.Region),
// 	config.WithEndpointResolver(customResolver),
// )
// if err != nil {
// 	return nil, err
// }
// s3client := s3.NewFromConfig(cfg, func(o *s3.Options) {
// 	o.UsePathStyle = (client.HostStyle == pkg.PathStyle)
// })

func NewS3Syncer(opt *Option) (Synchronizer, error) {
	if opt == nil {
		return nil, ErrorNilOption
	}

	ep := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{URL: opt.Endpoint, SigningRegion: region}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(opt.AccessID, opt.AccessSecret, "")),
		config.WithRegion(opt.Region),
		config.WithEndpointResolverWithOptions(ep),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3Syncer{client: client, bucket: opt.BucketName}, nil
}

func (s *S3Syncer) Persistent(ctx context.Context, key string, data io.Reader) (string, error) {
	uploader := manager.NewUploader(s.client)
	resp, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:   &s.bucket,
		Key:      aws.String(key),
		Body:     data,
		Metadata: util.MetadataFromContext(ctx),
	})
	if err != nil {
		return "", err
	}

	return resp.Location, nil
}

func (s *S3Syncer) PickOne(ctx context.Context) (string, error) {
	return "", nil
}
